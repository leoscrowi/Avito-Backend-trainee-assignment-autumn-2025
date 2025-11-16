package middleware

import (
	"net/http"
	"strings"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
)

var ErrorResponseUnauthorized = &domain.ErrorResponse{
	Code:    domain.UNAUTHORIZED,
	Message: "Unauthorized",
}

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isValidToken(r, cfg) {
				domain.WriteError(w, ErrorResponseUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isAdminToken(r, cfg) {
				domain.WriteError(w, ErrorResponseUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isAdminToken(r *http.Request, cfg *config.Config) bool {
	token := getTokenFromRequest(r)
	return token == cfg.AuthConfig.AdminToken
}

func isValidToken(r *http.Request, cfg *config.Config) bool {
	token := getTokenFromRequest(r)
	return token == cfg.AuthConfig.AdminToken || token == cfg.AuthConfig.UserToken
}

func getTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}
