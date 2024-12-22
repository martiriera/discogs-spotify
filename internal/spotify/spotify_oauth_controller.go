package spotify

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const oauthState = "AOWTCN2KHZ"

type OAuthController struct {
	config     *oauth2.Config
	oauthState string
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

func (o *OAuthController) SetToken(values url.Values) error {
	if err := values.Get("error"); err != "" {
		return errors.Wrap(errors.New(err), "spotify: error in callback")
	}
	code := values.Get("code")
	if code == "" {
		return errors.New("spotify: no code in callback")
	}
	actualState := values.Get("state")
	if actualState != o.oauthState {
		return errors.New("spotify: redirect state parameter doesn't match")
	}
	_, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		return errors.Wrap(err, "spotify: error exchanging code for token")
	}
	return nil
}
