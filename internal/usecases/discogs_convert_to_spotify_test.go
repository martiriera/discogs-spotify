package usecases

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"github.com/martiriera/discogs-spotify/util"
)

func BenchmarkGetAlbumUris(b *testing.B) {
	discogsResponses := entities.MotherNAlbums(300)
	spotifyServiceMock := &spotify.ServiceMock{
		SearchAlbumResponses: [][]entities.SpotifyAlbumItem{
			entities.MotherSpotifyAlbums()[0:2],
			entities.MotherSpotifyAlbums()[2:4],
		},
		SleepMillis: 600,
	}
	controller := NewDiscogsConvertToSpotify(spotifyServiceMock)
	ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

	b.ResetTimer()
	for i := range b.N {
		start := time.Now()
		_, err := controller.getSpotifyAlbumIDs(ctx, discogsResponses)
		if err != nil {
			b.Errorf("did not expect error, got %v", err)
		}
		elapsed := time.Since(start).Seconds()
		b.Logf("Iteration %d took %f seconds", i, elapsed)
	}
}
