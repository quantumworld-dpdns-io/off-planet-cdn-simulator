package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func apiPut(t *testing.T, path string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, controlAPIBase(t)+path, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Org-ID", demoOrgID)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func TestListPolicies(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/policies")
	require.Equal(t, 200, resp.StatusCode)
	var body struct{ Policies []map[string]any `json:"policies"` }
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Policies)
}

func TestCreateAndUpdatePolicy(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	siteID := getOrCreateSiteID(t)
	name := fmt.Sprintf("policy-%d", time.Now().UnixNano())

	// Create
	resp := apiPost(t, "/v1/policies", map[string]any{
		"site_id": siteID, "name": name, "description": "Test policy", "enabled": true,
	})
	require.Equal(t, 201, resp.StatusCode)
	var policy map[string]any
	decodeJSON(t, resp, &policy)
	assert.Equal(t, name, policy["name"])
	assert.Equal(t, true, policy["enabled"])
	policyID := policy["id"].(string)

	// Update
	newName := name + "-updated"
	updateResp := apiPut(t, "/v1/policies/"+policyID, map[string]any{
		"name": newName, "description": "Updated", "enabled": false,
	})
	require.Equal(t, 200, updateResp.StatusCode)
	var updated map[string]any
	decodeJSON(t, updateResp, &updated)
	assert.Equal(t, newName, updated["name"])
	assert.Equal(t, false, updated["enabled"])
}

func TestCreatePolicyMissingFields(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiPost(t, "/v1/policies", map[string]any{"name": "no-site"})
	assert.Equal(t, 400, resp.StatusCode)
}
