package filter

import "github.com/yourorg/driftwatch/internal/drift"

// Options holds criteria for filtering drift results.
type Options struct {
	// OnlyDrifted returns only results that have drift or are missing.
	OnlyDrifted bool
	// Services is an allowlist of service names; empty means all services.
	Services []string
}

// Apply filters a slice of drift.Result according to the given Options.
// It returns a new slice containing only the results that match.
func Apply(results []drift.Result, opts Options) []drift.Result {
	allowlist := buildAllowlist(opts.Services)

	filtered := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if len(allowlist) > 0 {
			if _, ok := allowlist[r.Service]; !ok {
				continue
			}
		}
		if opts.OnlyDrifted && !isDrifted(r) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// isDrifted reports whether a result represents any form of drift.
func isDrifted(r drift.Result) bool {
	return r.Missing || len(r.Diffs) > 0
}

// buildAllowlist converts a slice of service names into a set for O(1) lookup.
func buildAllowlist(services []string) map[string]struct{} {
	if len(services) == 0 {
		return nil
	}
	m := make(map[string]struct{}, len(services))
	for _, s := range services {
		m[s] = struct{}{}
	}
	return m
}
