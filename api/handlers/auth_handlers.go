package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"finala/api/auth"
	"finala/api/models"
	"finala/config"
	"finala/serverutil"
)

// LoginHandler handles user login requests.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		serverutil.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		serverutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		serverutil.RespondWithError(w, http.StatusBadRequest, "Username is required")
		return
	}

	if req.Password == "" {
		serverutil.RespondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if req.Username == config.AppCredentials.Username && req.Password == config.AppCredentials.Password {
		tokenString, _, err := auth.GenerateJWT(req.Username)
		if err != nil {
			log.Printf("ERROR: Generating JWT: %v", err)
			serverutil.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
			return
		}

		serverutil.RespondWithJSON(w, http.StatusOK, models.LoginResponse{
			Token:   tokenString,
			Message: "Login successful",
		})
	} else {
		log.Printf("WARN: Failed login attempt for username: %s", req.Username)
		serverutil.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
	}
}
