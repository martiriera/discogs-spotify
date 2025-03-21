package usecases

import (
	"crypto/rand"
	"encoding/base64"
	"context"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
)

var oauthState, _ = generateRandomState()

const (
	ErrNoCode                     = "spotify: no code in callback"
	ErrRedirectStateParamMismatch = "spotify: redirect state parameter doesn't match"
	ErrErrorInCallback            = "spotify: error in callback"
	ErrExchangingCode             = "spotify: error exchanging code"
	ErrSavingSession              = "spotify: error saving session"
)

var scopes = []string{
	"user-read-private",
	"user-read-email",
	"playlist-modify-public",
	"playlist-modify-private",
}

type SpotifyAuthenticate struct {
	config     *oauth2.Config
	oauthState string
}

func NewSpotifyAuthenticate(clientID, clientSecret, redirectURL string) *SpotifyAuthenticate {
	return &SpotifyAuthenticate{
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

func (o *SpotifyAuthenticate) GetAuthURL() string {
	return o.config.AuthCodeURL(o.oauthState, oauth2.AccessTypeOffline)
}

func (o *SpotifyAuthenticate) GenerateToken(ctx context.Context, code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New(ErrNoCode)
	}

	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, ErrExchangingCode)
	}
	return token, nil
}

func (o *SpotifyAuthenticate) GenerateTokenFromGin(ctx *gin.Context) (*oauth2.Token, error) {
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

	return o.GenerateToken(ctx, code)
}

func (o *SpotifyAuthenticate) StoreToken(ctx *gin.Context, s ports.SessionPort, token *oauth2.Token) error {
	err := s.SetData(ctx.Request, ctx.Writer, session.SpotifyTokenKey, token)

	if err != nil {
		return errors.Wrap(err, ErrSavingSession)
	}

	return nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
