// Polls control API for pending preload jobs and downloads content
use crate::config::Config;
use crate::telemetry;
use std::time::Duration;
use std::path::PathBuf;

#[derive(Debug, serde::Deserialize)]
struct PreloadJob {
    id: String,
    #[allow(dead_code)]
    status: String,
}

#[derive(Debug, serde::Deserialize)]
struct PreloadJobsResponse {
    jobs: Vec<PreloadJob>,
}

pub async fn start_sync_loop(cfg: Config) {
    let client = reqwest::Client::new();
    let org_id = std::env::var("ORG_ID")
        .unwrap_or_else(|_| "00000000-0000-0000-0000-000000000001".to_string());

    // Ensure cache directory exists
    if let Err(e) = std::fs::create_dir_all(&cfg.cache_dir) {
        telemetry::record_error("sync_loop_init", &e.to_string());
    }

    loop {
        tokio::time::sleep(Duration::from_secs(15)).await;

        match poll_pending_jobs(&client, &cfg, &org_id).await {
            Ok(count) => telemetry::record_sync_tick(count),
            Err(e) => telemetry::record_error("sync_poll", &e.to_string()),
        }
    }
}

async fn poll_pending_jobs(client: &reqwest::Client, cfg: &Config, org_id: &str) -> anyhow::Result<usize> {
    let url = format!("{}/v1/preload/jobs?status=PENDING", cfg.control_api_url);
    let resp = client
        .get(&url)
        .header("X-Org-ID", org_id)
        .send()
        .await?;

    if !resp.status().is_success() {
        return Ok(0);
    }

    let jobs_resp: PreloadJobsResponse = resp.json().await.unwrap_or(PreloadJobsResponse { jobs: vec![] });
    let count = jobs_resp.jobs.len();

    for job in &jobs_resp.jobs {
        if let Err(e) = process_job(client, cfg, org_id, &job.id).await {
            telemetry::record_error("process_job", &e.to_string());
        }
    }

    Ok(count)
}

async fn process_job(client: &reqwest::Client, cfg: &Config, org_id: &str, job_id: &str) -> anyhow::Result<()> {
    let url = format!("{}/v1/preload/jobs/{}", cfg.control_api_url, job_id);
    let resp = client.get(&url).header("X-Org-ID", org_id).send().await?;
    if !resp.status().is_success() {
        return Ok(());
    }

    // For now just log — full item processing implemented in S4/S5
    tracing::info!(job_id = job_id, "processing preload job");
    Ok(())
}

pub async fn fetch_object(cfg: &Config, object_id: &str, source_url: &str) -> anyhow::Result<PathBuf> {
    let client = reqwest::Client::builder()
        .timeout(Duration::from_secs(300))
        .build()?;

    let resp = client.get(source_url).send().await?;
    if !resp.status().is_success() {
        anyhow::bail!("fetch failed: {}", resp.status());
    }

    let bytes = resp.bytes().await?;

    // Store content-addressably under cache_dir/<object_id>
    let dest = PathBuf::from(&cfg.cache_dir).join(object_id);
    std::fs::write(&dest, &bytes)?;

    telemetry::record_fetch_complete(object_id, bytes.len() as u64, &cfg.cache_dir);
    Ok(dest)
}
