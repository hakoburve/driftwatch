package fetch_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/fetch"
)

func TestFetchConfig_Success(t *testing.T) {
	expected := map[string]string{"LOG_LEVEL": "info", "PORT": "8080"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	client := fetch.NewClient()
	got, err := client.FetchConfig(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("key %q: want %q, got %q", k, v, got[k])
		}
	}
}

func TestFetchConfig_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := fetch.NewClient()
	_, err := client.FetchConfig(context.Background(), ts.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestFetchConfig_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer ts.Close()

	client := fetch.NewClient()
	_, err := client.FetchConfig(context.Background(), ts.URL)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestFetchConfig_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := fetch.NewClient(fetch.WithTimeout(50 * time.Millisecond))
	_, err := client.FetchConfig(context.Background(), ts.URL)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestFetchConfig_BadURL(t *testing.T) {
	client := fetch.NewClient()
	_, err := client.FetchConfig(context.Background(), "://bad-url")
	if err == nil {
		t.Fatal("expected error for bad URL, got nil")
	}
}
