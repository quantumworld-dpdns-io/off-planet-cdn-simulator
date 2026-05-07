package optimizer

import (
	"sort"

	"github.com/off-planet-cdn/scheduler/internal/db"
)

// priorityOrder maps priority level strings to sort weight (lower = dispatched first).
var priorityOrder = map[string]int{
	"P0": 0,
	"P1": 1,
	"P2": 2,
	"P3": 3,
	"P4": 4,
	"P5": 5,
}

// Optimize sorts job items by priority (P0 first, P5 last).
// Items with the same priority are ordered by size ascending (smaller first = faster wins).
func Optimize(items []db.JobItem) []db.JobItem {
	sorted := make([]db.JobItem, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		pi := priorityOrder[sorted[i].Priority]
		pj := priorityOrder[sorted[j].Priority]
		if pi != pj {
			return pi < pj
		}
		return sorted[i].SizeBytes < sorted[j].SizeBytes
	})
	return sorted
}

// FitInBudget trims items to fit within a bandwidth budget (bytes).
// Returns as many highest-priority items as can fit.
// If budget is nil, all items are returned.
func FitInBudget(items []db.JobItem, budgetBytes *int64) []db.JobItem {
	if budgetBytes == nil {
		return items
	}
	var result []db.JobItem
	var used int64
	for _, item := range items {
		if used+item.SizeBytes > *budgetBytes {
			continue
		}
		result = append(result, item)
		used += item.SizeBytes
	}
	return result
}
