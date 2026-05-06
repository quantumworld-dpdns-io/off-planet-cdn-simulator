use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CacheManifest {
    pub manifest_version: u32,
    pub node_id: String,
    pub site_id: Uuid,
    pub generated_at: DateTime<Utc>,
    pub objects: Vec<ManifestEntry>,
    pub signature: Option<Vec<u8>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ManifestEntry {
    pub object_id: Uuid,
    pub name: String,
    pub content_hash: String,
    pub size_bytes: u64,
    pub priority: String,
    pub pinned: bool,
    pub storage_path: String,
}

impl CacheManifest {
    pub fn new(node_id: String, site_id: Uuid) -> Self {
        Self {
            manifest_version: 1,
            node_id,
            site_id,
            generated_at: chrono::Utc::now(),
            objects: vec![],
            signature: None,
        }
    }
}
