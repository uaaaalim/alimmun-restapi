package auth

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"uuid"

	"github.com/uaaaalim/alimmun-restapi/internal/repositories"
	"github.com/uaaaalim/alimmun-restapi/internal/response"
	"github.com/uaaaalim/alimmun-restapi/internal/utils"
)

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// POST example.com/auth/register
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	rl, err := h.services.Auth.IsRegisterIPRateLimited(
		r.Context(), w, r.Context().Value("real-ip").(string),
	)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	if rl {
		return
	}

	var input registerRequest
	if err := utils.ReadJSON(w, r, &input, true); err != nil {
		return
	}

	// checking values
	fields := make(map[string]string)

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = normalizeName(input.Name)

	address, err := mail.ParseAddress(input.Email)
	switch {
	case input.Email == "":
		fields["email"] = "required"
	case len(input.Email) > 254:
		fields["email"] = "maximum length is 254 bytes"
	case err != nil:
		fields["email"] = "invalid email"
	case address.Address != input.Email:
		fields["email"] = "invalid email"
	}

	switch {
	case input.Name == "":
		fields["name"] = "required"
	case utf8.RuneCountInString(input.Name) > 20:
		fields["name"] = "maximum length is 20 characters"
	}

	if input.Password == "" {
		fields["password"] = "required"
	} else {
		var hasDigit, hasUpper, hasSpecial, hasInvalid bool
		for _, char := range input.Password {
			switch {
			case char >= '0' && char <= '9':
				hasDigit = true
			case char >= 'A' && char <= 'Z':
				hasUpper = true
			case char >= 'a' && char <= 'z':
			case char >= '!' && char <= '~':
				hasSpecial = true
			default:
				hasInvalid = true
			}
		}

		switch {
		case utf8.RuneCountInString(input.Password) < 8:
			fields["password"] = "minimum length is 8 characters"
		case len(input.Password) > 1024:
			fields["password"] = "maximum length is 1024 bytes"
		case hasInvalid:
			fields["password"] = "only English letters, digits and ASCII special characters are allowed"
		case !hasDigit || !hasUpper || !hasSpecial:
			fields["password"] = "must contain at least one digit, one uppercase letter and one special character"
		}
	}

	if len(fields) > 0 {
		response.Fail(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Check the provided fields",
			fields,
		)
		return
	}

	rl, err = h.services.Auth.IsRegisterEmailRateLimited(
		r.Context(), w, input.Email,
	)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	if rl {
		return
	}

	user, err := h.services.Auth.Register(r.Context(), input.Email, input.Name, input.Password)
	if err != nil {
		if errors.Is(err, repositories.ErrEmailTaken) {
			response.Fail(w, http.StatusConflict, "EMAIL_TAKEN", "Email is already registered", nil)
			return
		}
		response.InternalError(w, err)
		return
	}

	res := struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		IsActivated bool      `json:"is_activated"`
		CreatedAt   time.Time `json:"created_at"`
	}{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		IsActivated: user.IsActivated,
		CreatedAt:   user.CreatedAt,
	}

	response.OK(w, http.StatusCreated, res)
}

func normalizeName(name string) string {
	chars := []rune(strings.ToLower(strings.TrimSpace(name)))
	if len(chars) > 0 {
		chars[0] = unicode.ToUpper(chars[0])
	}
	return string(chars)
}
