use shared::types::{CacheObject, PriorityClass};

pub fn can_evict(obj: &CacheObject) -> bool {
    if obj.pinned {
        return false;
    }
    if obj.priority == PriorityClass::P0 {
        return false;
    }
    true
}
