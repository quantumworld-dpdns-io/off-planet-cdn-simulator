mod config;
mod heartbeat;
mod sync;
mod api;
mod telemetry;

use anyhow::Result;

#[tokio::main]
async fn main() -> Result<()> {
    let cfg = config::Config::from_env()?;
    shared::telemetry::init("edge-agent");
    tracing::info!(port = cfg.port, "Edge agent starting");

    let app = api::router(cfg.clone());
    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", cfg.port)).await?;
    tracing::info!(port = cfg.port, "Listening");

    axum::serve(listener, app).await?;
    Ok(())
}
