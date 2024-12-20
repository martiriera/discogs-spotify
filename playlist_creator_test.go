package main

import (
	"testing"

	"github.com/martiriera/discogs-spotify/entities"
)

type DiscogsServiceMock struct{}

func (m *DiscogsServiceMock) GetReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return []entities.DiscogsRelease{}, nil
}

type SpotifyServiceMock struct{}

func (m *SpotifyServiceMock) GetAlbumUri(artist string, album string) (string, error) {
	return "spotify:album:1", nil
}

func TestPlaylistCreator(t *testing.T) {
	discogsServiceMock := &DiscogsServiceMock{}
	spotifyServiceMock := &SpotifyServiceMock{}
	playlistCreator := newPlaylistCreator(discogsServiceMock, spotifyServiceMock)

	_, err := playlistCreator.CreatePlaylist("username")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
