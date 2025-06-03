package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/usecases"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

type UserHandler struct {
	createUserUsecase *usecases.CreateUserUsecase
	logger            *log.Logger
	validate          *validator.Validate
}

func NewUserHandler(createUserUsecase *usecases.CreateUserUsecase, logger *log.Logger, validate *validator.Validate) *UserHandler {
	return &UserHandler{createUserUsecase: createUserUsecase, logger: logger, validate: validate}
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	type registerUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=20"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("failed to decode request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid request format"}, h.logger)
		return
	}

	err = h.validate.Struct(req)

	if err != nil {
		h.logger.Printf("failed to validate request body: %v", err)
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid request format"}, err, h.logger)
		return
	}

	user, err := h.createUserUsecase.Execute(r.Context(), req.Email, req.Username, req.Password)

	if err != nil {
		h.logger.Printf("failed to create user: %v", err)

		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			utils.WriteJSON(w, http.StatusConflict, utils.APIResponse{"error": apperrors.ErrUserAlreadyExists.Error()}, h.logger)
		} else if errors.Is(err, apperrors.ErrInvalidInput) {
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": apperrors.ErrInvalidInput.Error(), "details": err.Error()}, h.logger)
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal server error"}, h.logger)
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.APIResponse{"message": "user created successfully", "user": user}, h.logger)
}
