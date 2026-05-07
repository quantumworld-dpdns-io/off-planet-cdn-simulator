package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/off-planet-cdn/control-api/internal/db"
)

type AnalyticsHandler struct{ DB *db.Client }

func (h *AnalyticsHandler) CacheHits(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	points, err := h.DB.CacheHitTimeseries(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"points": points})
}

func (h *AnalyticsHandler) PriorityDistribution(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	buckets, err := h.DB.PriorityDistribution(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"buckets": buckets})
}

func (h *AnalyticsHandler) NodeFill(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	nodes, err := h.DB.NodeFillSummary(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"nodes": nodes})
}
