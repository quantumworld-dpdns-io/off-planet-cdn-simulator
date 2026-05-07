package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/control-api/internal/db"
	"github.com/off-planet-cdn/control-api/internal/models"
)

type NodeHandler struct{ DB *db.Client }

func (h *NodeHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	nodes, err := h.DB.ListNodes(r.Context(), orgID)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	if nodes == nil {
		nodes = []models.Node{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"nodes": nodes})
}

func (h *NodeHandler) Register(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	var body struct {
		SiteID        string `json:"site_id"`
		Name          string `json:"name"`
		CacheDir      string `json:"cache_dir"`
		CacheMaxBytes int64  `json:"cache_max_bytes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.SiteID == "" {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.CacheMaxBytes == 0 {
		body.CacheMaxBytes = 10737418240
	}
	if body.CacheDir == "" {
		body.CacheDir = "/tmp/offplanet-cache"
	}
	node, err := h.DB.RegisterNode(r.Context(), orgID, body.SiteID, body.Name, body.CacheDir, body.CacheMaxBytes)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

func (h *NodeHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	nodeID := chi.URLParam(r, "node_id")
	var body struct {
		Status         string `json:"status"`
		CacheUsedBytes int64  `json:"cache_used_bytes"`
		CacheMaxBytes  int64  `json:"cache_max_bytes"`
		AgentVersion   string `json:"agent_version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.Status == "" {
		body.Status = "ONLINE"
	}
	hb, err := h.DB.RecordHeartbeat(r.Context(), orgID, nodeID, body.Status, body.AgentVersion, body.CacheUsedBytes, body.CacheMaxBytes)
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hb)
}

func (h *NodeHandler) Status(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	nodeID := chi.URLParam(r, "node_id")
	node, err := h.DB.GetNodeStatus(r.Context(), orgID, nodeID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func (h *NodeHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := db.OrgIDFromContext(r.Context())
	nodeID := chi.URLParam(r, "node_id")
	node, err := h.DB.GetNodeStatus(r.Context(), orgID, nodeID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}
