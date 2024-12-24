package playlist

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/spotify"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUsername string) ([]string, error) {
	releases, err := c.discogsService.GetReleases(discogsUsername)
	if err != nil {
		return nil, err
	}
	albums := parseAlbumsFromReleases(releases)
	spotifyUris := []string{}
	for _, album := range albums {
		uri, err := c.spotifyService.GetAlbumUri(ctx, album)
		if err != nil {
			return nil, err
		}
		spotifyUris = append(spotifyUris, uri)
	}
	fmt.Printf("Spotify URIs: %v", spotifyUris)
	return filterNotFounds(spotifyUris), nil
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

func filterNotFounds(uris []string) []string {
	filtered := []string{}
	for _, uri := range uris {
		if uri != "" {
			filtered = append(filtered, uri)
		}
	}
	return filtered
}
