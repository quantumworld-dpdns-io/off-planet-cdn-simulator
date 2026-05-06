mod npm;
mod pypi;
mod crates_io;
mod oci;
mod model_registry;
mod manifest;

use anyhow::Result;
use axum::{routing::get, Router, Json};
use serde_json::{json, Value};

#[tokio::main]
async fn main() -> Result<()> {
    shared::telemetry::init("registry-mirror");

    let port: u16 = std::env::var("REGISTRY_MIRROR_PORT")
        .unwrap_or_else(|_| "8082".into())
        .parse()?;

    tracing::info!(port = port, "Registry mirror starting");

    let app = Router::new()
        .route("/local/health", get(health));

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", port)).await?;
    tracing::info!(port = port, "Listening");

    axum::serve(listener, app).await?;
    Ok(())
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "registry-mirror"}))
}
