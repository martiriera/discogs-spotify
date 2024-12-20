package main

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/entities"
)

type DiscogsServiceMock struct{}

func (m *DiscogsServiceMock) GetReleases(discogsUsername string) ([]entities.DiscogsRelease, error) {
	return []entities.DiscogsRelease{
		{
			Artists: []entities.DiscogsArtist{
				{Name: "artist1"},
			},
			BasicInformation: entities.DiscogsBasicInformation{
				Title: "title1",
			},
		},
		{
			Artists: []entities.DiscogsArtist{
				{Name: "artist2"},
			},
			BasicInformation: entities.DiscogsBasicInformation{
				Title: "title2",
			},
		},
	}, nil
}

type SpotifyServiceMock struct{}

func (m *SpotifyServiceMock) GetAlbumUri(artist string, title string) (string, error) {
	return "spotify:album:1", nil
}

func TestPlaylistCreator(t *testing.T) {
	discogsServiceMock := &DiscogsServiceMock{}
	spotifyServiceMock := &SpotifyServiceMock{}
	playlistCreator := newPlaylistCreator(discogsServiceMock, spotifyServiceMock)

	uris, err := playlistCreator.CreatePlaylist("username")
	if err != nil {
		t.Errorf("error is not nil")
	}
	if len(uris) != 2 {
		t.Errorf("got %d albums, want 2", len(uris))
	}
	want := []string{"spotify:album:1", "spotify:album:1"}
	if !reflect.DeepEqual(uris, want) {
		t.Errorf("got %v, want %v", uris, want)
	}
}
