package playlist

import (
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
	if err != nil {
		return "", err
	}
	albums := parseAlbumsFromReleases(releases)
	spotifyUris := []string{}
	for _, album := range albums {
		uri, err := c.spotifyService.GetAlbumUri(ctx, album)
		if err != nil {
			return "", errors.Wrap(err, "error getting album uri")
		}
		spotifyUris = append(spotifyUris, uri)
	}
	filterNotFounds(spotifyUris)
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
	return playlistId, nil
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
	// Remove duplicates
	return albums
}

// TODO: Necessary?
func joinArtists(artists []entities.DiscogsArtist) string {
	names := []string{}
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return strings.Join(names, ", ")
}

func filterNotFounds(uris []string) {
	j := 0
	for _, uri := range uris {
		if uri != "" {
			uris[j] = uri
			j++
		}
	}
	uris = uris[:j]
}
