package spotify

import (
	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const oauthState = "AOWTCN2KHZ"

var scopes = []string{
	"user-read-private",
	"user-read-email",
	"playlist-modify-public",
	"playlist-modify-private",
}

type OAuthController struct {
	config     *oauth2.Config
	oauthState string
}

func NewOAuthController(clientID, clientSecret, redirectURL string) *OAuthController {
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

func (o *OAuthController) GenerateToken(ctx *gin.Context) (*oauth2.Token, error) {
	values := ctx.Request.URL.Query()
	if err := values.Get("error"); err != "" {
		return nil, errors.Wrap(errors.New(err), "spotify: error in callback")
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("spotify: no code in callback")
	}
	actualState := values.Get("state")
	if actualState != o.oauthState {
		return nil, errors.New("spotify: redirect state parameter doesn't match")
	}

	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "spotify: error exchanging code for token")
	}
	return token, nil
}

func (o *OAuthController) StoreToken(ctx *gin.Context, s session.Session, token *oauth2.Token) error {
	err := s.SetData(ctx.Request, ctx.Writer, session.SpotifyTokenKey, token)

	if err != nil {
		return errors.Wrap(err, "spotify: error saving session")
	}

	return nil
}
