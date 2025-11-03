package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/usecases"
	"github.com/martiriera/discogs-spotify/internal/utils/env"
)

type AuthRouter struct {
	oauthController *usecases.SpotifyAuthenticate
	session         *ports.SessionPort
	config          *config.Config
}

func NewAuthRouter(c *usecases.SpotifyAuthenticate, session *ports.SessionPort, cfg *config.Config) *AuthRouter {
	router := &AuthRouter{oauthController: c, session: session, config: cfg}
	return router
}

func (router *AuthRouter) SetupRoutes(rg *gin.RouterGroup) {
	rg.GET("/login", router.handleLogin)
	rg.GET("/callback", router.handleLoginCallback)

	if router.config.Spotify.UseProxy {
		proxyGroup := rg.Group("/proxy")
		proxyGroup.GET("/callback/spotify", router.handleProxyCallback)
	}
}

func (router *AuthRouter) handleLogin(ctx *gin.Context) {
	url := router.oauthController.GetAuthURL()
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (router *AuthRouter) handleLoginCallback(ctx *gin.Context) {
	token, err := router.oauthController.GenerateTokenFromGin(ctx)
	if err != nil {
		handleError(ctx, err, http.StatusInternalServerError)
		return
	}
	err = router.oauthController.StoreToken(ctx, *router.session, token)
	if err != nil {
		handleError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/home")
}

func (router *AuthRouter) handleProxyCallback(ctx *gin.Context) {
	localDevURL := env.GetWithDefault("LOCAL_DEV_URL", "http://localhost:8080")

	// Parse the local development callback URL
	localCallbackURL, err := url.Parse(localDevURL + "/auth/callback")
	if err != nil {
		log.Printf("Error parsing local dev URL: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid local dev URL configuration"})
		return
	}

	// Copy all query parameters from the Spotify callback
	query := localCallbackURL.Query()
	for key, values := range ctx.Request.URL.Query() {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	localCallbackURL.RawQuery = query.Encode()

	log.Printf("Auth proxy redirecting to local development server: %s", localCallbackURL.String())

	// Redirect the user's browser to the local development server
	ctx.Redirect(http.StatusTemporaryRedirect, localCallbackURL.String())
}
