use thiserror::Error;

#[derive(Error, Debug)]
pub enum CdnError {
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    #[error("HTTP error: {0}")]
    Http(String),
    #[error("Serialization error: {0}")]
    Serde(#[from] serde_json::Error),
    #[error("Cache miss: {0}")]
    CacheMiss(String),
    #[error("Eviction blocked: object is pinned or P0")]
    EvictionBlocked,
    #[error("Insufficient space: needed {needed} bytes, can only free {available} bytes")]
    InsufficientSpace { needed: u64, available: u64 },
    #[error("Configuration error: {0}")]
    Config(String),
    #[error("Database error: {0}")]
    Database(String),
}
