//go:build nullauth

package daisy

import (
	"context"
	"net/http"
)

type nullAuth struct{}

func (n *nullAuth) Middleware(next http.Handler) http.Handler {
	return next
}

func (n *nullAuth) CurrentUserPrincipal(ctx context.Context) (string, error) {
	return "/jane@example.com/", nil
}

func NullAuth(_ AuthProvider) AuthProvider {
	return &nullAuth{}
}
