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
    let meta_url = format!("{}/pypi/{}/json", upstream.trim_end_matches('/'), name);
    let meta: serde_json::Value = client
        .get(&meta_url)
        .send()
        .await
        .context("Failed to fetch PyPI package metadata")?
        .error_for_status()
        .context("PyPI metadata request returned error status")?
        .json()
        .await
        .context("Failed to parse PyPI metadata JSON")?;

    // Resolve version
    let resolved_version = match version {
        Some(v) => v.to_string(),
        None => meta
            .get("info")
            .and_then(|i| i.get("version"))
            .and_then(|v| v.as_str())
            .ok_or_else(|| anyhow!("Could not find info.version for PyPI package {}", name))?
            .to_string(),
    };

    // Find best download URL: prefer .whl, fallback to .tar.gz
    let releases = meta
        .get("releases")
        .and_then(|r| r.get(&resolved_version))
        .and_then(|v| v.as_array())
        .ok_or_else(|| {
            anyhow!(
                "No releases found for PyPI package {}@{}",
                name,
                resolved_version
            )
        })?;

    let mut download_url: Option<String> = None;
    let mut ext = "whl";

    // First pass: look for .whl
    for item in releases.iter() {
        if let Some(url) = item.get("url").and_then(|u| u.as_str()) {
            if url.ends_with(".whl") {
                download_url = Some(url.to_string());
                ext = "whl";
                break;
            }
        }
    }

    // Second pass: fallback to .tar.gz
    if download_url.is_none() {
        for item in releases.iter() {
            if let Some(url) = item.get("url").and_then(|u| u.as_str()) {
                if url.ends_with(".tar.gz") {
                    download_url = Some(url.to_string());
                    ext = "tar.gz";
                    break;
                }
            }
        }
    }

    let download_url = download_url.ok_or_else(|| {
        anyhow!(
            "No .whl or .tar.gz found for PyPI package {}@{}",
            name,
            resolved_version
        )
    })?;

    // Download the file
    let bytes = client
        .get(&download_url)
        .send()
        .await
        .context("Failed to download PyPI package")?
        .error_for_status()
        .context("PyPI package download returned error status")?
        .bytes()
        .await
        .context("Failed to read PyPI package bytes")?;

    // Write to disk
    let out_dir = format!("{}/pypi/{}", mirror_dir, name);
    tokio::fs::create_dir_all(&out_dir)
        .await
        .context("Failed to create PyPI mirror directory")?;

    let out_path = format!("{}/{}.{}", out_dir, resolved_version, ext);
    let mut file = tokio::fs::File::create(&out_path)
        .await
        .context("Failed to create PyPI artifact file")?;
    file.write_all(&bytes)
        .await
        .context("Failed to write PyPI package to disk")?;

    Ok(out_path)
}
