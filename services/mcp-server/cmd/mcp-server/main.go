package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/off-planet-cdn/mcp-server/internal/tools"
)

// ToolCallRequest is the shape of a POST /tools/call body.
type ToolCallRequest struct {
	Tool  string          `json:"tool"`
	Input json.RawMessage `json:"input"`
}

func handleToolCall(w http.ResponseWriter, r *http.Request) {
	var req ToolCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch req.Tool {
	case "cache_status":
		var input tools.CacheStatusInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
			return
		}
		out, err := tools.CacheStatus(r.Context(), input)
		if err != nil {
			http.Error(w, `{"error":"tool error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(out)

	case "generate_preload_plan":
		out, err := tools.GeneratePreloadPlan(r.Context(), req.Input)
		if err != nil {
			http.Error(w, `{"error":"tool error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(out)

	case "inspect_node":
		out, err := tools.InspectNode(r.Context(), req.Input)
		if err != nil {
			http.Error(w, `{"error":"tool error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(out)

	case "simulate_eviction":
		out, err := tools.SimulateEviction(r.Context(), req.Input)
		if err != nil {
			http.Error(w, `{"error":"tool error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(out)

	case "summarize_incident":
		out, err := tools.SummarizeIncident(r.Context(), req.Input)
		if err != nil {
			http.Error(w, `{"error":"tool error"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"summary": out})

	default:
		http.Error(w, `{"error":"unknown tool"}`, http.StatusBadRequest)
	}
}

func main() {
	port := os.Getenv("MCP_SERVER_PORT")
	if port == "" {
		port = "8084"
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "mcp-server",
		})
	})

	r.Post("/tools/call", handleToolCall)

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("mcp-server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down mcp-server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("mcp-server stopped")
}
