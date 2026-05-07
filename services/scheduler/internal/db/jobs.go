package db

import (
	"context"
	"time"
)

type PendingJob struct {
	ID                   string
	OrgID                string
	SiteID               string
	Name                 string
	BandwidthBudgetBytes *int64
}

type JobItem struct {
	ID        string
	ObjectID  string
	SourceURL string
	Priority  string // P0..P5
	SizeBytes int64
}

type Node struct {
	ID   string
	Name string
}

// ListPendingJobs returns PENDING preload jobs ordered by created_at ASC (oldest first).
func (c *Client) ListPendingJobs(ctx context.Context) ([]PendingJob, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, site_id, name, bandwidth_budget_bytes
         FROM preload_jobs
         WHERE status = 'PENDING'
         ORDER BY created_at ASC
         LIMIT 50`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []PendingJob
	for rows.Next() {
		var j PendingJob
		if err := rows.Scan(&j.ID, &j.OrgID, &j.SiteID, &j.Name, &j.BandwidthBudgetBytes); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

// GetJobItems returns the pending items for a job, enriched with cache object metadata.
func (c *Client) GetJobItems(ctx context.Context, jobID string) ([]JobItem, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT ji.id, ji.object_id,
                COALESCE(co.source_url, ''),
                COALESCE(pc.level::text, 'P4'),
                COALESCE(co.size_bytes, 0)
         FROM preload_job_items ji
         LEFT JOIN cache_objects co ON co.id = ji.object_id
         LEFT JOIN priority_classes pc ON pc.id = co.priority_class_id
         WHERE ji.job_id = $1 AND ji.status = 'PENDING'`,
		jobID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JobItem
	for rows.Next() {
		var item JobItem
		if err := rows.Scan(&item.ID, &item.ObjectID, &item.SourceURL, &item.Priority, &item.SizeBytes); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetOnlineNodesForSite returns ONLINE nodes for a site.
func (c *Client) GetOnlineNodesForSite(ctx context.Context, siteID string) ([]Node, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, name FROM nodes
         WHERE site_id = $1 AND status = 'ONLINE'
         ORDER BY last_seen DESC`,
		siteID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []Node
	for rows.Next() {
		var n Node
		if err := rows.Scan(&n.ID, &n.Name); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

// MarkJobRunning sets a preload job's status to RUNNING.
func (c *Client) MarkJobRunning(ctx context.Context, jobID string) error {
	now := time.Now().UTC()
	_, err := c.Pool.Exec(ctx,
		`UPDATE preload_jobs SET status = 'RUNNING', started_at = $1, updated_at = $1 WHERE id = $2`,
		now, jobID,
	)
	return err
}

// MarkJobDone sets a preload job's status to DONE.
func (c *Client) MarkJobDone(ctx context.Context, jobID string) error {
	now := time.Now().UTC()
	_, err := c.Pool.Exec(ctx,
		`UPDATE preload_jobs SET status = 'DONE', completed_at = $1, updated_at = $1 WHERE id = $2`,
		now, jobID,
	)
	return err
}

// MarkJobFailed sets a preload job's status to FAILED.
func (c *Client) MarkJobFailed(ctx context.Context, jobID string) error {
	now := time.Now().UTC()
	_, err := c.Pool.Exec(ctx,
		`UPDATE preload_jobs SET status = 'FAILED', completed_at = $1, updated_at = $1 WHERE id = $2`,
		now, jobID,
	)
	return err
}

// IsJobCancelled returns true if the job is already CANCELLED.
func (c *Client) IsJobCancelled(ctx context.Context, jobID string) (bool, error) {
	var status string
	err := c.Pool.QueryRow(ctx,
		`SELECT status::text FROM preload_jobs WHERE id = $1`, jobID,
	).Scan(&status)
	if err != nil {
		return false, err
	}
	return status == "CANCELLED", nil
}

// GetActiveBandwidthWindow returns the currently active bandwidth window for a site, if any.
func (c *Client) GetActiveBandwidthWindow(ctx context.Context, siteID string) (*int64, error) {
	var bps int64
	err := c.Pool.QueryRow(ctx,
		`SELECT bandwidth_bps FROM bandwidth_windows
         WHERE site_id = $1
           AND window_start <= now() AT TIME ZONE 'UTC'
           AND window_end >= now() AT TIME ZONE 'UTC'
         ORDER BY reliability_score DESC
         LIMIT 1`,
		siteID,
	).Scan(&bps)
	if err != nil {
		return nil, nil // no active window — not an error
	}
	return &bps, nil
}
