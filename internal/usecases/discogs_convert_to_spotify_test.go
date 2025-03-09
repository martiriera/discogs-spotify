package usecases

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/adapters/session"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/util"
)

func BenchmarkGetAlbumUris(b *testing.B) {
	discogsResponses := entities.MotherNAlbums(300)
	spotifyServiceMock := &spotify.ServiceMock{
		Responses:   []string{"spotify:album:1", "spotify:album:2"},
		SleepMillis: 600,
	}
	controller := NewDiscogsConvertToSpotify(spotifyServiceMock)
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		_, err := controller.getSpotifyAlbumIDs(ctx, discogsResponses)
		if err != nil {
			b.Errorf("did not expect error, got %v", err)
		}
		elapsed := time.Since(start).Seconds()
		b.Logf("Iteration %d took %f seconds", i, elapsed)
	}
}
