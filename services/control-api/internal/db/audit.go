package db

import "context"

func (c *Client) WriteAuditLog(ctx context.Context, orgID, actorID, action, resourceType, resourceID string) error {
	_, err := c.Pool.Exec(ctx,
		`INSERT INTO audit_logs (org_id, actor_id, action, resource_type, resource_id)
         VALUES ($1, NULLIF($2,'')::uuid, $3, $4, NULLIF($5,'')::uuid)`,
		orgID, actorID, action, resourceType, resourceID,
	)
	return err
}
