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
				albums, err := c.spotifyService.SearchAlbum(ctx, album)
				if albums == nil {
					return
				}
				if err != nil {
					errChan <- errors.Wrap(err, "error getting album id")
					return
				}
				uri := getMatchingAlbumURI(album, albums)
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
	// David Bowie Hunky Dory = A Pedir De Boca = A Pedir De Boca
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

// compares album name and artist from Discogs with Spotify to discard unrelated albums
func getMatchingAlbumURI(album entities.Album, spotifyAlbums []entities.SpotifyAlbumItem) string {
	inputArtist := normalizeName(album.Artist)
	inputAlbumName := normalizeName(album.Title)

	// Case 1: Exact artist and album name match
	for _, spotifyAlbum := range spotifyAlbums {
		normalizedSpotifyAlbumName := normalizeName(spotifyAlbum.Name)

		for _, artist := range spotifyAlbum.Artists {
			normalizedSpotifyArtist := normalizeName(artist.Name)

			if strings.EqualFold(normalizedSpotifyArtist, inputArtist) &&
				strings.EqualFold(normalizedSpotifyAlbumName, inputAlbumName) {
				return spotifyAlbum.ID
			}
		}
	}

	// Case 2: Artist match and at least one word album name match
	for _, spotifyAlbum := range spotifyAlbums {
		normalizedSpotifyAlbumName := normalizeName(spotifyAlbum.Name)

		for _, artist := range spotifyAlbum.Artists {
			normalizedSpotifyArtist := normalizeName(artist.Name)

			if strings.EqualFold(normalizedSpotifyArtist, inputArtist) {
				inputWords := strings.Fields(inputAlbumName)
				spotifyWords := strings.Fields(normalizedSpotifyAlbumName)

				for _, inputWord := range inputWords {
					for _, spotifyWord := range spotifyWords {
						if strings.EqualFold(inputWord, spotifyWord) && len(inputWord) > 1 {
							return spotifyAlbum.ID
						}
					}
				}
			}
		}
	}

	// Case 3: All other cases are considered not OK
	return ""
}

// normalizeName removes common suffixes and normalizes the artist name
func normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.Split(name, "(")[0]
	name = strings.TrimSpace(name)
	return name
}
