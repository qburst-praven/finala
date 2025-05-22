package handlers_test

import (
	"bytes"

	"encoding/json"

	"net/http"

	"net/http/httptest"

	"strings"

	"testing"

	"finala/api/handlers"

	"finala/api/models"

	"finala/config"
	// "finala/serverutil" // Not directly used here but handlers depend on it
)

func TestLoginHandler(t *testing.T) {

	// Store original AppCredentials and JWT secret if any, and restore after test

	originalAppCreds := config.AppCredentials

	// originalJWTSecret := auth.jwtSecretKey // If jwtSecretKey was public, otherwise can't easily mock it without DI

	defer func() {

		config.AppCredentials = originalAppCreds

		// auth.jwtSecretKey = originalJWTSecret

	}()

	// Setup mock credentials for testing

	config.AppCredentials = config.AuthCredentialsConfig{

		Username: "testuser",

		Password: "testpassword",
	}

	tests := []struct {
		name string

		method string

		payload interface{}

		expectedStatusCode int

		expectToken bool

		expectErrorMsg string
	}{

		{

			name: "Successful Login",

			method: http.MethodPost,

			payload: models.LoginRequest{

				Username: "testuser",

				Password: "testpassword",
			},

			expectedStatusCode: http.StatusOK,

			expectToken: true,
		},

		{

			name: "Incorrect Password",

			method: http.MethodPost,

			payload: models.LoginRequest{

				Username: "testuser",

				Password: "wrongpassword",
			},

			expectedStatusCode: http.StatusUnauthorized,

			expectErrorMsg: "Invalid username or password",
		},

		{

			name: "Incorrect Username",

			method: http.MethodPost,

			payload: models.LoginRequest{

				Username: "wronguser",

				Password: "testpassword",
			},

			expectedStatusCode: http.StatusUnauthorized,

			expectErrorMsg: "Invalid username or password",
		},

		{

			name: "Method Not Allowed",

			method: http.MethodGet,

			expectedStatusCode: http.StatusMethodNotAllowed,

			expectErrorMsg: "Method not allowed",
		},

		{

			name: "Invalid Payload - Malformed JSON",

			method: http.MethodPost,

			payload: "not-json",

			expectedStatusCode: http.StatusBadRequest,

			expectErrorMsg: "Invalid request payload",
		},

		{

			name: "Missing Username",

			method: http.MethodPost,

			payload: models.LoginRequest{

				Password: "testpassword",
			},

			expectedStatusCode: http.StatusBadRequest,

			expectErrorMsg: "Username is required",
		},

		{

			name: "Missing Password",

			method: http.MethodPost,

			payload: models.LoginRequest{

				Username: "testuser",
			},

			expectedStatusCode: http.StatusBadRequest,

			expectErrorMsg: "Password is required",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			var reqBody []byte

			var err error

			if tt.payload != nil {

				if strPayload, ok := tt.payload.(string); ok {

					reqBody = []byte(strPayload)

				} else {

					reqBody, err = json.Marshal(tt.payload)

					if err != nil {

						t.Fatalf("Failed to marshal payload: %v", err)

					}

				}

			}

			req, err := http.NewRequest(tt.method, "/api/v1/auth/login", bytes.NewBuffer(reqBody))

			if err != nil {

				t.Fatal(err)

			}

			rr := httptest.NewRecorder()

			h := http.HandlerFunc(handlers.LoginHandler)

			h.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {

				t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",

					status, tt.expectedStatusCode, rr.Body.String())

			}

			if tt.expectToken {

				var resp models.LoginResponse

				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {

					t.Errorf("Failed to unmarshal response body: %v. Body: %s", err, rr.Body.String())

				}

				if resp.Token == "" {

					t.Errorf("Expected token in response, got none. Body: %s", rr.Body.String())

				}

				if !strings.Contains(resp.Message, "Login successful") {

					t.Errorf("Expected success message, got: %s", resp.Message)

				}

			} else if tt.expectErrorMsg != "" {

				var errResp struct { // Using anonymous struct as ErrorResponse is in serverutil

					Error string `json:"error"`
				}

				if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {

					t.Errorf("Failed to unmarshal error response body: %v. Body: %s", err, rr.Body.String())

				}

				if !strings.Contains(errResp.Error, tt.expectErrorMsg) {

					t.Errorf("Expected error message '%s', got '%s'. Body: %s", tt.expectErrorMsg, errResp.Error, rr.Body.String())

				}

			}

		})

	}

}
