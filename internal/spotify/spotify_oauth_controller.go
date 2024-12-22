package spotify

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const oauthState = "AOWTCN2KHZ"

type OAuthController struct {
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
		oauthState: oauthState,
	}
}

func (o *OAuthController) GetAuthUrl() string {
	return o.config.AuthCodeURL(o.oauthState, oauth2.AccessTypeOffline)
}

func (o *OAuthController) GetServiceFromCallback(state string, code string) (*HttpSpotifyService, error) {
	// Verify state parameter to prevent CSRF
	if state != o.oauthState {
		return nil, errors.New("state mismatch")
	}

	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange token")
	}

	// Optionally store the token
	if o.tokenStoreFunc != nil {
		o.tokenStoreFunc(token)
	}

	client := o.config.Client(context.Background(), token)
	return NewHttpSpotifyService(client), nil
}

func (o *OAuthController) SetTokenStoreFunc(fn func(token *oauth2.Token)) {
	o.tokenStoreFunc = fn
}
