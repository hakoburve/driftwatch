package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/cache"
)

func TestStoreAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := map[string]string{"LOG_LEVEL": "info", "TIMEOUT": "30s"}

	if err := cache.Store(dir, "my-service", cfg); err != nil {
		t.Fatalf("Store: unexpected error: %v", err)
	}

	entry, err := cache.Load(dir, "my-service")
	if err != nil {
		t.Fatalf("Load: unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("Load: expected entry, got nil")
	}
	if entry.Config["LOG_LEVEL"] != "info" {
		t.Errorf("Config mismatch: got %q, want %q", entry.Config["LOG_LEVEL"], "info")
	}
	if entry.Config["TIMEOUT"] != "30s" {
		t.Errorf("Config mismatch: got %q, want %q", entry.Config["TIMEOUT"], "30s")
	}
}

func TestLoad_MissingFile_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	entry, err := cache.Load(dir, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for missing cache, got %+v", entry)
	}
}

func TestIsExpired_Fresh(t *testing.T) {
	entry := &cache.Entry{
		FetchedAt: time.Now().UTC(),
		Config:    map[string]string{},
	}
	if entry.IsExpired(5 * time.Minute) {
		t.Error("expected fresh entry to not be expired")
	}
}

func TestIsExpired_Stale(t *testing.T) {
	entry := &cache.Entry{
		FetchedAt: time.Now().UTC().Add(-10 * time.Minute),
		Config:    map[string]string{},
	}
	if !entry.IsExpired(5 * time.Minute) {
		t.Error("expected stale entry to be expired")
	}
}

func TestStore_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/cache"

	if err := cache.Store(dir, "svc", map[string]string{"K": "V"}); err != nil {
		t.Fatalf("Store: unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected cache directory to be created")
	}
}

func TestStore_SanitizesServiceName(t *testing.T) {
	dir := t.TempDir()
	if err := cache.Store(dir, "ns/my:service", map[string]string{"A": "1"}); err != nil {
		t.Fatalf("Store: unexpected error: %v", err)
	}
	entry, err := cache.Load(dir, "ns/my:service")
	if err != nil {
		t.Fatalf("Load: unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry after sanitized store/load")
	}
}
