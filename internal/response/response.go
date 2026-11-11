package response

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"
)

type ErrorBody struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	RetryAfter int64             `json:"retry_after,omitempty"`
	Fields     map[string]string `json:"fields,omitempty"`
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
			`{"error":{"code":"internal_error","message":"An unknown error has happened. We've been notified and trying to fix it"}}`,
		))
		slog.Error("write json response", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(append(payload, '\n')); err != nil {
		slog.Error("write response", "error", err)
	}
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

func InternalError(
	w http.ResponseWriter,
	err error,
) {
	slog.Error("internal error", "error", err)
	write(w, http.StatusInternalServerError, Response{
		Error: &ErrorBody{
			Code:    "INTERNAL_ERROR",
			Message: "An unknown error has happened. We've been notified and trying to fix it",
		},
	})
}

func RateLimited(
	w http.ResponseWriter,
	retryAfter time.Duration,
) {
	seconds := max(int64(math.Ceil(retryAfter.Seconds())), 1)

	w.Header().Set("Retry-After", strconv.FormatInt(seconds, 10))

	write(w, http.StatusTooManyRequests, Response{
		Error: &ErrorBody{
			Code:       "TOO_MANY_REQUESTS",
			RetryAfter: seconds,
			Message:    "You have exceeded the rate limit",
		},
	})
}
