// Sends periodic heartbeat to control API
use crate::config::Config;
use crate::telemetry;
use std::time::Duration;

pub async fn start_heartbeat_loop(cfg: Config) {
    let client = reqwest::Client::new();
    let org_id = std::env::var("ORG_ID")
        .unwrap_or_else(|_| "00000000-0000-0000-0000-000000000001".to_string());

    loop {
        match send_heartbeat(&client, &cfg, &org_id).await {
            Ok(_) => telemetry::record_heartbeat_sent(&cfg.node_id, "ONLINE"),
            Err(e) => telemetry::record_error("heartbeat", &e.to_string()),
        }
        tokio::time::sleep(Duration::from_secs(30)).await;
    }
}

async fn send_heartbeat(client: &reqwest::Client, cfg: &Config, org_id: &str) -> anyhow::Result<()> {
    let cache_used = measure_cache_used_bytes(&cfg.cache_dir);

    let body = serde_json::json!({
        "status": "ONLINE",
        "cache_used_bytes": cache_used,
        "cache_max_bytes": cfg.cache_max_bytes,
        "agent_version": cfg.agent_version,
    });

    let url = format!("{}/v1/nodes/{}/heartbeat", cfg.control_api_url, cfg.node_id);
    let resp = client
        .post(&url)
        .header("Content-Type", "application/json")
        .header("X-Org-ID", org_id)
        .json(&body)
        .send()
        .await?;

    if !resp.status().is_success() {
        let status = resp.status();
        let text = resp.text().await.unwrap_or_default();
        anyhow::bail!("heartbeat failed: {} — {}", status, text);
    }
    Ok(())
}

fn measure_cache_used_bytes(cache_dir: &str) -> u64 {
    let path = std::path::Path::new(cache_dir);
    if !path.exists() {
        return 0;
    }
    walkdir_size(path)
}

fn walkdir_size(path: &std::path::Path) -> u64 {
    let mut total = 0u64;
    if let Ok(entries) = std::fs::read_dir(path) {
        for entry in entries.flatten() {
            let p = entry.path();
            if p.is_file() {
                total += p.metadata().map(|m| m.len()).unwrap_or(0);
            } else if p.is_dir() {
                total += walkdir_size(&p);
            }
        }
    }
    total
}
