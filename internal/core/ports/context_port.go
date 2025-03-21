package ports

import (
	"context"

	"golang.org/x/oauth2"
)

// ContextPort defines the interface for retrieving data from a context
type ContextPort interface {
	GetToken(ctx context.Context) (*oauth2.Token, error)
	GetUserID(ctx context.Context) (string, error)
	SetUserID(ctx context.Context, userID string) error
}
