package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashed_pw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return hashed_pw, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {

	matched, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error while comparing password and hash: %w", err)
	}

	return matched, nil
}
