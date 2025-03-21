package session

const (
	AuthSessionName = "auth-session"
)

type ContextKey string

const (
	SpotifyTokenKey  ContextKey = "spotify-token"
	SpotifyUserIDKey ContextKey = "spotify-user-id"
)
