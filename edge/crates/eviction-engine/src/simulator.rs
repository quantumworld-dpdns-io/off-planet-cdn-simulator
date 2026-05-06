use shared::types::CacheObject;
use crate::score::score_object;
use crate::constraints::can_evict;
use shared::errors::CdnError;

#[derive(Debug, Clone)]
pub struct EvictionCandidate {
    pub object_id: uuid::Uuid,
    pub name: String,
    pub score: f64,
    pub size_bytes: u64,
    pub reason: String,
}

pub fn simulate(
    objects: &[CacheObject],
    target_freed_bytes: u64,
) -> Result<Vec<EvictionCandidate>, CdnError> {
    // Score all evictable objects
    let mut candidates: Vec<(f64, &CacheObject)> = objects
        .iter()
        .filter(|o| can_evict(o))
        .map(|o| (score_object(o).score, o))
        .collect();

    // Sort ascending by score (lowest score evicted first)
    candidates.sort_by(|a, b| a.0.partial_cmp(&b.0).unwrap());

    let mut plan = Vec::new();
    let mut freed = 0u64;

    for (score, obj) in &candidates {
        if freed >= target_freed_bytes {
            break;
        }
        plan.push(EvictionCandidate {
            object_id: obj.id,
            name: obj.name.clone(),
            score: *score,
            size_bytes: obj.size_bytes,
            reason: format!("score={:.2}, priority={:?}", score, obj.priority),
        });
        freed += obj.size_bytes;
    }

    if freed < target_freed_bytes && freed < objects.iter().map(|o| o.size_bytes).sum::<u64>() {
        return Err(CdnError::InsufficientSpace {
            needed: target_freed_bytes,
            available: freed,
        });
    }

    Ok(plan)
}
