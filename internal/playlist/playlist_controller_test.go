package playlist

import (
	"reflect"
	"testing"

	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
	"golang.org/x/oauth2"
)

func TestPlaylistController(t *testing.T) {
	t.Run("create playlist", func(t *testing.T) {
		discogsServiceMock := &discogs.DiscogsServiceMock{
			Response: entities.MotherTwoAlbums(),
		}
		spotifyServiceMock := &spotify.SpotifyServiceMock{
			Responses: []string{"spotify:album:1", "spotify:album:2"},
		}
		controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

		playlistId, err := controller.CreatePlaylist(ctx, "discogs-digger")
		if err != nil {
			t.Errorf("error is not nil")
		}
		if playlistId == "" {
			t.Errorf("got empty playlist id, want not empty")
		}
	})

	t.Run("filter duplicates", func(t *testing.T) {
		discogsServiceMock := &discogs.DiscogsServiceMock{
			Response: entities.MotherTwoAlbums(),
		}
		uris := []string{"spotify:album:1", "", "spotify:album:3", "spotify:album:4"}
		spotifyServiceMock := &spotify.SpotifyServiceMock{
			Responses: uris,
		}

		controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		filteredUris := controller.filterNotFounds(uris)

		if len(filteredUris) != 3 {
			t.Errorf("got %d uris, want 3", len(filteredUris))
		}

		expectedUris := []string{"spotify:album:1", "spotify:album:3", "spotify:album:4"}
		if !reflect.DeepEqual(filteredUris, expectedUris) {
			t.Errorf("got %v, want %v", filteredUris, expectedUris)
		}
	})

	t.Run("add to playlist by batches", func(t *testing.T) {
		discogsServiceMock := &discogs.DiscogsServiceMock{}
		spotifyServiceMock := &spotify.SpotifyServiceMock{}
		uris := make([]string, 205)

		controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
		ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

		err := controller.addToSpotifyPlaylist(ctx, "6rqhFgbbKwnb9MLmUQDhG6", uris)
		if err != nil {
			t.Errorf("error is not nil")
		}

		if spotifyServiceMock.CalledCount != 3 {
			t.Errorf("got %d calls, want 3", spotifyServiceMock.CalledCount)
		}
	})
}

func BenchmarkGetAlbumUris(b *testing.B) {
	// TODO: Simulate slow response
	discogsResponses := entities.MotherNAlbums(300)
	discogsServiceMock := &discogs.DiscogsServiceMock{
		Response: discogsResponses,
	}
	spotifyServiceMock := &spotify.SpotifyServiceMock{
		Responses: []string{"spotify:album:1", "spotify:album:2"},
	}
	controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.getSpotifyAlbumUris(ctx, discogsResponses)
	}
}
