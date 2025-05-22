package serverutil

import (
	"encoding/json"

	"log"

	"net/http"
)

// ErrorResponse is a generic structure for JSON error responses.

type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithJSON sends a JSON response with a given status code and payload.

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {

		log.Printf("ERROR: Marshalling JSON response: %v, Payload: %+v", err, payload)

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)

		if _, writeErr := w.Write([]byte(`{"error":"Internal server error preparing response"}`)); writeErr != nil {

			log.Printf("ERROR: Writing fallback JSON error response: %v", writeErr)

		}

		return

	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	_, err = w.Write(response)

	if err != nil {

		log.Printf("ERROR: Writing JSON response: %v", err)

	}

}

// RespondWithError sends a JSON error response using the ErrorResponse struct.

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {

	RespondWithJSON(w, statusCode, ErrorResponse{Error: message})

}
