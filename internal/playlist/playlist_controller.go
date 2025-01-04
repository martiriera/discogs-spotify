package playlist

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/discogs"
	"github.com/martiriera/discogs-spotify/internal/entities"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUsername string) (*entities.Playlist, error) {
	// fetchReleases
	releases, err := c.discogsService.GetReleases(discogsUsername)

	if err != nil {
		return nil, err
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
	albums := parseAlbumsFromReleases(releases)
	uris := make([]string, len(albums))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(albums))

	for i, album := range albums {
		wg.Add(1)
		go func(i int, album entities.Album) {
			defer wg.Done()
			uri, err := c.spotifyService.GetAlbumId(ctx, album)
			if err != nil {
				errChan <- errors.Wrap(err, "error getting album id on channel")
				return
			}
			mu.Lock()
			uris[i] = uri
			mu.Unlock()
		}(i, album)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
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

// TODO: Necessary?
func joinArtists(artists []entities.DiscogsArtist) string {
	names := []string{}
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return strings.Join(names, ", ")
}
