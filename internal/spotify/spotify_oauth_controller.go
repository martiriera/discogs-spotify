package spotify

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type OAuthController struct {
	service        SpotifyService
	config         *oauth2.Config
	oauthState     string
	tokenStoreFunc func(token *oauth2.Token)
}

func NewOAuthController(clientID, clientSecret, redirectURL string, scopes []string) *OAuthController {
	return &OAuthController{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     spotify.Endpoint,
		},
		oauthState: "random-string",
	}
}

func (o *OAuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := o.config.AuthCodeURL(o.oauthState, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *OAuthController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter to prevent CSRF
	state := r.FormValue("state")
	if state != o.oauthState {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally store the token
	if o.tokenStoreFunc != nil {
		o.tokenStoreFunc(token)
	}

	client := o.config.Client(context.Background(), token)
	userInfo, err := o.service.GetSpotifyUserInfo(client)
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Display user info
	fmt.Fprintf(w, "User Info: %s\n", userInfo)
}

func (o *OAuthController) SetTokenStoreFunc(fn func(token *oauth2.Token)) {
	o.tokenStoreFunc = fn
}
