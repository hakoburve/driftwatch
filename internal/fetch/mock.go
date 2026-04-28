package fetch

import "context"

// MockClient is a test double for Client that returns pre-configured responses.
type MockClient struct {
	// Responses maps a URL to the config map that should be returned.
	Responses map[string]map[string]string
	// Errors maps a URL to an error that should be returned.
	Errors map[string]error
	// Calls records every URL that FetchConfig was called with.
	Calls []string
}

// NewMockClient initialises a MockClient with empty maps.
func NewMockClient() *MockClient {
	return &MockClient{
		Responses: make(map[string]map[string]string),
		Errors:    make(map[string]error),
	}
}

// FetchConfig returns the pre-configured response or error for the given URL.
func (m *MockClient) FetchConfig(_ context.Context, url string) (map[string]string, error) {
	m.Calls = append(m.Calls, url)
	if err, ok := m.Errors[url]; ok {
		return nil, err
	}
	if resp, ok := m.Responses[url]; ok {
		return resp, nil
	}
	return map[string]string{}, nil
}

// ConfigFetcher is the interface satisfied by both Client and MockClient.
type ConfigFetcher interface {
	FetchConfig(ctx context.Context, url string) (map[string]string, error)
}
