// Host function imports exposed to WASM policy plugins
// Full implementation in Phase 5

/// Returns the current UTC timestamp as seconds since epoch.
pub fn host_now_secs() -> u64 {
    // TODO: expose as WASM host import
    0
}

/// Logs a message from inside a WASM plugin.
pub fn host_log(_level: &str, _message: &str) {
    // TODO: expose as WASM host import
}

/// Returns the node's current cache fill ratio (0.0 – 1.0).
pub fn host_cache_fill_ratio() -> f64 {
    // TODO: expose as WASM host import
    0.0
}
