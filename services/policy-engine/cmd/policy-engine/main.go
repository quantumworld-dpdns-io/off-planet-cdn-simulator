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
	"github.com/off-planet-cdn/policy-engine/internal/scoring"
)

func main() {
	port := os.Getenv("POLICY_ENGINE_PORT")
	if port == "" {
		port = "8083"
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "policy-engine",
		})
	})

	r.Post("/score", func(w http.ResponseWriter, r *http.Request) {
		var ctx scoring.ScoringContext
		if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}
		score := scoring.Score(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"score": score})
	})

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("policy-engine listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down policy-engine...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("policy-engine stopped")
}
