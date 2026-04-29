package report

import (
	"time"

	"github.com/driftwatch/internal/drift"
)

// Status represents the overall drift status of a report.
type Status string

const (
	StatusClean   Status = "clean"
	StatusDrifted Status = "drifted"
)

// Report holds a complete drift analysis snapshot.
type Report struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Status      Status         `json:"status"`
	Summary     drift.Summary  `json:"summary"`
	Results     []drift.Result `json:"results"`
}

// New builds a Report from a slice of drift results.
func New(results []drift.Result) Report {
	summary := drift.Summarize(results)
	status := StatusClean
	if summary.Drifted > 0 || summary.Missing > 0 {
		status = StatusDrifted
	}
	return Report{
		GeneratedAt: time.Now().UTC(),
		Status:      status,
		Summary:     summary,
		Results:     results,
	}
}

// HasDrift returns true when the report contains any drift or missing services.
func (r Report) HasDrift() bool {
	return r.Status == StatusDrifted
}
