use std::path::Path;
use anyhow::{Context, Result};

/// Downloads the object from `source_url` and writes it to `{cache_dir}/{object_id}`.
/// Returns the raw bytes on success so the caller can serve them immediately without
/// a second disk read.
pub async fn fetch_and_store(
    client: &reqwest::Client,
    source_url: &str,
    object_id: &str,
    cache_dir: &Path,
) -> Result<Vec<u8>> {
    let resp = client
        .get(source_url)
        .send()
        .await
        .with_context(|| format!("GET {source_url}"))?;

    if !resp.status().is_success() {
        anyhow::bail!("upstream returned {}", resp.status());
    }

    let bytes = resp.bytes().await.context("read upstream body")?;

    // Ensure cache dir exists
    tokio::fs::create_dir_all(cache_dir)
        .await
        .context("create cache dir")?;

    let dest = cache_dir.join(object_id);
    tokio::fs::write(&dest, &bytes)
        .await
        .with_context(|| format!("write cache file {:?}", dest))?;

    tracing::info!(object_id, bytes = bytes.len(), "fetched and cached");
    Ok(bytes.to_vec())
}
