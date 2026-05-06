package ingest

import (
	"context"
	"encoding/json"
	"log"

	"github.com/off-planet-cdn/telemetry-ingest/internal/normalize"
	"github.com/off-planet-cdn/telemetry-ingest/internal/storage"
)

// Ingest normalizes and stores a batch of raw telemetry events.
// Returns the number of events successfully ingested.
func Ingest(ctx context.Context, events []json.RawMessage) (int, error) {
	log.Printf("ingest: received %d events", len(events))

	normalized := make([]map[string]any, 0, len(events))
	for _, raw := range events {
		m, err := normalize.Normalize(raw)
		if err != nil {
			log.Printf("ingest: normalize error: %v", err)
			continue
		}
		normalized = append(normalized, m)
	}

	if err := storage.Store(ctx, normalized); err != nil {
		return 0, err
	}

	return len(normalized), nil
}
