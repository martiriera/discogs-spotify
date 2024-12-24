package server

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"golang.org/x/oauth2"
)

func TestAcceptance(t *testing.T) {
	discogsServiceMock := &discogs.DiscogsServiceMock{
		Response: entities.MotherTwoAlbums(),
	}
	spotifyServiceMock := &spotify.SpotifyServiceMock{
		Responses: []string{"spotify:album:1", "spotify:album:2"},
	}
	oauthController := spotify.NewOAuthController(
		"client_id",
		"client_secret",
		"http://localhost:8080/callback",
	)
	sessionMock := session.NewInMemorySession()
	sessionMock.Init()

	t.Run("serve main", func(t *testing.T) {
		controller := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(controller, oauthController, sessionMock)
		request := httptest.NewRequest("GET", "/api/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
	})

	t.Run("login", func(t *testing.T) {
		controller := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(controller, oauthController, sessionMock)
		request := httptest.NewRequest("GET", "/auth/login", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 307)
	})

	t.Run("create playlist", func(t *testing.T) {
		controller := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		request := httptest.NewRequest("POST", "/api/playlist", strings.NewReader("username=test"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		token := &oauth2.Token{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			Expiry:       time.Now().Add(time.Hour),
			TokenType:    "token_type",
		}
		sessionMock.SetData(request, response, session.SpotifyTokenKey, token)

		server := NewServer(controller, oauthController, sessionMock)
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), `["spotify:album:1","spotify:album:2"]`)
	})

	t.Run("create playlist without username", func(t *testing.T) {
		controller := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(controller, oauthController, sessionMock)
		request := httptest.NewRequest("POST", "/api/playlist", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 400)
		assertResponseBody(t, response.Body.String(), "{\"error\":\"username is required\"}")
	})

	t.Run("server error from discogs", func(t *testing.T) {
		discogsServiceMock.Error = discogs.ErrUnexpectedStatus
		controller := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(controller, oauthController, sessionMock)
		request := httptest.NewRequest("POST", "/api/playlist", strings.NewReader("username=test"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 500)
		assertResponseBody(t, response.Body.String(), "{\"error\":\"discogs unexpected status error\"}")
	})
}

func assertResponseStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
