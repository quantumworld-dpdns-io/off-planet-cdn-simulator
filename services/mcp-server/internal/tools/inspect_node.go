package tools

import (
	"context"
	"encoding/json"
)

// NodeInfo holds diagnostic information about an edge node.
type NodeInfo struct {
	NodeID        string  `json:"node_id"`
	Status        string  `json:"status"`
	FillRatio     float64 `json:"fill_ratio"`
	CapacityBytes int64   `json:"capacity_bytes"`
	UsedBytes     int64   `json:"used_bytes"`
}

// InspectNode returns live diagnostic details for the specified node.
// Returns empty NodeInfo in this stub; Phase 2 will query the control API.
func InspectNode(_ context.Context, input json.RawMessage) (*NodeInfo, error) {
	return &NodeInfo{}, nil
}
