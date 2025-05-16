package usecases

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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

func (c *DiscogsConvertToSpotify) getSpotifyAlbumIDs(ctx context.Context, releases []entities.DiscogsRelease) ([]string, error) {
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
				uri, err := c.spotifyService.SearchAlbum(ctx, album)
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
	// TODO: Move this logic to domain
	artistName := release.BasicInformation.Artists[0].Name
	artistName = strings.TrimSpace(strings.Split(artistName, " (")[0])

	titleName := release.BasicInformation.Title
	titleName = strings.TrimSpace(strings.Split(titleName, " (")[0])

	album := entities.Album{
		Artist: artistName,
		Title:  titleName,
		Year:   release.BasicInformation.Year,
	}
	return album
}
