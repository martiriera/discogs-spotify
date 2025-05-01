package usecases

import (
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

func (c *DiscogsProcessURL) processDiscogsURL(discogsURL string) ([]entities.DiscogsRelease, *entities.DiscogsInputURL, error) {
	// fetch releases
	parsedDiscogsURL, err := parseDiscogsURL(discogsURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error parsing Discogs URL")
	}

	var releases []entities.DiscogsRelease
	if parsedDiscogsURL.Type == entities.CollectionType {
		releases, err = c.discogsService.GetCollectionReleases(parsedDiscogsURL.ID)
	} else if parsedDiscogsURL.Type == entities.WantlistType {
		releases, err = c.discogsService.GetWantlistReleases(parsedDiscogsURL.ID)
	} else if parsedDiscogsURL.Type == entities.ListType {
		releases, err = c.discogsService.GetListReleases(parsedDiscogsURL.ID)
	} else {
		return nil, nil, errors.New("unrecognized URL type")
	}

	if err != nil {
		return nil, nil, err
	}

	return releases, parsedDiscogsURL, nil
}

func parseDiscogsURL(inputURL string) (*entities.DiscogsInputURL, error) {
	if inputURL == "" {
		return nil, ErrInvalidDiscogsURL
	}

	// add scheme if missing to prevent url.Parse errors
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}

	// Handle wantlist URLs with query parameters
	if strings.Contains(parsedURL.Path, "/wantlist") {
		// Extract the user parameter from the query
		queryParams := parsedURL.Query()
		user := queryParams.Get("user")
		if user != "" {
			return &entities.DiscogsInputURL{ID: user, Type: entities.WantlistType}, nil
		}
		return nil, ErrInvalidDiscogsURL
	}

	// Handle collection and list URLs
	// Remove language code if present (e.g., /es/)
	path := parsedURL.Path
	re := regexp.MustCompile(`^/(?:[a-z]{2}/)?(.*)$`)
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return nil, ErrInvalidDiscogsURL
	}

	cleanPath := matches[1]

	// Check for collection URLs
	collectionRe := regexp.MustCompile(`^user/([^/]+)/collection$`)
	collectionMatches := collectionRe.FindStringSubmatch(cleanPath)
	if len(collectionMatches) > 1 {
		return &entities.DiscogsInputURL{ID: collectionMatches[1], Type: entities.CollectionType}, nil
	}

	// Check for list URLs
	listRe := regexp.MustCompile(`^lists/.+/(\d+)$`)
	listMatches := listRe.FindStringSubmatch(cleanPath)
	if len(listMatches) > 1 {
		return &entities.DiscogsInputURL{ID: listMatches[1], Type: entities.ListType}, nil
	}

	return nil, ErrInvalidDiscogsURL
}
