package usecases

import (
	"reflect"
	"testing"

	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/session"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/util"
)

func TestPlaylistController(t *testing.T) {
	t.Run("create playlist flow", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherTwoAlbums(),
		}
		spotifyServiceMock := &spotify.ServiceMock{
			Responses: []string{"spotify:album:1", "spotify:album:2"},
		}
		controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

		playlist, err := controller.CreatePlaylist(ctx, "https://www.discogs.com/user/digger/collection")
		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}
		if playlist.DiscogsReleases != 2 {
			t.Errorf("got %d releases, want 2", playlist.DiscogsReleases)
		}
		if playlist.SpotifyAlbums != 2 {
			t.Errorf("got %d albums, want 2", playlist.SpotifyAlbums)
		}
		if playlist.SpotifyPlaylist.ID != "6rqhFgbbKwnb9MLmUQDhG6" {
			t.Errorf("got %s, want 6rqhFgbbKwnb9MLmUQDhG6", playlist.SpotifyPlaylist.ID)
		}
		if playlist.SpotifyPlaylist.URL != "https://open.spotify.com/playlist/6rqhFgbbKwnb9MLmUQDhG6" {
			t.Errorf("got %s, want https://open.spotify.com/playlist/6rqhFgbbKwnb9MLmUQDhG6", playlist.SpotifyPlaylist.URL)
		}
	})

	t.Run("filter duplicates and not founds", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{}
		uris := []string{"spotify:album:1", "spotify:album:1", "spotify:album:2", "", "spotify:album:3"}
		spotifyServiceMock := &spotify.ServiceMock{
			Responses: uris,
		}

		controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		filteredUris := controller.filterValidUnique(uris)

		if len(filteredUris) != 3 {
			t.Errorf("got %d uris, want 3", len(filteredUris))
		}

		expectedUris := []string{"spotify:album:1", "spotify:album:2", "spotify:album:3"}
		if !reflect.DeepEqual(filteredUris, expectedUris) {
			t.Errorf("got %v, want %v", filteredUris, expectedUris)
		}
	})
}
