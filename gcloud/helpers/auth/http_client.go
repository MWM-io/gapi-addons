package auth

import (
	"context"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"google.golang.org/api/idtoken"
)

// GetS2SClient returns a http.Client that automatically adds an "Authorization" header
// audience can be the service account, url, or load balancer url of the destination
func GetS2SClient(ctx context.Context, audience string) (*http.Client, errors.Error) {
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, errors.Wrap(err).WithKind("init_client_failed")
	}

	return client, nil
}
