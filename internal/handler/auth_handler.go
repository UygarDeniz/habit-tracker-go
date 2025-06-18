package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uygardeniz/habit-tracker/internal/config"
	authUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/auth"
	userUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/user"
	"github.com/uygardeniz/habit-tracker/internal/utils"
)

var (
	frontendURL       = config.GetFrontendURL()
	googleOauthConfig = config.GetGoogleOauthConfig()
	oauthStateString  = "random"
)

type AuthHandler struct {
	logger                           *log.Logger
	loginOrRegisterGoogleUserUsecase *authUsecase.LoginOrRegisterGoogleUserUsecase
	getUserByIDUsecase               *userUsecase.GetUserByIDUsecase
}

func NewAuthHandler(logger *log.Logger, loginOrRegisterGoogleUserUsecase *authUsecase.LoginOrRegisterGoogleUserUsecase, getUserByIDUsecase *userUsecase.GetUserByIDUsecase) *AuthHandler {
	return &AuthHandler{
		logger:                           logger,
		loginOrRegisterGoogleUserUsecase: loginOrRegisterGoogleUserUsecase,
		getUserByIDUsecase:               getUserByIDUsecase,
	}
}

func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		h.logger.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		h.logger.Printf("no authorization code received")
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "no_authorization_code"}, h.logger)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.logger.Printf("code exchange failed with '%s'\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "code_exchange_failed"}, h.logger)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		h.logger.Printf("failed getting user info: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_get_user_info"}, h.logger)
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		h.logger.Printf("failed reading response body: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_read_user_info"}, h.logger)
		return
	}

	h.logger.Printf("user info: %s\n", string(contents))
	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		h.logger.Printf("failed to unmarshal user info: %s\n", err.Error())
		http.Redirect(w, r, fmt.Sprintf("%s/auth?auth_error=Internal server error", frontendURL), http.StatusTemporaryRedirect)
		return
	}

	user, err := h.loginOrRegisterGoogleUserUsecase.Execute(r.Context(), userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		h.logger.Printf("failed to login or register user: %s\n", err.Error())
		http.Redirect(w, r, fmt.Sprintf("%s/auth?auth_error=Internal server error", frontendURL), http.StatusTemporaryRedirect)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.Printf("failed to generate refresh token: %s\n", err.Error())
		http.Redirect(w, r, fmt.Sprintf("%s/auth?auth_error=Internal server error", frontendURL), http.StatusTemporaryRedirect)
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/api/auth",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &refreshTokenCookie)

	h.logger.Printf("Authentication successful. UserID: %s. Redirecting to frontend.", user.ID)
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "refresh token not found"}, h.logger)
		return
	}

	refreshTokenString := cookie.Value
	refreshToken, err := utils.ValidateToken(refreshTokenString, os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil || !refreshToken.Valid {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "invalid refresh token"}, h.logger)
		return
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to parse token claims"}, h.logger)
		return
	}

	userID := claims["sub"].(string)

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to generate new access token"}, h.logger)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"access_token": newAccessToken}, h.logger)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name: "refresh_token",

		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		Path:     "/api/auth",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"message": "logged out successfully"}, h.logger)
}

func (h *AuthHandler) HandleGetUserAndAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "refresh token not found"}, h.logger)
		return
	}

	refreshTokenString := cookie.Value
	refreshToken, err := utils.ValidateToken(refreshTokenString, os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil || !refreshToken.Valid {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{"error": "invalid refresh token"}, h.logger)
		return
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to parse token claims"}, h.logger)
		return
	}

	userID := claims["sub"].(string)

	user, err := h.getUserByIDUsecase.Execute(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to get user"}, h.logger)
		return
	}

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to generate access token"}, h.logger)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userID)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed to generate refresh token"}, h.logger)
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/api/auth",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &refreshTokenCookie)

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"user": user, "access_token": accessToken}, h.logger)
}
