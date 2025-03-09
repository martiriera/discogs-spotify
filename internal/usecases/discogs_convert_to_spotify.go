package usecases

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type DiscogsConvertToSpotify struct {
	spotifyService ports.SpotifyPort
}

func NewDiscogsConvertToSpotify(s ports.SpotifyPort) *DiscogsConvertToSpotify {
	return &DiscogsConvertToSpotify{spotifyService: s}
}

func (c *DiscogsConvertToSpotify) getSpotifyAlbumIDs(ctx *gin.Context, releases []entities.DiscogsRelease) ([]string, error) {
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

func getAlbumFromRelease(release entities.DiscogsRelease) entities.Album {
	album := entities.Album{
		Artist: release.BasicInformation.Artists[0].Name,
		Title:  strings.TrimSpace(release.BasicInformation.Title),
	}
	return album
}
