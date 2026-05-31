package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestCheckJWT(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	userUUID1 := uuid.New()
	tokenSecret1 := "mysecret123"
	expiresIn1 := time.Duration(time.Second * 10)

	signedToken1, _ := MakeJWT(userUUID1, tokenSecret1, expiresIn1)

	//userUUID2 := uuid.New()

	tests := []struct {
		name        string
		userUUID    uuid.UUID
		tokenSecret string
		signedToken string
		wantErr     bool
		match       bool
	}{
		{
			name:        "Correct decoding",
			tokenSecret: tokenSecret1,
			signedToken: signedToken1,
			userUUID:    userUUID1,
			wantErr:     false,
		},
		{
			name:        "Expired token",
			tokenSecret: tokenSecret1,
			signedToken: signedToken1,
			userUUID:    userUUID1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Expired token" {
				time.Sleep(time.Second * 12)
			}
			decodedUserUUID, err := ValidateJWT(tt.signedToken, tt.tokenSecret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v", err)
			}

			if (err != nil) != tt.wantErr && tt.name == "Expired token" {
				t.Errorf("ValidateJWT() error = %v", err)
			}
			if (decodedUserUUID != tt.userUUID) && !tt.wantErr {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.userUUID, decodedUserUUID)
			}

		})
	}
}
