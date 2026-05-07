package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/off-planet-cdn/scheduler/internal/contactwindows"
	"github.com/off-planet-cdn/scheduler/internal/db"
	"github.com/off-planet-cdn/scheduler/internal/optimizer"
	"github.com/off-planet-cdn/scheduler/internal/queues"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	redisURL := envOr("REDIS_URL", "redis://localhost:6379")
	dbURL := envOr("SUPABASE_DB_URL", "postgresql://postgres:postgres@localhost:54322/postgres")
	edgeAgentURL := envOr("EDGE_AGENT_URL", "http://localhost:8081")

	redisAddr := redisURL
	if after, found := strings.CutPrefix(redisAddr, "redis://"); found {
		redisAddr = after
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbClient, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("scheduler: connect to database: %v", err)
	}
	defer dbClient.Close()

	redisClient := queues.New(redisAddr)
	defer redisClient.Close()

	if err := redisClient.Ping(ctx); err != nil {
		log.Printf("scheduler: warn: redis not reachable: %v", err)
	}

	checker := contactwindows.New(dbClient)
	httpClient := &http.Client{Timeout: 30 * time.Second}

	log.Println("scheduler: started")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler: shutting down")
			return
		case <-ticker.C:
			if err := dispatchPendingJobs(ctx, dbClient, redisClient, checker, httpClient, edgeAgentURL); err != nil {
				log.Printf("scheduler: dispatch error: %v", err)
			}
		}
	}
}

func dispatchPendingJobs(
	ctx context.Context,
	dbClient *db.Client,
	redisClient *queues.Client,
	checker *contactwindows.Checker,
	httpClient *http.Client,
	edgeAgentURL string,
) error {
	jobs, err := dbClient.ListPendingJobs(ctx)
	if err != nil {
		return fmt.Errorf("list pending jobs: %w", err)
	}

	qLen, _ := redisClient.QueueLength(ctx)
	log.Printf("scheduler: tick — %d pending jobs in DB, %d in Redis queue", len(jobs), qLen)

	for _, job := range jobs {
		if err := dispatchJob(ctx, dbClient, redisClient, checker, httpClient, edgeAgentURL, job); err != nil {
			log.Printf("scheduler: job %s dispatch error: %v", job.ID, err)
			_ = dbClient.MarkJobFailed(ctx, job.ID)
		}
	}
	return nil
}

func dispatchJob(
	ctx context.Context,
	dbClient *db.Client,
	redisClient *queues.Client,
	checker *contactwindows.Checker,
	httpClient *http.Client,
	edgeAgentURL string,
	job db.PendingJob,
) error {
	// Skip if already cancelled
	cancelled, err := dbClient.IsJobCancelled(ctx, job.ID)
	if err != nil {
		return fmt.Errorf("check cancelled: %w", err)
	}
	if cancelled {
		log.Printf("scheduler: job %s is cancelled — skipping", job.ID)
		return nil
	}

	// Check Redis cancellation signal
	if redisCancelled, _ := redisClient.IsCancelled(ctx, job.ID); redisCancelled {
		log.Printf("scheduler: job %s cancelled via Redis — skipping", job.ID)
		return nil
	}

	// Check contact window
	open, err := checker.IsWindowOpen(ctx, job.SiteID)
	if err != nil {
		return fmt.Errorf("check window: %w", err)
	}
	if !open {
		log.Printf("scheduler: job %s — contact window closed for site %s, deferring", job.ID, job.SiteID)
		return nil
	}

	// Get job items
	items, err := dbClient.GetJobItems(ctx, job.ID)
	if err != nil {
		return fmt.Errorf("get job items: %w", err)
	}
	if len(items) == 0 {
		log.Printf("scheduler: job %s has no pending items — marking done", job.ID)
		return dbClient.MarkJobDone(ctx, job.ID)
	}

	// Sort by priority, trim to bandwidth budget
	items = optimizer.Optimize(items)
	items = optimizer.FitInBudget(items, job.BandwidthBudgetBytes)

	// Mark job running
	if err := dbClient.MarkJobRunning(ctx, job.ID); err != nil {
		return fmt.Errorf("mark running: %w", err)
	}

	// Dispatch to edge agent
	if err := dispatchToEdge(ctx, httpClient, edgeAgentURL, items); err != nil {
		return fmt.Errorf("dispatch to edge: %w", err)
	}

	log.Printf("scheduler: job %s dispatched %d items to edge agent", job.ID, len(items))
	return dbClient.MarkJobDone(ctx, job.ID)
}

type preloadObject struct {
	ObjectID  string `json:"object_id"`
	SourceURL string `json:"source_url"`
	Priority  string `json:"priority"`
}

type preloadRequest struct {
	Objects []preloadObject `json:"objects"`
}

func dispatchToEdge(ctx context.Context, client *http.Client, edgeURL string, items []db.JobItem) error {
	objects := make([]preloadObject, 0, len(items))
	for _, item := range items {
		if item.SourceURL == "" {
			continue // skip objects without a source URL
		}
		objects = append(objects, preloadObject{
			ObjectID:  item.ObjectID,
			SourceURL: item.SourceURL,
			Priority:  item.Priority,
		})
	}
	if len(objects) == 0 {
		return nil
	}

	body, err := json.Marshal(preloadRequest{Objects: objects})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		edgeURL+"/local/cache/preload", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("edge agent unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("edge agent returned %d", resp.StatusCode)
	}
	return nil
}
