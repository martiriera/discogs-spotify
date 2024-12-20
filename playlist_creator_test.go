package main

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/entities"
)

type DiscogsServiceMock struct {
	response []entities.DiscogsRelease
}

func (m *DiscogsServiceMock) GetReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return m.response, nil
}

type SpotifyServiceMock struct {
	responses []string
	index     int
}

func (m *SpotifyServiceMock) GetAlbumUri(artist string, title string) (string, error) {
	if m.index >= len(m.responses) {
		return "", nil
	}
	response := m.responses[m.index]
	m.index++
	return response, nil
}

func TestPlaylistCreator(t *testing.T) {
	discogsServiceMock := &DiscogsServiceMock{
		response: entities.MotherTwoAlbums(),
	}
	spotifyServiceMock := &SpotifyServiceMock{
		responses: []string{"spotify:album:1", "spotify:album:2"},
	}
	playlistCreator := newPlaylistCreator(discogsServiceMock, spotifyServiceMock)

	uris, err := playlistCreator.CreatePlaylist("digger")
	if err != nil {
		t.Errorf("error is not nil")
	}
	if len(uris) != 2 {
		t.Errorf("got %d albums, want 2", len(uris))
	}
	want := []string{"spotify:album:1", "spotify:album:2"}
	if !reflect.DeepEqual(uris, want) {
		t.Errorf("got %v, want %v", uris, want)
	}
}

// func TestRealPlaylistCreator(t *testing.T) {
// 	t.Setenv("SPOTIFY_CLIENT_ID", "0e0db614950547a9848c20f23c38ced3")
// 	t.Setenv("SPOTIFY_CLIENT_SECRET", "c9f69a062c9c482a95e77a200a9c04f4")
// 	playlistCreator := newPlaylistCreator(
// 		NewHttpDiscogsService(&http.Client{}),
// 		NewHttpSpotifyService(&http.Client{}, ""),
// 	)

// 	uris, err := playlistCreator.CreatePlaylist("martireir")
// 	if err != nil {
// 		t.Errorf("error is not nil")
// 	}
// 	if len(uris) == 0 {
// 		t.Errorf("got %d albums, want at least 1", len(uris))
// 	}
// }
