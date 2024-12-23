package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
)

type AuthRouter struct {
	oauthController *spotify.OAuthController
	session         *session.Session
}

func NewAuthRouter(c *spotify.OAuthController, session *session.Session) *AuthRouter {
	router := &AuthRouter{oauthController: c}
	router.oauthController = c
	router.session = session
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
	token, err := router.oauthController.GenerateToken(c)
	if err != nil {
		util.HandleError(c, err, http.StatusInternalServerError)
		return
	}
	err = router.oauthController.StoreToken(c, *router.session, token)
	if err != nil {
		util.HandleError(c, err, http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
