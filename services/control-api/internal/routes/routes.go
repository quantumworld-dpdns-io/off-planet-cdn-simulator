package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/handlers"
	"github.com/off-planet-cdn/control-api/internal/middleware"
)

func Register(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.OtelTrace)
	r.Get("/v1/health", handlers.Health)
	// Phase 1 routes registered here later
}
