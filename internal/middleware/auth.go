package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/andre/project-app-bioskop-golang/internal/services"
)

// AuthMiddleware is a middleware that checks for valid JWT token
func AuthMiddleware(userService *services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Verify token
			userID, err := userService.VerifyToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), "userID", userID)
			ctx = context.WithValue(ctx, "token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetTokenFromContext extracts token from context
func GetTokenFromContext(r *http.Request) (string, error) {
	token, ok := r.Context().Value("token").(string)
	if !ok {
		return "", errors.New("token not found in context")
	}
	return token, nil
}
