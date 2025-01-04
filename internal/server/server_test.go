package server

import (
	"fmt"
	"net/http/httptest"
	"os"
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
	userController := spotify.NewUserController(spotifyServiceMock)

	t.Run("api main get 200", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
	})

	t.Run("auth login post 200", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("GET", "/auth/login", nil)
		response := httptest.NewRecorder()
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 307)
	})

	t.Run("api playlist post 200", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("POST", "/playlist", strings.NewReader("discogs_username=test"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)

		token := &oauth2.Token{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			Expiry:       time.Now().Add(time.Hour),
			TokenType:    "token_type",
		}
		sessionMock.SetData(request, response, session.SpotifyTokenKey, token)

		server := NewServer(playlistController, oauthController, userController, sessionMock)
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
		want := "{\"discogs_releases\":2,\"id\":\"6rqhFgbbKwnb9MLmUQDhG6\",\"spotify_albums\":2,\"url\":\"https://open.spotify.com/playlist/6rqhFgbbKwnb9MLmUQDhG6\"}"
		assertResponseBody(t, response.Body.String(), want)
	})

	t.Run("api playlist post 400 no username", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("POST", "/playlist", nil)
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Minute)})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 400)
		assertResponseBody(t, response.Body.String(), "{\"error\":\"username is required\"}")
	})

	t.Run("api playlist post 500 discogs error", func(t *testing.T) {
		discogsServiceMock.Error = discogs.ErrUnexpectedStatus
		sessionMock := initSessionMock()

		request := httptest.NewRequest("POST", "/playlist", strings.NewReader("discogs_username=test"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Minute)})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 500)
		assertResponseBody(t, response.Body.String(), "{\"error\":\"discogs unexpected status error\"}")
	})

	t.Run("api get home 200", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("GET", "/home", nil)
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Minute)})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)
		fmt.Println(os.Getwd())
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
	})

	t.Run("api get home 302 expired token", func(t *testing.T) {
		sessionMock := initSessionMock()
		request := httptest.NewRequest("GET", "/home", nil)
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Second)})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		time.Sleep(1 * time.Second)
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 302)
	})

	t.Run("api get home 302 expired session", func(t *testing.T) {
		sessionMock := initSessionMock()
		sessionMock.Init(1)
		request := httptest.NewRequest("GET", "/home", nil)
		response := httptest.NewRecorder()
		sessionMock.SetData(request, response, session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test", Expiry: time.Now().Add(time.Minute)})
		playlistController := playlist.NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistController, oauthController, userController, sessionMock)

		// TODO: Find a way to avoid sleep
		time.Sleep(2 * time.Second)
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 302)
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

func initSessionMock() *session.InMemorySession {
	sessionMock := session.NewInMemorySession()
	sessionMock.Init(90)
	return sessionMock
}
