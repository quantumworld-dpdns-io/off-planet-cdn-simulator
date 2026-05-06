mod chunk;
mod embeddings;
mod qdrant;
mod duckdb_export;

use anyhow::Result;
use axum::{routing::get, Router, Json};
use serde_json::{json, Value};

#[tokio::main]
async fn main() -> Result<()> {
    shared::telemetry::init("content-indexer");

    let port: u16 = std::env::var("CONTENT_INDEXER_PORT")
        .unwrap_or_else(|_| "8083".into())
        .parse()?;

    tracing::info!(port = port, "Content indexer starting");

    let app = Router::new()
        .route("/local/health", get(health));

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", port)).await?;
    tracing::info!(port = port, "Listening");

    axum::serve(listener, app).await?;
    Ok(())
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "content-indexer"}))
}
