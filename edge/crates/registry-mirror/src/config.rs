pub struct Config {
    pub port: u16,
    pub mirror_dir: String,
    pub upstream_npm: String,
    pub upstream_pypi: String,
    pub upstream_crates: String,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            port: std::env::var("REGISTRY_MIRROR_PORT")
                .unwrap_or_else(|_| "8082".into())
                .parse()
                .unwrap_or(8082),
            mirror_dir: std::env::var("MIRROR_DIR")
                .unwrap_or_else(|_| "/tmp/registry-mirror".into()),
            upstream_npm: std::env::var("NPM_UPSTREAM")
                .unwrap_or_else(|_| "https://registry.npmjs.org".into()),
            upstream_pypi: std::env::var("PYPI_UPSTREAM")
                .unwrap_or_else(|_| "https://pypi.org".into()),
            upstream_crates: std::env::var("CRATES_UPSTREAM")
                .unwrap_or_else(|_| "https://crates.io".into()),
        }
    }
}
