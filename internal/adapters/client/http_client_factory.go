package client

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HTTPClientFactory struct{}

func NewHTTPClientFactory() *HTTPClientFactory {
	return &HTTPClientFactory{}
}

func (f *HTTPClientFactory) CreateClient(timeout time.Duration, retryAttempts int, retryDelay time.Duration) HTTPClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retryAttempts
	retryClient.RetryWaitMin = retryDelay
	retryClient.RetryWaitMax = retryDelay * 10

	retryClient.HTTPClient.Timeout = timeout

	retryClient.Logger = nil // Disable default logger

	return retryClient.StandardClient()
}

func (f *HTTPClientFactory) CreateDiscogsClient(timeout time.Duration, retryAttempts int, retryDelay time.Duration) HTTPClient {
	client := f.CreateClient(timeout, retryAttempts, retryDelay)


	transport := client.(*http.Client).Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	client.(*http.Client).Transport = &userAgentTransport{
		base: transport,
		userAgent: "DiscogsSpotify/1.0",
	}

	return client
}

func (f *HTTPClientFactory) CreateSpotifyClient(timeout time.Duration, retryAttempts int, retryDelay time.Duration) HTTPClient {
	client := f.CreateClient(timeout, retryAttempts, retryDelay)

	return client
}

// userAgentTransport is a custom transport that adds a User-Agent header
type userAgentTransport struct {
	base      http.RoundTripper
	userAgent string
}

// RoundTrip implements the http.RoundTripper interface
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}