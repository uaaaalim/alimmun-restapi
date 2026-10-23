package response

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, Response{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func Fault(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, Response{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
