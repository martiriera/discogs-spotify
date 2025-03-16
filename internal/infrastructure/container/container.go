package container

import (
	"net/http"

	"github.com/martiriera/discogs-spotify/internal/adapters/client"
	"github.com/martiriera/discogs-spotify/internal/adapters/discogs"
	"github.com/martiriera/discogs-spotify/internal/adapters/spotify"
	"github.com/martiriera/discogs-spotify/internal/core/ports"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/server"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/session"
	"github.com/martiriera/discogs-spotify/internal/usecases"
)

type Container struct {
	Config            *config.Config
	Server            *server.Server
	HTTPServer        *http.Server
	Session           ports.SessionPort
	DiscogsService    ports.DiscogsPort
	SpotifyService    ports.SpotifyPort
	PlaylistController *usecases.Controller
	OAuthController   *usecases.SpotifyAuthenticate
	UserController    *usecases.GetSpotifyUser
	HTTPClientFactory *client.HTTPClientFactory
}

func NewContainer(cfg *config.Config) *Container {
	c := &Container{
		Config: cfg,
		HTTPClientFactory: client.NewHTTPClientFactory(),
	}

	c.initSession()
	c.initServices()
	c.initControllers()
	c.initServer()

	return c
}

func (c *Container) initSession() {
	s := session.NewGorillaSession()
	s.Init(c.Config.Session.MaxAgeSec)
	c.Session = s
}

func (c *Container) initServices() {
	discogsClient := c.HTTPClientFactory.CreateDiscogsClient(
		c.Config.HTTP.DiscogsTimeout,
		c.Config.HTTP.RetryAttempts,
		c.Config.HTTP.RetryDelay,
	)

	spotifyClient := c.HTTPClientFactory.CreateSpotifyClient(
		c.Config.HTTP.SpotifyTimeout,
		c.Config.HTTP.RetryAttempts,
		c.Config.HTTP.RetryDelay,
	)

	c.DiscogsService = discogs.NewHTTPService(discogsClient)
	c.SpotifyService = spotify.NewHTTPService(spotifyClient)
}

func (c *Container) initControllers() {
	c.PlaylistController = usecases.NewPlaylistController(
		c.DiscogsService,
		c.SpotifyService,
	)

	c.OAuthController = usecases.NewSpotifyAuthenticate(
		c.Config.Spotify.ClientID,
		c.Config.Spotify.ClientSecret,
		c.Config.Spotify.RedirectURI,
	)

	c.UserController = usecases.NewGetSpotifyUser(c.SpotifyService)
}

func (c *Container) initServer() {
	c.Server = server.NewServer(
		c.PlaylistController,
		c.OAuthController,
		c.UserController,
		c.Session,
	)

	c.HTTPServer = &http.Server{
		Addr:         ":" + c.Config.Server.Port,
		Handler:      c.Server,
		ReadTimeout:  c.Config.Server.ReadTimeout,
		WriteTimeout: c.Config.Server.WriteTimeout,
		IdleTimeout:  c.Config.Server.IdleTimeout,
	}
}

func (c *Container) GetHTTPServer() *http.Server {
	return c.HTTPServer
}