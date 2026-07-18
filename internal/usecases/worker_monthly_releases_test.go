package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/entities"
)

// errorSpotifyMock returns an error on SearchAlbum calls.
type errorSpotifyMock struct{}

func (e *errorSpotifyMock) SearchAlbum(_ context.Context, _ entities.Album) ([]entities.SpotifyAlbumItem, error) {
	return nil, errors.New("spotify unavailable")
}
func (e *errorSpotifyMock) GetUserID(_ context.Context) (string, error)                    { return "", nil }
func (e *errorSpotifyMock) CreatePlaylist(_ context.Context, _, _ string) (entities.SpotifyPlaylist, error) {
	return entities.SpotifyPlaylist{}, nil
}
func (e *errorSpotifyMock) AddToPlaylist(_ context.Context, _ string, _ []string) error { return nil }
func (e *errorSpotifyMock) GetAlbumsTrackUris(_ context.Context, _ []string) ([]string, error) {
	return nil, nil
}

// --- parseReleaseMonth ---

func TestParseReleaseMonth(t *testing.T) {
	tests := []struct {
		name        string
		releaseDate string
		precision   string
		wantMonth   time.Month
		wantOK      bool
	}{
		{
			name:        "day precision - valid date",
			releaseDate: "2024-08-15",
			precision:   "day",
			wantMonth:   time.August,
			wantOK:      true,
		},
		{
			name:        "month precision - valid date",
			releaseDate: "2024-03",
			precision:   "month",
			wantMonth:   time.March,
			wantOK:      true,
		},
		{
			name:        "year precision - skipped",
			releaseDate: "2024",
			precision:   "year",
			wantMonth:   0,
			wantOK:      false,
		},
		{
			name:        "unknown precision - skipped",
			releaseDate: "2024-08-15",
			precision:   "unknown",
			wantMonth:   0,
			wantOK:      false,
		},
		{
			name:        "day precision - malformed date",
			releaseDate: "not-a-date",
			precision:   "day",
			wantMonth:   0,
			wantOK:      false,
		},
		{
			name:        "month precision - malformed date",
			releaseDate: "2024/08",
			precision:   "month",
			wantMonth:   0,
			wantOK:      false,
		},
		{
			name:        "day precision - January",
			releaseDate: "2023-01-01",
			precision:   "day",
			wantMonth:   time.January,
			wantOK:      true,
		},
		{
			name:        "month precision - December",
			releaseDate: "2023-12",
			precision:   "month",
			wantMonth:   time.December,
			wantOK:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMonth, gotOK := parseReleaseMonth(tc.releaseDate, tc.precision)
			if gotOK != tc.wantOK {
				t.Errorf("parseReleaseMonth(%q, %q) ok = %v, want %v", tc.releaseDate, tc.precision, gotOK, tc.wantOK)
			}
			if gotMonth != tc.wantMonth {
				t.Errorf("parseReleaseMonth(%q, %q) month = %v, want %v", tc.releaseDate, tc.precision, gotMonth, tc.wantMonth)
			}
		})
	}
}

// --- filterByMonth ---

