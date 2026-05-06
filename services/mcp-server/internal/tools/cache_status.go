package tools

import "context"

type CacheStatusInput struct {
	NodeID string `json:"node_id"`
}

type CacheStatusOutput struct {
	FillRatio   float64  `json:"fill_ratio"`
	PinnedCount int      `json:"pinned_count"`
	TopObjects  []string `json:"top_objects"`
}

func CacheStatus(ctx context.Context, input CacheStatusInput) (*CacheStatusOutput, error) {
	// TODO: query control API
	return &CacheStatusOutput{FillRatio: 0.0, PinnedCount: 0, TopObjects: []string{}}, nil
}
