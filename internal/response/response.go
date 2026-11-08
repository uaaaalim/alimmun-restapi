package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Response struct {
	Data  any        `json:"data,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

func write(w http.ResponseWriter, status int, response Response) {
	payload, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(
			`{"error":{"code":"internal_error","message":"An unknown error has happen. We've been notified and trying to fix it'"}}`,
		))
		slog.Error("write json response", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(append(payload, '\n'))
	slog.Error("write response", "error", err)
}

func OK(w http.ResponseWriter, status int, data any) {
	write(w, status, Response{Data: data})
}

func Fail(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	fields map[string]string,
) {
	write(w, status, Response{
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}
