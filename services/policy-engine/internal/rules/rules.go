package rules

// Evaluate applies the policy identified by policyID to objectID and returns
// the recommended action. Returns "PREFETCH" unconditionally in this stub;
// Phase 2 will load and execute policy rules from the database.
func Evaluate(policyID string, objectID string) (string, error) {
	return "PREFETCH", nil
}
