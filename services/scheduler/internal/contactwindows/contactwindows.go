package contactwindows

import "context"

// IsWindowOpen reports whether a communication window is currently open for
// the given site. Returns true unconditionally in this stub — Phase 2 will
// integrate real orbital-mechanics calculations.
func IsWindowOpen(_ context.Context, siteID string) (bool, error) {
	return true, nil
}
