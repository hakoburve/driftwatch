package report_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func driftedResults() []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Diffs: []drift.Diff{
				{Key: "LOG_LEVEL", Declared: "info", Live: "debug"},
			},
		},
		{
			Service:  "worker",
			Missing:  true,
		},
	}
}

func TestNew_CleanReport(t *testing.T) {
	r := report.New([]drift.Result{})
	if r.Status != report.StatusClean {
		t.Errorf("expected clean, got %s", r.Status)
	}
	if r.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
	if r.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestNew_DriftedReport(t *testing.T) {
	r := report.New(driftedResults())
	if r.Status != report.StatusDrifted {
		t.Errorf("expected drifted, got %s", r.Status)
	}
	if !r.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestNew_SummaryPopulated(t *testing.T) {
	r := report.New(driftedResults())
	if r.Summary.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Summary.Total)
	}
	if r.Summary.Drifted != 1 {
		t.Errorf("expected drifted 1, got %d", r.Summary.Drifted)
	}
	if r.Summary.Missing != 1 {
		t.Errorf("expected missing 1, got %d", r.Summary.Missing)
	}
}

func TestNew_ResultsPreserved(t *testing.T) {
	results := driftedResults()
	r := report.New(results)
	if len(r.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(r.Results))
	}
}
