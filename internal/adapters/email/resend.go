package email

import (
	"fmt"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"

	"github.com/martiriera/discogs-spotify/internal/usecases"
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

// htmlEscape escapes the minimal set of characters needed for safe HTML content.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
