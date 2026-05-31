package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String()},
	)
	log.Printf("userID: %s", userID)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		fmt.Printf("error while signing token: %v", err)
	}
	return signedToken, nil
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		fmt.Printf("error while parsing token: %v\n", err)
		return uuid.Nil, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Printf("error while getting subject: %v", err)
		return uuid.Nil, err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	parsedSubject, err := uuid.Parse(subject)
	if err != nil {
		fmt.Printf("error while decoding subject: %v", err)
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	log.Printf("parsedSubject: %s", parsedSubject)
	return parsedSubject, nil
}
