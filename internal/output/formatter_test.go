package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func TestWriteText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	results := []drift.Result{
		{
			Service: "api",
			Missing: false,
			Diffs: []drift.Diff{
				{Key: "LOG_LEVEL", Expected: "info", Actual: "debug"},
			},
			ExtraKeys: []string{"UNKNOWN_VAR"},
		},
	}
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFTED] api") {
		t.Errorf("expected DRIFTED label, got: %q", out)
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected diff key LOG_LEVEL, got: %q", out)
	}
	if !strings.Contains(out, "UNKNOWN_VAR") {
		t.Errorf("expected extra key UNKNOWN_VAR, got: %q", out)
	}
}

func TestWriteText_MissingService(t *testing.T) {
	results := []drift.Result{
		{Service: "worker", Missing: true},
	}
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[MISSING] worker") {
		t.Errorf("expected MISSING label, got: %q", buf.String())
	}
}

func TestWriteJSON_NoResults(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"status":"ok"`) {
		t.Errorf("expected ok status in JSON, got: %q", buf.String())
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	results := []drift.Result{
		{Service: "cache", Missing: false, Diffs: []drift.Diff{{Key: "TTL", Expected: "300", Actual: "600"}}},
	}
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"status":"drifted"`) {
		t.Errorf("expected drifted status in JSON, got: %q", buf.String())
	}
}
