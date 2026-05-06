package offcdn

import "time"

// Site represents a Moon/Mars habitat deployment site.
type Site struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Node is a CDN edge node within a site.
type Node struct {
	ID            string    `json:"id"`
	SiteID        string    `json:"site_id"`
	Hostname      string    `json:"hostname"`
	Status        string    `json:"status"`
	CapacityBytes int64     `json:"capacity_bytes"`
	UsedBytes     int64     `json:"used_bytes"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CacheObject represents a cached content item on a node.
type CacheObject struct {
	ID          string    `json:"id"`
	NodeID      string    `json:"node_id"`
	Key         string    `json:"key"`
	URL         string    `json:"url"`
	SizeBytes   int64     `json:"size_bytes"`
	Priority    int       `json:"priority"`
	Pinned      bool      `json:"pinned"`
	Score       float64   `json:"score"`
	ContentType string    `json:"content_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PreloadJob is a scheduled job to prefetch content to a node.
type PreloadJob struct {
	ID          string     `json:"id"`
	SiteID      string     `json:"site_id"`
	NodeID      string     `json:"node_id"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Policy defines caching rules for a site.
type Policy struct {
	ID               string    `json:"id"`
	SiteID           string    `json:"site_id"`
	Name             string    `json:"name"`
	Rules            string    `json:"rules"`
	DefaultTTL       int       `json:"default_ttl"`
	MaxObjectSize    int64     `json:"max_object_size"`
	EvictionStrategy string    `json:"eviction_strategy"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
