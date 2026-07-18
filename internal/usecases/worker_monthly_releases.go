package usecases

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

// WorkerMonthlyReleases fetches a Discogs collection, searches each release on
// Spotify, and filters albums whose release month matches the target month.
type WorkerMonthlyReleases struct {
	discogsService ports.DiscogsPort
	spotifyService ports.SpotifyPort
}

func NewWorkerMonthlyReleases(d ports.DiscogsPort, s ports.SpotifyPort) *WorkerMonthlyReleases {
	return &WorkerMonthlyReleases{
		discogsService: d,
		spotifyService: s,
	}
}

// MonthlyRelease holds the data needed for the email digest.
type MonthlyRelease struct {
	Artist      string
	Title       string
	ReleaseDate string
	SpotifyURL  string
}

// Run fetches the collection at discogsURL and returns albums released in targetMonth.
func (w *WorkerMonthlyReleases) Run(ctx context.Context, discogsURL string, targetMonth time.Month) ([]MonthlyRelease, error) {
	parsed, err := parseDiscogsURL(discogsURL)
	if err != nil {
		return nil, fmt.Errorf("parsing discogs URL: %w", err)
	}

	processor := NewDiscogsProcessURL(w.discogsService)
	releases, err := processor.processDiscogsURL(ctx, parsed)
	if err != nil {
		return nil, fmt.Errorf("fetching discogs releases: %w", err)
	}

	albumItems, err := w.searchSpotifyAlbums(ctx, releases)
	if err != nil {
		return nil, fmt.Errorf("searching spotify albums: %w", err)
	}

	return filterByMonth(albumItems, targetMonth), nil
}

func (w *WorkerMonthlyReleases) searchSpotifyAlbums(
	ctx context.Context,
	releases []entities.DiscogsRelease,
) ([]entities.SpotifyAlbumItem, error) {
	type result struct {
		item entities.SpotifyAlbumItem
		err  error
	}

	resultsChan := make(chan result, len(releases))
	var wg sync.WaitGroup
	rateLimiter := time.Tick(spotifyAPIRateLimit)

	for _, release := range releases {
		album := getAlbumFromRelease(&release)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-rateLimiter:
			wg.Add(1)
			go func(album entities.Album) {
				defer wg.Done()
				items, err := w.spotifyService.SearchAlbum(ctx, album)
				if err != nil {
					resultsChan <- result{err: errors.Wrap(err, "searching album")}
					return
				}
				matched := getMatchingAlbumItem(album, items)
				if matched != nil {
					resultsChan <- result{item: *matched}
				}
			}(album)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var items []entities.SpotifyAlbumItem
	var errs []error
	for r := range resultsChan {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		if r.item.ID != "" {
			items = append(items, r.item)
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("encountered errors during spotify search: %v", errs)
	}

	return items, nil
}

// getMatchingAlbumItem returns the full SpotifyAlbumItem (instead of just the ID)
// using the same two-pass fuzzy matching as getMatchingAlbumURI.
func getMatchingAlbumItem(album entities.Album, spotifyAlbums []entities.SpotifyAlbumItem) *entities.SpotifyAlbumItem {
	matchedID := getMatchingAlbumURI(album, spotifyAlbums)
	if matchedID == "" {
		return nil
	}
	for i := range spotifyAlbums {
		if spotifyAlbums[i].ID == matchedID {
			return &spotifyAlbums[i]
		}
	}
	return nil
}

// filterByMonth keeps albums whose Spotify release_date month matches targetMonth.
// Albums with year-only precision are skipped.
func filterByMonth(items []entities.SpotifyAlbumItem, targetMonth time.Month) []MonthlyRelease {
	seen := make(map[string]struct{})
	var result []MonthlyRelease

	for i := range items {
		item := &items[i]
		if _, ok := seen[item.ID]; ok {
			continue
		}

		releaseMonth, ok := parseReleaseMonth(item.ReleaseDate, item.ReleaseDatePrecision)
		if !ok {
			continue
		}

		if releaseMonth == targetMonth {
			seen[item.ID] = struct{}{}
			artist := ""
			if len(item.Artists) > 0 {
				artist = item.Artists[0].Name
			}
			result = append(result, MonthlyRelease{
				Artist:      artist,
				Title:       item.Name,
				ReleaseDate: item.ReleaseDate,
				SpotifyURL:  item.ExternalURLs.Spotify,
			})
		}
	}

	return result
}

// parseReleaseMonth extracts the month from a Spotify release_date string.
// Returns (month, true) for "day" and "month" precision; (0, false) otherwise.
func parseReleaseMonth(releaseDate, precision string) (time.Month, bool) {
	switch precision {
	case "day":
		t, err := time.Parse("2006-01-02", releaseDate)
		if err != nil {
			return 0, false
		}
		return t.Month(), true
	case "month":
		t, err := time.Parse("2006-01", releaseDate)
		if err != nil {
			return 0, false
		}
		return t.Month(), true
	default:
		return 0, false
	}
}
