package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

type AuthMiddleware struct {
	logger *log.Logger
}

func NewAuthMiddleware(logger *log.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
	}
}

// RequireAuth validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		// First, check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// If token is not in the header, check for the cookie
		if tokenString == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			m.logger.Println("Authorization token not found")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, m.logger)
			return
		}

		token, err := utils.ValidateToken(tokenString, os.Getenv("JWT_ACCESS_SECRET"))
		if err != nil || !token.Valid {
			m.logger.Printf("Invalid token: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "invalid_token"}, m.logger)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Printf("Failed to parse token claims")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "invalid_token_claims"}, m.logger)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			m.logger.Printf("Missing or invalid user ID in token claims")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "invalid_user_id_in_token"}, m.logger)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func (m *AuthMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
