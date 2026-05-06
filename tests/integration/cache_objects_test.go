package integration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getPriorityClassID fetches the first priority class ID for the demo org
func getPriorityClassID(t *testing.T) string {
	t.Helper()
	pool := newPool(t)
	var id string
	err := pool.QueryRow(t.Context(),
		`SELECT id FROM priority_classes WHERE org_id = $1 ORDER BY created_at LIMIT 1`,
		demoOrgID,
	).Scan(&id)
	require.NoError(t, err, "need at least one priority class — run supabase seed")
	return id
}

func TestListCacheObjects(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/cache/objects")
	require.Equal(t, 200, resp.StatusCode)
	var body struct{ Objects []map[string]any `json:"objects"` }
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Objects)
}

func TestCreateAndGetCacheObject(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	pcID := getPriorityClassID(t)
	name := fmt.Sprintf("test-object-%d", time.Now().UnixNano())

	resp := apiPost(t, "/v1/cache/objects", map[string]any{
		"site_id":           siteID,
		"priority_class_id": pcID,
		"name":              name,
		"source_url":        "https://example.com/file.pdf",
		"size_bytes":        1048576,
		"tags":              []string{"medical", "p0"},
	})
	require.Equal(t, 201, resp.StatusCode)

	var obj map[string]any
	decodeJSON(t, resp, &obj)
	assert.Equal(t, name, obj["name"])
	assert.Equal(t, false, obj["pinned"])
	assert.Equal(t, "ACTIVE", obj["status"])
	objID := obj["id"].(string)

	// Get by ID
	getResp := apiGet(t, "/v1/cache/objects/"+objID)
	require.Equal(t, 200, getResp.StatusCode)
	var got map[string]any
	decodeJSON(t, getResp, &got)
	assert.Equal(t, objID, got["id"])
}

func TestPinAndUnpinCacheObject(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	pcID := getPriorityClassID(t)

	// Create
	resp := apiPost(t, "/v1/cache/objects", map[string]any{
		"site_id": siteID, "priority_class_id": pcID,
		"name":       fmt.Sprintf("pin-test-%d", time.Now().UnixNano()),
		"size_bytes": 512,
	})
	require.Equal(t, 201, resp.StatusCode)
	var obj map[string]any
	decodeJSON(t, resp, &obj)
	objID := obj["id"].(string)

	// Pin
	pinResp := apiPost(t, "/v1/cache/objects/"+objID+"/pin", nil)
	require.Equal(t, 200, pinResp.StatusCode)
	var pinned map[string]any
	decodeJSON(t, pinResp, &pinned)
	assert.Equal(t, true, pinned["pinned"])

	// Verify in DB
	pool := newPool(t)
	var isPinned bool
	pool.QueryRow(t.Context(), `SELECT pinned FROM cache_objects WHERE id = $1`, objID).Scan(&isPinned)
	assert.True(t, isPinned)

	// Unpin
	unpinResp := apiPost(t, "/v1/cache/objects/"+objID+"/unpin", nil)
	require.Equal(t, 200, unpinResp.StatusCode)
	var unpinned map[string]any
	decodeJSON(t, unpinResp, &unpinned)
	assert.Equal(t, false, unpinned["pinned"])
}

func TestCacheObjectNotFound(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/cache/objects/00000000-0000-0000-0000-000000000000")
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCreateCacheObjectMissingFields(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiPost(t, "/v1/cache/objects", map[string]any{"name": "no-site"})
	assert.Equal(t, 400, resp.StatusCode)
}

func TestListCacheObjectsFilterBySite(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	resp := apiGet(t, "/v1/cache/objects?site_id="+siteID)
	require.Equal(t, 200, resp.StatusCode)
	var body struct{ Objects []map[string]any `json:"objects"` }
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Objects)
}

func TestPinAuditLogCreated(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	pcID := getPriorityClassID(t)

	resp := apiPost(t, "/v1/cache/objects", map[string]any{
		"site_id": siteID, "priority_class_id": pcID,
		"name":       fmt.Sprintf("audit-test-%d", time.Now().UnixNano()),
		"size_bytes": 1024,
	})
	require.Equal(t, 201, resp.StatusCode)
	var obj map[string]any
	decodeJSON(t, resp, &obj)
	objID := obj["id"].(string)

	apiPost(t, "/v1/cache/objects/"+objID+"/pin", nil)

	pool := newPool(t)
	var count int
	pool.QueryRow(t.Context(),
		`SELECT COUNT(*) FROM audit_logs WHERE resource_id = $1 AND action = 'pin'`,
		objID,
	).Scan(&count)
	assert.Equal(t, 1, count, "pin action must create an audit log entry")
}
