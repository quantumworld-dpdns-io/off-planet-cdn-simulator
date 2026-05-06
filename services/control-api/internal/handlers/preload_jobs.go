package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
	cdnredis "github.com/off-planet-cdn/control-api/internal/redis"
)

type PreloadJobHandler struct {
	DB    *db.Client
	Redis *cdnredis.Client
}

func (h *PreloadJobHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	var body struct {
		SiteID               string   `json:"site_id"`
		Name                 string   `json:"name"`
		ObjectIDs            []string `json:"object_ids"`
		BandwidthBudgetBytes *int64   `json:"bandwidth_budget_bytes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.SiteID == "" {
		http.Error(w, `{"error":"name and site_id are required"}`, http.StatusBadRequest)
		return
	}
	job, err := h.DB.CreatePreloadJob(r.Context(), orgID, body.SiteID, body.Name, body.BandwidthBudgetBytes)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if len(body.ObjectIDs) > 0 {
		if err := h.DB.AddPreloadJobItems(r.Context(), orgID, job.ID, body.ObjectIDs); err != nil {
			http.Error(w, `{"error":"failed to add job items"}`, http.StatusInternalServerError)
			return
		}
	}
	if h.Redis != nil {
		_ = h.Redis.EnqueuePreloadJob(r.Context(), job.ID)
	}
	_ = h.DB.WriteAuditLog(r.Context(), orgID, "", "create", "preload_job", job.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func (h *PreloadJobHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	jobs, err := h.DB.ListPreloadJobs(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if jobs == nil {
		jobs = []models.PreloadJob{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"jobs": jobs})
}

func (h *PreloadJobHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	jobID := chi.URLParam(r, "job_id")
	job, err := h.DB.GetPreloadJob(r.Context(), orgID, jobID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *PreloadJobHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	jobID := chi.URLParam(r, "job_id")
	if err := h.DB.CancelPreloadJob(r.Context(), orgID, jobID); err != nil {
		http.Error(w, `{"error":"job not found or not cancellable"}`, http.StatusNotFound)
		return
	}
	if h.Redis != nil {
		_ = h.Redis.CancelPreloadJob(r.Context(), jobID)
	}
	_ = h.DB.WriteAuditLog(r.Context(), orgID, "", "cancel", "preload_job", jobID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled", "job_id": jobID})
}
