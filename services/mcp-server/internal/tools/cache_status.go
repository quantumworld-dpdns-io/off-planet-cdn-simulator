package tools

import (
	"context"
	"fmt"
)

type CacheStatusInput struct {
	NodeID string `json:"node_id"`
}

type CacheStatusOutput struct {
	FillRatio   float64  `json:"fill_ratio"`
	PinnedCount int      `json:"pinned_count"`
	TopObjects  []string `json:"top_objects"`
}

// nodeResponse is the minimal shape returned by GET /v1/nodes/{id}.
type nodeResponse struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	SiteID         string `json:"site_id"`
	CacheUsedBytes int64  `json:"cache_used_bytes"`
	CacheMaxBytes  int64  `json:"cache_max_bytes"`
}

// objectsListResponse is the minimal shape returned by GET /v1/objects.
type objectsListResponse struct {
	Objects []objectItem `json:"objects"`
}

type objectItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SiteID    string `json:"site_id"`
	Pinned    bool   `json:"pinned"`
	SizeBytes int64  `json:"size_bytes"`
	Status    string `json:"status"`
	SourceURL string `json:"source_url"`
}

func CacheStatus(ctx context.Context, input CacheStatusInput) (*CacheStatusOutput, error) {
	// 1. Fetch node info.
	var node nodeResponse
	if err := apiGet(ctx, fmt.Sprintf("/v1/nodes/%s", input.NodeID), &node); err != nil {
		return nil, err
	}

	// 2. Fetch objects list.
	var objList objectsListResponse
	if err := apiGet(ctx, "/v1/objects?limit=100", &objList); err != nil {
		return nil, err
	}

	// 3. Filter objects belonging to this node's site.
	var fillRatio float64
	if node.CacheMaxBytes > 0 {
		fillRatio = float64(node.CacheUsedBytes) / float64(node.CacheMaxBytes)
	}

	var pinnedCount int
	var topObjects []string
	for _, obj := range objList.Objects {
		if obj.SiteID != node.SiteID {
			continue
		}
		if obj.Pinned {
			pinnedCount++
		}
		if len(topObjects) < 5 {
			topObjects = append(topObjects, obj.Name)
		}
	}

	if topObjects == nil {
		topObjects = []string{}
	}

	return &CacheStatusOutput{
		FillRatio:   fillRatio,
		PinnedCount: pinnedCount,
		TopObjects:  topObjects,
	}, nil
}
