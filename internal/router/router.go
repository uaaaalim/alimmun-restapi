package router

import "net/http"

func CreateRouter() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
