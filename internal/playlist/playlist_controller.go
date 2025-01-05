package playlist

import (
	"fmt"
	"net/url"
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

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context, discogsUrl string) (*entities.Playlist, error) {
	stop := util.StartTimer("CreatePlaylist")
	defer stop()

	// fetchReleases
	discogsUrlParsed, err := parseDiscogsUrl(discogsUrl)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing Discogs URL")
	}

	var releases []entities.DiscogsRelease
	if discogsUrlParsed.UrlType == entities.CollectionType {
		releases, err = c.discogsService.GetReleases(discogsUrl)
	} else if discogsUrlParsed.UrlType == entities.WantlistType {
		releases, err = c.discogsService.GetWantlistReleases(discogsUrl)
	} else if discogsUrlParsed.UrlType == entities.ListType {
		// releases, err = c.discogsService.GetListReleases(discogsUrl)
	} else {
		return nil, errors.New("unrecognized URL type")
	}

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

func parseDiscogsUrl(urlStr string) (*entities.DiscogsInputUrl, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL")
	}

	if parsedUrl.Host != "www.discogs.com" && !strings.Contains(parsedUrl.Path, "discogs.com") {
		return nil, errors.New("URL must be from discogs.com")
	}

	path := parsedUrl.Path
	query := parsedUrl.Query()

	if strings.Contains(path, "/collection") {
		user := strings.Split(path, "/")[3]
		return &entities.DiscogsInputUrl{Id: user, UrlType: entities.CollectionType}, nil
	}

	if wantlistUser := query.Get("user"); wantlistUser != "" {
		return &entities.DiscogsInputUrl{Id: wantlistUser, UrlType: entities.WantlistType}, nil
	}

	if strings.Contains(path, "/lists/") {
		listId := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
		return &entities.DiscogsInputUrl{Id: listId, UrlType: entities.ListType}, nil
	}

	return nil, errors.New("unrecognized Discogs URL")
}
