package security

import "github.com/alexedwards/argon2id"

func HashPassword(value string) (string, error) {
	return argon2id.CreateHash(value, &argon2id.Params{
		Memory:      7 * 1024, // 7 МиБ
		Iterations:  5,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
}

func VerifyPassword(value, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(value, hash)
}
