package usecases

import (
	"context"

	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type GetSpotifyUser struct {
	spotifyService ports.SpotifyPort
}

func NewGetSpotifyUser(s ports.SpotifyPort) *GetSpotifyUser {
	return &GetSpotifyUser{spotifyService: s}
}

func (c *GetSpotifyUser) GetUserID(ctx context.Context) (string, error) {
	userID, err := c.spotifyService.GetUserID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}
	return userID, nil
}
