package dto

// GoogleUserInfo represents the data received from Google OAuth
type GoogleUserInfo struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

// AuthResponse represents the response data after successful authentication
type AuthResponse struct {
	AccessToken string
	User        UserData
	Message     string
}

// UserData represents the user information in the response
type UserData struct {
	ID      string
	Email   string
	Name    string
	Picture string
}
