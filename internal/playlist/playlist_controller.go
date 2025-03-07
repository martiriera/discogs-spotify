package playlist

import (
	"fmt"
	"net/url"
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUrl string) (*entities.Playlist, error) {
	stop := util.StartTimer("CreatePlaylist")
	defer stop()

	// fetch releases
	parsedDiscogsUrl, err := parseDiscogsUrl(discogsUrl)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing Discogs URL")
	}

	var releases []entities.DiscogsRelease
	if parsedDiscogsUrl.Type == entities.CollectionType {
		releases, err = c.discogsService.GetCollectionReleases(parsedDiscogsUrl.ID)
	} else if parsedDiscogsUrl.Type == entities.WantlistType {
		releases, err = c.discogsService.GetWantlistReleases(parsedDiscogsUrl.ID)
	} else if parsedDiscogsUrl.Type == entities.ListType {
		releases, err = c.discogsService.GetListReleases(parsedDiscogsUrl.ID)
	} else {
		return nil, errors.New("unrecognized URL type")
	}

	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("no releases found on Discogs list")
	}

	// process album IDs
	albumIDs, err := c.getSpotifyAlbumIDs(ctx, releases)
	if err != nil {
		return nil, errors.Wrap(err, "error getting spotify album uris")
	}
	albumIDs = c.filterValidUnique(albumIDs)

	// create playlist
	playlistBuilder := NewPlaylistBuilder(c.spotifyService)
	err = playlistBuilder.AddAlbums(ctx, albumIDs)
	if err != nil {
		return nil, errors.Wrap(err, "error adding albums to playlist builder")
	}
	playlist, err := playlistBuilder.CreateAndPopulate(
		ctx,
		"Discogs "+cases.Title(language.English).String(parsedDiscogsUrl.Type.String())+" by "+parsedDiscogsUrl.ID,
		"Created from: "+discogsUrl,
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

func (c *PlaylistController) getSpotifyAlbumIDs(ctx *gin.Context, releases []entities.DiscogsRelease) ([]string, error) {
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
				uri, err := c.spotifyService.GetAlbumID(ctx, album)
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

var ErrInvalidDiscogsUrl = errors.New("invalid Discogs URL")

func parseDiscogsUrl(urlStr string) (*entities.DiscogsInputUrl, error) {
	// validate host
	parsedUrl, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	pathWithQuery := parsedUrl.Path
	if parsedUrl.RawQuery != "" && strings.Contains(parsedUrl.RawQuery, "user=") {
		pathWithQuery += "?" + parsedUrl.RawQuery
	}

	matchingUrl := ""
	if parsedUrl.Host == "www.discogs.com" {
		matchingUrl = pathWithQuery
	} else if parsedUrl.Host == "" {
		matchingUrl = "/" + strings.SplitN(pathWithQuery, "/", 2)[1]
	} else {
		return nil, ErrInvalidDiscogsUrl
	}

	// validate path
	re := regexp.MustCompile(`^/(?:[a-z]{2}/)?(?:user/(.+)/collection|wantlist\?user=(.+))$|lists/.+/(\d+)`)
	matches := re.FindStringSubmatch(matchingUrl)
	if matches == nil {
		return nil, ErrInvalidDiscogsUrl
	}

	for i, match := range matches {
		// https://www.discogs.com/es/user/digger/collection
		if i == 1 && match != "" {
			return &entities.DiscogsInputUrl{ID: match, Type: entities.CollectionType}, nil
		}
		// https://www.discogs.com/es/wantlist?user=digger
		if i == 2 && match != "" {
			return &entities.DiscogsInputUrl{ID: match, Type: entities.WantlistType}, nil
		}

		// https://www.discogs.com/es/lists/MyList/1545836
		if i == 3 && match != "" {
			return &entities.DiscogsInputUrl{ID: match, Type: entities.ListType}, nil
		}
	}
	return nil, ErrInvalidDiscogsUrl
}
