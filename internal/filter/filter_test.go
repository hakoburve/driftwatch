package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/filter"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "alpha", Diffs: []drift.Diff{{Key: "port", Expected: "8080", Actual: "9090"}}},
		{Service: "beta", Diffs: nil},
		{Service: "gamma", Missing: true},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{})
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_OnlyDrifted_ExcludesClean(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{OnlyDrifted: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(got))
	}
	for _, r := range got {
		if r.Service == "beta" {
			t.Errorf("clean service 'beta' should have been filtered out")
		}
	}
}

func TestApply_ServiceAllowlist_FiltersCorrectly(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Services: []string{"alpha", "gamma"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Service == "beta" {
			t.Errorf("service 'beta' should have been excluded by allowlist")
		}
	}
}

func TestApply_CombinedOptions(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{
		OnlyDrifted: true,
		Services:    []string{"alpha", "beta"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Service != "alpha" {
		t.Errorf("expected 'alpha', got %q", got[0].Service)
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	got := filter.Apply(nil, filter.Options{OnlyDrifted: true})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestApply_UnknownServiceInAllowlist_ReturnsEmpty(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{Services: []string{"nonexistent"}})
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}
