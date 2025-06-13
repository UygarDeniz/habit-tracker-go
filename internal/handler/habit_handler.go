package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/middleware"
	habitUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/habit"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

type HabitHandler struct {
	createHabitUsecase *habitUsecase.CreateHabitUsecase
	logger             *log.Logger
	v                  *validator.Validate
}

func NewHabitHandler(createHabitUsecase *habitUsecase.CreateHabitUsecase, logger *log.Logger, v *validator.Validate) *HabitHandler {
	return &HabitHandler{
		createHabitUsecase: createHabitUsecase,
		logger:             logger,
		v:                  v,
	}
}

func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	var req dto.CreateHabitDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Failed to decode request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_request_format"}, h.logger)
		return
	}

	if err := h.v.Struct(&req); err != nil {
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "validation_failed"}, err, h.logger)
		return
	}

	habit, err := h.createHabitUsecase.Execute(r.Context(), userID, req)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidInput:
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_input"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "not_found"}, h.logger)
		default:
			h.logger.Printf("Error creating habit: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	// Convert entity to response DTO
	response := dto.HabitResponseDTO{
		ID:               habit.ID,
		UserID:           habit.UserID,
		Name:             habit.Name,
		Description:      habit.Description,
		Motivation:       habit.Motivation,
		Color:            habit.Color,
		Category:         habit.Category,
		Frequency:        habit.Frequency,
		TargetCount:      habit.TargetCount,
		CurrentStreak:    habit.CurrentStreak,
		BestStreak:       habit.BestStreak,
		TotalCompletions: habit.TotalCompletions,
		IsActive:         habit.IsActive,
		CreatedAt:        habit.CreatedAt,
		UpdatedAt:        habit.UpdatedAt,
	}

	// Handle TargetDays conversion
	if habit.TargetDays != nil {
		targetDaysStr, err := entity.ConvertTargetDaysToJSON(habit.TargetDays)
		if err != nil {
			h.logger.Printf("Error converting target days to JSON: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
			return
		}
		response.TargetDays = targetDaysStr
	}

	h.logger.Printf("Habit created successfully. HabitID: %s, UserID: %s", habit.ID, userID)
	utils.WriteJSON(w, http.StatusCreated, utils.APIResponse{"habit": response}, h.logger)
}
