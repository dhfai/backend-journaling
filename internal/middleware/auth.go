package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend-journaling/internal/handlers"
	"backend-journaling/pkg/jwt"
)

func Authenticate(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.WriteError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handlers.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]
			claims, err := jwtManager.Verify(token)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					handlers.WriteError(w, http.StatusUnauthorized, "Token has expired")
					return
				}
				handlers.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value("user").(*jwt.Claims)

			if claims.Role != role {
				handlers.WriteError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
