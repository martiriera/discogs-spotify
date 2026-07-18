package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpClient "github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

const tokenURL = "https://accounts.spotify.com/api/token" //nolint:gosec // not a credential, it's a public endpoint URL

// tokenExpiryBuffer is subtracted from the token's declared expiry to ensure
// we refresh before the token actually expires.
const tokenExpiryBuffer = 30 * time.Second

// ClientCredentialsProvider fetches and caches a Spotify access token
// using the Client Credentials OAuth2 flow (machine-to-machine, no user).
type ClientCredentialsProvider struct {
	clientID     string
	clientSecret string
	httpClient   httpClient.HTTPClient
	token        string
	expiry       time.Time
}

func NewClientCredentialsProvider(clientID, clientSecret string, c httpClient.HTTPClient) *ClientCredentialsProvider {
	return &ClientCredentialsProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   c,
	}
}

// Token returns a valid access token, refreshing if necessary.
func (p *ClientCredentialsProvider) Token(ctx context.Context) (string, error) {
	if p.token != "" && time.Now().Before(p.expiry) {
		return p.token, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("building token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.clientID, p.clientSecret)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching spotify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify token endpoint returned %d", resp.StatusCode)
	}

	var tokenResp entities.SpotifyAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding token response: %w", err)
	}

	p.token = tokenResp.AccessToken
	// subtract 30s buffer so we refresh before actual expiry
	p.expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn)*time.Second - tokenExpiryBuffer)

	return p.token, nil
}
