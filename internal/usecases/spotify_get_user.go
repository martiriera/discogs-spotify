package usecases

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type GetSpotifyUser struct {
	spotifyService ports.SpotifyPort
}

func NewGetSpotifyUser(s ports.SpotifyPort) *GetSpotifyUser {
	return &GetSpotifyUser{spotifyService: s}
}

func (c *GetSpotifyUser) GetSpotifyUserID(ctx *gin.Context) (string, error) {
	userID, err := c.spotifyService.GetSpotifyUserID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}
	return userID, nil
}
