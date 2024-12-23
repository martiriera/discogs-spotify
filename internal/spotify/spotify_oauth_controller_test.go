package spotify

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"golang.org/x/oauth2"
)

func TestSpotifyOauthController(t *testing.T) {
	t.Run("get auth url", func(t *testing.T) {
		controller := NewOAuthController("client_id", "client_secret", "redirect_uri", []string{"scope1", "scope2"})
		redirectURL := controller.GetAuthUrl()

		want := "https://accounts.spotify.com/authorize?access_type=offline&client_id=client_id&redirect_uri=redirect_uri&response_type=code&scope=scope1+scope2&state=" + oauthState

		if redirectURL != want {
			t.Errorf("got %s, want %s", redirectURL, want)
		}
	})

	t.Run("store token", func(t *testing.T) {
		t.Setenv("SESSION_KEY", "session_key")
		s := session.NewGorillaSession()
		s.Init()
		controller := NewOAuthController("client_id", "client_secret", "redirect_uri", []string{"scope1", "scope2"})
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("POST", "/", nil)

		token := &oauth2.Token{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			Expiry:       time.Now().Add(time.Hour),
			TokenType:    "token_type",
		}
		err := controller.StoreToken(c, s, token)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		tokenJson, err := s.GetData(c.Request, session.SpotifyTokenKey)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		var storedToken oauth2.Token
		err = json.Unmarshal([]byte(tokenJson.(string)), &storedToken)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if storedToken.AccessToken != token.AccessToken {
			t.Errorf("got %s, want %s", storedToken.AccessToken, token.AccessToken)
		}
	})
}
