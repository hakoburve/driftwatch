package fetch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/driftwatch/internal/fetch"
)

func TestMockClient_ReturnsConfiguredResponse(t *testing.T) {
	m := fetch.NewMockClient()
	url := "http://svc-a/config"
	m.Responses[url] = map[string]string{"KEY": "value"}

	got, err := m.FetchConfig(context.Background(), url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "value" {
		t.Errorf("want %q, got %q", "value", got["KEY"])
	}
	if len(m.Calls) != 1 || m.Calls[0] != url {
		t.Errorf("expected call recorded for %q", url)
	}
}

func TestMockClient_ReturnsConfiguredError(t *testing.T) {
	m := fetch.NewMockClient()
	url := "http://svc-b/config"
	expectedErr := errors.New("connection refused")
	m.Errors[url] = expectedErr

	_, err := m.FetchConfig(context.Background(), url)
	if !errors.Is(err, expectedErr) {
		t.Errorf("want %v, got %v", expectedErr, err)
	}
}

func TestMockClient_ReturnsEmptyMapForUnknownURL(t *testing.T) {
	m := fetch.NewMockClient()
	got, err := m.FetchConfig(context.Background(), "http://unknown/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMockClient_RecordsMultipleCalls(t *testing.T) {
	m := fetch.NewMockClient()
	urls := []string{"http://a/config", "http://b/config", "http://a/config"}
	for _, u := range urls {
		m.FetchConfig(context.Background(), u)
	}
	if len(m.Calls) != len(urls) {
		t.Errorf("want %d calls recorded, got %d", len(urls), len(m.Calls))
	}
}

// TestMockClient_ErrorTakesPrecedenceOverResponse verifies that when both an
// error and a response are configured for the same URL, the error is returned.
func TestMockClient_ErrorTakesPrecedenceOverResponse(t *testing.T) {
	m := fetch.NewMockClient()
	url := "http://svc-c/config"
	expectedErr := errors.New("timeout")
	m.Responses[url] = map[string]string{"KEY": "value"}
	m.Errors[url] = expectedErr

	_, err := m.FetchConfig(context.Background(), url)
	if !errors.Is(err, expectedErr) {
		t.Errorf("want error %v, got %v", expectedErr, err)
	}
}

// Compile-time check that both types satisfy the interface.
var _ fetch.ConfigFetcher = (*fetch.Client)(nil)
var _ fetch.ConfigFetcher = (*fetch.MockClient)(nil)
