package simulator

// Simulate predicts which cache objects on nodeID would be evicted to free
// targetBytes of space. Returns an empty slice in this stub; Phase 2 will
// implement the full scoring-based eviction simulation.
func Simulate(nodeID string, targetBytes int64) ([]string, error) {
	return []string{}, nil
}
