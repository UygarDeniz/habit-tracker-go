package routes

import (
	"net/http"

	"github.com/uygardeniz/habit-tracker/internal/app"
	"github.com/uygardeniz/habit-tracker/internal/middleware"
)

func SetupRoutes(app *app.Application) http.Handler {
	router := http.NewServeMux()
	protectedMux := http.NewServeMux()
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(app.Logger)

	// Authentication routes
	router.HandleFunc("GET /api/auth/google/login", app.AuthHandler.HandleGoogleLogin)
	router.HandleFunc("GET /api/auth/google/callback", app.AuthHandler.HandleGoogleCallback)
	router.HandleFunc("POST /api/auth/refresh_token", app.AuthHandler.HandleRefreshToken)
	router.HandleFunc("POST /api/auth/logout", app.AuthHandler.HandleLogout)

	// Habit routes
	protectedMux.HandleFunc("POST /api/habits", app.HabitHandler.CreateHabit)

	// Apply auth middleware to protected routes
	router.Handle("/api/habits", authMiddleware.RequireAuth(protectedMux))

	handler := authMiddleware.Logging(router)

	return handler
}
