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

func (c *DiscogsConvertToSpotify) getSpotifyAlbums(ctx context.Context, releases []entities.DiscogsRelease) ([]entities.Album, error) {
	albumsChan := make(chan entities.Album, len(releases))
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
				albumsChan <- uri
			}(album)
		}
	}

	go func() {
		wg.Wait()
		close(albumsChan)
		close(errChan)
	}()

	var albums []entities.Album
	for uri := range albumsChan {
		albums = append(albums, uri)
	}

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("encountered errors: %v", errs)
	}

	return albums, nil
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
	}
	return album
}
