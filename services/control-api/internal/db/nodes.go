package db

import (
	"context"
	"time"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) ListNodes(ctx context.Context, orgID string) ([]models.Node, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, site_id, name, status::text, cache_dir, cache_max_bytes, cache_used_bytes, last_seen, created_at, updated_at
         FROM nodes WHERE org_id = $1 ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []models.Node
	for rows.Next() {
		var n models.Node
		if err := rows.Scan(&n.ID, &n.OrgID, &n.SiteID, &n.Name, &n.Status, &n.CacheDir, &n.CacheMaxBytes, &n.CacheUsedBytes, &n.LastSeen, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

func (c *Client) RegisterNode(ctx context.Context, orgID, siteID, name, cacheDir string, cacheMaxBytes int64) (*models.Node, error) {
	var n models.Node
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO nodes (org_id, site_id, name, cache_dir, cache_max_bytes, status)
         VALUES ($1, $2, $3, $4, $5, 'UNKNOWN')
         RETURNING id, org_id, site_id, name, status::text, cache_dir, cache_max_bytes, cache_used_bytes, last_seen, created_at, updated_at`,
		orgID, siteID, name, cacheDir, cacheMaxBytes,
	).Scan(&n.ID, &n.OrgID, &n.SiteID, &n.Name, &n.Status, &n.CacheDir, &n.CacheMaxBytes, &n.CacheUsedBytes, &n.LastSeen, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (c *Client) GetNodeStatus(ctx context.Context, orgID, nodeID string) (*models.Node, error) {
	var n models.Node
	err := c.Pool.QueryRow(ctx,
		`SELECT id, org_id, site_id, name, status::text, cache_dir, cache_max_bytes, cache_used_bytes, last_seen, created_at, updated_at
         FROM nodes WHERE id = $1 AND org_id = $2`,
		nodeID, orgID,
	).Scan(&n.ID, &n.OrgID, &n.SiteID, &n.Name, &n.Status, &n.CacheDir, &n.CacheMaxBytes, &n.CacheUsedBytes, &n.LastSeen, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (c *Client) RecordHeartbeat(ctx context.Context, orgID, nodeID, status, agentVersion string, cacheUsedBytes, cacheMaxBytes int64) (*models.NodeHeartbeat, error) {
	now := time.Now().UTC()

	// Insert heartbeat row
	var hb models.NodeHeartbeat
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO node_heartbeats (org_id, node_id, status, cache_used_bytes, cache_max_bytes, agent_version, created_at)
         VALUES ($1, $2, $3::node_status, $4, $5, $6, $7)
         RETURNING id, org_id, node_id, status::text, cache_used_bytes, cache_max_bytes, COALESCE(agent_version,''), created_at`,
		orgID, nodeID, status, cacheUsedBytes, cacheMaxBytes, agentVersion, now,
	).Scan(&hb.ID, &hb.OrgID, &hb.NodeID, &hb.Status, &hb.CacheUsedBytes, &hb.CacheMaxBytes, &hb.AgentVersion, &hb.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update node last_seen, status, cache_used_bytes
	_, err = c.Pool.Exec(ctx,
		`UPDATE nodes SET last_seen = $1, status = $2::node_status, cache_used_bytes = $3, updated_at = $1
         WHERE id = $4 AND org_id = $5`,
		now, status, cacheUsedBytes, nodeID, orgID,
	)
	if err != nil {
		return nil, err
	}

	return &hb, nil
}
