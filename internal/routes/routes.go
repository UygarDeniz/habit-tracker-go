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
	router.HandleFunc("GET /api/auth/google/login", http.HandlerFunc(app.AuthHandler.HandleGoogleLogin))
	router.HandleFunc("GET /api/auth/google/callback", app.AuthHandler.HandleGoogleCallback)
	router.HandleFunc("GET /api/auth/refresh_token", app.AuthHandler.HandleRefreshToken)
	router.HandleFunc("GET /api/auth/session", app.AuthHandler.HandleGetUserAndAccessToken)
	router.HandleFunc("POST /api/auth/logout", app.AuthHandler.HandleLogout)

	// User routes
	protectedMux.HandleFunc("GET /api/user/me", app.UserHandler.GetMe)

	// Habit routes
	protectedMux.HandleFunc("GET /api/habits", app.HabitHandler.GetHabitsByUserID)
	protectedMux.HandleFunc("POST /api/habits", app.HabitHandler.CreateHabit)
	protectedMux.HandleFunc("GET /api/habits/{habitID}", app.HabitHandler.GetHabit)
	protectedMux.HandleFunc("PUT /api/habits/{habitID}", app.HabitHandler.UpdateHabit)
	protectedMux.HandleFunc("DELETE /api/habits/{habitID}", app.HabitHandler.DeleteHabit)

	// Completion routes
	protectedMux.HandleFunc("GET /api/completions", app.CompletionHandler.GetCompletions)
	protectedMux.HandleFunc("POST /api/habits/{habitID}/completions", app.CompletionHandler.CreateCompletion)
	protectedMux.HandleFunc("GET /api/completions/{completionID}", app.CompletionHandler.GetCompletion)
	protectedMux.HandleFunc("PUT /api/completions/{completionID}", app.CompletionHandler.UpdateCompletion)
	protectedMux.HandleFunc("DELETE /api/completions/{completionID}", app.CompletionHandler.DeleteCompletion)

	// Apply auth middleware to protected routes
	router.Handle("/api/user/me", authMiddleware.RequireAuth(protectedMux))
	router.Handle("/api/habits", authMiddleware.RequireAuth(protectedMux))
	router.Handle("/api/habits/", authMiddleware.RequireAuth(protectedMux))
	router.Handle("/api/completions", authMiddleware.RequireAuth(protectedMux))
	router.Handle("/api/completions/", authMiddleware.RequireAuth(protectedMux))

	handler := authMiddleware.Logging(router)

	return handler
}
