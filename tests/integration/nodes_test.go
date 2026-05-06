package integration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getOrCreateSiteID uses the Lunar Habitat Alpha site from seed.sql.
// Falls back to creating one via API if not present.
func getOrCreateSiteID(t *testing.T) string {
	t.Helper()
	pool := newPool(t)
	var siteID string
	err := pool.QueryRow(
		t.Context(),
		`SELECT id FROM sites WHERE org_id = $1 LIMIT 1`,
		demoOrgID,
	).Scan(&siteID)
	if err == nil {
		return siteID
	}
	// Create a site via API
	resp := apiPost(t, "/v1/sites", map[string]string{"name": "Test Site for Nodes"})
	require.Equal(t, 201, resp.StatusCode)
	var body map[string]any
	decodeJSON(t, resp, &body)
	return body["id"].(string)
}

func TestListNodes(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/nodes")
	require.Equal(t, 200, resp.StatusCode)

	var body struct {
		Nodes []map[string]any `json:"nodes"`
	}
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Nodes)
}

func TestRegisterNode(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	name := fmt.Sprintf("test-node-%d", time.Now().UnixNano())

	resp := apiPost(t, "/v1/nodes/register", map[string]any{
		"site_id":         siteID,
		"name":            name,
		"cache_dir":       "/tmp/test-cache",
		"cache_max_bytes": 1073741824,
	})
	require.Equal(t, 201, resp.StatusCode)

	var node map[string]any
	decodeJSON(t, resp, &node)
	assert.Equal(t, name, node["name"])
	assert.Equal(t, siteID, node["site_id"])
	assert.Equal(t, "UNKNOWN", node["status"])
	nodeID, ok := node["id"].(string)
	require.True(t, ok)

	// Verify node appears in DB
	pool := newPool(t)
	var dbName string
	err := pool.QueryRow(t.Context(), `SELECT name FROM nodes WHERE id = $1`, nodeID).Scan(&dbName)
	require.NoError(t, err)
	assert.Equal(t, name, dbName)
}

func TestHeartbeatUpdatesLastSeen(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)

	// Register a node
	regResp := apiPost(t, "/v1/nodes/register", map[string]any{
		"site_id": siteID,
		"name":    fmt.Sprintf("hb-node-%d", time.Now().UnixNano()),
	})
	require.Equal(t, 201, regResp.StatusCode)
	var regNode map[string]any
	decodeJSON(t, regResp, &regNode)
	nodeID := regNode["id"].(string)

	// Send heartbeat
	hbResp := apiPost(t, "/v1/nodes/"+nodeID+"/heartbeat", map[string]any{
		"status":           "ONLINE",
		"cache_used_bytes": 1073741824,
		"cache_max_bytes":  10737418240,
		"agent_version":    "0.1.0",
	})
	require.Equal(t, 200, hbResp.StatusCode)

	var hb map[string]any
	decodeJSON(t, hbResp, &hb)
	assert.Equal(t, nodeID, hb["node_id"])
	assert.Equal(t, "ONLINE", hb["status"])

	// Verify last_seen updated in DB and status changed
	pool := newPool(t)
	var lastSeen *time.Time
	var status string
	err := pool.QueryRow(t.Context(),
		`SELECT last_seen, status::text FROM nodes WHERE id = $1`, nodeID,
	).Scan(&lastSeen, &status)
	require.NoError(t, err)
	assert.NotNil(t, lastSeen, "last_seen must be set after heartbeat")
	assert.Equal(t, "ONLINE", status)

	// Send second heartbeat and verify last_seen advances
	time.Sleep(50 * time.Millisecond)
	firstSeen := *lastSeen
	apiPost(t, "/v1/nodes/"+nodeID+"/heartbeat", map[string]any{
		"status": "ONLINE", "cache_used_bytes": 2147483648, "cache_max_bytes": 10737418240,
	})
	var lastSeen2 *time.Time
	pool.QueryRow(t.Context(), `SELECT last_seen FROM nodes WHERE id = $1`, nodeID).Scan(&lastSeen2)
	assert.True(t, lastSeen2 != nil && !lastSeen2.Before(firstSeen), "last_seen must advance on second heartbeat")
}

func TestNodeStatusEndpoint(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)

	regResp := apiPost(t, "/v1/nodes/register", map[string]any{
		"site_id": siteID,
		"name":    fmt.Sprintf("status-node-%d", time.Now().UnixNano()),
	})
	require.Equal(t, 201, regResp.StatusCode)
	var regNode map[string]any
	decodeJSON(t, regResp, &regNode)
	nodeID := regNode["id"].(string)

	resp := apiGet(t, "/v1/nodes/"+nodeID+"/status")
	require.Equal(t, 200, resp.StatusCode)
	var node map[string]any
	decodeJSON(t, resp, &node)
	assert.Equal(t, nodeID, node["id"])
}

func TestHeartbeatAppearsInNodeHeartbeatsTable(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)

	regResp := apiPost(t, "/v1/nodes/register", map[string]any{
		"site_id": siteID,
		"name":    fmt.Sprintf("hbt-node-%d", time.Now().UnixNano()),
	})
	require.Equal(t, 201, regResp.StatusCode)
	var regNode map[string]any
	decodeJSON(t, regResp, &regNode)
	nodeID := regNode["id"].(string)

	apiPost(t, "/v1/nodes/"+nodeID+"/heartbeat", map[string]any{
		"status": "ONLINE", "cache_used_bytes": 0, "cache_max_bytes": 10737418240,
	})

	pool := newPool(t)
	var count int
	pool.QueryRow(t.Context(), `SELECT COUNT(*) FROM node_heartbeats WHERE node_id = $1`, nodeID).Scan(&count)
	assert.Equal(t, 1, count, "heartbeat row must exist in node_heartbeats")
}

func TestRegisterNodeMissingSiteID(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiPost(t, "/v1/nodes/register", map[string]any{"name": "bad-node"})
	assert.Equal(t, 400, resp.StatusCode)
}
