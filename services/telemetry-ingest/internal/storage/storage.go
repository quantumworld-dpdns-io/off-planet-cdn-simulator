package storage

import (
	"context"
	"log"
)

// Store persists a batch of normalized telemetry events.
// This stub logs the count; Phase 2 will write to TimescaleDB/Supabase.
func Store(_ context.Context, events []map[string]any) error {
	log.Printf("storage: storing %d events", len(events))
	return nil
}
