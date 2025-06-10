package routes

import (
	"net/http"

	"github.com/uygardeniz/habit-tracker/internal/app"
)

func SetupRoutes(app *app.Application) *http.ServeMux {
	router := http.NewServeMux()

	// Authentication routes
	router.HandleFunc("GET /api/auth/google/login", app.AuthHandler.HandleGoogleLogin)
	router.HandleFunc("GET /api/auth/google/callback", app.AuthHandler.HandleGoogleCallback)
	router.HandleFunc("POST /api/auth/refresh_token", app.AuthHandler.HandleRefreshToken)
	router.HandleFunc("POST /api/auth/logout", app.AuthHandler.HandleLogout)

	return router
}
