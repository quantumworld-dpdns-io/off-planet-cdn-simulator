package tools

import (
	"context"
	"encoding/json"
)

// GeneratePreloadPlanOutput is the result of a preload plan generation.
type GeneratePreloadPlanOutput struct {
	URLs []string `json:"urls"`
}

// GeneratePreloadPlan produces a prioritized list of URLs to prefetch for a node.
// Returns an empty plan in this stub; Phase 2 will integrate the scheduler and
// policy engine to produce an optimised plan.
func GeneratePreloadPlan(_ context.Context, input json.RawMessage) (*GeneratePreloadPlanOutput, error) {
	return &GeneratePreloadPlanOutput{URLs: []string{}}, nil
}
