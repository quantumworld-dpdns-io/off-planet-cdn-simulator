package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// GeneratePreloadPlanOutput is the result of a preload plan generation.
type GeneratePreloadPlanOutput struct {
	URLs []string `json:"urls"`
}

type generatePreloadPlanInput struct {
	SiteID string `json:"site_id"`
	Limit  int    `json:"limit"`
}

// GeneratePreloadPlan produces a prioritized list of URLs to prefetch for a site.
func GeneratePreloadPlan(ctx context.Context, input json.RawMessage) (*GeneratePreloadPlanOutput, error) {
	var in generatePreloadPlanInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}

	path := "/v1/objects?status=ACTIVE"
	if in.SiteID != "" {
		path += fmt.Sprintf("&site_id=%s", in.SiteID)
	}

	var objList objectsListResponse
	if err := apiGet(ctx, path, &objList); err != nil {
		return nil, err
	}

	urls := make([]string, 0, in.Limit)
	for _, obj := range objList.Objects {
		if len(urls) >= in.Limit {
			break
		}
		if obj.SourceURL == "" {
			continue
		}
		urls = append(urls, obj.SourceURL)
	}

	return &GeneratePreloadPlanOutput{URLs: urls}, nil
}
