package usecases

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

func TestSpotifyAuthenticate(t *testing.T) {
	t.Run("get auth url", func(t *testing.T) {
		controller := NewSpotifyAuthenticate("client_id", "client_secret", "redirect_uri")
		redirectURL := controller.GetAuthURL()

		want := "https://accounts.spotify.com/authorize?access_type=offline&client_id=client_id&redirect_uri=redirect_uri&response_type=code&scope=user-read-private+user-read-email+playlist-modify-public+playlist-modify-private&state=" +
			url.QueryEscape(oauthState)

		if redirectURL != want {
			t.Errorf("got %s, want %s", redirectURL, want)
		}
	})

	t.Run("store token on gorilla session", func(t *testing.T) {
		t.Setenv("SESSION_KEY", "session_key")
		s := session.NewGorillaSession()
		s.Init(60)
		controller := NewSpotifyAuthenticate("client_id", "client_secret", "redirect_uri")
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

	t.Run("generate token with error in callback", func(t *testing.T) {
		controller := NewSpotifyAuthenticate("client_id", "client_secret", "redirect_uri")
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?error=access_denied", nil)
		const expectedError = ErrErrorInCallback + ": access_denied"

		_, err := controller.GenerateToken(ctx)
		if err == nil || err.Error() != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("generate token with no code in callback", func(t *testing.T) {
		controller := NewSpotifyAuthenticate("client_id", "client_secret", "redirect_uri")
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?state="+oauthState, nil)

		_, err := controller.GenerateToken(ctx)
		if err == nil || err.Error() != ErrNoCode {
			t.Errorf("expected error %v, got %v", ErrNoCode, err)
		}
	})

	t.Run("generate token with state mismatch", func(t *testing.T) {
		controller := NewSpotifyAuthenticate("client_id", "client_secret", "redirect_uri")
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest("GET", "/callback?code=auth_code&state=wrong_state", nil)

		_, err := controller.GenerateToken(ctx)
		if err == nil || err.Error() != ErrRedirectStateParamMismatch {
			t.Errorf("expected error %v, got %v", ErrRedirectStateParamMismatch, err)
		}
	})
}
