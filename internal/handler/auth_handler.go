package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	userUsecase "github.com/uygardeniz/habit-tracker/internal/usecases/user"
	"github.com/uygardeniz/habit-tracker/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func getGoogleOauthConfig() *oauth2.Config {
	if googleOauthConfig == nil {
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	}
	return googleOauthConfig
}

type AuthHandler struct {
	logger                           *log.Logger
	loginOrRegisterGoogleUserUsecase *userUsecase.LoginOrRegisterGoogleUserUsecase
}

func NewAuthHandler(logger *log.Logger, loginOrRegisterGoogleUserUsecase *userUsecase.LoginOrRegisterGoogleUserUsecase) *AuthHandler {
	return &AuthHandler{logger: logger, loginOrRegisterGoogleUserUsecase: loginOrRegisterGoogleUserUsecase}
}

func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := getGoogleOauthConfig()
	url := config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	config := getGoogleOauthConfig()

	state := r.FormValue("state")
	if state != oauthStateString {
		h.logger.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "invalid_oauth_state"}, h.logger)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		h.logger.Printf("no authorization code received")
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{"error": "no_authorization_code"}, h.logger)
		return
	}

	token, err := config.Exchange(context.Background(), code)
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
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_parse_user_info"}, h.logger)
		return
	}

	user, err := h.loginOrRegisterGoogleUserUsecase.Execute(r.Context(), userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		h.logger.Printf("failed to login or register user: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_process_user"}, h.logger)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		h.logger.Printf("failed to generate access token: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_generate_access_token"}, h.logger)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.Printf("failed to generate refresh token: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{"error": "failed_to_generate_refresh_token"}, h.logger)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/api/auth",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	userData := dto.AuthResponse{
		AccessToken: accessToken,
		User: dto.UserData{
			ID:      user.ID,
			Email:   user.Email,
			Name:    user.Name,
			Picture: user.Picture,
		},
		Message: "authentication_successful",
	}

	h.logger.Printf("Authentication successful. UserID: %s", user.ID)
	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"data": userData}, h.logger)
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
	// Clear refresh token cookie
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, // Set to false for localhost development
		Path:     "/api/auth",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	utils.WriteJSON(w, http.StatusOK, utils.APIResponse{"message": "logged out successfully"}, h.logger)
}
