package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// NodeInfo holds diagnostic information about an edge node.
type NodeInfo struct {
	NodeID        string  `json:"node_id"`
	Status        string  `json:"status"`
	FillRatio     float64 `json:"fill_ratio"`
	CapacityBytes int64   `json:"capacity_bytes"`
	UsedBytes     int64   `json:"used_bytes"`
}

type inspectNodeInput struct {
	NodeID string `json:"node_id"`
}

// InspectNode returns live diagnostic details for the specified node.
func InspectNode(ctx context.Context, input json.RawMessage) (*NodeInfo, error) {
	var in inspectNodeInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	var node nodeResponse
	if err := apiGet(ctx, fmt.Sprintf("/v1/nodes/%s", in.NodeID), &node); err != nil {
		return nil, err
	}

	var fillRatio float64
	if node.CacheMaxBytes > 0 {
		fillRatio = float64(node.CacheUsedBytes) / float64(node.CacheMaxBytes)
	}

	return &NodeInfo{
		NodeID:        node.ID,
		Status:        node.Status,
		FillRatio:     fillRatio,
		CapacityBytes: node.CacheMaxBytes,
		UsedBytes:     node.CacheUsedBytes,
	}, nil
}
