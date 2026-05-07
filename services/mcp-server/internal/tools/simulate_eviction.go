package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

// EvictionSimulationOutput lists cache object IDs that would be evicted.
type EvictionSimulationOutput struct {
	Candidates []string `json:"candidates"`
}

type simulateEvictionInput struct {
	SiteID   string  `json:"site_id"`
	TargetMB float64 `json:"target_mb"`
}

// SimulateEviction predicts which objects would be removed from cache to meet a space target.
func SimulateEviction(ctx context.Context, input json.RawMessage) (*EvictionSimulationOutput, error) {
	var in simulateEvictionInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}
	if in.TargetMB <= 0 {
		in.TargetMB = 100
	}

	path := "/v1/objects?status=ACTIVE"
	if in.SiteID != "" {
		path += fmt.Sprintf("&site_id=%s", in.SiteID)
	}

	var objList objectsListResponse
	if err := apiGet(ctx, path, &objList); err != nil {
		return nil, err
	}

	// Filter out pinned objects.
	var evictable []objectItem
	for _, obj := range objList.Objects {
		if !obj.Pinned {
			evictable = append(evictable, obj)
		}
	}

	// Sort by size_bytes descending.
	sort.Slice(evictable, func(i, j int) bool {
		return evictable[i].SizeBytes > evictable[j].SizeBytes
	})

	// Accumulate until target bytes met.
	targetBytes := int64(in.TargetMB * 1024 * 1024)
	var accumulated int64
	candidates := []string{}
	for _, obj := range evictable {
		if accumulated >= targetBytes {
			break
		}
		candidates = append(candidates, obj.ID)
		accumulated += obj.SizeBytes
	}

	return &EvictionSimulationOutput{Candidates: candidates}, nil
}
