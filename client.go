package deadmanssnitch

import (
	"fmt"
	"net/http"
)

const (
	apiEndpoint = "https://api.deadmanssnitch.com"
	apiVersion  = "v1"
)

// Client is the Dead Man's Snitch API client
type Client struct {
	httpClient *http.Client
	apiBaseURL string
	apiKey     string
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	client := &Client{
		httpClient: &http.Client{},
		apiBaseURL: fmt.Sprintf("%s/%s", apiEndpoint, apiVersion),
		apiKey:     apiKey,
	}

	return client
}
