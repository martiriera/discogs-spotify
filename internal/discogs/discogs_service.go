package discogs

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/client"
	"github.com/martiriera/discogs-spotify/internal/entities"
	"github.com/pkg/errors"
)

var ErrUnauthorized = errors.New("discogs unauthorized error")
var ErrUnexpectedStatus = errors.New("discogs unexpected status error")
var ErrRequest = errors.New("discogs request error")
var ErrResponse = errors.New("discogs response error")

type DiscogsService interface {
	GetReleases(username string) ([]entities.DiscogsRelease, error)
}

type HttpDiscogsService struct {
	client client.HttpClient
}

const basePath = "https://api.discogs.com"

func NewHttpDiscogsService(client client.HttpClient) *HttpDiscogsService {
	return &HttpDiscogsService{client: client}
}

func (s *HttpDiscogsService) GetReleases(username string) ([]entities.DiscogsRelease, error) {
	url := basePath + "/users/" + username + "/collection/folders/0/releases?per_page=100"
	result := make([]entities.DiscogsRelease, 0)
	response, err := doRequest(s.client, url)
	if err != nil {
		return nil, err
	}
	result = append(result, response.Releases...)
	for response.Pagination.Urls.Next != "" {
		response, err = doRequest(s.client, response.Pagination.Urls.Next)
		if err != nil {
			return nil, err
		}
		result = append(result, response.Releases...)
	}
	return result, nil
}

func doRequest(client client.HttpClient, url string) (*entities.DiscogsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(ErrRequest, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.Wrap(ErrUnauthorized, "private collection")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.Wrapf(ErrUnexpectedStatus, "status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response entities.DiscogsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(ErrResponse, err.Error())
	}
	return &response, nil
}
