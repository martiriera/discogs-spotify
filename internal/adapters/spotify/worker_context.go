package spotify

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// WorkerContextProvider implements ports.ContextPort for headless worker use.
// It fetches tokens via Client Credentials and is a no-op for user ID ops.
type WorkerContextProvider struct {
	tokenProvider *ClientCredentialsProvider
}

func NewWorkerContextProvider(p *ClientCredentialsProvider) *WorkerContextProvider {
	return &WorkerContextProvider{tokenProvider: p}
}

func (w *WorkerContextProvider) GetToken(ctx context.Context) (*oauth2.Token, error) {
	accessToken, err := w.tokenProvider.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("worker context: getting spotify token: %w", err)
	}
	return &oauth2.Token{AccessToken: accessToken}, nil
}

func (*WorkerContextProvider) GetUserID(_ context.Context) (string, error) {
	return "", nil
}

func (*WorkerContextProvider) SetUserID(_ context.Context, _ string) error {
	return nil
}
