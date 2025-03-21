package usecases

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

type mockOAuthConfig struct {
	*oauth2.Config
	exchangeFunc func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

func (m *mockOAuthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if m.exchangeFunc != nil {
		return m.exchangeFunc(ctx, code, opts...)
	}
	return nil, nil
}

type mockSession struct {
	ports.SessionPort
	setDataError error
}

func (m *mockSession) SetData(r *http.Request, w http.ResponseWriter, key session.ContextKey, value any) error {
	return m.setDataError
}

func createMockConfig() *mockOAuthConfig {
	return &mockOAuthConfig{
		Config: &oauth2.Config{
			ClientID:     "test_client_id",
			ClientSecret: "test_client_secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       scopes,
			Endpoint:     spotify.Endpoint,
		},
	}
}

func TestSpotifyAuthenticate(t *testing.T) {
	t.Run("get auth url", func(t *testing.T) {
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		redirectURL := controller.GetAuthURL()

		want := "https://accounts.spotify.com/authorize?access_type=offline&client_id=test_client_id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=user-read-private+user-read-email+playlist-modify-public+playlist-modify-private&state=" +
			url.QueryEscape(oauthState)

		if redirectURL != want {
			t.Errorf("got %s, want %s", redirectURL, want)
		}
	})

	t.Run("store token on gorilla session", func(t *testing.T) {
		t.Setenv("SESSION_KEY", "session_key")
		s := session.NewGorillaSession()
		s.Init(60)
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("POST", "/", nil)

		token := &oauth2.Token{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			Expiry:       time.Now().Add(time.Hour),
			TokenType:    "token_type",
		}
		err := controller.StoreToken(ctx, s, token)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		storedTokenData, err := s.GetData(ctx.Request, session.SpotifyTokenKey)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		storedToken, ok := storedTokenData.(*oauth2.Token)
		if !ok {
			t.Errorf("stored token is not of type *oauth2.Token")
		}

		if storedToken.AccessToken != token.AccessToken {
			t.Errorf("got %s, want %s", storedToken.AccessToken, token.AccessToken)
		}
	})

	t.Run("store token with session error", func(t *testing.T) {
		mockSession := &mockSession{
			setDataError: errors.New("session error"),
		}
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("POST", "/", nil)

		token := &oauth2.Token{
			AccessToken: "access_token",
		}
		err := controller.StoreToken(ctx, mockSession, token)

		if err == nil || !strings.Contains(err.Error(), ErrSavingSession) {
			t.Errorf("expected error containing %s, got %v", ErrSavingSession, err)
		}
	})

	t.Run("generate token from gin with error in callback", func(t *testing.T) {
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?error=access_denied", nil)
		const expectedError = ErrErrorInCallback + ": access_denied"

		_, err := controller.GenerateTokenFromGin(ctx)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("generate token from gin with no code in callback", func(t *testing.T) {
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?state="+oauthState, nil)

		_, err := controller.GenerateTokenFromGin(ctx)
		if err == nil || err.Error() != ErrNoCode {
			t.Errorf("expected error %v, got %v", ErrNoCode, err)
		}
	})

	t.Run("generate token from gin with state mismatch", func(t *testing.T) {
		mockConfig := createMockConfig()
		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?code=auth_code&state=wrong_state", nil)

		_, err := controller.GenerateTokenFromGin(ctx)
		if err == nil || err.Error() != ErrRedirectStateParamMismatch {
			t.Errorf("expected error %v, got %v", ErrRedirectStateParamMismatch, err)
		}
	})

	t.Run("generate token from gin successfully", func(t *testing.T) {
		expectedToken := &oauth2.Token{
			AccessToken: "test_access_token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		}

		mockConfig := createMockConfig()
		mockConfig.exchangeFunc = func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
			if code != "test_auth_code" {
				t.Errorf("got code %s, want test_auth_code", code)
			}
			return expectedToken, nil
		}

		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?code=test_auth_code&state="+oauthState, nil)

		token, err := controller.GenerateTokenFromGin(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if token.AccessToken != expectedToken.AccessToken {
			t.Errorf("got access token %s, want %s", token.AccessToken, expectedToken.AccessToken)
		}
		if token.TokenType != expectedToken.TokenType {
			t.Errorf("got token type %s, want %s", token.TokenType, expectedToken.TokenType)
		}
	})

	t.Run("generate token with exchange error", func(t *testing.T) {
		mockConfig := createMockConfig()
		mockConfig.exchangeFunc = func(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
			return nil, errors.New("exchange error")
		}

		controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
		ctx := context.Background()

		_, err := controller.GenerateToken(ctx, "test_auth_code")
		if err == nil || !strings.Contains(err.Error(), ErrExchangingCode) {
			t.Errorf("expected error containing %s, got %v", ErrExchangingCode, err)
		}
	})
}