func TestFilterByMonth(t *testing.T) {
	baseArtist := entities.SpotifyAlbumArtist{Name: "Test Artist"}

	tests := []struct {
		name        string
		items       []entities.SpotifyAlbumItem
		targetMonth time.Month
		wantLen     int
		wantTitles  []string
	}{
		{
			name: "filters matching month - day precision",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "id1",
					Name:                 "August Album",
					ReleaseDate:          "2024-08-10",
					ReleaseDatePrecision: "day",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
					ExternalURLs:         entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/id1"},
				},
				{
					ID:                   "id2",
					Name:                 "March Album",
					ReleaseDate:          "2024-03-01",
					ReleaseDatePrecision: "day",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
					ExternalURLs:         entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/id2"},
				},
			},
			targetMonth: time.August,
			wantLen:     1,
			wantTitles:  []string{"August Album"},
		},
		{
			name: "filters matching month - month precision",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "id3",
					Name:                 "May Album",
					ReleaseDate:          "2024-05",
					ReleaseDatePrecision: "month",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
				},
				{
					ID:                   "id4",
					Name:                 "June Album",
					ReleaseDate:          "2024-06",
					ReleaseDatePrecision: "month",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
				},
			},
			targetMonth: time.May,
			wantLen:     1,
			wantTitles:  []string{"May Album"},
		},
		{
			name: "skips year-only precision",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "id5",
					Name:                 "Year Album",
					ReleaseDate:          "2024",
					ReleaseDatePrecision: "year",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
				},
			},
			targetMonth: time.January,
			wantLen:     0,
			wantTitles:  []string{},
		},
		{
			name: "deduplicates by Spotify ID",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "dup1",
					Name:                 "Duplicate Album",
					ReleaseDate:          "2024-07-01",
					ReleaseDatePrecision: "day",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
				},
				{
					ID:                   "dup1",
					Name:                 "Duplicate Album",
					ReleaseDate:          "2024-07-01",
					ReleaseDatePrecision: "day",
					Artists:              []entities.SpotifyAlbumArtist{baseArtist},
				},
			},
			targetMonth: time.July,
			wantLen:     1,
			wantTitles:  []string{"Duplicate Album"},
		},
		{
			name:        "empty input returns nil",
			items:       []entities.SpotifyAlbumItem{},
			targetMonth: time.January,
			wantLen:     0,
			wantTitles:  []string{},
		},
		{
			name: "uses first artist name",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "id6",
					Name:                 "Multi Artist Album",
					ReleaseDate:          "2024-11-15",
					ReleaseDatePrecision: "day",
					Artists: []entities.SpotifyAlbumArtist{
						{Name: "First Artist"},
						{Name: "Second Artist"},
					},
					ExternalURLs: entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/id6"},
				},
			},
			targetMonth: time.November,
			wantLen:     1,
			wantTitles:  []string{"Multi Artist Album"},
		},
		{
			name: "empty artist list uses empty string",
			items: []entities.SpotifyAlbumItem{
				{
					ID:                   "id7",
					Name:                 "No Artist Album",
					ReleaseDate:          "2024-02",
					ReleaseDatePrecision: "month",
					Artists:              []entities.SpotifyAlbumArtist{},
				},
			},
			targetMonth: time.February,
			wantLen:     1,
			wantTitles:  []string{"No Artist Album"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterByMonth(tc.items, tc.targetMonth)
			if len(got) != tc.wantLen {
				t.Errorf("filterByMonth() returned %d items, want %d", len(got), tc.wantLen)
			}
			for i, title := range tc.wantTitles {
				if i >= len(got) {
					break
				}
				if got[i].Title != title {
					t.Errorf("filterByMonth()[%d].Title = %q, want %q", i, got[i].Title, title)
				}
			}
		})
	}
}

func TestFilterByMonth_PopulatesFields(t *testing.T) {
	items := []entities.SpotifyAlbumItem{
		{
			ID:          "id1",
			Name:        "Test Album",
			ReleaseDate: "2024-04-20",
			ReleaseDatePrecision: "day",
			Artists: []entities.SpotifyAlbumArtist{
				{Name: "Test Artist"},
			},
			ExternalURLs: entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/id1"},
		},
	}

	result := filterByMonth(items, time.April)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	r := result[0]
	if r.Artist != "Test Artist" {
		t.Errorf("Artist = %q, want %q", r.Artist, "Test Artist")
	}
	if r.Title != "Test Album" {
		t.Errorf("Title = %q, want %q", r.Title, "Test Album")
	}
	if r.ReleaseDate != "2024-04-20" {
		t.Errorf("ReleaseDate = %q, want %q", r.ReleaseDate, "2024-04-20")
	}
	if r.SpotifyURL != "https://open.spotify.com/album/id1" {
		t.Errorf("SpotifyURL = %q, want %q", r.SpotifyURL, "https://open.spotify.com/album/id1")
	}
}

// --- getMatchingAlbumItem ---

func TestGetMatchingAlbumItem(t *testing.T) {
	albums := entities.MotherSpotifyAlbums()

	tests := []struct {
		name       string
		album      entities.Album
		candidates []entities.SpotifyAlbumItem
		wantID     string
		wantNil    bool
	}{
		{
			name:       "exact match returns item",
			album:      entities.Album{Artist: "Descendents", Title: "Milo Goes to College"},
			candidates: albums,
			wantID:     entities.SpotifyAlbumIDMiloGoesToCollege,
		},
		{
			name:       "no match returns nil",
			album:      entities.Album{Artist: "Unknown Artist", Title: "Unknown Album"},
			candidates: albums,
			wantNil:    true,
		},
		{
			name:       "empty candidates returns nil",
			album:      entities.Album{Artist: "Descendents", Title: "Milo Goes to College"},
			candidates: []entities.SpotifyAlbumItem{},
			wantNil:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getMatchingAlbumItem(tc.album, tc.candidates)
			if tc.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected non-nil result")
			}
			if got.ID != tc.wantID {
				t.Errorf("got ID %q, want %q", got.ID, tc.wantID)
			}
		})
	}
}

// --- WorkerMonthlyReleases.Run ---

