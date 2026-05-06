use anyhow::{anyhow, Result};

#[derive(Clone, Debug)]
pub struct Config {
    pub port: u16,
    pub node_id: String,
    pub control_api_url: String,
    pub cache_dir: String,
    pub cache_max_bytes: u64,
    pub agent_version: String,
}

impl Config {
    pub fn from_env() -> Result<Self> {
        Ok(Self {
            port: std::env::var("EDGE_AGENT_PORT")
                .unwrap_or_else(|_| "8081".into())
                .parse()?,
            node_id: std::env::var("EDGE_NODE_ID")
                .unwrap_or_else(|_| "local-dev-node".into()),
            control_api_url: std::env::var("CONTROL_API_URL")
                .unwrap_or_else(|_| "http://localhost:8080".into()),
            cache_dir: std::env::var("CACHE_DIR")
                .unwrap_or_else(|_| "/tmp/offplanet-cache".into()),
            cache_max_bytes: std::env::var("CACHE_MAX_BYTES")
                .unwrap_or_else(|_| "10737418240".into())
                .parse()?,
            agent_version: "0.1.0".into(),
        })
    }
}
