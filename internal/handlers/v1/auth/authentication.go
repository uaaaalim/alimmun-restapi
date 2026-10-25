package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/services"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_BODY", "Incorrect JSON type has been given")
		return
	}
	body.Email = strings.TrimSpace(body.Email)

	address, err := mail.ParseAddress(body.Email)
	if err != nil || address.Address != body.Email || len(body.Email) > 255 {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_EMAIL", "Specified email is incorrect")
		return
	} else if len(body.Password) > 32 || len(body.Password) < 8 {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_PASSWORD", "\"Password\" is too long or too short (requires from 8 up to 32 chars)")
		return
	} else if len(body.Name) > 100 || len(body.Name) < 2 {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_NAME", "\"Name\" is too long or too short (requires from 8 up to 32 chars)")
		return
	}

	user, err := h.AuthService.Register(r.Context(), body.Name, body.Email, body.Password)
	if errors.Is(err, services.AuthUserAlreadyExists) {
		response.Fault(w, http.StatusConflict,
			"USER_ALREADY_EXISTS", "Пользователь с таким email уже существует")
		return
	} else if err != nil {
		response.Fault(w, http.StatusInternalServerError,
			"INTERNAL_ERROR", "Не удалось создать пользователя")
		return
	}

	data := struct {
		ID        int64           `json:"id"`
		Name      string          `json:"name"`
		Email     string          `json:"email"`
		Role      models.UserRole `json:"role"`
		CreatedAt time.Time       `json:"created_at"`
	}{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	response.Success(w, http.StatusOK, data)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_BODY", "Incorrect JSON type has been given")
		return
	}
	body.Email = strings.TrimSpace(body.Email)

	address, err := mail.ParseAddress(body.Email)
	if err != nil || address.Address != body.Email || len(body.Email) > 255 {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_EMAIL", "Specified email is incorrect")
		return
	} else if len(body.Password) > 32 || len(body.Password) < 8 {
		response.Fault(w, http.StatusBadRequest,
			"INVALID_PASSWORD", "Password is too long or too short (requires from 8 up to 32 chars)")
		return
	}

	user, err := h.AuthService.Login(r.Context(), body.Email, body.Password)

	if errors.Is(err, services.AuthInvalidCredentials) {
		response.Fault(w, http.StatusUnauthorized,
			"INVALID_CREDENTIALS", "Invalid email or password")
		return
	} else if err != nil {
		response.Fault(w, http.StatusInternalServerError,
			"INTERNAL_ERROR", "Failed to login")
		return
	}

	data := struct {
		ID        int64           `json:"id"`
		Name      string          `json:"name"`
		Email     string          `json:"email"`
		Role      models.UserRole `json:"role"`
		CreatedAt time.Time       `json:"created_at"`
	}{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	response.Success(w, http.StatusOK, data)
}
