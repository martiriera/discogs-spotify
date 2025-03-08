// README: Move to user package when growing
package spotify

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
)

type UserController struct {
	spotifyService ports.SpotifyPort
}

func NewUserController(s ports.SpotifyPort) *UserController {
	return &UserController{spotifyService: s}
}

func (c *UserController) GetSpotifyUserID(ctx *gin.Context) (string, error) {
	userID, err := c.spotifyService.GetSpotifyUserID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}
	return userID, nil
}
