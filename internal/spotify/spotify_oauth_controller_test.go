package spotify

import (
	"testing"
)

func TestSpotifyOauthController(t *testing.T) {
	t.Run("get redirection url", func(t *testing.T) {
		controller := NewOAuthController("client_id", "client_secret", "redirect_uri", []string{"scope1", "scope2"})
		redirectURL := controller.GetRedirectionUrl()

		want := "https://accounts.spotify.com/authorize?access_type=offline&client_id=client_id&redirect_uri=redirect_uri&response_type=code&scope=scope1+scope2&state=random-string"

		if redirectURL != want {
			t.Errorf("got %s, want %s", redirectURL, want)
		}
	})

	// TODO: Add test for GetServiceFromCallback
}
