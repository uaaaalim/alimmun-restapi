package security

import "golang.org/x/crypto/bcrypt"

// HashPassword creates password hash
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPassword Return nil as error if password matches
func CheckPassword(hash, password string) error {
	if len(password) > 72 {
		return bcrypt.ErrPasswordTooLong
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
