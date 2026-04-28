package drift

import (
	"testing"
)

func TestDetect_InSync(t *testing.T) {
	states := []ServiceState{
		{
			Name:     "api",
			Declared: map[string]string{"IMAGE": "nginx:1.25", "REPLICAS": "3"},
			Live:     map[string]string{"IMAGE": "nginx:1.25", "REPLICAS": "3"},
		},
	}

	results := Detect(states)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusInSync {
		t.Errorf("expected in-sync, got %s", results[0].Status)
	}
	if len(results[0].Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(results[0].Changes))
	}
}

func TestDetect_Drifted(t *testing.T) {
	states := []ServiceState{
		{
			Name:     "worker",
			Declared: map[string]string{"IMAGE": "app:v2", "REPLICAS": "5"},
			Live:     map[string]string{"IMAGE": "app:v1", "REPLICAS": "5"},
		},
	}

	results := Detect(states)
	if results[0].Status != StatusDrifted {
		t.Errorf("expected drifted, got %s", results[0].Status)
	}
	if len(results[0].Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(results[0].Changes))
	}
	c := results[0].Changes[0]
	if c.Key != "IMAGE" || c.Declared != "app:v2" || c.Live != "app:v1" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDetect_Missing(t *testing.T) {
	states := []ServiceState{
		{
			Name:     "scheduler",
			Declared: map[string]string{"IMAGE": "sched:latest"},
			Live:     nil,
		},
	}

	results := Detect(states)
	if results[0].Status != StatusMissing {
		t.Errorf("expected missing, got %s", results[0].Status)
	}
}

func TestDetect_ExtraLiveKey(t *testing.T) {
	states := []ServiceState{
		{
			Name:     "proxy",
			Declared: map[string]string{"IMAGE": "proxy:1.0"},
			Live:     map[string]string{"IMAGE": "proxy:1.0", "DEBUG": "true"},
		},
	}

	results := Detect(states)
	if results[0].Status != StatusDrifted {
		t.Errorf("expected drifted due to extra live key, got %s", results[0].Status)
	}
}

func TestSummary(t *testing.T) {
	cases := []struct {
		result   DriftResult
		contains string
	}{
		{DriftResult{Service: "api", Status: StatusInSync}, "in-sync"},
		{DriftResult{Service: "api", Status: StatusMissing}, "not found"},
		{
			DriftResult{
				Service: "api",
				Status:  StatusDrifted,
				Changes: []Change{{Key: "IMAGE", Declared: "a", Live: "b"}},
			},
			"1 key(s)",
		},
	}

	for _, tc := range cases {
		s := Summary(tc.result)
		if !contains(s, tc.contains) {
			t.Errorf("Summary(%v) = %q, want substring %q", tc.result.Status, s, tc.contains)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
