package email

import (
	"fmt"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"

	"github.com/martiriera/discogs-spotify/internal/usecases"
)

const (
	icsEventHour     = 18
	icsEventDuration = 1 // hours
	icsAlarmMinutes  = -5
)

// ResendSender sends the monthly digest email via the Resend API.
type ResendSender struct {
	client *resend.Client
	from   string
	to     string
}

func NewResendSender(apiKey, from, to string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   from,
		to:     to,
	}
}

// SendMonthlyDigest sends the monthly release digest email.
func (s *ResendSender) SendMonthlyDigest(month time.Month, year int, releases []usecases.MonthlyRelease) error {
	subject := fmt.Sprintf("Upcoming releases for %s %d", month.String(), year)
	html := buildEmailHTML(month, year, releases)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{s.to},
		Subject: subject,
		Html:    html,
	}

	if len(releases) > 0 {
		icsContent := buildICS(month, year, releases)
		filename := fmt.Sprintf("releases-%s-%d.ics", strings.ToLower(month.String()), year)
		params.Attachments = []*resend.Attachment{
			{
				Filename:    filename,
				Content:     []byte(icsContent),
				ContentType: "text/calendar; charset=utf-8",
			},
		}
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("sending email via resend: %w", err)
	}

	return nil
}

func buildEmailHTML(month time.Month, year int, releases []usecases.MonthlyRelease) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body { font-family: sans-serif; color: #222; max-width: 600px; margin: 0 auto; padding: 24px; }
  h1 { font-size: 1.4rem; margin-bottom: 4px; }
  p.subtitle { color: #666; margin-top: 0; margin-bottom: 24px; }
  table { width: 100%; border-collapse: collapse; }
  th { text-align: left; border-bottom: 2px solid #ddd; padding: 8px 4px; font-size: 0.85rem; color: #555; }
  td { padding: 10px 4px; border-bottom: 1px solid #eee; vertical-align: top; }
  a { color: #1DB954; text-decoration: none; }
  a:hover { text-decoration: underline; }
  .empty { color: #888; font-style: italic; }
</style>
</head>
<body>
`)

	fmt.Fprintf(&sb, "<h1>Releases for %s</h1>\n", month.String())
	fmt.Fprintf(&sb, "<p class=\"subtitle\">Albums in your Discogs collection released in %s (any year)</p>\n", month.String())

	_ = year // year is in the subject; not repeated in body to keep it timeless

	if len(releases) == 0 {
		sb.WriteString("<p class=\"empty\">No releases found for this month in your collection.</p>\n")
	} else {
		sb.WriteString("<table>\n<thead>\n<tr><th>Artist</th><th>Album</th><th>Released</th></tr>\n</thead>\n<tbody>\n")
		for _, r := range releases {
			fmt.Fprintf(&sb,
				"<tr><td>%s</td><td><a href=\"%s\">%s</a></td><td>%s</td></tr>\n",
				htmlEscape(r.Artist),
				r.SpotifyURL,
				htmlEscape(r.Title),
				htmlEscape(r.ReleaseDate),
			)
		}
		sb.WriteString("</tbody>\n</table>\n")
	}

	sb.WriteString("</body>\n</html>")
	return sb.String()
}

// buildICS generates an iCalendar (.ics) file containing one VEVENT per release.
// Each event starts at 18:00 and lasts one hour, with a 5-minute VALARM.
// The release date is used as-is from Discogs (format "YYYY-MM-DD" or "YYYY-MM").
func buildICS(month time.Month, year int, releases []usecases.MonthlyRelease) string {
	now := time.Now().UTC().Format("20060102T150405Z")
	var sb strings.Builder

	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//discogs-spotify//Monthly Digest//EN\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")

	for i, r := range releases {
		dtstart, ok := parseDateToICS(r.ReleaseDate)
		if !ok {
			continue
		}
		dtend := shiftICSHour(dtstart, icsEventDuration)
		uid := fmt.Sprintf("release-%d-%s-%d-%d@discogs-spotify", i, strings.ToLower(month.String()), year, i)
		summary := icsEscape(fmt.Sprintf("%s - %s", r.Artist, r.Title))

		sb.WriteString("BEGIN:VEVENT\r\n")
		fmt.Fprintf(&sb, "UID:%s\r\n", uid)
		fmt.Fprintf(&sb, "DTSTAMP:%s\r\n", now)
		fmt.Fprintf(&sb, "DTSTART:%s\r\n", dtstart)
		fmt.Fprintf(&sb, "DTEND:%s\r\n", dtend)
		fmt.Fprintf(&sb, "SUMMARY:%s\r\n", summary)
		sb.WriteString("BEGIN:VALARM\r\n")
		sb.WriteString("ACTION:DISPLAY\r\n")
		fmt.Fprintf(&sb, "TRIGGER:-PT%dM\r\n", -icsAlarmMinutes)
		sb.WriteString("DESCRIPTION:Reminder\r\n")
		sb.WriteString("END:VALARM\r\n")
		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")
	return sb.String()
}

// parseDateToICS converts a Discogs release date string to an iCalendar
// DTSTART value at 18:00 local time (no timezone).
// Supports "YYYY-MM-DD" and "YYYY-MM" formats.
func parseDateToICS(releaseDate string) (string, bool) {
	var t time.Time
	var err error

	switch len(releaseDate) {
	case 10: // YYYY-MM-DD
		t, err = time.Parse("2006-01-02", releaseDate)
	case 7: // YYYY-MM
		t, err = time.Parse("2006-01", releaseDate)
	default:
		return "", false
	}

	if err != nil {
		return "", false
	}

	result := time.Date(t.Year(), t.Month(), t.Day(), icsEventHour, 0, 0, 0, time.UTC)
	return result.Format("20060102T150405"), true
}

// shiftICSHour adds hours to an iCalendar datetime string (format: 20060102T150405).
func shiftICSHour(icsDatetime string, hours int) string {
	t, err := time.Parse("20060102T150405", icsDatetime)
	if err != nil {
		return icsDatetime
	}
	return t.Add(time.Duration(hours) * time.Hour).Format("20060102T150405")
}

// icsEscape escapes special characters for iCalendar text values.
func icsEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, ";", `\;`)
	s = strings.ReplaceAll(s, ",", `\,`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

// htmlEscape escapes the minimal set of characters needed for safe HTML content.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
