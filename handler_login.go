package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/naive-pixel/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type userLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading incoming request", err)
		return
	}
	userParams := userLogin{}
	err = json.Unmarshal(data, &userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error unmarshaling request", err)
		return
	}
	user, err := apiCfg.db.GetPasswordByUserID(req.Context(), userParams.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user by email", err)
	}
	log.Printf("userParams.Password: %s", userParams.Password)
	log.Printf("user.HashedPassword: %s", user.HashedPassword)

	matched, err := auth.CheckPasswordHash(userParams.Password, user.HashedPassword)

	if !matched || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	userResponse := User{ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
}
