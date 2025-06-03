package main

import (
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/uygardeniz/habit-tracker/internal/app"
	"github.com/uygardeniz/habit-tracker/internal/routes"
)

func main() {
	port := ":8080"

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	application, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer application.DB.Close()

	router := routes.SetupRoutes(application)

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	application.Logger.Printf("Server is running on port %s", port)

	err = server.ListenAndServe()

	if err != nil {
		application.Logger.Fatalf("Server failed to start: %v", err)
	}
}
