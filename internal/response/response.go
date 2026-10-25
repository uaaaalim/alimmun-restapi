package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"uuid"
)

const internalErrorMessage = "An unknown error has happen. We've been notified about it and we're trying to fix it"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ErrorID string `json:"id,omitempty"`
}

type Response struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body Response) {
	data, err := json.Marshal(body)
	if err != nil {
		internalError(w, context.Background(), fmt.Errorf("encode JSON response: %w", err), "original_status", status)
		return
	}
	writeJSON(w, status, data)
}

func writeJSON(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(append(data, '\n')); err != nil {
		slog.Error("failed to write HTTP response",
			"error", err,
			"status", status,
			"error_id", w.Header().Get("X-Error-ID"),
		)
	}
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

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	internalError(w, r.Context(), err, "method", r.Method, "path", r.URL.Path)
}

func internalError(w http.ResponseWriter, ctx context.Context, err error, attrs ...any) {
	errorID := uuid.New().String()

	attrs = append(attrs, "error_id", errorID, "error", err)
	slog.ErrorContext(ctx, "internal server error", attrs...)

	w.Header().Set("X-Error-ID", errorID)
	w.Header().Set("Cache-Control", "no-store")

	// Both values are ASCII strings controlled by the server. Build the fallback
	// directly so an encoding failure never calls JSON recursively.
	data := []byte(fmt.Sprintf(
		`{"success":false,"data":null,"error":{"code":"INTERNAL_ERROR","message":%q,"id":%q}}`,
		internalErrorMessage, errorID,
	))
	writeJSON(w, http.StatusInternalServerError, data)
}
