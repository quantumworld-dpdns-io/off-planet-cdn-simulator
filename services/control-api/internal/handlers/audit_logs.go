package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

type AuditLogHandler struct{ DB *db.Client }

func (h *AuditLogHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	logs, err := h.DB.ListAuditLogs(r.Context(), orgID,
		q.Get("actor_id"), q.Get("action"), q.Get("resource_type"),
		limit, offset)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if logs == nil {
		logs = []models.AuditLog{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"logs": logs, "limit": limit, "offset": offset})
}
