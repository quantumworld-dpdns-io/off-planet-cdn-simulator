mod config;
mod npm;
mod pypi;
mod crates_io;
mod oci;
mod model_registry;
mod manifest;

use std::sync::Arc;

use anyhow::Result;
use axum::{
    extract::State,
    http::StatusCode,
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};

use config::Config;

#[derive(Clone)]
pub struct AppState {
    pub config: Arc<Config>,
    pub client: reqwest::Client,
}

#[derive(Deserialize)]
struct SyncRequest {
    registry_type: String,
    name: String,
    version: Option<String>,
}

#[derive(Serialize)]
struct SyncResponse {
    ok: bool,
    artifact_path: String,
    message: String,
}

#[derive(Serialize)]
struct ArtifactEntry {
    path: String,
    size_bytes: u64,
    registry_type: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    shared::telemetry::init("registry-mirror");

    let config = Arc::new(Config::from_env());

    let client = reqwest::Client::builder()
        .user_agent("off-planet-cdn-mirror/0.1")
        .build()
        .expect("Failed to build reqwest client");

    let state = AppState {
        config: config.clone(),
        client,
    };

    tracing::info!(port = config.port, "Registry mirror starting");

    let app = Router::new()
        .route("/local/health", get(health))
        .route("/local/mirrors/artifacts", get(list_artifacts))
        .route("/local/mirrors/sync", post(sync_handler))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", config.port)).await?;
    tracing::info!(port = config.port, "Listening");

    axum::serve(listener, app).await?;
    Ok(())
}

async fn health() -> Json<Value> {
    Json(json!({"status": "ok", "service": "registry-mirror"}))
}

async fn list_artifacts(
    State(state): State<AppState>,
) -> Result<Json<Vec<ArtifactEntry>>, StatusCode> {
    let mirror_dir = &state.config.mirror_dir;
    let mut entries = Vec::new();

    if let Ok(()) = collect_artifacts(mirror_dir, mirror_dir, &mut entries).await {
        Ok(Json(entries))
    } else {
        // If mirror_dir doesn't exist yet, return empty list
        Ok(Json(vec![]))
    }
}

async fn collect_artifacts(
    mirror_dir: &str,
    dir: &str,
    entries: &mut Vec<ArtifactEntry>,
) -> Result<()> {
    let mut read_dir = tokio::fs::read_dir(dir).await?;
    while let Some(entry) = read_dir.next_entry().await? {
        let path = entry.path();
        let meta = entry.metadata().await?;
        if meta.is_dir() {
            let path_str = path.to_string_lossy().to_string();
            Box::pin(collect_artifacts(mirror_dir, &path_str, entries)).await?;
        } else if meta.is_file() {
            let full_path = path.to_string_lossy().to_string();
            let size_bytes = meta.len();

            // Determine registry_type from first path component after mirror_dir
            let relative = full_path
                .strip_prefix(mirror_dir)
                .unwrap_or(&full_path)
                .trim_start_matches('/');
            let registry_type = relative
                .split('/')
                .next()
                .unwrap_or("unknown")
                .to_string();

            entries.push(ArtifactEntry {
                path: full_path,
                size_bytes,
                registry_type,
            });
        }
    }
    Ok(())
}

async fn sync_handler(
    State(state): State<AppState>,
    Json(req): Json<SyncRequest>,
) -> Result<Json<SyncResponse>, StatusCode> {
    let name = &req.name;
    let version = req.version.as_deref();
    let mirror_dir = &state.config.mirror_dir;
    let client = &state.client;
    let config = &state.config;

    let result = match req.registry_type.as_str() {
        "npm" => {
            npm::sync(client, &config.upstream_npm, mirror_dir, name, version).await
        }
        "pypi" => {
            pypi::sync(client, &config.upstream_pypi, mirror_dir, name, version).await
        }
        "crates_io" => {
            crates_io::sync(client, &config.upstream_crates, mirror_dir, name, version).await
        }
        "oci" => {
            oci::sync(client, mirror_dir, name, version).await
        }
        "model" => {
            model_registry::sync(client, mirror_dir, name, version).await
        }
        other => {
            tracing::warn!(registry_type = other, "Unknown registry type requested");
            return Ok(Json(SyncResponse {
                ok: false,
                artifact_path: String::new(),
                message: format!("Unknown registry_type: {}", other),
            }));
        }
    };

    match result {
        Ok(artifact_path) => {
            tracing::info!(
                registry_type = req.registry_type.as_str(),
                name = name.as_str(),
                path = artifact_path.as_str(),
                "Sync completed"
            );
            Ok(Json(SyncResponse {
                ok: true,
                artifact_path,
                message: "sync successful".to_string(),
            }))
        }
        Err(err) => {
            tracing::error!(
                registry_type = req.registry_type.as_str(),
                name = name.as_str(),
                error = %err,
                "Sync failed"
            );
            Ok(Json(SyncResponse {
                ok: false,
                artifact_path: String::new(),
                message: format!("sync failed: {}", err),
            }))
        }
    }
}
