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

func (router *AuthRouter) handleLogin(ctx *gin.Context) {
	url := router.oauthController.GetAuthUrl()
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (router *AuthRouter) handleLoginCallback(ctx *gin.Context) {
	token, err := router.oauthController.GenerateToken(ctx)
	if err != nil {
		util.HandleError(ctx, err, http.StatusInternalServerError)
		return
	}
	err = router.oauthController.StoreToken(ctx, *router.session, token)
	if err != nil {
		util.HandleError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/api/playlist"+ctx.Request.URL.RawQuery)
}
