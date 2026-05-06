use shared::types::{CacheObject, CacheScore, ScoreBreakdown};
use chrono::Utc;

pub fn score_object(obj: &CacheObject) -> CacheScore {
    let priority_weight = obj.priority.weight();

    // Staleness: penalize objects not accessed in > 30 days
    let days_since_access = (Utc::now() - obj.last_accessed_at).num_days();
    let staleness_penalty = (days_since_access as f64 * 10.0).min(5_000.0);

    // Size: penalize large objects (1 point per MB above 100 MB)
    let size_mb = obj.size_bytes as f64 / 1_048_576.0;
    let size_penalty = (size_mb - 100.0).max(0.0);

    let breakdown = ScoreBreakdown {
        priority_weight,
        mission_relevance: 0.0,   // populated by policy engine
        predicted_demand: 0.0,    // populated by scheduler
        offline_criticality: 0.0,
        revalidation_cost: 0.0,
        fetch_latency_cost: 0.0,
        package_dependency_score: 0.0,
        size_penalty,
        staleness_penalty,
        redundancy_penalty: 0.0,
    };

    let score = breakdown.priority_weight
        + breakdown.mission_relevance
        + breakdown.predicted_demand
        + breakdown.offline_criticality
        + breakdown.revalidation_cost
        + breakdown.fetch_latency_cost
        + breakdown.package_dependency_score
        - breakdown.size_penalty
        - breakdown.staleness_penalty
        - breakdown.redundancy_penalty;

    CacheScore {
        object_id: obj.id,
        score,
        breakdown,
    }
}
