package discogs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	httpClient "github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/core/entities"

	"github.com/pkg/errors"
)

var ErrUnauthorized = errors.New("discogs unauthorized error")
var ErrUnexpectedStatus = errors.New("discogs unexpected status error")
var ErrRequest = errors.New("discogs request error")
var ErrResponse = errors.New("discogs response error")

type HTTPService struct {
	client httpClient.HTTPClient
}

const basePath = "https://api.discogs.com"

func NewHTTPService(client httpClient.HTTPClient) *HTTPService {
	return &HTTPService{client: client}
}

func (s *HTTPService) GetCollectionReleases(ctx context.Context, username string) ([]entities.DiscogsRelease, error) {
	url := basePath + "/users/" + username + "/collection/folders/0/releases?per_page=100&sort=artist&sort_order=asc"
	return paginate(ctx, s.client, url)
}

func (s *HTTPService) GetWantlistReleases(ctx context.Context, username string) ([]entities.DiscogsRelease, error) {
	url := basePath + "/users/" + username + "/wants?per_page=100&sort=artist&sort_order=asc"
	return paginate(ctx, s.client, url)
}

func (s *HTTPService) GetListReleases(ctx context.Context, listID string) ([]entities.DiscogsRelease, error) {
	url := basePath + "/lists/" + listID
	response, err := doRequest(ctx, s.client, url)
	if err != nil {
		return nil, err
	}
	return response.GetReleases(), nil
}

func paginate(ctx context.Context, client httpClient.HTTPClient, url string) ([]entities.DiscogsRelease, error) {
	result := make([]entities.DiscogsRelease, 0)
	response, err := doRequest(ctx, client, url)
	if err != nil {
		return nil, err
	}
	result = append(result, response.GetReleases()...)
	for response.GetPagination().Urls.Next != "" {
		response, err = doRequest(ctx, client, response.GetPagination().Urls.Next)
		if err != nil {
			return nil, err
		}
		result = append(result, response.GetReleases()...)
	}
	return result, nil
}

func doRequest(ctx context.Context, client httpClient.HTTPClient, url string) (entities.DiscogsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.Wrap(ErrUnauthorized, "private resource")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrapf(ErrUnexpectedStatus, "status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if strings.Contains(url, "collection") {
		var response entities.DiscogsCollectionResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, errors.Wrap(ErrResponse, err.Error())
		}
		return &response, nil
	} else if strings.Contains(url, "wants") {
		var response entities.DiscogsWantlistResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, errors.Wrap(ErrResponse, err.Error())
		}
		return &response, nil
	} else if strings.Contains(url, "lists") {
		var response entities.DiscogsListResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, errors.Wrap(ErrResponse, err.Error())
		}
		return &response, nil
	}

	return nil, errors.Wrapf(ErrResponse, "unknown response type for URL: %s", url)
}
