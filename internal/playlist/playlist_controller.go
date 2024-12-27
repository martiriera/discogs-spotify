package playlist

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/pkg/errors"
)

type PlaylistController struct {
	discogsService discogs.DiscogsService
	spotifyService spotify.SpotifyService
}

func NewPlaylistController(discogsService discogs.DiscogsService, spotifyService spotify.SpotifyService) *PlaylistController {
	return &PlaylistController{
		discogsService: discogsService,
		spotifyService: spotifyService,
	}
}

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUsername string) (string, error) {
	releases, err := c.discogsService.GetReleases(discogsUsername)
	log.Println("Releases: ", releases)

	if err != nil {
		return "", err
	}

	uris, err := c.getSpotifyAlbumUris(ctx, releases)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify album uris")
	}
	log.Println("URIs: ", uris)

	userId, err := c.spotifyService.GetSpotifyUserId(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}

	// TODO: Store also on session?
	ctx.Set(session.SpotifyUserIdKey, userId)

	playlistId, err := c.spotifyService.CreatePlaylist(ctx, "Discogs Playlist", "Playlist created from Discogs")
	if err != nil {
		return "", errors.Wrap(err, "error creating playlist")
	}

	err = c.spotifyService.AddToPlaylist(ctx, playlistId, uris)
	if err != nil {
		return "", errors.Wrap(err, "error adding to playlist")
	}

	return playlistId, nil
}

func (c *PlaylistController) getSpotifyAlbumUris(ctx *gin.Context, releases []entities.DiscogsRelease) ([]string, error) {
	albums := parseAlbumsFromReleases(releases)
	uris := []string{}
	for _, album := range albums {
		uri, err := c.spotifyService.GetAlbumUri(ctx, album)
		if err != nil {
			return nil, errors.Wrap(err, "error getting album uri")
		}
		uris = append(uris, uri)
	}
	filteredUris := c.filterNotFounds(uris)
	return filteredUris, nil
}

func parseAlbumsFromReleases(releases []entities.DiscogsRelease) []entities.Album {
	albums := []entities.Album{}
	for _, release := range releases {
		album := entities.Album{
			Artist: joinArtists(release.BasicInformation.Artists),
			Title:  strings.TrimSpace(release.BasicInformation.Title),
		}
		albums = append(albums, album)
	}
	return albums
}

func (c *PlaylistController) filterNotFounds(uris []string) []string {
	filtered := []string{}
	for _, uri := range uris {
		if uri != "" {
			filtered = append(filtered, uri)
		}
	}
	return filtered
}

// TODO: Necessary?
func joinArtists(artists []entities.DiscogsArtist) string {
	names := []string{}
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return strings.Join(names, ", ")
}
