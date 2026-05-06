package tools

import (
	"context"
	"encoding/json"
)

// EvictionSimulationOutput lists cache object IDs that would be evicted.
type EvictionSimulationOutput struct {
	Candidates []string `json:"candidates"`
}

// SimulateEviction predicts which objects would be removed from cache to meet a
// space target. Returns an empty candidate list in this stub; Phase 2 will
// delegate to the policy-engine's simulate endpoint.
func SimulateEviction(_ context.Context, input json.RawMessage) (*EvictionSimulationOutput, error) {
	return &EvictionSimulationOutput{Candidates: []string{}}, nil
}
