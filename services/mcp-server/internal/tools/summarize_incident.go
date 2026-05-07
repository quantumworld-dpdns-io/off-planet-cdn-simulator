package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type summarizeIncidentInput struct {
	Limit int `json:"limit"`
}

type auditLogsResponse struct {
	Logs []auditLogEntry `json:"logs"`
}

type auditLogEntry struct {
	ID           string `json:"id"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	CreatedAt    string `json:"created_at"`
}

// SummarizeIncident produces a human-readable summary of recent audit log events.
func SummarizeIncident(_ context.Context, input json.RawMessage) (string, error) {
	var in summarizeIncidentInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}

	var logs auditLogsResponse
	if err := apiGet(context.Background(), fmt.Sprintf("/v1/audit-logs?limit=%d", in.Limit), &logs); err != nil {
		return "", err
	}

	if len(logs.Logs) == 0 {
		return "No recent audit events found.", nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Recent %d audit events:\n", len(logs.Logs))
	for _, entry := range logs.Logs {
		fmt.Fprintf(&sb, "- %s %s %s at %s\n", entry.Action, entry.ResourceType, entry.ResourceID, entry.CreatedAt)
	}
	return sb.String(), nil
}
