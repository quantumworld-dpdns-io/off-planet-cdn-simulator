use anyhow::{anyhow, Context, Result};
use tokio::io::AsyncWriteExt;

pub async fn sync(
    client: &reqwest::Client,
    upstream: &str,
    mirror_dir: &str,
    name: &str,
    version: Option<&str>,
) -> Result<String> {
    // Fetch crate metadata
    let meta_url = format!("{}/api/v1/crates/{}", upstream.trim_end_matches('/'), name);
    let meta: serde_json::Value = client
        .get(&meta_url)
        .send()
        .await
        .context("Failed to fetch crates.io metadata")?
        .error_for_status()
        .context("crates.io metadata request returned error status")?
        .json()
        .await
        .context("Failed to parse crates.io metadata JSON")?;

    // Resolve version
    let resolved_version = match version {
        Some(v) => v.to_string(),
        None => meta
            .get("crate")
            .and_then(|c| c.get("newest_version"))
            .and_then(|v| v.as_str())
            .ok_or_else(|| {
                anyhow!("Could not find crate.newest_version for crates.io package {}", name)
            })?
            .to_string(),
    };

    // Download the .crate file (reqwest follows redirects automatically)
    let download_url = format!(
        "{}/api/v1/crates/{}/{}/download",
        upstream.trim_end_matches('/'),
        name,
        resolved_version
    );

    let bytes = client
        .get(&download_url)
        .send()
        .await
        .context("Failed to download crates.io package")?
        .error_for_status()
        .context("crates.io package download returned error status")?
        .bytes()
        .await
        .context("Failed to read crates.io package bytes")?;

    // Write to disk
    let out_dir = format!("{}/crates_io/{}", mirror_dir, name);
    tokio::fs::create_dir_all(&out_dir)
        .await
        .context("Failed to create crates_io mirror directory")?;

    let out_path = format!("{}/{}.crate", out_dir, resolved_version);
    let mut file = tokio::fs::File::create(&out_path)
        .await
        .context("Failed to create crates.io artifact file")?;
    file.write_all(&bytes)
        .await
        .context("Failed to write crates.io package to disk")?;

    Ok(out_path)
}
