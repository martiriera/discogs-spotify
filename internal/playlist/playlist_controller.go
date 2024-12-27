package playlist

import (
	"log"
	"strings"
	"sync"

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
	filteredUris := c.filterNotFounds(uris)
	log.Println("URIs: ", filteredUris)

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

	err = c.spotifyService.AddToPlaylist(ctx, playlistId, filteredUris)
	if err != nil {
		return "", errors.Wrap(err, "error adding to playlist")
	}

	return playlistId, nil
}

func (c *PlaylistController) getSpotifyAlbumUris(ctx *gin.Context, releases []entities.DiscogsRelease) ([]string, error) {
	albums := parseAlbumsFromReleases(releases)
	uris := make([]string, len(albums))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(albums))

	for i, album := range albums {
		wg.Add(1)
		go func(i int, album entities.Album) {
			defer wg.Done()
			uri, err := c.spotifyService.GetAlbumUri(ctx, album)
			if err != nil {
				errChan <- errors.Wrap(err, "error getting album uri")
				return
			}
			mu.Lock()
			uris[i] = uri
			mu.Unlock()
		}(i, album)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return uris, nil
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
