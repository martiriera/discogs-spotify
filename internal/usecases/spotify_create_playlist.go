package usecases

import (
	"context"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type SpotifyCreatePlaylist struct {
	spotifyService ports.SpotifyPort
	tracks         []string
}

func NewSpotifyCreatePlaylist(spotifyService ports.SpotifyPort) *SpotifyCreatePlaylist {
	return &SpotifyCreatePlaylist{
		spotifyService: spotifyService,
	}
}

func (u *SpotifyCreatePlaylist) AppendAlbumsTracks(ctx context.Context, albums []string) error {
	trackUris, err := u.getSpotifyTrackUris(ctx, albums)
	if err != nil {
		return err
	}
	u.tracks = append(u.tracks, trackUris...)
	return nil
}

func (u *SpotifyCreatePlaylist) CreateAndPopulate(ctx context.Context, name, description string) (*entities.SpotifyPlaylist, error) {
	playlist, err := u.spotifyService.CreatePlaylist(ctx, name, description)
	if err != nil {
		return nil, errors.Wrap(err, "error creating playlist")
	}
	err = u.addToSpotifyPlaylist(ctx, playlist.ID, u.tracks)
	if err != nil {
		return nil, errors.Wrap(err, "error adding to playlist")
	}
	return &playlist, nil
}

func (u *SpotifyCreatePlaylist) getSpotifyTrackUris(ctx context.Context, albums []string) ([]string, error) {
	batchSize := 20
	uris := []string{}
	err := batchRequests(ctx, albums, batchSize, func(ctx context.Context, batch []string) error {
		tracks, err := u.spotifyService.GetAlbumsTrackUris(ctx, batch)
		if err != nil {
			return errors.Wrap(err, "error getting album track uris")
		}
		uris = append(uris, tracks...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return uris, nil
}

func (u *SpotifyCreatePlaylist) addToSpotifyPlaylist(ctx context.Context, playlistID string, tracks []string) error {
	batchSize := 100
	return batchRequests(ctx, tracks, batchSize, func(ctx context.Context, batch []string) error {
		err := u.spotifyService.AddToPlaylist(ctx, playlistID, batch)
		if err != nil {
			return errors.Wrap(err, "error adding to playlist")
		}
		return nil
	})
}

func batchRequests(ctx context.Context, totalItems []string, batchSize int, fn func(ctx context.Context, batch []string) error) error {
	for i := 0; i < len(totalItems); i += batchSize {
		end := i + batchSize
		if end > len(totalItems) {
			end = len(totalItems)
		}
		batch := totalItems[i:end]
		if err := fn(ctx, batch); err != nil {
			return err
		}
	}
	return nil
}
