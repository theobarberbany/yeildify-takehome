package common

import (
	"encoding/json"
	"net/http"
)

// The common package contains shared structures and functions

// Result is something that may be returned as a result
type Result interface {
	// a dummy method to limit the possible types that can
	// satisfy this interface
	isResult()
}

// Message is a struct containing the message field
type Message struct {
	Message      string  `json:"message"`
	RandomNumber float64 `json:"rand,omitempty"`
}

func (s Message) isResult() {}

// ErrorResponse is an error type
type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func (s ErrorResponse) isResult() {}

// WriteJSONResponse is a helper function that sets the 'Content-Type' header to
// 'application/json' and encodes a Message into a http.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, response Result, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
