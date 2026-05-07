pub const X_CACHE: &str = "X-Cache";
pub const X_PRIORITY_CLASS: &str = "X-Priority-Class";
pub const X_OFFLINE_AVAILABLE: &str = "X-Offline-Available";

use axum::http::HeaderMap;

/// Sets X-Cache: HIT or X-Cache: MISS
pub fn set_cache_status(headers: &mut HeaderMap, hit: bool) {
    let val = if hit { "HIT" } else { "MISS" };
    headers.insert(X_CACHE, val.parse().unwrap());
}

/// Sets X-Priority-Class header if priority is known
pub fn set_priority_class(headers: &mut HeaderMap, priority: Option<&str>) {
    if let Some(p) = priority {
        if let Ok(v) = p.parse() {
            headers.insert(X_PRIORITY_CLASS, v);
        }
    }
}

/// Sets X-Offline-Available: true when served from local cache
pub fn set_offline_available(headers: &mut HeaderMap, available: bool) {
    let val = if available { "true" } else { "false" };
    headers.insert(X_OFFLINE_AVAILABLE, val.parse().unwrap());
}
