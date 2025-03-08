package playlist

import (
	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/spotify"
)

type Builder struct {
	spotifyService spotify.Service
	tracks         []string
}

func NewPlaylistBuilder(spotifyService spotify.Service) *Builder {
	return &Builder{
		spotifyService: spotifyService,
	}
}

func (pb *Builder) AddAlbums(ctx *gin.Context, albums []string) error {
	trackUris, err := pb.getSpotifyTrackUris(ctx, albums)
	if err != nil {
		return err
	}
	pb.tracks = append(pb.tracks, trackUris...)
	return nil
}

func (pb *Builder) CreateAndPopulate(ctx *gin.Context, name, description string) (*entities.SpotifyPlaylist, error) {
	playlist, err := pb.spotifyService.CreatePlaylist(ctx, name, description)
	if err != nil {
		return nil, errors.Wrap(err, "error creating playlist")
	}
	err = pb.addToSpotifyPlaylist(ctx, playlist.ID, pb.tracks)
	if err != nil {
		return nil, errors.Wrap(err, "error adding to playlist")
	}
	return &playlist, nil
}

func (pb *Builder) getSpotifyTrackUris(ctx *gin.Context, albums []string) ([]string, error) {
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

func (pb *Builder) addToSpotifyPlaylist(ctx *gin.Context, playlistID string, tracks []string) error {
	batchSize := 100
	return batchRequests(ctx, tracks, batchSize, func(ctx *gin.Context, batch []string) error {
		err := pb.spotifyService.AddToPlaylist(ctx, playlistID, batch)
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
