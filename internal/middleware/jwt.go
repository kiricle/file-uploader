package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kiricle/file-uploader/internal/services"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JwtMiddleware(jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Ensure format is "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader { // If "Bearer " is not found, authHeader remains unchanged
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// Parse and validate the token
			claims, err := jwtService.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Store user ID in request context)
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
