package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
version: "1"
services:
  - name: api
    type: kubernetes
    source: ./manifests/api.yaml
    target: default/deployment/api
    ignore_keys:
      - metadata.resourceVersion
`
	path := writeTemp(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Name != "api" {
		t.Errorf("expected service name 'api', got %q", cfg.Services[0].Name)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoServices(t *testing.T) {
	path := writeTemp(t, "version: \"1\"\nservices: []\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty services")
	}
}

func TestLoad_DuplicateServiceName(t *testing.T) {
	content := `
version: "1"
services:
  - name: api
    type: kubernetes
    source: ./manifests/api.yaml
    target: default/deployment/api
  - name: api
    type: docker
    source: ./docker-compose.yml
    target: api_container
`
	path := writeTemp(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate service name")
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	cases := []struct {
		name    string
		content string
	}{
		{"missing name", "version: \"1\"\nservices:\n  - type: kubernetes\n    source: s\n    target: t\n"},
		{"missing type", "version: \"1\"\nservices:\n  - name: svc\n    source: s\n    target: t\n"},
		{"missing source", "version: \"1\"\nservices:\n  - name: svc\n    type: kubernetes\n    target: t\n"},
		{"missing target", "version: \"1\"\nservices:\n  - name: svc\n    type: kubernetes\n    source: s\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTemp(t, tc.content)
			_, err := Load(path)
			if err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
		})
	}
}
