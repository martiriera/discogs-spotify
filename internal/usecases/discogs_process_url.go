package usecases

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/entities"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/util"
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
	stop := util.StartTimer("CreatePlaylist")
	defer stop()

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

func parseDiscogsURL(urlStr string) (*entities.DiscogsInputURL, error) {
	// validate host
	parsedURL, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	pathWithQuery := parsedURL.Path
	if parsedURL.RawQuery != "" && strings.Contains(parsedURL.RawQuery, "user=") {
		pathWithQuery += "?" + parsedURL.RawQuery
	}

	matchingURL := ""
	if parsedURL.Host == "www.discogs.com" {
		matchingURL = pathWithQuery
	} else if parsedURL.Host == "" {
		matchingURL = "/" + strings.SplitN(pathWithQuery, "/", 2)[1]
	} else {
		return nil, ErrInvalidDiscogsURL
	}

	// validate path
	re := regexp.MustCompile(`^/(?:[a-z]{2}/)?(?:user/(.+)/collection|wantlist\?user=(.+))$|lists/.+/(\d+)`)
	matches := re.FindStringSubmatch(matchingURL)
	if matches == nil {
		return nil, ErrInvalidDiscogsURL
	}

	for i, match := range matches {
		// https://www.discogs.com/es/user/digger/collection
		if i == 1 && match != "" {
			return &entities.DiscogsInputURL{ID: match, Type: entities.CollectionType}, nil
		}
		// https://www.discogs.com/es/wantlist?user=digger
		if i == 2 && match != "" {
			return &entities.DiscogsInputURL{ID: match, Type: entities.WantlistType}, nil
		}

		// https://www.discogs.com/es/lists/MyList/1545836
		if i == 3 && match != "" {
			return &entities.DiscogsInputURL{ID: match, Type: entities.ListType}, nil
		}
	}
	return nil, ErrInvalidDiscogsURL
}
