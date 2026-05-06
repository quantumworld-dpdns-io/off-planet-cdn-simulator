package scoring

type ScoringContext struct {
	PriorityWeight     float64 `json:"priority_weight"`
	MissionRelevance   float64 `json:"mission_relevance"`
	PredictedDemand    float64 `json:"predicted_demand"`
	OfflineCriticality float64 `json:"offline_criticality"`
	RevalidationCost   float64 `json:"revalidation_cost"`
	FetchLatencyCost   float64 `json:"fetch_latency_cost"`
	PackageDependency  float64 `json:"package_dependency"`
	SizePenalty        float64 `json:"size_penalty"`
	StalenessPenalty   float64 `json:"staleness_penalty"`
	RedundancyPenalty  float64 `json:"redundancy_penalty"`
}

func Score(ctx ScoringContext) float64 {
	return ctx.PriorityWeight +
		ctx.MissionRelevance +
		ctx.PredictedDemand +
		ctx.OfflineCriticality +
		ctx.RevalidationCost +
		ctx.FetchLatencyCost +
		ctx.PackageDependency -
		ctx.SizePenalty -
		ctx.StalenessPenalty -
		ctx.RedundancyPenalty
}
