package config

import (
	"os"

	"golang.org/x/oauth2"
)

var (
	googleOauthConfig *oauth2.Config
	frontendURL       string
)

func GetGoogleOauthConfig() *oauth2.Config {
	if googleOauthConfig == nil {
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		}
	}

	return googleOauthConfig
}

func GetFrontendURL() string {
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
	}
	return frontendURL
}
