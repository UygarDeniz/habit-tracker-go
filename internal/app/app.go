package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/handler"
	"github.com/uygardeniz/habit-tracker/internal/repository"
	"github.com/uygardeniz/habit-tracker/internal/usecases"
)

type Application struct {
	Logger      *log.Logger
	DB          *sql.DB
	UserHandler *handler.UserHandler
	Validate    *validator.Validate
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
	createUserUsecase := usecases.NewCreateUserUsecase(userRepository)

	v := validator.New()
	// Initialize handlers
	userHandler := handler.NewUserHandler(createUserUsecase, logger, v)

	app := &Application{
		Logger:      logger,
		DB:          db,
		UserHandler: userHandler,
	}

	return app, nil
}
