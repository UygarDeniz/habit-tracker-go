package routes

import (
	"net/http"

	"github.com/uygardeniz/habit-tracker/internal/app"
)

func SetupRoutes(app *app.Application) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/auth/register", app.UserHandler.HandleRegisterUser)

	return router
}
