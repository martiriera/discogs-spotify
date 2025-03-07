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

func (c *UserController) GetSpotifyUserID(ctx *gin.Context) (string, error) {
	userID, err := c.spotifyService.GetSpotifyUserID(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error getting spotify user id")
	}
	return userID, nil
}
