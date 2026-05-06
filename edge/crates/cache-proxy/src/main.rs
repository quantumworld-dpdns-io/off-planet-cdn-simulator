mod proxy;
mod fetch;
mod range_requests;
mod revalidate;
mod headers;

use anyhow::Result;
use axum::{routing::get, Router, Json};
use serde_json::{json, Value};

#[tokio::main]
async fn main() -> Result<()> {
    shared::telemetry::init("cache-proxy");

    let port: u16 = std::env::var("CACHE_PROXY_PORT")
        .unwrap_or_else(|_| "3128".into())
        .parse()?;

    tracing::info!(port = port, "Cache proxy starting");

    let app = Router::new()
        .route("/local/health", get(health));

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", port)).await?;
    tracing::info!(port = port, "Listening");

    axum::serve(listener, app).await?;
    Ok(())
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "cache-proxy"}))
}
