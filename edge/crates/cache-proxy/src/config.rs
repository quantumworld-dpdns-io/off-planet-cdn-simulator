use std::path::PathBuf;

#[derive(Clone, Debug)]
pub struct Config {
    pub port: u16,
    pub cache_dir: PathBuf,
    pub edge_agent_url: String,
    pub upstream_url: Option<String>,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            port: std::env::var("CACHE_PROXY_PORT")
                .unwrap_or_else(|_| "3128".into())
                .parse()
                .unwrap_or(3128),
            cache_dir: PathBuf::from(
                std::env::var("CACHE_DIR").unwrap_or_else(|_| "/tmp/edge-cache".into())
            ),
            edge_agent_url: std::env::var("EDGE_AGENT_URL")
                .unwrap_or_else(|_| "http://localhost:8081".into()),
            upstream_url: std::env::var("UPSTREAM_URL").ok(),
        }
    }
}
