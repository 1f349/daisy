//go:build nullauth

package daisy

import (
	"context"
	"github.com/charmbracelet/log"
	"net/http"
)

func NewAuth(db string, logger *log.Logger) AuthProvider {
	return &nullAuth{}
}

type nullAuth struct{}

func (n *nullAuth) Middleware(next http.Handler) http.Handler {
	return next
}

func (n *nullAuth) CurrentUserPrincipal(ctx context.Context) (string, error) {
	return "/jane@example.com/", nil
}
