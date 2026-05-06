#[cfg(test)]
mod tests {
    use crate::{score, constraints, simulator};
    use shared::types::{CacheObject, PriorityClass};
    use uuid::Uuid;
    use chrono::{Utc, Duration};

    fn make_object(priority: PriorityClass, size_bytes: u64, pinned: bool, days_old: i64) -> CacheObject {
        CacheObject {
            id: Uuid::new_v4(),
            name: format!("{:?}-object", priority),
            priority,
            size_bytes,
            pinned,
            content_hash: "abc123".into(),
            source_url: "https://example.com/file".into(),
            last_accessed_at: Utc::now() - Duration::days(days_old),
            created_at: Utc::now() - Duration::days(days_old + 1),
        }
    }

    // Score tests
    #[test]
    fn p0_score_always_above_threshold() {
        let obj = make_object(PriorityClass::P0, 1024, false, 0);
        let s = score::score_object(&obj);
        assert!(s.score >= 10_000.0, "P0 score must be >= 10000, got {}", s.score);
    }

    #[test]
    fn p4_score_below_p0() {
        let p0 = make_object(PriorityClass::P0, 1024, false, 0);
        let p4 = make_object(PriorityClass::P4, 1024, false, 0);
        assert!(score::score_object(&p0).score > score::score_object(&p4).score);
    }

    #[test]
    fn large_object_penalized() {
        let small = make_object(PriorityClass::P2, 1024, false, 0);             // 1 KB
        let large = make_object(PriorityClass::P2, 10 * 1024 * 1024 * 1024, false, 0); // 10 GB
        assert!(score::score_object(&small).score > score::score_object(&large).score);
    }

    #[test]
    fn stale_object_penalized() {
        let fresh = make_object(PriorityClass::P2, 1024, false, 0);
        let stale = make_object(PriorityClass::P2, 1024, false, 90);
        assert!(score::score_object(&fresh).score > score::score_object(&stale).score);
    }

    #[test]
    fn pinned_object_immune_from_eviction() {
        let pinned = make_object(PriorityClass::P4, 1024, true, 0);
        assert!(!constraints::can_evict(&pinned));
    }

    #[test]
    fn p0_object_immune_from_eviction() {
        let p0 = make_object(PriorityClass::P0, 1024, false, 0);
        assert!(!constraints::can_evict(&p0));
    }

    #[test]
    fn score_is_deterministic() {
        let obj = make_object(PriorityClass::P1, 1024 * 1024, false, 5);
        let scores: Vec<f64> = (0..100).map(|_| score::score_object(&obj).score).collect();
        let first = scores[0];
        assert!(scores.iter().all(|&s| (s - first).abs() < 0.001));
    }

    // Simulator tests
    #[test]
    fn simulate_respects_capacity_target() {
        let objects = vec![
            make_object(PriorityClass::P4, 100 * 1024 * 1024, false, 30),
            make_object(PriorityClass::P4, 200 * 1024 * 1024, false, 60),
            make_object(PriorityClass::P3, 50 * 1024 * 1024, false, 10),
        ];
        let target = 150 * 1024 * 1024u64;
        let plan = simulator::simulate(&objects, target).unwrap();
        let freed: u64 = plan.iter().map(|c| c.size_bytes).sum();
        assert!(freed >= target);
    }

    #[test]
    fn simulate_never_evicts_pinned() {
        let objects = vec![
            make_object(PriorityClass::P4, 100 * 1024 * 1024, true, 30),
            make_object(PriorityClass::P4, 200 * 1024 * 1024, false, 60),
        ];
        let plan = simulator::simulate(&objects, 50 * 1024 * 1024).unwrap();
        for candidate in &plan {
            let obj = objects.iter().find(|o| o.id == candidate.object_id).unwrap();
            assert!(!obj.pinned, "Pinned object should not be in eviction plan");
        }
    }

    #[test]
    fn simulate_evicts_p4_before_p3() {
        let p4 = make_object(PriorityClass::P4, 100 * 1024 * 1024, false, 30);
        let p3 = make_object(PriorityClass::P3, 100 * 1024 * 1024, false, 30);
        let objects = vec![p4.clone(), p3.clone()];
        let plan = simulator::simulate(&objects, 50 * 1024 * 1024).unwrap();
        assert_eq!(plan[0].object_id, p4.id, "P4 should be evicted first");
    }

    #[test]
    fn simulate_empty_cache_returns_empty() {
        let result = simulator::simulate(&[], 1024).unwrap();
        assert!(result.is_empty());
    }
}
