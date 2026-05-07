use anyhow::{anyhow, Context, Result};
use tokio::io::AsyncWriteExt;

pub async fn sync(
    client: &reqwest::Client,
    upstream: &str,
    mirror_dir: &str,
    name: &str,
    version: Option<&str>,
) -> Result<String> {
    // Fetch package metadata
    let meta_url = format!("{}/{}", upstream.trim_end_matches('/'), name);
    let meta: serde_json::Value = client
        .get(&meta_url)
        .send()
        .await
        .context("Failed to fetch npm package metadata")?
        .error_for_status()
        .context("npm metadata request returned error status")?
        .json()
        .await
        .context("Failed to parse npm metadata JSON")?;

    // Resolve version
    let resolved_version = match version {
        Some(v) => v.to_string(),
        None => meta
            .get("dist-tags")
            .and_then(|dt| dt.get("latest"))
            .and_then(|v| v.as_str())
            .ok_or_else(|| anyhow!("Could not find dist-tags.latest for npm package {}", name))?
            .to_string(),
    };

    // Get tarball URL
    let tarball_url = meta
        .get("versions")
        .and_then(|vs| vs.get(&resolved_version))
        .and_then(|v| v.get("dist"))
        .and_then(|d| d.get("tarball"))
        .and_then(|t| t.as_str())
        .ok_or_else(|| {
            anyhow!(
                "Could not find tarball URL for npm package {}@{}",
                name,
                resolved_version
            )
        })?
        .to_string();

    // Download tarball
    let bytes = client
        .get(&tarball_url)
        .send()
        .await
        .context("Failed to download npm tarball")?
        .error_for_status()
        .context("npm tarball download returned error status")?
        .bytes()
        .await
        .context("Failed to read npm tarball bytes")?;

    // Write to disk
    let out_dir = format!("{}/npm/{}", mirror_dir, name);
    tokio::fs::create_dir_all(&out_dir)
        .await
        .context("Failed to create npm mirror directory")?;

    let out_path = format!("{}/{}.tgz", out_dir, resolved_version);
    let mut file = tokio::fs::File::create(&out_path)
        .await
        .context("Failed to create npm artifact file")?;
    file.write_all(&bytes)
        .await
        .context("Failed to write npm tarball to disk")?;

    Ok(out_path)
}
