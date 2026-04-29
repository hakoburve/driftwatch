package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a cached live config snapshot with metadata.
type Entry struct {
	FetchedAt time.Time         `json:"fetched_at"`
	Config    map[string]string `json:"config"`
}

// IsExpired returns true if the entry is older than ttl.
func (e *Entry) IsExpired(ttl time.Duration) bool {
	return time.Since(e.FetchedAt) > ttl
}

// Store persists a config snapshot for the given service to the cache directory.
func Store(dir, service string, cfg map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cache: create dir: %w", err)
	}
	entry := Entry{
		FetchedAt: time.Now().UTC(),
		Config:    cfg,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	path := filepath.Join(dir, sanitize(service)+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("cache: write file: %w", err)
	}
	return nil
}

// Load reads a cached entry for the given service. Returns (nil, nil) when no
// cache file exists yet.
func Load(dir, service string) (*Entry, error) {
	path := filepath.Join(dir, sanitize(service)+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read file: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return &entry, nil
}

// sanitize replaces characters that are unsafe in filenames.
func sanitize(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '/' || c == '\\' || c == ':' || c == '*' || c == '?' {
			out[i] = '_'
		} else {
			out[i] = c
		}
	}
	return string(out)
}
