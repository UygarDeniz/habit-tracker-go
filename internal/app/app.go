package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/handler"
	"github.com/uygardeniz/habit-tracker/internal/repository"
	habitUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/habit"
	userUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/user"
)

type Application struct {
	Logger       *log.Logger
	DB           *sql.DB
	AuthHandler  *handler.AuthHandler
	HabitHandler *handler.HabitHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	v := validator.New()
	db, err := repository.OpenDB()

	if err != nil {
		return nil, fmt.Errorf("failed to prepare database connection: %w", err)
	}

	// Initialize repositories
	userRepository := repository.NewPostgresUserRepository(db)
	habitRepository := repository.NewPostgresHabitRepository(db)
	// Initialize usecases
	loginOrRegisterGoogleUserUsecase := userUsecase.NewLoginOrRegisterGoogleUserUsecase(userRepository)
	createHabitUsecase := habitUsecase.NewCreateHabitUsecase(habitRepository)
	// Initialize handlers

	authHandler := handler.NewAuthHandler(logger, loginOrRegisterGoogleUserUsecase)
	habitHandler := handler.NewHabitHandler(createHabitUsecase, logger, v)
	app := &Application{
		Logger:       logger,
		DB:           db,
		AuthHandler:  authHandler,
		HabitHandler: habitHandler,
	}

	return app, nil
}
