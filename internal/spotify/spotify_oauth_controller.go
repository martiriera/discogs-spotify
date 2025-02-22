package spotify

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martiriera/discogs-spotify/internal/session"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var oauthState = generateRandomState()

const (
	ErrNoCode                     = "spotify: no code in callback"
	ErrRedirectStateParamMismatch = "spotify: redirect state parameter doesn't match"
	ErrErrorInCallback            = "spotify: error in callback"
	ErrrExchangingCodeForToken    = "spotify: error exchanging code for token"
	ErrSavingSession              = "spotify: error saving session"
)

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
		return nil, errors.Wrap(errors.New(err), ErrErrorInCallback)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New(ErrNoCode)
	}
	actualState := values.Get("state")
	if actualState != o.oauthState {
		return nil, errors.New(ErrRedirectStateParamMismatch)
	}

	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, ErrrExchangingCodeForToken)
	}
	return token, nil
}

func (o *OAuthController) StoreToken(ctx *gin.Context, s session.Session, token *oauth2.Token) error {
	err := s.SetData(ctx.Request, ctx.Writer, session.SpotifyTokenKey, token)

	if err != nil {
		return errors.Wrap(err, ErrSavingSession)
	}

	return nil
}

func generateRandomState() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
