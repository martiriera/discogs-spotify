package playlist

import (
	"reflect"
	"testing"
	"time"

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

		playlist, err := controller.CreatePlaylist(ctx, "discogs-digger")
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
		discogsServiceMock := &discogs.DiscogsServiceMock{}
		uris := []string{"spotify:album:1", "spotify:album:1", "spotify:album:2", "", "spotify:album:3"}
		spotifyServiceMock := &spotify.SpotifyServiceMock{
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

func BenchmarkGetAlbumUris(b *testing.B) {
	discogsResponses := entities.MotherNAlbums(300)
	discogsServiceMock := &discogs.DiscogsServiceMock{
		Response: discogsResponses,
	}
	spotifyServiceMock := &spotify.SpotifyServiceMock{
		Responses:   []string{"spotify:album:1", "spotify:album:2"},
		SleepMillis: 600,
	}
	controller := NewPlaylistController(discogsServiceMock, spotifyServiceMock)
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		controller.getSpotifyAlbumIds(ctx, discogsResponses)
		elapsed := time.Since(start).Seconds()
		b.Logf("Iteration %d took %f seconds", i, elapsed)
	}
}
