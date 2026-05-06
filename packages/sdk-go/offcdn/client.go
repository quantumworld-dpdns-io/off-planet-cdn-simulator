package offcdn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func New(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, httpClient: &http.Client{}, token: token}
}

func (c *Client) Health(ctx context.Context) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/health", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed: %d", resp.StatusCode)
	}
	var result map[string]string
	return result, json.NewDecoder(resp.Body).Decode(&result)
}
