package playlist

import (
	"log"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUsername string) (string, error) {
	// fetchReleases
	releases, err := c.discogsService.GetReleases(discogsUsername)
	log.Println("Releases: ", len(releases))

	if err != nil {
		return "", err
	}

	// processAlbumIds
	albumIds, err := c.getSpotifyAlbumIds(ctx, releases)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify album uris")
	}
	albumIds = c.filterValidUnique(albumIds)
	log.Println("IDs: ", len(albumIds))

	// playlistBuilder
	tracks, err := c.getSpotifyTrackUris(ctx, albumIds)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify track uris")
	}
	playlist, err := c.spotifyService.CreatePlaylist(ctx, "Discogs Playlist", "Playlist created from Discogs")
	if err != nil {
		return "", errors.Wrap(err, "error creating playlist")
	}
	err = c.addToSpotifyPlaylist(ctx, playlist.ID, tracks)
	if err != nil {
		return "", errors.Wrap(err, "error adding to playlist")
	}

	return playlist.URL, nil
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

func (c *PlaylistController) addToSpotifyPlaylist(ctx *gin.Context, playlistId string, uris []string) error {
	batchSize := 100
	return batchRequests(ctx, uris, batchSize, func(ctx *gin.Context, batch []string) error {
		err := c.spotifyService.AddToPlaylist(ctx, playlistId, batch)
		if err != nil {
			return errors.Wrap(err, "error adding to playlist")
		}
		return nil
	})
}

func (c *PlaylistController) getSpotifyTrackUris(ctx *gin.Context, albums []string) ([]string, error) {
	batckSize := 20
	uris := []string{}
	err := batchRequests(ctx, albums, batckSize, func(ctx *gin.Context, batch []string) error {
		tracks, err := c.spotifyService.GetAlbumsTrackUris(ctx, batch)
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
