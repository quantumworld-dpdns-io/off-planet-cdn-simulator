package integration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSites(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/sites")
	require.Equal(t, 200, resp.StatusCode)

	var body struct {
		Sites []map[string]any `json:"sites"`
	}
	decodeJSON(t, resp, &body)
	assert.NotNil(t, body.Sites, "sites must not be null")
}

func TestCreateAndGetSite(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	name := fmt.Sprintf("Test Site %d", time.Now().UnixNano())

	// Create
	resp := apiPost(t, "/v1/sites", map[string]string{
		"name":        name,
		"location":    "Test Location",
		"description": "Integration test site",
	})
	require.Equal(t, 201, resp.StatusCode)

	var created map[string]any
	decodeJSON(t, resp, &created)
	assert.Equal(t, name, created["name"])
	siteID, ok := created["id"].(string)
	require.True(t, ok, "id must be a string")

	// Get by ID
	getResp := apiGet(t, "/v1/sites/"+siteID)
	require.Equal(t, 200, getResp.StatusCode)

	var got map[string]any
	decodeJSON(t, getResp, &got)
	assert.Equal(t, siteID, got["id"])
	assert.Equal(t, name, got["name"])
}

func TestGetSiteNotFound(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiGet(t, "/v1/sites/00000000-0000-0000-0000-000000000000")
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCreateSiteMissingName(t *testing.T) {
	waitForAPI(t, 10*time.Second)
	resp := apiPost(t, "/v1/sites", map[string]string{"location": "Nowhere"})
	assert.Equal(t, 400, resp.StatusCode)
}
