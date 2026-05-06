package routes

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/handlers"
	"github.com/off-planet-cdn/control-api/internal/middleware"
)

func Register(r *chi.Mux, dbClient *db.Client) {
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.OtelTrace)
	r.Use(middleware.OrgID) // reads X-Org-ID header and sets in context

	r.Get("/v1/health", handlers.Health)

	sites := &handlers.SiteHandler{DB: dbClient}
	r.Get("/v1/sites", sites.List)
	r.Post("/v1/sites", sites.Create)
	r.Get("/v1/sites/{site_id}", sites.Get)

	nodes := &handlers.NodeHandler{DB: dbClient}
	r.Get("/v1/nodes", nodes.List)
	r.Post("/v1/nodes/register", nodes.Register)
	r.Post("/v1/nodes/{node_id}/heartbeat", nodes.Heartbeat)
	r.Get("/v1/nodes/{node_id}/status", nodes.Status)
}
