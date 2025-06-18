package handler

import (
	"log"
	"net/http"

	userUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/user"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

type UserHandler struct {
	logger       *log.Logger
	getMeUsecase *userUsecase.GetMeUsecase
}

func NewUserHandler(logger *log.Logger, getMeUsecase *userUsecase.GetMeUsecase) *UserHandler {
	return &UserHandler{logger: logger, getMeUsecase: getMeUsecase}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.getMeUsecase.Execute(r.Context())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to find user"}, h.logger)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"user": user}, h.logger)
}
