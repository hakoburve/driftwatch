package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client fetches live configuration from a remote service endpoint.
type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

// Option configures a Client.
type Option func(*Client)

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
		c.httpClient.Timeout = d
	}
}

// NewClient creates a new fetch Client with sensible defaults.
func NewClient(opts ...Option) *Client {
	c := &Client{
		timeout: 10 * time.Second,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// FetchConfig retrieves the live key/value config map from the given URL.
// The endpoint is expected to return a flat JSON object.
func (c *Client) FetchConfig(ctx context.Context, url string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch: building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: executing request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch: unexpected status %d from %s", resp.StatusCode, url)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("fetch: decoding response from %s: %w", url, err)
	}
	return result, nil
}
