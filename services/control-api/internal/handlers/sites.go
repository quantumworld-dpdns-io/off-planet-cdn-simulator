package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

type SiteHandler struct{ DB *db.Client }

func (h *SiteHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	if orgID == "" {
		http.Error(w, `{"error":"missing org_id"}`, http.StatusUnauthorized)
		return
	}
	sites, err := h.DB.ListSites(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if sites == nil {
		sites = []models.Site{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"sites": sites})
}

func (h *SiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	if orgID == "" {
		http.Error(w, `{"error":"missing org_id"}`, http.StatusUnauthorized)
		return
	}
	var body struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	site, err := h.DB.CreateSite(r.Context(), orgID, body.Name, body.Location, body.Description)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

func (h *SiteHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	siteID := chi.URLParam(r, "site_id")
	site, err := h.DB.GetSite(r.Context(), orgID, siteID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}
