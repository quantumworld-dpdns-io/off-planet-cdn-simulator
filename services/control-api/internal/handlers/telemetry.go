package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/off-planet-cdn/control-api/internal/db"
)

type TelemetryHandler struct{ DB *db.Client }

func (h *TelemetryHandler) IngestEvents(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	var events []struct {
		SiteID    string                 `json:"site_id"`
		NodeID    string                 `json:"node_id"`
		EventType string                 `json:"event_type"`
		Payload   map[string]interface{} `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	accepted := 0
	for _, e := range events {
		if e.EventType == "" {
			continue
		}
		if e.Payload == nil {
			e.Payload = map[string]interface{}{}
		}
		if err := h.DB.WritetelemetryEvent(r.Context(), orgID, e.SiteID, e.NodeID, e.EventType, e.Payload); err == nil {
			accepted++
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"accepted": accepted})
}
