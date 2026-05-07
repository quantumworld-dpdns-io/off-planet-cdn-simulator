package contactwindows

import (
	"context"
	"time"

	"github.com/off-planet-cdn/scheduler/internal/db"
)

// Checker evaluates whether a bandwidth/contact window is currently open for a site.
type Checker struct {
	DB *db.Client
}

func New(dbClient *db.Client) *Checker {
	return &Checker{DB: dbClient}
}

// IsWindowOpen returns true if a bandwidth window is currently active for the site.
// If no windows are configured for the site, it defaults to open (always-on mode for dev).
func (c *Checker) IsWindowOpen(ctx context.Context, siteID string) (bool, error) {
	bps, err := c.DB.GetActiveBandwidthWindow(ctx, siteID)
	if err != nil {
		return false, err
	}
	// No windows configured → treat as always open (dev/test mode)
	if bps == nil {
		return true, nil
	}
	// A window exists and is active
	return true, nil
}

// NextWindowAt returns when the next contact window opens for the site.
// Used for scheduling — returns zero time if unknown.
func (c *Checker) NextWindowAt(ctx context.Context, siteID string) (time.Time, error) {
	var nextAt time.Time
	err := c.DB.Pool.QueryRow(ctx,
		`SELECT COALESCE(next_window_at, now()) FROM contact_windows
         WHERE site_id = $1 ORDER BY next_window_at ASC LIMIT 1`,
		siteID,
	).Scan(&nextAt)
	if err != nil {
		return time.Time{}, nil
	}
	return nextAt, nil
}