func TestWorkerMonthlyReleasesRun(t *testing.T) {
	t.Run("returns filtered releases for target month", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherTwoDiscogsAlbums(),
		}
		// Return albums with release dates in August and March
		spotifyServiceMock := &spotify.ServiceMock{
			SearchAlbumResponses: [][]entities.SpotifyAlbumItem{
				{
					{
						ID:                   entities.SpotifyAlbumIDMiloGoesToCollege,
						Name:                 "Milo Goes to College",
						ReleaseDate:          "1982-08-01",
						ReleaseDatePrecision: "day",
						Artists: []entities.SpotifyAlbumArtist{
							{Name: "Descendents"},
						},
						ExternalURLs: entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/" + entities.SpotifyAlbumIDMiloGoesToCollege},
					},
				},
				{
					{
						ID:                   entities.SpotifyAlbumIDCatholicBoy,
						Name:                 "Catholic Boy",
						ReleaseDate:          "1980-03-01",
						ReleaseDatePrecision: "day",
						Artists: []entities.SpotifyAlbumArtist{
							{Name: "The Jim Carroll Band"},
						},
						ExternalURLs: entities.SpotifyExternalURLs{Spotify: "https://open.spotify.com/album/" + entities.SpotifyAlbumIDCatholicBoy},
					},
				},
			},
		}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		results, err := worker.Run(context.Background(), "https://www.discogs.com/user/digger/collection", time.August)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("got %d results, want 1", len(results))
		}
		if results[0].Title != "Milo Goes to College" {
			t.Errorf("Title = %q, want %q", results[0].Title, "Milo Goes to College")
		}
		if results[0].Artist != "Descendents" {
			t.Errorf("Artist = %q, want %q", results[0].Artist, "Descendents")
		}
	})

	t.Run("returns empty slice when no releases match target month", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherTwoDiscogsAlbums(),
		}
		spotifyServiceMock := &spotify.ServiceMock{
			SearchAlbumResponses: [][]entities.SpotifyAlbumItem{
				{
					{
						ID:                   entities.SpotifyAlbumIDMiloGoesToCollege,
						Name:                 "Milo Goes to College",
						ReleaseDate:          "1982-08-01",
						ReleaseDatePrecision: "day",
						Artists:              []entities.SpotifyAlbumArtist{{Name: "Descendents"}},
					},
				},
			},
		}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		results, err := worker.Run(context.Background(), "https://www.discogs.com/user/digger/collection", time.January)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("got %d results, want 0", len(results))
		}
	})

	t.Run("returns error on invalid discogs URL", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{}
		spotifyServiceMock := &spotify.ServiceMock{}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		_, err := worker.Run(context.Background(), "https://www.notdiscogs.com/something", time.August)
		if err == nil {
			t.Fatal("expected error for invalid URL, got nil")
		}
	})

	t.Run("returns error when discogs fetch fails", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Error: errors.New("discogs service unavailable"),
		}
		spotifyServiceMock := &spotify.ServiceMock{}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		_, err := worker.Run(context.Background(), "https://www.discogs.com/user/digger/collection", time.August)
		if err == nil {
			t.Fatal("expected error when discogs fetch fails, got nil")
		}
	})

	t.Run("returns error when spotify search fails", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherTwoDiscogsAlbums(),
		}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, &errorSpotifyMock{})
		_, err := worker.Run(context.Background(), "https://www.discogs.com/user/digger/collection", time.August)
		if err == nil {
			t.Fatal("expected error when spotify search fails, got nil")
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherNAlbums(10),
		}
		spotifyServiceMock := &spotify.ServiceMock{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		_, err := worker.Run(ctx, "https://www.discogs.com/user/digger/collection", time.August)
		if err == nil {
			t.Fatal("expected error for cancelled context, got nil")
		}
	})

	t.Run("works with wantlist URL", func(t *testing.T) {
		discogsServiceMock := &discogs.ServiceMock{
			Response: entities.MotherTwoDiscogsAlbums(),
		}
		spotifyServiceMock := &spotify.ServiceMock{
			SearchAlbumResponses: [][]entities.SpotifyAlbumItem{
				{
					{
						ID:                   entities.SpotifyAlbumIDMiloGoesToCollege,
						Name:                 "Milo Goes to College",
						ReleaseDate:          "1982-08",
						ReleaseDatePrecision: "month",
						Artists:              []entities.SpotifyAlbumArtist{{Name: "Descendents"}},
					},
				},
			},
		}

		worker := NewWorkerMonthlyReleases(discogsServiceMock, spotifyServiceMock)
		results, err := worker.Run(context.Background(), "https://www.discogs.com/wantlist?user=digger", time.August)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("got %d results, want 1", len(results))
		}
	})
}
