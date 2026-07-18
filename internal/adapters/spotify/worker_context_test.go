package spotify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

func TestWorkerContextProvider_GetToken(t *testing.T) {
	t.Run("returns token on successful client credentials fetch", func(t *testing.T) {
		stubClient := &StubSpotifyHTTPClient{
			Responses: []*http.Response{
				{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"access_token": "worker-access-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}`)),
				},
			},
		}
		provider := NewClientCredentialsProvider("client-id", "client-secret", stubClient)
		workerCtx := NewWorkerContextProvider(provider)

		token, err := workerCtx.GetToken(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == nil {
			t.Fatal("expected non-nil token")
		}
		if token.AccessToken != "worker-access-token" {
			t.Errorf("AccessToken = %q, want %q", token.AccessToken, "worker-access-token")
		}
	})

	t.Run("returns error when token endpoint fails", func(t *testing.T) {
		stubClient := &StubSpotifyHTTPClient{
			Responses: []*http.Response{
				{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "invalid_client"}`)),
				},
			},
		}
		provider := NewClientCredentialsProvider("bad-id", "bad-secret", stubClient)
		workerCtx := NewWorkerContextProvider(provider)

		_, err := workerCtx.GetToken(context.Background())
		if err == nil {
			t.Fatal("expected error when token endpoint returns 401, got nil")
		}
	})

	t.Run("caches token and avoids second HTTP call", func(t *testing.T) {
		stubClient := &StubSpotifyHTTPClient{
			Responses: []*http.Response{
				{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"access_token": "cached-token",
						"token_type": "Bearer",
						"expires_in": 3600
					}`)),
				},
			},
		}
		provider := NewClientCredentialsProvider("client-id", "client-secret", stubClient)
		workerCtx := NewWorkerContextProvider(provider)

		token1, err := workerCtx.GetToken(context.Background())
		if err != nil {
			t.Fatalf("first call failed: %v", err)
		}

		// Second call should use the cached token (stub has no more responses).
		token2, err := workerCtx.GetToken(context.Background())
		if err != nil {
			t.Fatalf("second call failed: %v", err)
		}
		if token1.AccessToken != token2.AccessToken {
			t.Errorf("expected same cached token, got %q and %q", token1.AccessToken, token2.AccessToken)
		}
	})
}

func TestWorkerContextProvider_GetUserID(t *testing.T) {
	provider := NewClientCredentialsProvider("id", "secret", &StubSpotifyHTTPClient{})
	workerCtx := NewWorkerContextProvider(provider)

	userID, err := workerCtx.GetUserID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "" {
		t.Errorf("GetUserID() = %q, want empty string", userID)
	}
}

func TestWorkerContextProvider_SetUserID(t *testing.T) {
	provider := NewClientCredentialsProvider("id", "secret", &StubSpotifyHTTPClient{})
	workerCtx := NewWorkerContextProvider(provider)

	err := workerCtx.SetUserID(context.Background(), "any-user")
	if err != nil {
		t.Errorf("SetUserID() returned unexpected error: %v", err)
	}
}
