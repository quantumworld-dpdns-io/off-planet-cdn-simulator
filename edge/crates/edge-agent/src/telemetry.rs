// Edge-agent specific OTel span helpers

pub fn record_heartbeat_sent(node_id: &str, status: &str) {
    tracing::info!(node_id = node_id, status = status, "heartbeat sent");
}

pub fn record_fetch_complete(object_id: &str, bytes: u64, cache_dir: &str) {
    tracing::info!(object_id = object_id, bytes = bytes, cache_dir = cache_dir, "content fetch complete");
}

pub fn record_sync_tick(pending_jobs: usize) {
    tracing::debug!(pending_jobs = pending_jobs, "sync tick");
}

pub fn record_error(context: &str, err: &str) {
    tracing::error!(context = context, error = err, "edge agent error");
}
