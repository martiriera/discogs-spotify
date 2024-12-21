package server

import (
	"net/http/httptest"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/playlist"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

func TestServer(t *testing.T) {
	discogsServiceMock := &discogs.DiscogsServiceMock{
		Response: entities.MotherTwoAlbums(),
	}
	spotifyServiceMock := &spotify.SpotifyServiceMock{
		Responses: []string{"spotify:album:1", "spotify:album:2"},
	}

	t.Run("create playlist", func(t *testing.T) {
		playlistCreator := playlist.NewPlaylistCreator(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistCreator)
		request := httptest.NewRequest("POST", "/create-playlist?username=test", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), `["spotify:album:1","spotify:album:2"]`)
	})

	t.Run("create playlist without username", func(t *testing.T) {
		playlistCreator := playlist.NewPlaylistCreator(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistCreator)
		request := httptest.NewRequest("POST", "/create-playlist", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 400)
		assertResponseBody(t, response.Body.String(), "username is required\n")
	})

	t.Run("server error from discogs", func(t *testing.T) {
		discogsServiceMock.Error = discogs.ErrUnexpectedStatus
		playlistCreator := playlist.NewPlaylistCreator(discogsServiceMock, spotifyServiceMock)
		server := NewServer(playlistCreator)
		request := httptest.NewRequest("POST", "/create-playlist?username=test", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, 500)
		assertResponseBody(t, response.Body.String(), "{\"error\":\"discogs unexpected status error\"}\n")
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
