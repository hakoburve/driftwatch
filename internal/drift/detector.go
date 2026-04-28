package drift

import (
	"fmt"
	"sort"
)

// DriftStatus represents the drift state of a single service.
type DriftStatus int

const (
	StatusInSync DriftStatus = iota
	StatusDrifted
	StatusMissing
)

func (s DriftStatus) String() string {
	switch s {
	case StatusInSync:
		return "in-sync"
	case StatusDrifted:
		return "drifted"
	case StatusMissing:
		return "missing"
	default:
		return "unknown"
	}
}

// ServiceState holds the declared and live configuration for a service.
type ServiceState struct {
	Name     string
	Declared map[string]string
	Live     map[string]string
}

// DriftResult describes the outcome of comparing a service's declared vs live state.
type DriftResult struct {
	Service  string
	Status   DriftStatus
	Changes  []Change
}

// Change captures a single key-level difference.
type Change struct {
	Key      string
	Declared string
	Live     string
}

// Detect compares declared configuration against live configuration for each
// ServiceState and returns a DriftResult per service.
func Detect(states []ServiceState) []DriftResult {
	results := make([]DriftResult, 0, len(states))

	for _, svc := range states {
		result := DriftResult{Service: svc.Name}

		if svc.Live == nil {
			result.Status = StatusMissing
			results = append(results, result)
			continue
		}

		keys := unionKeys(svc.Declared, svc.Live)
		for _, key := range keys {
			declaredVal := svc.Declared[key]
			liveVal := svc.Live[key]
			if declaredVal != liveVal {
				result.Changes = append(result.Changes, Change{
					Key:      key,
					Declared: declaredVal,
					Live:     liveVal,
				})
			}
		}

		if len(result.Changes) > 0 {
			result.Status = StatusDrifted
		} else {
			result.Status = StatusInSync
		}

		results = append(results, result)
	}

	return results
}

// Summary returns a human-readable summary line for a DriftResult.
func Summary(r DriftResult) string {
	switch r.Status {
	case StatusMissing:
		return fmt.Sprintf("[%s] %s — service not found in live environment", r.Status, r.Service)
	case StatusDrifted:
		return fmt.Sprintf("[%s] %s — %d key(s) differ", r.Status, r.Service, len(r.Changes))
	default:
		return fmt.Sprintf("[%s] %s", r.Status, r.Service)
	}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
