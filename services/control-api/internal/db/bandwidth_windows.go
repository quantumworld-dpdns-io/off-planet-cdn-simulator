package db

import (
	"context"

	"github.com/off-planet-cdn/control-api/internal/models"
)

func (c *Client) ListBandwidthWindows(ctx context.Context, orgID, siteID string) ([]models.BandwidthWindow, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT id, org_id, COALESCE(site_id::text,''), COALESCE(label,''),
		       window_start, window_end, bandwidth_bps, reliability_score, created_at
		FROM bandwidth_windows
		WHERE org_id = $1
		  AND ($2 = '' OR site_id::text = $2)
		ORDER BY window_start ASC`,
		orgID, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var windows []models.BandwidthWindow
	for rows.Next() {
		var w models.BandwidthWindow
		if err := rows.Scan(&w.ID, &w.OrgID, &w.SiteID, &w.Label,
			&w.WindowStart, &w.WindowEnd, &w.BandwidthBps, &w.ReliabilityScore, &w.CreatedAt); err != nil {
			return nil, err
		}
		windows = append(windows, w)
	}
	return windows, rows.Err()
}

func (c *Client) CreateBandwidthWindow(ctx context.Context, orgID, siteID, label string, windowStart, windowEnd string, bandwidthBps int64, reliabilityScore float64) (*models.BandwidthWindow, error) {
	var w models.BandwidthWindow
	err := c.Pool.QueryRow(ctx, `
		INSERT INTO bandwidth_windows (org_id, site_id, label, window_start, window_end, bandwidth_bps, reliability_score)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,''), $4::timestamptz, $5::timestamptz, $6, $7)
		RETURNING id, org_id, COALESCE(site_id::text,''), COALESCE(label,''),
		          window_start, window_end, bandwidth_bps, reliability_score, created_at`,
		orgID, siteID, label, windowStart, windowEnd, bandwidthBps, reliabilityScore,
	).Scan(&w.ID, &w.OrgID, &w.SiteID, &w.Label,
		&w.WindowStart, &w.WindowEnd, &w.BandwidthBps, &w.ReliabilityScore, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
