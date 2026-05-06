package models

import "time"

// Org represents a tenant organisation.
type Org struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Plan      string    `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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

// CachePolicy defines caching rules for a site.
type CachePolicy struct {
	ID              string    `json:"id"`
	SiteID          string    `json:"site_id"`
	Name            string    `json:"name"`
	Rules           string    `json:"rules"` // JSON-encoded rule set
	DefaultTTL      int       `json:"default_ttl"`
	MaxObjectSize   int64     `json:"max_object_size"`
	EvictionStrategy string   `json:"eviction_strategy"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PreloadJob is a scheduled job to prefetch content to a node.
type PreloadJob struct {
	ID          string    `json:"id"`
	SiteID      string    `json:"site_id"`
	NodeID      string    `json:"node_id"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PreloadJobItem is a single URL within a PreloadJob.
type PreloadJobItem struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	SizeBytes  int64     `json:"size_bytes"`
	Error      string    `json:"error,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EvictionRun records an eviction pass on a node.
type EvictionRun struct {
	ID            string    `json:"id"`
	NodeID        string    `json:"node_id"`
	Strategy      string    `json:"strategy"`
	BytesFreed    int64     `json:"bytes_freed"`
	ObjectsEvicted int      `json:"objects_evicted"`
	Duration      int       `json:"duration_ms"`
	CreatedAt     time.Time `json:"created_at"`
}

// AuditLog records user-initiated actions for compliance.
type AuditLog struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"org_id"`
	ActorID    string    `json:"actor_id"`
	Action     string    `json:"action"`
	ResourceID string    `json:"resource_id"`
	Resource   string    `json:"resource"`
	Meta       string    `json:"meta"` // JSON-encoded metadata
	CreatedAt  time.Time `json:"created_at"`
}
