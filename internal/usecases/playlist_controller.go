package usecases

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/database"
)

type Controller struct {
	importer  *DiscogsProcessURL
	builder   *SpotifyCreatePlaylist
	converter *DiscogsConvertToSpotify
	store     *StorePlaylistUseCase
}

func NewPlaylistController(discogsService ports.DiscogsPort, spotifyService ports.SpotifyPort, repository *database.PlaylistRepository) *Controller {
	return &Controller{
		importer:  NewDiscogsProcessURL(discogsService),
		builder:   NewSpotifyCreatePlaylist(spotifyService),
		converter: NewDiscogsConvertToSpotify(spotifyService),
		store:     NewStorePlaylistUseCase(spotifyService, repository),
	}
}

func (c *Controller) CreatePlaylist(ctx context.Context, discogsURL string) (*entities.Playlist, error) {
	stop := StartTimer("CreatePlaylist")
	defer stop()

	parsedDiscogsURL, err := parseDiscogsURL(discogsURL)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing Discogs URL")
	}

	// fetch releases
	releases, err := c.importer.processDiscogsURL(parsedDiscogsURL)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("no releases found on Discogs list")
	}

	// process album IDs
	albums, err := c.converter.getSpotifyAlbums(ctx, releases)
	if err != nil {
		return nil, errors.Wrap(err, "error getting spotify album uris")
	}
	albums = c.filterValidUnique(albums)

	// create playlist
	albumURIs := make([]string, len(albums))
	for i, album := range albums {
		albumURIs[i] = album.SpotifyURI
	}
	err = c.builder.AppendAlbumsTracks(ctx, albumURIs)
	if err != nil {
		return nil, errors.Wrap(err, "error adding albums to playlist builder")
	}
	playlist, err := c.builder.CreateAndPopulate(
		ctx,
		"Discogs "+cases.Title(language.English).String(parsedDiscogsURL.Type.String())+" by "+parsedDiscogsURL.ID,
		"Created from: "+discogsURL,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating and populating playlist")
	}

	// Store playlist data asynchronously
	var wg sync.WaitGroup
	c.store.ExecuteAsync(ctx, albums, &wg)
	// Note: We don't wait for the store operation to complete

	return &entities.Playlist{
		DiscogsReleases: len(releases),
		SpotifyAlbums:   len(albums),
		SpotifyPlaylist: *playlist,
	}, nil
}

func (c *Controller) filterValidUnique(albums []entities.Album) []entities.Album {
	seen := map[entities.Album]bool{}
	filtered := []entities.Album{}
	for _, album := range albums {
		if album.SpotifyURI != "" && !seen[album] {
			filtered = append(filtered, album)
			seen[album] = true
		}
	}
	return filtered
}
