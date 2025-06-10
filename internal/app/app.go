package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/uygardeniz/habit-tracker/internal/handler"
	"github.com/uygardeniz/habit-tracker/internal/repository"
	"github.com/uygardeniz/habit-tracker/internal/usecases"
)

type Application struct {
	Logger      *log.Logger
	DB          *sql.DB
	AuthHandler *handler.AuthHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := repository.OpenDB()

	if err != nil {
		return nil, fmt.Errorf("failed to prepare database connection: %w", err)
	}

	// Initialize repositories
	userRepository := repository.NewPostgresUserRepository(db)

	// Initialize usecases
	loginOrRegisterGoogleUserUsecase := usecases.NewLoginOrRegisterGoogleUserUsecase(userRepository)

	// Initialize handlers

	authHandler := handler.NewAuthHandler(logger, loginOrRegisterGoogleUserUsecase)

	app := &Application{
		Logger:      logger,
		DB:          db,
		AuthHandler: authHandler,
	}

	return app, nil
}
