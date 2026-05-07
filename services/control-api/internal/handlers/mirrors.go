package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

type MirrorHandler struct{ DB *db.Client }

func (h *MirrorHandler) ListSources(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	sources, err := h.DB.ListMirrorSources(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if sources == nil {
		sources = []models.MirrorSource{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"sources": sources})
}

func (h *MirrorHandler) CreateSource(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())

	var body struct {
		RegistryType string `json:"registry_type"`
		UpstreamURL  string `json:"upstream_url"`
		Label        string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if body.RegistryType == "" || body.UpstreamURL == "" {
		http.Error(w, `{"error":"registry_type and upstream_url are required"}`, http.StatusBadRequest)
		return
	}

	source, err := h.DB.CreateMirrorSource(r.Context(), orgID, body.RegistryType, body.UpstreamURL, body.Label)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(source)
}

func (h *MirrorHandler) ListArtifacts(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	sourceID := r.URL.Query().Get("source_id")

	artifacts, err := h.DB.ListMirrorArtifacts(r.Context(), orgID, sourceID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if artifacts == nil {
		artifacts = []models.MirrorArtifact{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"artifacts": artifacts})
}
