use axum::{routing::get, Router, Json};
use serde_json::{json, Value};
use crate::config::Config;

pub fn router(_cfg: Config) -> Router {
    Router::new()
        .route("/local/health", get(health))
        .route("/local/cache/status", get(cache_status))
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "edge-agent", "version": "0.1.0"}))
}

async fn cache_status() -> Json<Value> {
    Json(json!({
        "cache_used_bytes": 0,
        "cache_max_bytes": 10737418240_u64,
        "object_count": 0,
        "pinned_count": 0
    }))
}
