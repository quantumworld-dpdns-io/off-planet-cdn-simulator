use anyhow::{Context, Result};
use tokio::io::AsyncWriteExt;

pub async fn sync(
    _client: &reqwest::Client,
    mirror_dir: &str,
    name: &str,
    version: Option<&str>,
) -> Result<String> {
    let tag = version.unwrap_or("latest");

    let synced_at = chrono::Utc::now().to_rfc3339();
    let metadata = serde_json::json!({
        "image": name,
        "tag": tag,
        "status": "metadata-only",
        "synced_at": synced_at,
    });

    let out_dir = format!("{}/oci/{}", mirror_dir, name);
    tokio::fs::create_dir_all(&out_dir)
        .await
        .context("Failed to create OCI mirror directory")?;

    let out_path = format!("{}/{}.json", out_dir, tag);
    let mut file = tokio::fs::File::create(&out_path)
        .await
        .context("Failed to create OCI metadata file")?;
    file.write_all(metadata.to_string().as_bytes())
        .await
        .context("Failed to write OCI metadata to disk")?;

    Ok(out_path)
}
