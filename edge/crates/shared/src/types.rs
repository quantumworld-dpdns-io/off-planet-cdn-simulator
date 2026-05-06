use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub enum PriorityClass {
    P0, P1, P2, P3, P4, P5,
}

impl PriorityClass {
    pub fn weight(&self) -> f64 {
        match self {
            PriorityClass::P0 => 10_000.0,
            PriorityClass::P1 => 5_000.0,
            PriorityClass::P2 => 2_000.0,
            PriorityClass::P3 => 1_000.0,
            PriorityClass::P4 => 100.0,
            PriorityClass::P5 => 0.0,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CacheObject {
    pub id: Uuid,
    pub name: String,
    pub priority: PriorityClass,
    pub size_bytes: u64,
    pub pinned: bool,
    pub content_hash: String,
    pub source_url: String,
    pub last_accessed_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CacheScore {
    pub object_id: Uuid,
    pub score: f64,
    pub breakdown: ScoreBreakdown,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ScoreBreakdown {
    pub priority_weight: f64,
    pub mission_relevance: f64,
    pub predicted_demand: f64,
    pub offline_criticality: f64,
    pub revalidation_cost: f64,
    pub fetch_latency_cost: f64,
    pub package_dependency_score: f64,
    pub size_penalty: f64,
    pub staleness_penalty: f64,
    pub redundancy_penalty: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NodeStatus {
    pub node_id: String,
    pub cache_used_bytes: u64,
    pub cache_max_bytes: u64,
    pub object_count: u64,
    pub pinned_count: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Heartbeat {
    pub node_id: String,
    pub status: String,
    pub cache_used_bytes: u64,
    pub cache_max_bytes: u64,
    pub agent_version: String,
}
