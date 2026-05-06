package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

type CacheObjectHandler struct{ DB *db.Client }

func (h *CacheObjectHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	q := r.URL.Query()
	objects, err := h.DB.ListObjects(r.Context(), orgID,
		q.Get("site_id"), q.Get("priority_class_id"), q.Get("tag"))
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if objects == nil {
		objects = []models.CacheObject{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"objects": objects})
}

func (h *CacheObjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	var body struct {
		SiteID          string   `json:"site_id"`
		PriorityClassID string   `json:"priority_class_id"`
		Name            string   `json:"name"`
		ContentType     string   `json:"content_type"`
		SourceURL       string   `json:"source_url"`
		SizeBytes       int64    `json:"size_bytes"`
		Tags            []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.SiteID == "" || body.PriorityClassID == "" {
		http.Error(w, `{"error":"name, site_id, and priority_class_id are required"}`, http.StatusBadRequest)
		return
	}
	obj, err := h.DB.CreateObject(r.Context(), orgID, body.SiteID, body.PriorityClassID,
		body.Name, body.ContentType, body.SourceURL, body.SizeBytes, body.Tags)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(obj)
}

func (h *CacheObjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	objectID := chi.URLParam(r, "object_id")
	obj, err := h.DB.GetObject(r.Context(), orgID, objectID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func (h *CacheObjectHandler) Pin(w http.ResponseWriter, r *http.Request) {
	h.setPinned(w, r, true)
}

func (h *CacheObjectHandler) Unpin(w http.ResponseWriter, r *http.Request) {
	h.setPinned(w, r, false)
}

func (h *CacheObjectHandler) setPinned(w http.ResponseWriter, r *http.Request, pinned bool) {
	orgID := db.OrgIDFromContext(r.Context())
	objectID := chi.URLParam(r, "object_id")
	obj, err := h.DB.SetPinned(r.Context(), orgID, objectID, pinned)
	if err != nil {
		http.Error(w, `{"error":"not found or database error"}`, http.StatusNotFound)
		return
	}
	action := "unpin"
	if pinned {
		action = "pin"
	}
	_ = h.DB.WriteAuditLog(r.Context(), orgID, "", action, "cache_object", objectID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}
