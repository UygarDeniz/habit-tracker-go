package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken generates a JWT access token for the given user ID
func GenerateAccessToken(userID string) (string, error) {
	secretKey := os.Getenv("JWT_ACCESS_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_ACCESS_SECRET environment variable is not set")
	}
	if len(secretKey) < 32 {
		return "", errors.New("JWT_ACCESS_SECRET must be at least 32 characters")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a JWT refresh token for the given user ID
func GenerateRefreshToken(userID string) (string, error) {
	secretKey := os.Getenv("JWT_REFRESH_SECRET")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token with the given secret
func ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
