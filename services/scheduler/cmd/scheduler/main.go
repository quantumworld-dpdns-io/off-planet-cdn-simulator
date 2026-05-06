package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	redisURL := os.Getenv("REDIS_URL")
	dbURL := os.Getenv("SUPABASE_DB_URL")
	controlAPIURL := getEnv("CONTROL_API_URL", "http://control-api:8080")

	log.Printf("scheduler starting: redis=%s db=%s controlAPI=%s", redisURL, dbURL, controlAPIURL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("scheduler running — tick interval 10s")

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler shutting down")
			return
		case t := <-ticker.C:
			log.Printf("scheduler tick at %s", t.Format(time.RFC3339))
		}
	}
}
