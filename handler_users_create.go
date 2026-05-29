package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/naive-pixel/chirpy/internal/auth"
	"github.com/naive-pixel/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (apiCfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading incoming request", err)
		return
	}
	params := parameter{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error unmarshaling incoming request", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while hashing password", err)
		return
	}

	userParams := database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword}
	user, err := apiCfg.db.CreateUser(req.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	userResponse := User{ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, userResponse)

}

/*
type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}
*/
