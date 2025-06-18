package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/middleware"
	completionUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/completion"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

type CompletionHandler struct {
	createCompletionUsecase *completionUsecase.CreateCompletionUsecase
	getCompletionUsecase    *completionUsecase.GetCompletionUsecase
	getCompletionsUsecase   *completionUsecase.GetCompletionsUsecase
	updateCompletionUsecase *completionUsecase.UpdateCompletionUsecase
	deleteCompletionUsecase *completionUsecase.DeleteCompletionUsecase
	logger                  *log.Logger
	v                       *validator.Validate
}

func NewCompletionHandler(
	createCompletionUsecase *completionUsecase.CreateCompletionUsecase,
	getCompletionUsecase *completionUsecase.GetCompletionUsecase,
	getCompletionsUsecase *completionUsecase.GetCompletionsUsecase,
	updateCompletionUsecase *completionUsecase.UpdateCompletionUsecase,
	deleteCompletionUsecase *completionUsecase.DeleteCompletionUsecase,
	logger *log.Logger,
	v *validator.Validate,
) *CompletionHandler {
	return &CompletionHandler{
		createCompletionUsecase: createCompletionUsecase,
		getCompletionUsecase:    getCompletionUsecase,
		getCompletionsUsecase:   getCompletionsUsecase,
		updateCompletionUsecase: updateCompletionUsecase,
		deleteCompletionUsecase: deleteCompletionUsecase,
		logger:                  logger,
		v:                       v,
	}
}

func (h *CompletionHandler) CreateCompletion(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	habitID := r.PathValue("habitID")
	if habitID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "habit ID is required"}, h.logger)
		return
	}

	var req dto.CreateCompletionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Failed to decode request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_request_format"}, h.logger)
		return
	}

	if err := h.v.Struct(&req); err != nil {
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "validation_failed"}, err, h.logger)
		return
	}

	completion, err := h.createCompletionUsecase.Execute(r.Context(), habitID, userID, req)
	if err != nil {
		switch err {
		case apperrors.ErrAlreadyExists:
			utils.WriteJSON(w, http.StatusConflict, utils.APIResponse{"error": "completion already exists for this date"}, h.logger)
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrInvalidInput:
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_input"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "habit not found"}, h.logger)
		default:
			h.logger.Printf("Error creating completion: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	response := toCompletionResponseDTO(completion)
	h.logger.Printf("Completion created successfully. CompletionID: %s, HabitID: %s, UserID: %s", completion.ID, habitID, userID)
	utils.WriteJSON(w, http.StatusCreated, utils.APIResponse{"completion": response}, h.logger)
}

func (h *CompletionHandler) GetCompletion(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	completionID := r.PathValue("completionID")

	completion, err := h.getCompletionUsecase.Execute(r.Context(), completionID, userID)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "completion not found"}, h.logger)
		default:
			h.logger.Printf("Error getting completion: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	response := toCompletionResponseDTO(completion)
	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"completion": response}, h.logger)
}

func (h *CompletionHandler) GetCompletions(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	query := dto.GetCompletionsQueryDTO{}

	if habitID := r.URL.Query().Get("habit_id"); habitID != "" {
		query.HabitID = &habitID
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		query.StartDate = &startDate
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		query.EndDate = &endDate
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.Limit = &limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			query.Offset = &offset
		}
	}

	if err := h.v.Struct(&query); err != nil {
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "validation_failed"}, err, h.logger)
		return
	}

	completions, err := h.getCompletionsUsecase.Execute(r.Context(), userID, query)
	if err != nil {
		h.logger.Printf("Error getting completions for user %s: %v", userID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		return
	}

	var responses []dto.CompletionResponseDTO
	for _, completion := range completions {
		responses = append(responses, toCompletionResponseDTO(completion))
	}

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"completions": responses}, h.logger)
}

func (h *CompletionHandler) UpdateCompletion(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	completionID := r.PathValue("completionID")

	var req dto.UpdateCompletionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Failed to decode request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_request_format"}, h.logger)
		return
	}

	if err := h.v.Struct(&req); err != nil {
		utils.WriteValidationErrorResponse(w, http.StatusBadRequest, utils.APIResponse{"error": "validation_failed"}, err, h.logger)
		return
	}

	completion, err := h.updateCompletionUsecase.Execute(r.Context(), completionID, userID, req)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrInvalidInput:
			utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_input"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "completion not found"}, h.logger)
		default:
			h.logger.Printf("Error updating completion: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	response := toCompletionResponseDTO(completion)
	h.logger.Printf("Completion updated successfully. CompletionID: %s, UserID: %s", completion.ID, userID)
	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"completion": response}, h.logger)
}

func (h *CompletionHandler) DeleteCompletion(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.logger.Printf("Failed to get user ID from context: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "unauthorized"}, h.logger)
		return
	}

	completionID := r.PathValue("completionID")

	err = h.deleteCompletionUsecase.Execute(r.Context(), completionID, userID)
	if err != nil {
		switch err {
		case apperrors.ErrForbidden:
			utils.WriteJSON(w, http.StatusForbidden, utils.APIResponse{"error": "forbidden"}, h.logger)
		case apperrors.ErrNotFound:
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{"error": "completion not found"}, h.logger)
		default:
			h.logger.Printf("Error deleting completion: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "internal_server_error"}, h.logger)
		}
		return
	}

	h.logger.Printf("Completion deleted successfully. CompletionID: %s, UserID: %s", completionID, userID)
	utils.WriteJSON(w, http.StatusNoContent, nil, h.logger)
}

func toCompletionResponseDTO(completion *entity.HabitCompletion) dto.CompletionResponseDTO {
	return dto.CompletionResponseDTO{
		ID:             completion.ID,
		HabitID:        completion.HabitID,
		UserID:         completion.UserID,
		CompletedAt:    completion.CompletedAt,
		CompletionDate: completion.CompletionDate,
		Count:          completion.Count,
		Notes:          completion.Notes,
		CreatedAt:      completion.CreatedAt,
	}
}
