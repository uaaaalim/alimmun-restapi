package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			duration := time.Since(start)
			level := slog.LevelDebug
			if duration > time.Second {
				level = slog.LevelWarn
			}

			slog.Log(r.Context(), level, "http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
			)
		}()

		next.ServeHTTP(w, r)
	})
}
