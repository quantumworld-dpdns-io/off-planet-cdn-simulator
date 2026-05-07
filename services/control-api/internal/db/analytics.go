package db

import (
	"context"
	"time"
)

type CacheHitPoint struct {
	Hour    string `json:"hour"`
	Hits    int    `json:"hits"`
	Misses  int    `json:"misses"`
}

type PriorityBucket struct {
	Level      string `json:"level"`
	Count      int    `json:"count"`
	TotalBytes int64  `json:"total_bytes"`
}

type NodeFill struct {
	NodeID    string `json:"node_id"`
	NodeName  string `json:"node_name"`
	UsedBytes int64  `json:"used_bytes"`
	MaxBytes  int64  `json:"max_bytes"`
}

func (c *Client) CacheHitTimeseries(ctx context.Context, orgID string) ([]CacheHitPoint, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT date_trunc('hour', created_at) AS hour, event_type, COUNT(*) AS cnt
		FROM telemetry_events
		WHERE org_id = $1
		  AND event_type IN ('cache_hit', 'cache_miss')
		  AND created_at > NOW() - INTERVAL '24 hours'
		GROUP BY hour, event_type
		ORDER BY hour ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	interim := map[time.Time]*CacheHitPoint{}
	var order []time.Time

	for rows.Next() {
		var hour time.Time
		var eventType string
		var cnt int
		if err := rows.Scan(&hour, &eventType, &cnt); err != nil {
			return nil, err
		}
		if _, ok := interim[hour]; !ok {
			interim[hour] = &CacheHitPoint{Hour: hour.UTC().Format(time.RFC3339)}
			order = append(order, hour)
		}
		switch eventType {
		case "cache_hit":
			interim[hour].Hits = cnt
		case "cache_miss":
			interim[hour].Misses = cnt
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	points := make([]CacheHitPoint, 0, len(order))
	for _, h := range order {
		points = append(points, *interim[h])
	}
	return points, nil
}

func (c *Client) PriorityDistribution(ctx context.Context, orgID string) ([]PriorityBucket, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT COALESCE(pc.level::text, 'UNKNOWN'), COUNT(*), COALESCE(SUM(co.size_bytes), 0)
		FROM cache_objects co
		LEFT JOIN priority_classes pc ON co.priority_class_id = pc.id
		WHERE co.org_id = $1
		  AND co.status = 'ACTIVE'
		GROUP BY pc.level
		ORDER BY pc.level ASC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buckets := []PriorityBucket{}
	for rows.Next() {
		var b PriorityBucket
		if err := rows.Scan(&b.Level, &b.Count, &b.TotalBytes); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}

func (c *Client) NodeFillSummary(ctx context.Context, orgID string) ([]NodeFill, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT id, name, cache_used_bytes, cache_max_bytes
		FROM nodes
		WHERE org_id = $1
		  AND status = 'ONLINE'
		ORDER BY cache_used_bytes DESC
		LIMIT 20`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fills := []NodeFill{}
	for rows.Next() {
		var f NodeFill
		if err := rows.Scan(&f.NodeID, &f.NodeName, &f.UsedBytes, &f.MaxBytes); err != nil {
			return nil, err
		}
		fills = append(fills, f)
	}
	return fills, rows.Err()
}
