package server

import (
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/spotify"
	"github.com/martiriera/discogs-spotify/util"
)

type AuthRouter struct {
	oauthController *spotify.OAuthController
}

func NewAuthRouter(c *spotify.OAuthController) *http.ServeMux {
	router := &AuthRouter{oauthController: c}
	router.oauthController = c

	mux := http.NewServeMux()
	mux.Handle("/login", http.HandlerFunc(router.handleLogin))
	mux.Handle("/callback", http.HandlerFunc(router.handleLoginCallback))
	return mux
}

func (router *AuthRouter) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := router.oauthController.GetAuthUrl()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (router *AuthRouter) handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter to prevent CSRF
	state := r.FormValue("state")
	code := r.FormValue("code")
	_, err := router.oauthController.GetServiceFromCallback(state, code)
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	// router.playlistCreator.SetSpotifyService(service)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
