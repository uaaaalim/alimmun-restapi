package auth

import "net/http"

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("{\"hello\": \"world!\"}"))
	if err != nil {
		return
	}
}
