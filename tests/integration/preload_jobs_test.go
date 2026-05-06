package integration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetPreloadJob(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	name := fmt.Sprintf("job-%d", time.Now().UnixNano())

	resp := apiPost(t, "/v1/preload/jobs", map[string]any{
		"site_id": siteID, "name": name,
	})
	require.Equal(t, 201, resp.StatusCode)
	var job map[string]any
	decodeJSON(t, resp, &job)
	assert.Equal(t, name, job["name"])
	assert.Equal(t, "PENDING", job["status"])
	jobID := job["id"].(string)

	// Get by ID
	getResp := apiGet(t, "/v1/preload/jobs/"+jobID)
	require.Equal(t, 200, getResp.StatusCode)
	var got map[string]any
	decodeJSON(t, getResp, &got)
	assert.Equal(t, jobID, got["id"])
}

func TestListPreloadJobs(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/preload/jobs")
	require.Equal(t, 200, resp.StatusCode)
	var body struct{ Jobs []map[string]any `json:"jobs"` }
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Jobs)
}

func TestCancelPreloadJob(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)

	// Create a job
	resp := apiPost(t, "/v1/preload/jobs", map[string]any{
		"site_id": siteID, "name": fmt.Sprintf("cancel-job-%d", time.Now().UnixNano()),
	})
	require.Equal(t, 201, resp.StatusCode)
	var job map[string]any
	decodeJSON(t, resp, &job)
	jobID := job["id"].(string)

	// Cancel it
	cancelResp := apiPost(t, "/v1/preload/jobs/"+jobID+"/cancel", nil)
	require.Equal(t, 200, cancelResp.StatusCode)
	var cancelBody map[string]any
	decodeJSON(t, cancelResp, &cancelBody)
	assert.Equal(t, "cancelled", cancelBody["status"])

	// Verify DB status
	pool := newPool(t)
	var status string
	pool.QueryRow(t.Context(), `SELECT status::text FROM preload_jobs WHERE id = $1`, jobID).Scan(&status)
	assert.Equal(t, "CANCELLED", status)
}

func TestPreloadJobNotFound(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/preload/jobs/00000000-0000-0000-0000-000000000000")
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCreatePreloadJobMissingFields(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiPost(t, "/v1/preload/jobs", map[string]any{"name": "no-site"})
	assert.Equal(t, 400, resp.StatusCode)
}

func TestPreloadJobAuditLog(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)

	resp := apiPost(t, "/v1/preload/jobs", map[string]any{
		"site_id": siteID, "name": fmt.Sprintf("audit-job-%d", time.Now().UnixNano()),
	})
	require.Equal(t, 201, resp.StatusCode)
	var job map[string]any
	decodeJSON(t, resp, &job)
	jobID := job["id"].(string)

	pool := newPool(t)
	var count int
	pool.QueryRow(t.Context(),
		`SELECT COUNT(*) FROM audit_logs WHERE resource_id = $1 AND action = 'create'`,
		jobID,
	).Scan(&count)
	assert.Equal(t, 1, count, "create preload job must produce an audit log")
}
