package main

import (
	"net/http/httptest"
	"testing"

	"github.com/martiriera/discogs-spotify/entities"
)

func TestServer(t *testing.T) {
	discogsServiceMock := &DiscogsServiceMock{
		response: entities.MotherTwoAlbums(),
	}
	spotifyServiceMock := &SpotifyServiceMock{
		responses: []string{"spotify:album:1", "spotify:album:2"},
	}
	playlistCreator := newPlaylistCreator(discogsServiceMock, spotifyServiceMock)
	server := NewServer(playlistCreator)

	request := httptest.NewRequest("POST", "/create-playlist?username=test", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseStatus(t, response.Code, 200)
	assertResponseBody(t, response.Body.String(), `["spotify:album:1","spotify:album:2"]`)
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
