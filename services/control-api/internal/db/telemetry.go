package db

import (
	"context"
	"encoding/json"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) WritetelemetryEvent(ctx context.Context, orgID, siteID, nodeID, eventType string, payload map[string]interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		payloadBytes = []byte("{}")
	}
	_, err = c.Pool.Exec(ctx,
		`INSERT INTO telemetry_events (org_id, site_id, node_id, event_type, payload)
         VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid, $4, $5)`,
		orgID, siteID, nodeID, eventType, payloadBytes,
	)
	return err
}

func (c *Client) ListAuditLogs(ctx context.Context, orgID, actorID, action, resourceType string, limit, offset int) ([]models.AuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, COALESCE(actor_id::text,''), action, resource_type,
                COALESCE(resource_id::text,''), created_at
         FROM audit_logs
         WHERE org_id = $1
           AND ($2 = '' OR actor_id::text = $2)
           AND ($3 = '' OR action = $3)
           AND ($4 = '' OR resource_type = $4)
         ORDER BY created_at DESC
         LIMIT $5 OFFSET $6`,
		orgID, actorID, action, resourceType, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.OrgID, &l.ActorID, &l.Action, &l.ResourceType, &l.ResourceID, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
