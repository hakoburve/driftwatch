package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func singleDriftReport() report.Report {
	return report.New([]drift.Result{
		{
			Service: "api",
			Diffs: []drift.Diff{
				{Key: "TIMEOUT", Declared: "30s", Live: "60s"},
			},
		},
	})
}

func TestWriteText_ContainsStatus(t *testing.T) {
	var buf bytes.Buffer
	r := singleDriftReport()
	if err := report.WriteText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "drifted") {
		t.Errorf("expected 'drifted' in output, got:\n%s", buf.String())
	}
}

func TestWriteText_ContainsDiffKey(t *testing.T) {
	var buf bytes.Buffer
	r := singleDriftReport()
	report.WriteText(&buf, r)
	if !strings.Contains(buf.String(), "TIMEOUT") {
		t.Errorf("expected 'TIMEOUT' in output, got:\n%s", buf.String())
	}
}

func TestWriteText_MissingService(t *testing.T) {
	var buf bytes.Buffer
	r := report.New([]drift.Result{{Service: "ghost", Missing: true}})
	report.WriteText(&buf, r)
	if !strings.Contains(buf.String(), "[MISSING]") {
		t.Errorf("expected '[MISSING]' in output, got:\n%s", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := singleDriftReport()
	if err := report.WriteJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteJSON_ContainsStatus(t *testing.T) {
	var buf bytes.Buffer
	r := singleDriftReport()
	report.WriteJSON(&buf, r)
	if !strings.Contains(buf.String(), `"status"`) {
		t.Errorf("expected 'status' field in JSON output")
	}
}

func TestWriteJSON_CleanReport(t *testing.T) {
	var buf bytes.Buffer
	r := report.New([]drift.Result{})
	report.WriteJSON(&buf, r)
	if !strings.Contains(buf.String(), `"clean"`) {
		t.Errorf("expected 'clean' status in JSON output")
	}
}
