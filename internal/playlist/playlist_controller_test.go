package playlist

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

func TestPlaylistController(t *testing.T) {
	discogsServiceMock := &discogs.DiscogsServiceMock{
		Response: entities.MotherTwoAlbums(),
	}
	spotifyServiceMock := &spotify.SpotifyServiceMock{
		Responses: []string{"spotify:album:1", "spotify:album:2"},
	}
	controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)

	uris, err := controller.CreatePlaylist("digger")
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
// 	playlistCreator := NewPlaylistCreator(
// 		discogs.NewHttpDiscogsService(&http.Client{}),
// 		spotify.NewHttpSpotifyService(&http.Client{}, ""),
// 	)

// 	uris, err := playlistCreator.CreatePlaylist("martireir")
// 	if err != nil {
// 		t.Errorf("error is not nil")
// 	}
// 	if len(uris) == 0 {
// 		t.Errorf("got %d albums, want at least 1", len(uris))
// 	}
// }
