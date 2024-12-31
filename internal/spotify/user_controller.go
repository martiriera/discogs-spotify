// README: Move to user package when growing
package spotify

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserController struct {
	spotifyService SpotifyService
}

func NewUserController(s SpotifyService) *UserController {
	return &UserController{spotifyService: s}
}

func (c *UserController) GetSpotifyUserId(ctx *gin.Context) (string, error) {
	userId, err := c.spotifyService.GetSpotifyUserId(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}
	return userId, nil
}
