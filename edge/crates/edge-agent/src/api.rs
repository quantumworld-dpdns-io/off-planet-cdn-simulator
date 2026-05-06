use axum::{
    extract::{Json, State},
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use std::sync::Arc;
use crate::config::Config;

pub fn router(cfg: Config) -> Router {
    let state = Arc::new(cfg);
    Router::new()
        .route("/local/health", get(health))
        .route("/local/cache/status", get(cache_status))
        .route("/local/cache/fetch", post(cache_fetch))
        .route("/local/cache/preload", post(cache_preload))
        .route("/local/policy/reload", post(policy_reload))
        .with_state(state)
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "edge-agent", "version": "0.1.0"}))
}

async fn cache_status(State(cfg): State<Arc<Config>>) -> Json<Value> {
    let used = measure_cache_bytes(&cfg.cache_dir);
    let object_count = count_cache_objects(&cfg.cache_dir);
    Json(json!({
        "node_id": cfg.node_id,
        "cache_used_bytes": used,
        "cache_max_bytes": cfg.cache_max_bytes,
        "object_count": object_count,
        "pinned_count": 0,
        "fill_ratio": if cfg.cache_max_bytes > 0 { used as f64 / cfg.cache_max_bytes as f64 } else { 0.0 },
    }))
}

#[derive(Deserialize)]
struct FetchRequest {
    object_id: String,
    source_url: String,
}

#[derive(Serialize)]
struct FetchResponse {
    object_id: String,
    stored_path: String,
    bytes: u64,
}

async fn cache_fetch(
    State(cfg): State<Arc<Config>>,
    Json(req): Json<FetchRequest>,
) -> Result<Json<Value>, (axum::http::StatusCode, Json<Value>)> {
    if req.object_id.is_empty() || req.source_url.is_empty() {
        return Err((
            axum::http::StatusCode::BAD_REQUEST,
            Json(json!({"error": "object_id and source_url are required"})),
        ));
    }

    match crate::sync::fetch_object(&cfg, &req.object_id, &req.source_url).await {
        Ok(path) => {
            let bytes = std::fs::metadata(&path).map(|m| m.len()).unwrap_or(0);
            Ok(Json(json!({
                "object_id": req.object_id,
                "stored_path": path.to_string_lossy(),
                "bytes": bytes,
            })))
        }
        Err(e) => Err((
            axum::http::StatusCode::INTERNAL_SERVER_ERROR,
            Json(json!({"error": e.to_string()})),
        )),
    }
}

#[derive(Deserialize)]
struct PreloadRequest {
    objects: Vec<PreloadItem>,
}

#[derive(Deserialize)]
struct PreloadItem {
    object_id: String,
    source_url: String,
    #[allow(dead_code)]
    priority: Option<String>,
}

async fn cache_preload(
    State(cfg): State<Arc<Config>>,
    Json(req): Json<PreloadRequest>,
) -> Json<Value> {
    let count = req.objects.len();
    tracing::info!(count = count, "preload request received");

    // Spawn individual fetches in background (non-blocking)
    for item in req.objects {
        let cfg_clone = cfg.clone();
        tokio::spawn(async move {
            if let Err(e) = crate::sync::fetch_object(&cfg_clone, &item.object_id, &item.source_url).await {
                tracing::error!(object_id = %item.object_id, error = %e, "preload fetch failed");
            }
        });
    }

    Json(json!({
        "accepted": count,
        "status": "queued",
    }))
}

async fn policy_reload() -> Json<Value> {
    tracing::info!("policy reload requested");
    Json(json!({"reloaded": true}))
}

fn measure_cache_bytes(cache_dir: &str) -> u64 {
    let path = std::path::Path::new(cache_dir);
    if !path.exists() {
        return 0;
    }
    dir_size(path)
}

fn dir_size(path: &std::path::Path) -> u64 {
    let mut total = 0u64;
    if let Ok(entries) = std::fs::read_dir(path) {
        for entry in entries.flatten() {
            let p = entry.path();
            if p.is_file() {
                total += p.metadata().map(|m| m.len()).unwrap_or(0);
            } else if p.is_dir() {
                total += dir_size(&p);
            }
        }
    }
    total
}

fn count_cache_objects(cache_dir: &str) -> u64 {
    let path = std::path::Path::new(cache_dir);
    if !path.exists() {
        return 0;
    }
    std::fs::read_dir(path)
        .map(|entries| entries.flatten().filter(|e| e.path().is_file()).count() as u64)
        .unwrap_or(0)
}
