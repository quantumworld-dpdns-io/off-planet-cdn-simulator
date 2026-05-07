mod config;
mod fetch;
mod headers;
mod proxy;
mod range_requests;
mod revalidate;

use std::sync::Arc;
use anyhow::Result;
use axum::{routing::get, Router};
use axum::extract::{Path, State};
use axum::http::{HeaderMap, StatusCode};
use axum::response::{IntoResponse, Response};

use proxy::AppState;
use config::Config;

#[tokio::main]
async fn main() -> Result<()> {
    shared::telemetry::init("cache-proxy");

    let cfg = Config::from_env();
    tracing::info!(port = cfg.port, cache_dir = ?cfg.cache_dir, "cache-proxy starting");

    // Ensure cache dir exists at startup
    tokio::fs::create_dir_all(&cfg.cache_dir).await?;

    let state = Arc::new(AppState {
        config: Arc::new(cfg.clone()),
        http: reqwest::Client::new(),
    });

    let app = Router::new()
        .route("/health", get(health))
        .route("/cache/{object_id}", get(handle_with_range))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", cfg.port)).await?;
    tracing::info!(port = cfg.port, "listening");
    axum::serve(listener, app).await?;
    Ok(())
}

async fn health() -> axum::Json<serde_json::Value> {
    axum::Json(serde_json::json!({"status": "ok", "service": "cache-proxy"}))
}

/// Wraps proxy_handler to add Range + revalidation support on top of the base hit/miss logic.
async fn handle_with_range(
    State(state): State<Arc<AppState>>,
    Path(object_id): Path<String>,
    request_headers: HeaderMap,
) -> Response {
    let cache_path = state.config.cache_dir.join(&object_id);

    // Check ETag / Last-Modified before reading the file
    if cache_path.exists() {
        if let Some((etag, mtime)) = revalidate::file_etag(&cache_path).await {
            if revalidate::is_not_modified_etag(&request_headers, &etag)
                || revalidate::is_not_modified_since(&request_headers, mtime)
            {
                let mut h = HeaderMap::new();
                revalidate::apply_cache_control_headers(&mut h, &etag, mtime);
                headers::set_cache_status(&mut h, true);
                return (StatusCode::NOT_MODIFIED, h).into_response();
            }
        }
    }

    // Delegate to base proxy handler (returns full file body on HIT or fetches on MISS)
    // We reconstruct the response with range support applied
    let base = proxy::proxy_handler(State(state), Path(object_id)).await;

    // Apply Range header to successful 200 responses
    let range_header = request_headers
        .get(axum::http::header::RANGE)
        .and_then(|v| v.to_str().ok())
        .map(|s| s.to_string());

    if base.status() == StatusCode::OK {
        if let Some(ref range) = range_header {
            // Extract body bytes from the response — we need to buffer it
            let (mut parts, body) = base.into_parts();
            let bytes = axum::body::to_bytes(body, usize::MAX).await.unwrap_or_default();
            let rr = range_requests::apply_range(bytes.to_vec(), Some(range));
            range_requests::inject_range_headers(&mut parts.headers, &rr);
            parts.status = rr.status;
            return (parts.status, parts.headers, rr.body).into_response();
        }
    }

    base
}
