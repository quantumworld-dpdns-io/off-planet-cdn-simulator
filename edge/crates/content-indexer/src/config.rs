pub struct Config {
    pub port: u16,
    pub qdrant_url: String,
    pub collection: String,
    pub vector_size: u64,
}

impl Config {
    pub fn from_env() -> Self {
        let port = std::env::var("CONTENT_INDEXER_PORT")
            .unwrap_or_else(|_| "8083".into())
            .parse::<u16>()
            .expect("CONTENT_INDEXER_PORT must be a valid port number");

        let qdrant_url = std::env::var("QDRANT_URL")
            .unwrap_or_else(|_| "http://localhost:6333".into());

        let collection = std::env::var("COLLECTION_NAME")
            .unwrap_or_else(|_| "cdn_content".into());

        Config {
            port,
            qdrant_url,
            collection,
            vector_size: 128,
        }
    }
}
