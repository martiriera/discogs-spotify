package usecases

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type Controller struct {
	importer  *DiscogsProcessURL
	builder   *SpotifyCreatePlaylist
	converter *DiscogsConvertToSpotify
}

func NewPlaylistController(discogsService ports.DiscogsPort, spotifyService ports.SpotifyPort) *Controller {
	return &Controller{
		importer:  NewDiscogsProcessURL(discogsService),
		builder:   NewSpotifyCreatePlaylist(spotifyService),
		converter: NewDiscogsConvertToSpotify(spotifyService),
	}
}

func (c *Controller) CreatePlaylist(ctx *gin.Context, discogsURL string) (*entities.Playlist, error) {
	stop := StartTimer("CreatePlaylist")
	defer stop()

	// fetch releases
	releases, parsedDiscogsURL, err := c.importer.processDiscogsURL(discogsURL)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("no releases found on Discogs list")
	}

	// process album IDs
	albumIDs, err := c.converter.getSpotifyAlbumIDs(ctx, releases)
	if err != nil {
		return nil, errors.Wrap(err, "error getting spotify album uris")
	}
	albumIDs = c.filterValidUnique(albumIDs)

	// create playlist
	err = c.builder.AppendAlbumsTracks(ctx, albumIDs)
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

	return &entities.Playlist{
		DiscogsReleases: len(releases),
		SpotifyAlbums:   len(albumIDs),
		SpotifyPlaylist: *playlist,
	}, nil
}

func (c *Controller) filterValidUnique(uris []string) []string {
	seen := map[string]bool{}
	filtered := []string{}
	for _, uri := range uris {
		if uri != "" && !seen[uri] {
			filtered = append(filtered, uri)
			seen[uri] = true
		}
	}
	return filtered
}

func getAlbumFromRelease(release entities.DiscogsRelease) entities.Album {
	album := entities.Album{
		Artist: release.BasicInformation.Artists[0].Name,
		Title:  strings.TrimSpace(release.BasicInformation.Title),
	}
	return album
}
