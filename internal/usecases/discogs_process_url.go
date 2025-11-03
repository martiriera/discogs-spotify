package usecases

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

var ErrInvalidDiscogsURL = errors.New("invalid Discogs URL")

type DiscogsProcessURL struct {
	discogsService ports.DiscogsPort
}

func NewDiscogsProcessURL(discogsService ports.DiscogsPort) *DiscogsProcessURL {
	return &DiscogsProcessURL{
		discogsService: discogsService,
	}
}

func (c *DiscogsProcessURL) processDiscogsURL(ctx context.Context, parsedDiscogsURL *entities.ParsedDiscogsURL) ([]entities.DiscogsRelease, error) {
	var releases []entities.DiscogsRelease
	var err error
	switch parsedDiscogsURL.Type {
	case entities.CollectionType:
		releases, err = c.discogsService.GetCollectionReleases(ctx, parsedDiscogsURL.ID)
	case entities.WantlistType:
		releases, err = c.discogsService.GetWantlistReleases(ctx, parsedDiscogsURL.ID)
	case entities.ListType:
		releases, err = c.discogsService.GetListReleases(ctx, parsedDiscogsURL.ID)
	default:
		return nil, errors.New("unrecognized URL type")
	}

	if err != nil {
		return nil, err
	}

	return releases, nil
}

func parseDiscogsURL(inputURL string) (*entities.ParsedDiscogsURL, error) {
	if inputURL == "" {
		return nil, ErrInvalidDiscogsURL
	}

	parsedURL, err := ensureAndParseInputURL(inputURL)
	if err != nil {
		return nil, err
	}

	// Try to parse as a wantlist URL
	// https://www.discogs.com/es/wantlist?user=digger
	if wantlistURL := parseWantlistURL(parsedURL); wantlistURL != nil {
		return wantlistURL, nil
	}

	// Clean the path by removing language code if present
	cleanPath := cleanURLPath(parsedURL.Path)

	// Try to parse as a collection URL
	// https://www.discogs.com/es/user/digger/collection
	if collectionURL := parseCollectionURL(cleanPath); collectionURL != nil {
		return collectionURL, nil
	}

	// Try to parse as a list URL
	// https://www.discogs.com/es/lists/MyList/1545836
	if listURL := parseListURL(cleanPath); listURL != nil {
		return listURL, nil
	}

	return nil, ErrInvalidDiscogsURL
}

func ensureAndParseInputURL(inputURL string) (*url.URL, error) {
	// add scheme if missing to prevent url.Parse errors
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	return url.Parse(inputURL)
}

func parseWantlistURL(parsedURL *url.URL) *entities.ParsedDiscogsURL {
	if strings.Contains(parsedURL.Path, "/wantlist") {
		user := parsedURL.Query().Get("user")
		if user != "" {
			return &entities.ParsedDiscogsURL{ID: user, Type: entities.WantlistType}
		}
	}
	return nil
}

// clean language code from the URL path if present
func cleanURLPath(path string) string {
	re := regexp.MustCompile(`^/(?:[a-z]{2}/)?(.*)$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func parseCollectionURL(cleanPath string) *entities.ParsedDiscogsURL {
	collectionRe := regexp.MustCompile(`^user/([^/]+)/collection$`)
	collectionMatches := collectionRe.FindStringSubmatch(cleanPath)
	if len(collectionMatches) > 1 {
		return &entities.ParsedDiscogsURL{ID: collectionMatches[1], Type: entities.CollectionType}
	}
	return nil
}

func parseListURL(cleanPath string) *entities.ParsedDiscogsURL {
	listRe := regexp.MustCompile(`^lists/.+/(\d+)$`)
	listMatches := listRe.FindStringSubmatch(cleanPath)
	if len(listMatches) > 1 {
		return &entities.ParsedDiscogsURL{ID: listMatches[1], Type: entities.ListType}
	}
	return nil
}
