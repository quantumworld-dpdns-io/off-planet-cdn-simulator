package routes

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/handlers"
	"github.com/off-planet-cdn/control-api/internal/middleware"
	cdnredis "github.com/off-planet-cdn/control-api/internal/redis"
)

func Register(r *chi.Mux, dbClient *db.Client, redisClient *cdnredis.Client) {
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.OtelTrace)
	r.Use(middleware.OrgID)

	r.Get("/v1/health", handlers.Health)

	// Sites
	sites := &handlers.SiteHandler{DB: dbClient}
	r.Get("/v1/sites", sites.List)
	r.Post("/v1/sites", sites.Create)
	r.Get("/v1/sites/{site_id}", sites.Get)

	// Nodes
	nodes := &handlers.NodeHandler{DB: dbClient}
	r.Get("/v1/nodes", nodes.List)
	r.Get("/v1/nodes/{node_id}", nodes.Get)
	r.Post("/v1/nodes/register", nodes.Register)
	r.Post("/v1/nodes/{node_id}/heartbeat", nodes.Heartbeat)
	r.Get("/v1/nodes/{node_id}/status", nodes.Status)

	// Cache objects — frontend uses /v1/objects
	objects := &handlers.CacheObjectHandler{DB: dbClient}
	r.Get("/v1/objects", objects.List)
	r.Post("/v1/objects", objects.Create)
	r.Get("/v1/objects/{object_id}", objects.Get)
	r.Post("/v1/objects/{object_id}/pin", objects.Pin)
	r.Post("/v1/objects/{object_id}/unpin", objects.Unpin)

	// Policies
	policies := &handlers.PolicyHandler{DB: dbClient}
	r.Get("/v1/policies", policies.List)
	r.Post("/v1/policies", policies.Create)
	r.Get("/v1/policies/{policy_id}", policies.Get)
	r.Put("/v1/policies/{policy_id}", policies.Update)

	// Preload jobs — frontend uses /v1/preload-jobs (hyphen, not slash)
	preloadJobs := &handlers.PreloadJobHandler{DB: dbClient, Redis: redisClient}
	r.Get("/v1/preload-jobs", preloadJobs.List)
	r.Post("/v1/preload-jobs", preloadJobs.Create)
	r.Get("/v1/preload-jobs/{job_id}", preloadJobs.Get)
	r.Post("/v1/preload-jobs/{job_id}/cancel", preloadJobs.Cancel)

	// Bandwidth windows
	bwWindows := &handlers.BandwidthWindowHandler{DB: dbClient}
	r.Get("/v1/bandwidth-windows", bwWindows.List)
	r.Post("/v1/bandwidth-windows", bwWindows.Create)

	// Analytics
	analytics := &handlers.AnalyticsHandler{DB: dbClient}
	r.Get("/v1/analytics/cache-hits", analytics.CacheHits)
	r.Get("/v1/analytics/priority-distribution", analytics.PriorityDistribution)
	r.Get("/v1/analytics/node-fill", analytics.NodeFill)

	// Telemetry & audit
	telemetry := &handlers.TelemetryHandler{DB: dbClient}
	r.Post("/v1/telemetry/events", telemetry.IngestEvents)

	auditLogs := &handlers.AuditLogHandler{DB: dbClient}
	r.Get("/v1/audit-logs", auditLogs.List)
}
