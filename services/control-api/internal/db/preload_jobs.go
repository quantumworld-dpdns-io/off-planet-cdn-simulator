package db

import (
	"context"
	"fmt"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) CreatePreloadJob(ctx context.Context, orgID, siteID, name string, bandwidthBudgetBytes *int64) (*models.PreloadJob, error) {
	var j models.PreloadJob
	err := c.Pool.QueryRow(ctx,
		`INSERT INTO preload_jobs (org_id, site_id, name, status, bandwidth_budget_bytes)
         VALUES ($1, $2, $3, 'PENDING', $4)
         RETURNING id, org_id, site_id, name, status::text, bandwidth_budget_bytes,
                   started_at, completed_at, created_at, updated_at`,
		orgID, siteID, name, bandwidthBudgetBytes,
	).Scan(&j.ID, &j.OrgID, &j.SiteID, &j.Name, &j.Status, &j.BandwidthBudgetBytes,
		&j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (c *Client) ListPreloadJobs(ctx context.Context, orgID string) ([]models.PreloadJob, error) {
	rows, err := c.Pool.Query(ctx,
		`SELECT id, org_id, site_id, name, status::text, bandwidth_budget_bytes,
                started_at, completed_at, created_at, updated_at
         FROM preload_jobs WHERE org_id = $1 ORDER BY created_at DESC LIMIT 200`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []models.PreloadJob
	for rows.Next() {
		var j models.PreloadJob
		if err := rows.Scan(&j.ID, &j.OrgID, &j.SiteID, &j.Name, &j.Status, &j.BandwidthBudgetBytes,
			&j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (c *Client) GetPreloadJob(ctx context.Context, orgID, jobID string) (*models.PreloadJob, error) {
	var j models.PreloadJob
	err := c.Pool.QueryRow(ctx,
		`SELECT id, org_id, site_id, name, status::text, bandwidth_budget_bytes,
                started_at, completed_at, created_at, updated_at
         FROM preload_jobs WHERE id = $1 AND org_id = $2`,
		jobID, orgID,
	).Scan(&j.ID, &j.OrgID, &j.SiteID, &j.Name, &j.Status, &j.BandwidthBudgetBytes,
		&j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (c *Client) CancelPreloadJob(ctx context.Context, orgID, jobID string) error {
	result, err := c.Pool.Exec(ctx,
		`UPDATE preload_jobs SET status = 'CANCELLED', updated_at = now()
         WHERE id = $1 AND org_id = $2 AND status = 'PENDING'`,
		jobID, orgID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("job not found or not cancellable")
	}
	return nil
}

func (c *Client) AddPreloadJobItems(ctx context.Context, orgID, jobID string, objectIDs []string) error {
	for _, objID := range objectIDs {
		_, err := c.Pool.Exec(ctx,
			`INSERT INTO preload_job_items (org_id, job_id, object_id, status)
             VALUES ($1, $2, $3, 'PENDING')`,
			orgID, jobID, objID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
