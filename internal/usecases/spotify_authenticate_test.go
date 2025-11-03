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

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

type mockOauth2Config struct {
	authCodeURL  func(string, ...oauth2.AuthCodeOption) string
	exchangeFunc func(string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

func (m *mockOauth2Config) AuthCodeURL(state string, _ ...oauth2.AuthCodeOption) string {
	if m.authCodeURL != nil {
		return m.authCodeURL(state)
	}
	return "https://accounts.spotify.com/authorize?access_type=offline&client_id=test_client_id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=user-read-private+user-read-email+playlist-modify-public+playlist-modify-private&state=" + url.QueryEscape(state)
}

func (m *mockOauth2Config) Exchange(_ context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if m.exchangeFunc != nil {
		return m.exchangeFunc(code, opts...)
	}
	return &oauth2.Token{
		AccessToken: "test_access_token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}, nil
}

type mockSession struct {
	data         map[session.ContextKey]any
	setDataError error
}

func (ms *mockSession) Init(_ int) {}

func (ms *mockSession) Get(_ *http.Request, _ string) (map[any]any, error) {
	result := make(map[any]any)
	for k, v := range ms.data {
		result[k] = v
	}
	return result, nil
}

func (ms *mockSession) GetData(_ *http.Request, key session.ContextKey) (any, error) {
	return ms.data[key], nil
}

func (ms *mockSession) SetData(_ *http.Request, _ http.ResponseWriter, key session.ContextKey, value any) error {
	if ms.setDataError != nil {
		return ms.setDataError
	}
	if ms.data == nil {
		ms.data = make(map[session.ContextKey]any)
	}
	ms.data[key] = value
	return nil
}

func createMockConfig() *mockOauth2Config {
	return &mockOauth2Config{}
}

func TestSpotifyAuthenticate_GetAuthURL(t *testing.T) {
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	redirectURL := controller.GetAuthURL()

	want := "https://accounts.spotify.com/authorize?access_type=offline&client_id=test_client_id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=user-read-private+user-read-email+playlist-modify-public+playlist-modify-private&state=" +
		url.QueryEscape(oauthState)

	if redirectURL != want {
		t.Errorf("got %s, want %s", redirectURL, want)
	}
}

func TestSpotifyAuthenticate_StoreTokenOnGorillaSession(t *testing.T) {
	t.Setenv("SESSION_KEY", "session_key")
	s := session.NewGorillaSession()
	s.Init(60)
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", http.NoBody)

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
}

func TestSpotifyAuthenticate_StoreTokenWithSessionError(t *testing.T) {
	mockSession := &mockSession{
		setDataError: errors.New("session error"),
	}
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/", http.NoBody)

	token := &oauth2.Token{
		AccessToken: "access_token",
	}
	err := controller.StoreToken(ctx, mockSession, token)

	if err == nil || !strings.Contains(err.Error(), ErrSavingSession) {
		t.Errorf("expected error containing %s, got %v", ErrSavingSession, err)
	}
}

func TestSpotifyAuthenticate_GenerateTokenFromGinWithErrorInCallback(t *testing.T) {
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/callback?error=access_denied", http.NoBody)
	const expectedError = ErrErrorInCallback + ": access_denied"

	_, err := controller.GenerateTokenFromGin(ctx)
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestSpotifyAuthenticate_GenerateTokenFromGinWithNoCodeInCallback(t *testing.T) {
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/callback?state="+oauthState, http.NoBody)

	_, err := controller.GenerateTokenFromGin(ctx)
	if err == nil || err.Error() != ErrNoCode {
		t.Errorf("expected error %v, got %v", ErrNoCode, err)
	}
}

func TestSpotifyAuthenticate_GenerateTokenFromGinWithStateMismatch(t *testing.T) {
	mockConfig := createMockConfig()
	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/callback?code=auth_code&state=wrong_state", http.NoBody)

	_, err := controller.GenerateTokenFromGin(ctx)
	if err == nil || err.Error() != ErrRedirectStateParamMismatch {
		t.Errorf("expected error %v, got %v", ErrRedirectStateParamMismatch, err)
	}
}

func TestSpotifyAuthenticate_GenerateTokenFromGinSuccessfully(t *testing.T) {
	expectedToken := &oauth2.Token{
		AccessToken: "test_access_token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockConfig := createMockConfig()
	mockConfig.exchangeFunc = func(code string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
		if code != "test_auth_code" {
			t.Errorf("got code %s, want test_auth_code", code)
		}
		return expectedToken, nil
	}

	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/callback?code=test_auth_code&state="+oauthState, http.NoBody)

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
}

func TestSpotifyAuthenticate_GenerateTokenWithExchangeError(t *testing.T) {
	mockConfig := createMockConfig()
	mockConfig.exchangeFunc = func(_ string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
		return nil, errors.New("exchange error")
	}

	controller := NewSpotifyAuthenticateWithConfig(mockConfig, oauthState)
	ctx := context.Background()

	_, err := controller.GenerateToken(ctx, "test_auth_code")
	if err == nil || !strings.Contains(err.Error(), ErrExchangingCode) {
		t.Errorf("expected error containing %s, got %v", ErrExchangingCode, err)
	}
}
