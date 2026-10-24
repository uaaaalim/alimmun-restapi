package router

import (
	"net/http"

	"github.com/uaaaalim/alimmun-restapi/internal/handlers"
)

func CreateRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handlers.Home)

	return mux
}
