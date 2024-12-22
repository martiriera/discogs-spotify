package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
)

type AuthRouter struct {
	oauthController *spotify.OAuthController
}

func NewAuthRouter(c *spotify.OAuthController) *AuthRouter {
	router := &AuthRouter{oauthController: c}
	router.oauthController = c
	return router
}

func (router *AuthRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", router.handleLogin)
	rg.POST("/callback", router.handleLoginCallback)
}

func (router *AuthRouter) handleLogin(c *gin.Context) {
	url := router.oauthController.GetAuthUrl()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (router *AuthRouter) handleLoginCallback(c *gin.Context) {
	err := router.oauthController.SetToken(c.Request.URL.Query())
	if err != nil {
		util.HandleError(c, err, http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
