package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "control-api",
		"version": "0.1.0",
	})
}

// BandwidthWindowHandler handles bandwidth window CRUD endpoints.
type BandwidthWindowHandler struct{ DB *db.Client }

func (h *BandwidthWindowHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	siteID := r.URL.Query().Get("site_id")
	windows, err := h.DB.ListBandwidthWindows(r.Context(), orgID, siteID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if windows == nil {
		windows = []models.BandwidthWindow{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"windows": windows})
}

func (h *BandwidthWindowHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	var body struct {
		SiteID           string  `json:"site_id"`
		Label            string  `json:"label"`
		WindowStart      string  `json:"window_start"`
		WindowEnd        string  `json:"window_end"`
		BandwidthBps     int64   `json:"bandwidth_bps"`
		ReliabilityScore float64 `json:"reliability_score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.WindowStart == "" || body.WindowEnd == "" || body.BandwidthBps <= 0 {
		http.Error(w, `{"error":"window_start, window_end, and bandwidth_bps are required"}`, http.StatusBadRequest)
		return
	}
	if body.ReliabilityScore == 0 {
		body.ReliabilityScore = 1.0
	}
	win, err := h.DB.CreateBandwidthWindow(r.Context(), orgID, body.SiteID, body.Label, body.WindowStart, body.WindowEnd, body.BandwidthBps, body.ReliabilityScore)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(win)
}
