package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	ControlAPIURL = getEnv("CONTROL_API_URL", "http://localhost:8080")
	OrgID         = getEnv("ORG_ID", "00000000-0000-0000-0000-000000000001")
	HTTPClient    = &http.Client{Timeout: 15 * time.Second}
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// apiGet fetches JSON from the control-api and decodes it into dest.
func apiGet(ctx context.Context, path string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ControlAPIURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Org-ID", OrgID)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dest)
}
