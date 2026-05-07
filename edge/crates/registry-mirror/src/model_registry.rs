use anyhow::{Context, Result};
use tokio::io::AsyncWriteExt;

pub async fn sync(
    client: &reqwest::Client,
    mirror_dir: &str,
    name: &str,
    version: Option<&str>,
) -> Result<String> {
    let version_str = version.unwrap_or("main");

    // Sanitize name for filesystem (replace / with __)
    let name_sanitized = name.replace('/', "__");

    // Attempt to fetch HuggingFace model card metadata
    let api_url = format!("https://huggingface.co/api/models/{}", name);
    let metadata = match client.get(&api_url).send().await {
        Ok(resp) if resp.status().is_success() => {
            match resp.json::<serde_json::Value>().await {
                Ok(json) => {
                    serde_json::json!({
                        "modelId": json.get("modelId").unwrap_or(&serde_json::Value::Null),
                        "pipeline_tag": json.get("pipeline_tag").unwrap_or(&serde_json::Value::Null),
                        "tags": json.get("tags").unwrap_or(&serde_json::Value::Null),
                        "lastModified": json.get("lastModified").unwrap_or(&serde_json::Value::Null),
                    })
                }
                Err(_) => {
                    serde_json::json!({
                        "model": name,
                        "status": "metadata-unavailable",
                    })
                }
            }
        }
        _ => {
            serde_json::json!({
                "model": name,
                "status": "metadata-unavailable",
            })
        }
    };

    let out_dir = format!("{}/model/{}", mirror_dir, name_sanitized);
    tokio::fs::create_dir_all(&out_dir)
        .await
        .context("Failed to create model mirror directory")?;

    let out_path = format!("{}/{}.json", out_dir, version_str);
    let mut file = tokio::fs::File::create(&out_path)
        .await
        .context("Failed to create model metadata file")?;
    file.write_all(metadata.to_string().as_bytes())
        .await
        .context("Failed to write model metadata to disk")?;

    Ok(out_path)
}
