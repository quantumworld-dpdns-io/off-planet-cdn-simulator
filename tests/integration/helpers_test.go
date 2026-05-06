package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const demoOrgID = "00000000-0000-0000-0000-000000000001"

func dbURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("SUPABASE_DB_URL")
	if url == "" {
		url = "postgresql://postgres:postgres@localhost:54322/postgres"
	}
	return url
}

func newPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dbURL(t))
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	return pool
}

func controlAPIBase(t *testing.T) string {
	t.Helper()
	base := os.Getenv("CONTROL_API_URL")
	if base == "" {
		base = "http://localhost:8080"
	}
	return base
}

func apiGet(t *testing.T, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, controlAPIBase(t)+path, nil)
	require.NoError(t, err)
	req.Header.Set("X-Org-ID", demoOrgID)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func apiPost(t *testing.T, path string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, controlAPIBase(t)+path, strings.NewReader(string(b)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Org-ID", demoOrgID)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(v))
}

func waitForAPI(t *testing.T, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(controlAPIBase(t) + "/v1/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("control API not reachable within timeout")
}
