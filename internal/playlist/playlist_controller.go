package playlist

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUsername string) (*entities.Playlist, error) {
	stop := util.StartTimer("CreatePlaylist")
	defer stop()

	// fetchReleases
	releases, err := c.discogsService.GetReleases(discogsUsername)

	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("no releases found on Discogs list")
	}

	// processAlbumIds
	albumIds, err := c.getSpotifyAlbumIds(ctx, releases)
	if err != nil {
		return nil, errors.Wrap(err, "error getting spotify album uris")
	}
	albumIds = c.filterValidUnique(albumIds)

	// createPlaylist
	playlistBuilder := NewPlaylistBuilder(c.spotifyService)
	err = playlistBuilder.AddAlbums(ctx, albumIds)
	if err != nil {
		return nil, errors.Wrap(err, "error adding albums to playlist builder")
	}
	playlist, err := playlistBuilder.CreateAndPopulate(ctx, "Discogs Playlist", "Playlist created from Discogs")
	if err != nil {
		return nil, errors.Wrap(err, "error creating and populating playlist")
	}

	return &entities.Playlist{
		DiscogsReleases: len(releases),
		SpotifyAlbums:   len(albumIds),
		SpotifyPlaylist: *playlist,
	}, nil
}

func (c *PlaylistController) getSpotifyAlbumIds(ctx *gin.Context, releases []entities.DiscogsRelease) ([]string, error) {
	urisChan := make(chan string, len(releases))
	errChan := make(chan error, len(releases))

	var wg sync.WaitGroup
	rateLimiter := time.Tick(200 * time.Millisecond)

	for _, release := range releases {
		album := getAlbumFromRelease(release)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-rateLimiter:
			wg.Add(1)
			go func(album entities.Album) {
				defer wg.Done()
				uri, err := c.spotifyService.GetAlbumId(ctx, album)
				if err != nil {
					errChan <- errors.Wrap(err, "error getting album id")
					return
				}
				urisChan <- uri
			}(album)
		}
	}

	go func() {
		wg.Wait()
		close(urisChan)
		close(errChan)
	}()

	var uris []string
	for uri := range urisChan {
		uris = append(uris, uri)
	}

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("encountered errors: %v", errs)
	}

	return uris, nil
}

func (c *PlaylistController) filterValidUnique(uris []string) []string {
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

func parseDiscogsUrl(url string) string {
	re := regexp.MustCompile(`^(?:https://)?www\.discogs\.com/[^/]+/(?:user/(.+)/collection|lists/.+/(\d+)|wantlist\?user=(.+))$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 0 {
		for _, match := range matches[1:] {
			if match != "" {
				return match
			}
		}
	}
	return ""
}
