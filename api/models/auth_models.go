package models

// LoginRequest defines the structure for the JSON body expected in login requests.

type LoginRequest struct {
	Username string `json:"username"`

	Password string `json:"password"`
}

// LoginResponse defines the structure for the JSON response on successful login.

type LoginResponse struct {
	Token string `json:"token"`

	Message string `json:"message,omitempty"`
}
