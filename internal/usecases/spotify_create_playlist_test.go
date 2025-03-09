package usecases

import (
	"testing"

	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"

	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"github.com/martiriera/discogs-spotify/util"
)

func TestSpotifyCreatePlaylist(t *testing.T) {
	t.Run("add to playlist by batches", func(t *testing.T) {
		spotifyServiceMock := &spotify.ServiceMock{}
		uris := make([]string, 205)

		builder := NewSpotifyCreatePlaylist(spotifyServiceMock)
		ctx := util.NewTestContextWithToken(session.SpotifyTokenKey, &oauth2.Token{AccessToken: "test"})

		err := builder.addToSpotifyPlaylist(ctx, "6rqhFgbbKwnb9MLmUQDhG6", uris)
		if err != nil {
			t.Errorf("error is not nil")
		}

		// Additions are done in batches of 100
		if spotifyServiceMock.CalledCount != 3 {
			t.Errorf("got %d calls, want 3", spotifyServiceMock.CalledCount)
		}
	})

	t.Run("batch requests function", func(t *testing.T) {
		tcs := []struct {
			name      string
			uris      []string
			batchSize int
			wantCalls int
		}{
			{
				name:      "empty uris",
				uris:      []string{},
				batchSize: 10,
				wantCalls: 0,
			},
			{
				name:      "less than batch size",
				uris:      make([]string, 5),
				batchSize: 10,
				wantCalls: 1,
			},
			{
				name:      "equal to batch size",
				uris:      make([]string, 10),
				batchSize: 10,
				wantCalls: 1,
			},
			{
				name:      "more than batch size",
				uris:      make([]string, 15),
				batchSize: 10,
				wantCalls: 2,
			},
		}

		gotCalls := 0
		testFunc := func(_ *gin.Context, _ []string) error {
			gotCalls++
			return nil
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				gotCalls = 0
				err := batchRequests(nil, tc.uris, tc.batchSize, testFunc)
				if err != nil {
					t.Errorf("error is not nil")
				}

				if gotCalls != tc.wantCalls {
					t.Errorf("got %d calls, want 3", gotCalls)
				}
			})
		}
	})

}
