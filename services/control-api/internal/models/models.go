package models

import "time"

type Site struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Node struct {
	ID             string     `json:"id"`
	OrgID          string     `json:"org_id"`
	SiteID         string     `json:"site_id"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	CacheDir       string     `json:"cache_dir"`
	CacheMaxBytes  int64      `json:"cache_max_bytes"`
	CacheUsedBytes int64      `json:"cache_used_bytes"`
	LastSeen       *time.Time `json:"last_seen,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type NodeHeartbeat struct {
	ID             string    `json:"id"`
	OrgID          string    `json:"org_id"`
	NodeID         string    `json:"node_id"`
	Status         string    `json:"status"`
	CacheUsedBytes int64     `json:"cache_used_bytes"`
	CacheMaxBytes  int64     `json:"cache_max_bytes"`
	AgentVersion   string    `json:"agent_version"`
	CreatedAt      time.Time `json:"created_at"`
}

type CacheObject struct {
	ID              string    `json:"id"`
	OrgID           string    `json:"org_id"`
	SiteID          string    `json:"site_id"`
	PriorityClassID string    `json:"priority_class_id"`
	Name            string    `json:"name"`
	ContentType     string    `json:"content_type,omitempty"`
	SourceURL       string    `json:"source_url,omitempty"`
	ContentHash     string    `json:"content_hash,omitempty"`
	SizeBytes       int64     `json:"size_bytes"`
	Pinned          bool      `json:"pinned"`
	Status          string    `json:"status"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CachePolicy struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	SiteID      string    `json:"site_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PreloadJob struct {
	ID                   string     `json:"id"`
	OrgID                string     `json:"org_id"`
	SiteID               string     `json:"site_id"`
	Name                 string     `json:"name"`
	Status               string     `json:"status"`
	BandwidthBudgetBytes *int64     `json:"bandwidth_budget_bytes,omitempty"`
	StartedAt            *time.Time `json:"started_at,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type PreloadJobItem struct {
	ID               string    `json:"id"`
	OrgID            string    `json:"org_id"`
	JobID            string    `json:"job_id"`
	ObjectID         string    `json:"object_id"`
	Status           string    `json:"status"`
	BytesTransferred int64     `json:"bytes_transferred"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type EvictionRun struct {
	ID               string    `json:"id"`
	OrgID            string    `json:"org_id"`
	NodeID           string    `json:"node_id"`
	Status           string    `json:"status"`
	TargetFreedBytes int64     `json:"target_freed_bytes"`
	ActualFreedBytes int64     `json:"actual_freed_bytes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID           string    `json:"id"`
	OrgID        string    `json:"org_id"`
	ActorID      string    `json:"actor_id,omitempty"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type BandwidthWindow struct {
	ID               string    `json:"id"`
	OrgID            string    `json:"org_id"`
	SiteID           string    `json:"site_id,omitempty"`
	Label            string    `json:"label,omitempty"`
	WindowStart      time.Time `json:"window_start"`
	WindowEnd        time.Time `json:"window_end"`
	BandwidthBps     int64     `json:"bandwidth_bps"`
	ReliabilityScore float64   `json:"reliability_score"`
	CreatedAt        time.Time `json:"created_at"`
}
