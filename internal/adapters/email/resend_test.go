package email

import (
	"strings"
	"testing"
	"time"

	"github.com/martiriera/discogs-spotify/internal/usecases"
)

func TestBuildICS_ContainsRequiredCalendarFields(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Radiohead", Title: "OK Computer", ReleaseDate: "1997-06-16", SpotifyURL: "https://open.spotify.com/album/123"},
	}

	ics := buildICS(time.June, 1997, releases)

	requiredFields := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:",
		"BEGIN:VEVENT",
		"END:VEVENT",
		"END:VCALENDAR",
		"DTSTART:",
		"DTEND:",
		"SUMMARY:Radiohead - OK Computer",
		"BEGIN:VALARM",
		"TRIGGER:-PT5M",
		"END:VALARM",
	}

	for _, field := range requiredFields {
		if !strings.Contains(ics, field) {
			t.Errorf("expected ICS to contain %q, but it did not\nFull output:\n%s", field, ics)
		}
	}
}

func TestBuildICS_EventTimesCorrect(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Artist", Title: "Album", ReleaseDate: "2024-08-15", SpotifyURL: "https://open.spotify.com/album/abc"},
	}

	ics := buildICS(time.August, 2024, releases)

	if !strings.Contains(ics, "DTSTART:20240815T180000") {
		t.Errorf("expected DTSTART at 18:00 on release date, got:\n%s", ics)
	}
	if !strings.Contains(ics, "DTEND:20240815T190000") {
		t.Errorf("expected DTEND at 19:00 on release date, got:\n%s", ics)
	}
}

func TestBuildICS_MultipleReleases(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Artist A", Title: "Album A", ReleaseDate: "2024-08-01", SpotifyURL: "https://open.spotify.com/album/a"},
		{Artist: "Artist B", Title: "Album B", ReleaseDate: "2024-08-10", SpotifyURL: "https://open.spotify.com/album/b"},
		{Artist: "Artist C", Title: "Album C", ReleaseDate: "2024-08-20", SpotifyURL: "https://open.spotify.com/album/c"},
	}

	ics := buildICS(time.August, 2024, releases)

	count := strings.Count(ics, "BEGIN:VEVENT")
	if count != 3 {
		t.Errorf("expected 3 VEVENT blocks, got %d", count)
	}
}

func TestBuildICS_MonthPrecisionDate(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Artist", Title: "Album", ReleaseDate: "2024-08", SpotifyURL: "https://open.spotify.com/album/abc"},
	}

	ics := buildICS(time.August, 2024, releases)

	if !strings.Contains(ics, "BEGIN:VEVENT") {
		t.Error("expected a VEVENT for month-precision release date, got none")
	}
	if !strings.Contains(ics, "DTSTART:20240801T180000") {
		t.Errorf("expected DTSTART at 18:00 on first of month, got:\n%s", ics)
	}
}

func TestBuildICS_SkipsInvalidDate(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Artist", Title: "Album", ReleaseDate: "1997", SpotifyURL: "https://open.spotify.com/album/abc"},
	}

	ics := buildICS(time.June, 1997, releases)

	if strings.Contains(ics, "BEGIN:VEVENT") {
		t.Error("expected no VEVENT for year-only release date, but found one")
	}
}

func TestBuildICS_ICSEscapesSpecialCharacters(t *testing.T) {
	releases := []usecases.MonthlyRelease{
		{Artist: "Artist, Jr.", Title: "Album; Deluxe", ReleaseDate: "2024-08-01", SpotifyURL: "https://open.spotify.com/album/abc"},
	}

	ics := buildICS(time.August, 2024, releases)

	if !strings.Contains(ics, `Artist\, Jr. - Album\; Deluxe`) {
		t.Errorf("expected escaped commas and semicolons in SUMMARY, got:\n%s", ics)
	}
}

func TestParseDateToICS_DayPrecision(t *testing.T) {
	result, ok := parseDateToICS("1997-06-16")
	if !ok {
		t.Fatal("expected ok=true for YYYY-MM-DD input")
	}
	if result != "19970616T180000" {
		t.Errorf("expected 19970616T180000, got %s", result)
	}
}

func TestParseDateToICS_MonthPrecision(t *testing.T) {
	result, ok := parseDateToICS("1997-06")
	if !ok {
		t.Fatal("expected ok=true for YYYY-MM input")
	}
	if result != "19970601T180000" {
		t.Errorf("expected 19970601T180000, got %s", result)
	}
}

func TestParseDateToICS_YearOnlyReturnsFalse(t *testing.T) {
	_, ok := parseDateToICS("1997")
	if ok {
		t.Error("expected ok=false for year-only input")
	}
}

func TestParseDateToICS_InvalidReturnsFalse(t *testing.T) {
	_, ok := parseDateToICS("not-a-date")
	if ok {
		t.Error("expected ok=false for invalid date input")
	}
}
