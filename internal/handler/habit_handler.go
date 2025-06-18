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
	createHabitUsecase     *habitUsecase.CreateHabitUsecase
	getHabitUsecase        *habitUsecase.GetHabitUsecase
	updateHabitUsecase     *habitUsecase.UpdateHabitUsecase
	getHabitsByUserUsecase *habitUsecase.GetHabitsByUserUsecase
	deleteHabitUsecase     *habitUsecase.DeleteHabitUsecase
	logger                 *log.Logger
	v                      *validator.Validate
}

func NewHabitHandler(
	createHabitUsecase *habitUsecase.CreateHabitUsecase,
	getHabitUsecase *habitUsecase.GetHabitUsecase,
	updateHabitUsecase *habitUsecase.UpdateHabitUsecase,
	getHabitsByUserUsecase *habitUsecase.GetHabitsByUserUsecase,
	deleteHabitUsecase *habitUsecase.DeleteHabitUsecase,
	logger *log.Logger,
	v *validator.Validate,
) *HabitHandler {
	return &HabitHandler{
		createHabitUsecase:     createHabitUsecase,
		getHabitUsecase:        getHabitUsecase,
		updateHabitUsecase:     updateHabitUsecase,
		getHabitsByUserUsecase: getHabitsByUserUsecase,
		deleteHabitUsecase:     deleteHabitUsecase,
		logger:                 logger,
		v:                      v,
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

	response, err := toHabitResponseDTO(habit)
	if err != nil {
		h.logger.Printf("Failed to map habit to response DTO for habit %s: %v", habit.ID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		return
	}

	h.logger.Printf("Habit created successfully. HabitID: %s, UserID: %s", habit.ID, userID)
	utils.WriteJSON(w, http.StatusCreated, utils.APIResponse{"habit": response}, h.logger)
}

func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	habitID := r.PathValue("habitID")

	habit, err := h.getHabitUsecase.Execute(r.Context(), habitID, userID)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "you are not allowed to access this habit"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "habit not found"}, h.logger)
		default:
			h.logger.Printf("Error getting habit: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal server error"}, h.logger)
		}
		return
	}

	response, err := toHabitResponseDTO(habit)
	if err != nil {
		h.logger.Printf("Failed to map habit to response DTO for habit %s: %v", habit.ID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"habit": response}, h.logger)
}

func (h *HabitHandler) GetHabitsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	habits, err := h.getHabitsByUserUsecase.Execute(r.Context(), userID)
	if err != nil {
		h.logger.Printf("Error getting habits for user %s: %v", userID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		return
	}

	var responses []dto.HabitResponseDTO
	for _, habit := range habits {
		response, err := toHabitResponseDTO(habit)
		if err != nil {
			h.logger.Printf("Failed to map habit to response DTO for habit %s: %v", habit.ID, err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
			return
		}
		responses = append(responses, response)
	}

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"habits": responses}, h.logger)
}

func (h *HabitHandler) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	habitID := r.PathValue("habitID")

	var req dto.UpdateHabitDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Failed to decode request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_request_format"}, h.logger)
		return
	}

	if err := h.v.Struct(&req); err != nil {
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "validation_failed"}, err, h.logger)
		return
	}

	habit, err := h.updateHabitUsecase.Execute(r.Context(), habitID, userID, req)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrInvalidInput:
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_input"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "not_found"}, h.logger)
		default:
			h.logger.Printf("Error updating habit: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	response, err := toHabitResponseDTO(habit)
	if err != nil {
		h.logger.Printf("Failed to map habit to response DTO for habit %s: %v", habit.ID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		return
	}

	h.logger.Printf("Habit updated successfully. HabitID: %s, UserID: %s", habit.ID, userID)
	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"habit": response}, h.logger)
}

func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	habitID := r.PathValue("habitID")

	err = h.deleteHabitUsecase.Execute(r.Context(), habitID, userID)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "habit not found"}, h.logger)
		default:
			h.logger.Printf("Error deleting habit: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	h.logger.Printf("Habit deleted successfully. HabitID: %s, UserID: %s", habitID, userID)
	w.WriteHeader(http.StatusNoContent)
}

// toHabitResponseDTO converts an entity.Habit to a dto.HabitResponseDTO
func toHabitResponseDTO(habit *entity.Habit) (dto.HabitResponseDTO, error) {
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

	if habit.TargetDays != nil {
		targetDaysStr, err := dto.ConvertTargetDaysToJSON(habit.TargetDays)
		if err != nil {
			return dto.HabitResponseDTO{}, err
		}
		response.TargetDays = targetDaysStr
	}

	return response, nil
}
