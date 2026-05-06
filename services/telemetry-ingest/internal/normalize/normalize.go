package normalize

import "encoding/json"

// Normalize parses a raw JSON telemetry event into a generic map.
// Phase 2 will apply field-level normalization, unit conversion, and schema
// validation before this data reaches storage.
func Normalize(raw json.RawMessage) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}
