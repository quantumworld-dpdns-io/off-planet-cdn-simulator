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
    tracing::info!(port = cfg.port, node_id = %cfg.node_id, "Edge agent starting");

    // Spawn background tasks
    tokio::spawn(heartbeat::start_heartbeat_loop(cfg.clone()));
    tokio::spawn(sync::start_sync_loop(cfg.clone()));

    // Start HTTP server
    let app = api::router(cfg.clone());
    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", cfg.port)).await?;
    tracing::info!(port = cfg.port, "Listening");
    axum::serve(listener, app).await?;
    Ok(())
}
