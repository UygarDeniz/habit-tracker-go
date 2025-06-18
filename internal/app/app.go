package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/handler"
	"github.com/uygardeniz/habit-tracker/internal/repository"
	authUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/auth"
	completionUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/completion"
	habitUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/habit"
	userUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/user"
)

type Application struct {
	Logger            *log.Logger
	DB                *sql.DB
	AuthHandler       *handler.AuthHandler
	HabitHandler      *handler.HabitHandler
	CompletionHandler *handler.CompletionHandler
	UserHandler       *handler.UserHandler
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
	completionRepository := repository.NewPostgresCompletionRepository(db)

	// Initialize user usecases
	getMeUsecase := userUsecase.NewGetMeUsecase(userRepository)
	getUserByIDUsecase := userUsecase.NewGetUserByIDUsecase(userRepository)

	// Initialize auth usecases
	loginOrRegisterGoogleUserUsecase := authUsecase.NewLoginOrRegisterGoogleUserUsecase(userRepository)

	// Initialize habit usecases
	createHabitUsecase := habitUsecase.NewCreateHabitUsecase(habitRepository)
	getHabitUsecase := habitUsecase.NewGetHabitUsecase(habitRepository)
	getHabitsByUserUsecase := habitUsecase.NewGetHabitsByUserUsecase(habitRepository)
	updateHabitUsecase := habitUsecase.NewUpdateHabitUsecase(habitRepository)
	deleteHabitUsecase := habitUsecase.NewDeleteHabitUsecase(habitRepository)

	// Initialize completion usecases
	createCompletionUsecase := completionUsecase.NewCreateCompletionUsecase(completionRepository, habitRepository)
	getCompletionUsecase := completionUsecase.NewGetCompletionUsecase(completionRepository)
	getCompletionsUsecase := completionUsecase.NewGetCompletionsUsecase(completionRepository)
	updateCompletionUsecase := completionUsecase.NewUpdateCompletionUsecase(completionRepository, habitRepository)
	deleteCompletionUsecase := completionUsecase.NewDeleteCompletionUsecase(completionRepository, habitRepository)

	// Initialize handlers
	userHandler := handler.NewUserHandler(logger, getMeUsecase)
	authHandler := handler.NewAuthHandler(logger, loginOrRegisterGoogleUserUsecase, getUserByIDUsecase)
	habitHandler := handler.NewHabitHandler(createHabitUsecase, getHabitUsecase, updateHabitUsecase, getHabitsByUserUsecase, deleteHabitUsecase, logger, v)
	completionHandler := handler.NewCompletionHandler(createCompletionUsecase, getCompletionUsecase, getCompletionsUsecase, updateCompletionUsecase, deleteCompletionUsecase, logger, v)

	app := &Application{
		Logger:            logger,
		DB:                db,
		AuthHandler:       authHandler,
		HabitHandler:      habitHandler,
		CompletionHandler: completionHandler,
		UserHandler:       userHandler,
	}

	return app, nil
}
