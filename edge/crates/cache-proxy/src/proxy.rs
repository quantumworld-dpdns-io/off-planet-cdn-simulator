use std::sync::Arc;
use axum::{
    extract::{Path, State},
    http::{HeaderMap, StatusCode},
    response::{IntoResponse, Response},
};
use serde::Deserialize;

use crate::{config::Config, fetch, headers};

#[derive(Deserialize)]
struct CacheStatusResponse {
    objects: Vec<CacheObjectInfo>,
}

#[derive(Deserialize)]
struct CacheObjectInfo {
    object_id: String,
    source_url: Option<String>,
    priority: Option<String>,
    #[allow(dead_code)]
    size_bytes: Option<u64>,
    #[allow(dead_code)]
    pinned: Option<bool>,
}

pub struct AppState {
    pub config: Arc<Config>,
    pub http: reqwest::Client,
}

pub async fn proxy_handler(
    State(state): State<Arc<AppState>>,
    Path(object_id): Path<String>,
) -> Response {
    let cache_path = state.config.cache_dir.join(&object_id);

    // --- HIT path ---
    if cache_path.exists() {
        match tokio::fs::read(&cache_path).await {
            Ok(data) => {
                let mut h = HeaderMap::new();
                headers::set_cache_status(&mut h, true);
                headers::set_offline_available(&mut h, true);
                // Try to get priority from edge-agent asynchronously (best-effort)
                if let Some(info) = get_object_info(&state, &object_id).await {
                    headers::set_priority_class(&mut h, info.priority.as_deref());
                }
                return (StatusCode::OK, h, data).into_response();
            }
            Err(e) => {
                tracing::error!(object_id, error = %e, "failed to read cached file");
                return StatusCode::INTERNAL_SERVER_ERROR.into_response();
            }
        }
    }

    // --- MISS path ---
    let mut response_headers = HeaderMap::new();
    headers::set_cache_status(&mut response_headers, false);
    headers::set_offline_available(&mut response_headers, false);

    // Look up object info from edge-agent
    let info = get_object_info(&state, &object_id).await;
    let source_url = info
        .as_ref()
        .and_then(|i| i.source_url.clone())
        .or_else(|| state.config.upstream_url.as_ref().map(|u| format!("{u}/{object_id}")));

    let Some(url) = source_url else {
        tracing::warn!(object_id, "cache miss and no source_url available");
        return (StatusCode::NOT_FOUND, response_headers, b"object not found".to_vec()).into_response();
    };

    match fetch::fetch_and_store(&state.http, &url, &object_id, &state.config.cache_dir).await {
        Ok(data) => {
            if let Some(ref i) = info {
                headers::set_priority_class(&mut response_headers, i.priority.as_deref());
            }
            (StatusCode::OK, response_headers, data).into_response()
        }
        Err(e) => {
            tracing::error!(object_id, error = %e, "upstream fetch failed");
            (StatusCode::BAD_GATEWAY, response_headers, b"upstream fetch failed".to_vec()).into_response()
        }
    }
}

async fn get_object_info(state: &AppState, object_id: &str) -> Option<CacheObjectInfo> {
    let url = format!("{}/local/cache/status", state.config.edge_agent_url);
    let resp = state.http.get(&url).send().await.ok()?;
    let status_resp: CacheStatusResponse = resp.json().await.ok()?;
    status_resp.objects.into_iter().find(|o| o.object_id == object_id)
}
