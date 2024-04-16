package utils

import (
	"fmt"
	"net/http"
	"encoding/json"
)


type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{StatusCode: statusCode, Message: message}
}


func JsonHttpResponse(w http.ResponseWriter, statusCode int, data map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Log the error and return a generic error message if JSON marshaling fails
		fmt.Fprintf(w, `{"error": "Internal server error"}`)
		return
	}
	w.Write(jsonData)
}

