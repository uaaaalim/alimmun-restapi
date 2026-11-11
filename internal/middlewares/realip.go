package middlewares

import (
	"context"
	"net/http"
	"net/netip"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: X-Real-IP is set by parent proxy server!
		ip, err := netip.ParseAddr(r.Header.Get("X-Real-IP"))
		if err != nil {
			response.Fail(w, http.StatusBadRequest, "INVALID_HEADERS",
				"Failed to resolve origin remove client IP", nil)
			return
		}

		ctx := context.WithValue(r.Context(), "real-ip", ip.Unmap().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
