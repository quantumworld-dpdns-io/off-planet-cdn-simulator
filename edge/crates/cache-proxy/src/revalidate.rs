use axum::http::{HeaderMap, header};
use std::time::UNIX_EPOCH;

/// Generates a simple ETag from file size + mtime (weak ETag).
pub fn make_etag(size: u64, mtime_secs: u64) -> String {
    format!("W/\"{size}-{mtime_secs}\"")
}

/// Returns true if the client's If-None-Match matches our ETag (304 Not Modified).
pub fn is_not_modified_etag(request_headers: &HeaderMap, our_etag: &str) -> bool {
    if let Some(val) = request_headers.get(header::IF_NONE_MATCH) {
        if let Ok(s) = val.to_str() {
            return s == our_etag || s == "*";
        }
    }
    false
}

/// Returns true if If-Modified-Since header indicates content is fresh.
pub fn is_not_modified_since(request_headers: &HeaderMap, mtime_secs: u64) -> bool {
    if let Some(val) = request_headers.get(header::IF_MODIFIED_SINCE) {
        if let Ok(s) = val.to_str() {
            // Very simple: compare as epoch seconds encoded in the header string.
            // In prod you'd parse HTTP date; for the simulator we accept epoch seconds.
            if let Ok(client_time) = s.trim().parse::<u64>() {
                return mtime_secs <= client_time;
            }
        }
    }
    false
}

/// Checks cache file metadata and returns (etag, mtime_secs) if available.
pub async fn file_etag(path: &std::path::Path) -> Option<(String, u64)> {
    let meta = tokio::fs::metadata(path).await.ok()?;
    let size = meta.len();
    let mtime = meta
        .modified()
        .ok()?
        .duration_since(UNIX_EPOCH)
        .ok()?
        .as_secs();
    Some((make_etag(size, mtime), mtime))
}

/// Sets ETag and Last-Modified headers on response.
pub fn apply_cache_control_headers(headers: &mut HeaderMap, etag: &str, mtime_secs: u64) {
    if let Ok(v) = etag.parse() {
        headers.insert(header::ETAG, v);
    }
    if let Ok(v) = mtime_secs.to_string().parse() {
        headers.insert(header::LAST_MODIFIED, v);
    }
    // Cache-Control: no-store for mission-critical; clients should always revalidate
    if let Ok(v) = "no-store".parse() {
        headers.insert(header::CACHE_CONTROL, v);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn make_etag_format() {
        let etag = make_etag(512, 1700000000);
        assert_eq!(etag, "W/\"512-1700000000\"");
    }

    #[test]
    fn etag_deterministic() {
        assert_eq!(make_etag(1234, 9999), make_etag(1234, 9999));
    }

    #[test]
    fn etag_differs_on_size_change() {
        assert_ne!(make_etag(100, 9999), make_etag(200, 9999));
    }

    #[test]
    fn etag_differs_on_mtime_change() {
        assert_ne!(make_etag(100, 1000), make_etag(100, 2000));
    }

    #[test]
    fn is_not_modified_etag_match() {
        let mut headers = HeaderMap::new();
        headers.insert(header::IF_NONE_MATCH, "W/\"100-9999\"".parse().unwrap());
        assert!(is_not_modified_etag(&headers, "W/\"100-9999\""));
    }

    #[test]
    fn is_not_modified_etag_wildcard() {
        let mut headers = HeaderMap::new();
        headers.insert(header::IF_NONE_MATCH, "*".parse().unwrap());
        assert!(is_not_modified_etag(&headers, "W/\"100-9999\""));
    }

    #[test]
    fn is_not_modified_etag_mismatch() {
        let mut headers = HeaderMap::new();
        headers.insert(header::IF_NONE_MATCH, "W/\"999-9999\"".parse().unwrap());
        assert!(!is_not_modified_etag(&headers, "W/\"100-9999\""));
    }

    #[test]
    fn is_not_modified_since_fresh() {
        let mut headers = HeaderMap::new();
        // client says they have version from time 1000, file mtime is 900 (not newer)
        headers.insert(header::IF_MODIFIED_SINCE, "1000".parse().unwrap());
        assert!(is_not_modified_since(&headers, 900));
    }

    #[test]
    fn is_not_modified_since_stale() {
        let mut headers = HeaderMap::new();
        // client says they have version from time 500, file mtime is 1000 (newer)
        headers.insert(header::IF_MODIFIED_SINCE, "500".parse().unwrap());
        assert!(!is_not_modified_since(&headers, 1000));
    }
}
