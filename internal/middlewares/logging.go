package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *ResponseWriter) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	return r.ResponseWriter.Write(b)
}

func LoggerMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		writer := ResponseWriter{ResponseWriter: w}

		next.ServeHTTP(&writer, r)

		duration := time.Since(now)
		level := slog.LevelDebug
		status := writer.status
		if status == 0 {
			status = http.StatusOK
		}

		switch {
		case status >= 500:
			level = slog.LevelError
		case duration > 500*time.Millisecond:
			level = slog.LevelWarn
		}

		log.Log(r.Context(), level, "HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration", duration,
		)
	})
}
