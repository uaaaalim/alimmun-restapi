package handlers

import (
	"net/http"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/response"
)

func Home(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, struct {
		Message string `json:"message"`
		Time    string `json:"time"`
	}{
		Message: "Hello World",
		Time:    time.Now().Format(time.RFC3339),
	})
}
