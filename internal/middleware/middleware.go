package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/helper"
	jwttoken "github.com/codepnw/stdlib-ticket-system/pkg/jwt"
)

type AuthMiddleware struct {
	token jwttoken.JWTToken
}

func NewMiddleware(token jwttoken.JWTToken) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.ErrorResponse(w, http.StatusUnauthorized, "header is missing")
			return
		}

		args := strings.Fields(authHeader)
		if len(args) != 2 || args[0] != "Bearer" {
			helper.ErrorResponse(w, http.StatusUnauthorized, "invalid header format")
			return
		}

		claims, err := m.token.VerifyAccessToken(args[1])
		if err != nil {
			helper.ErrorResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.ContextUserClaimsKey, claims)
		ctx = context.WithValue(ctx, config.ContextUserIDKey, claims.UserID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
