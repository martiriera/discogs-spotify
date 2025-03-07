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
	t.Run("create playlist flow", func(t *testing.T) {
		discogsServiceMock := &discogs.DiscogsServiceMock{
			Response: entities.MotherTwoAlbums(),
		}
		spotifyServiceMock := &spotify.SpotifyServiceMock{
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

	t.Run("parse discogs url", func(t *testing.T) {
		tcs := []struct {
			name        string
			url         string
			expected    *entities.DiscogsInputUrl
			expectError bool
		}{
			{
				"short",
				"discogs.com/user/digger/collection",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.CollectionType},
				false,
			},
			{
				"https",
				"https://www.discogs.com/user/digger/collection",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.CollectionType},
				false,
			},
			{
				"https es",
				"https://www.discogs.com/es/user/digger/collection",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.CollectionType},
				false,
			},
			{
				"www es",
				"www.discogs.com/es/user/digger/collection",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.CollectionType},
				false,
			},
			{
				"https other user",
				"https://www.discogs.com/user/johndoe/collection",
				&entities.DiscogsInputUrl{Id: "johndoe", Type: entities.CollectionType},
				false,
			},
			{
				"https with header",
				"https://www.discogs.com/user/digger/collection?header=1",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.CollectionType},
				false,
			},
			{
				"https list",
				"https://www.discogs.com/lists/MyList/1545836",
				&entities.DiscogsInputUrl{Id: "1545836", Type: entities.ListType},
				false,
			},
			{
				"www list",
				"www.discogs.com/lists/MyList/1545836",
				&entities.DiscogsInputUrl{Id: "1545836", Type: entities.ListType},
				false,
			},
			{
				"https wantlist",
				"https://www.discogs.com/wantlist?user=digger",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.WantlistType},
				false,
			},
			{
				"www wantlist",
				"www.discogs.com/wantlist?user=digger",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.WantlistType},
				false,
			},
			{
				"short wantlist es",
				"discogs.com/es/wantlist?user=digger",
				&entities.DiscogsInputUrl{Id: "digger", Type: entities.WantlistType},
				false,
			},
			{
				"wrong",
				"https://www.discogs.com",
				nil,
				true,
			},
			{
				"wrong collection",
				"www.discogs.com/user/digger",
				nil,
				true,
			},
			{
				"incomplete collection",
				"www.discogs.com/user/digger/collectio",
				nil,
				true,
			},
			{
				"wrong query",
				"https://www.discogs.com/wantlist?digger",
				nil,
				true,
			},
			{
				"wrong lists",
				"https://www.discogs.com/user/digger/lists",
				nil,
				true,
			},
		}
		for _, tc := range tcs {
			got, err := parseDiscogsUrl(tc.url)
			if (err != nil) != tc.expectError {
				t.Errorf("error = %v, expectError = %v", err, tc.expectError)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("got %v, want %v: %s", got, tc.expected, tc.name)
			}
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
		_, err := controller.getSpotifyAlbumIds(ctx, discogsResponses)
		if err != nil {
			b.Errorf("did not expect error, got %v", err)
		}
		elapsed := time.Since(start).Seconds()
		b.Logf("Iteration %d took %f seconds", i, elapsed)
	}
}
