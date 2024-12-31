package playlist

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/pkg/errors"
)

type PlaylistBuilder struct {
	spotifyService spotify.SpotifyService
	tracks         []string
}

func NewPlaylistBuilder(spotifyService spotify.SpotifyService) *PlaylistBuilder {
	return &PlaylistBuilder{
		spotifyService: spotifyService,
	}
}

func (pb *PlaylistBuilder) AddAlbums(ctx *gin.Context, albums []string) error {
	trackUris, err := pb.getSpotifyTrackUris(ctx, albums)
	if err != nil {
		return err
	}
	pb.tracks = append(pb.tracks, trackUris...)
	return nil
}

func (pb *PlaylistBuilder) CreateAndPopulate(ctx *gin.Context, name, description string) (string, error) {
	playlist, err := pb.spotifyService.CreatePlaylist(ctx, name, description)
	if err != nil {
		return "", errors.Wrap(err, "error creating playlist")
	}
	err = pb.addToSpotifyPlaylist(ctx, playlist.ID, pb.tracks)
	if err != nil {
		return "", errors.Wrap(err, "error adding to playlist")
	}
	return playlist.URL, nil
}

func (pb *PlaylistBuilder) getSpotifyTrackUris(ctx *gin.Context, albums []string) ([]string, error) {
	batckSize := 20
	uris := []string{}
	err := batchRequests(ctx, albums, batckSize, func(ctx *gin.Context, batch []string) error {
		tracks, err := pb.spotifyService.GetAlbumsTrackUris(ctx, batch)
		if err != nil {
			return errors.Wrap(err, "error getting album track uris")
		}
		uris = append(uris, tracks...)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting album track uris")
	}
	return uris, nil
}

func (pb *PlaylistBuilder) addToSpotifyPlaylist(ctx *gin.Context, playlistId string, uris []string) error {
	batchSize := 100
	return batchRequests(ctx, uris, batchSize, func(ctx *gin.Context, batch []string) error {
		err := pb.spotifyService.AddToPlaylist(ctx, playlistId, batch)
		if err != nil {
			return errors.Wrap(err, "error adding to playlist")
		}
		return nil
	})
}

func batchRequests(ctx *gin.Context, totalItems []string, batchSize int, fn func(ctx *gin.Context, batch []string) error) error {
	for i := 0; i < len(totalItems); i += batchSize {
		end := i + batchSize
		if end > len(totalItems) {
			end = len(totalItems)
		}
		batch := totalItems[i:end]
		err := fn(ctx, batch)
		if err != nil {
			return errors.Wrap(err, "error processing batch")
		}
	}
	return nil
}
