package middleware

import (
	"arch/internal/server/entity"
	"arch/internal/server/service/auth"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

type authService interface {
	Login(ctx context.Context, credentials entity.UserCredentials, resource entity.ResourceInfo) error
}

type basicAuth struct {
	authService authService
	logger      *zerolog.Logger
}

func NewBasicAuth(service authService, logger *zerolog.Logger) func(http.Handler) http.Handler {
	b := basicAuth{
		authService: service,
		logger:      logger,
	}

	return b.basicAuthMiddleware
}

func (b *basicAuth) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		creds := entity.UserCredentials{
			Login:    username,
			Password: password,
		}
		resource := entity.ResourceInfo{
			Name:   r.RequestURI,
			Method: strings.ToUpper(r.Method),
		}
		if err := b.authService.Login(r.Context(), creds, resource); err != nil {
			switch {
			case errors.Is(err, auth.ErrUnauthorized):
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return

			case errors.Is(err, auth.ErrBadPermission):
				http.Error(w, "Bad permissions", http.StatusForbidden)
				return

			default:
				http.Error(w, fmt.Sprintf("b.authService.Login: %v", err), http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
