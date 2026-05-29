package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/naive-pixel/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (apiCfg *apiConfig) handlerCreateChrip(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	params := parameters{}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while reading []byte stream", err)
		return

	}
	err = json.Unmarshal(data, &params)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "error while converting request from []byte to struct", err)
		return

	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	cleanedChirp := censorChirp(params.Body)
	chripParam := database.CreateChripParams{Body: cleanedChirp,
		UserID: params.UserId}

	createdChrip, err := apiCfg.db.CreateChrip(req.Context(), chripParam)
	if err != nil {
		respondWithError(w, 500, "error while creating chirp", err)
		return
	}
	returnChirp := Chirp{ID: createdChrip.ID,
		CreatedAt: createdChrip.CreatedAt,
		UpdatedAt: createdChrip.UpdatedAt,
		Body:      createdChrip.Body,
		UserID:    createdChrip.UserID,
	}

	respondWithJSON(w, http.StatusCreated,
		returnChirp)

}

func (apiCfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	allChirps, err := apiCfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "error while getting all chirps", err)
		return
	}
	returnChirps := []Chirp{}
	for _, chrip := range allChirps {
		returnChirps = append(returnChirps, Chirp{ID: chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body:      chrip.Body,
			UserID:    chrip.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK,
		returnChirps)

}
func (apiCfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpID")
	if chirpId == "" {
		log.Println("Missing chirp URL")
		respondWithError(w, http.StatusInternalServerError, "Missing chirp URL", nil)
	}
	chripIDUUID, err := uuid.Parse(chirpId)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "Could not convert chripID from string to UUID", err)
		return
	}
	chrip, err := apiCfg.db.GetChirp(req.Context(), chripIDUUID)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusNotFound, "error while getting chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK,
		Chirp{ID: chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body:      chrip.Body,
			UserID:    chrip.UserID,
		})

}

func censorChirp(chirp string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(chirp)
	var cleanChirp []string
	for _, word := range words {
		if slices.Contains(bannedWords, strings.ToLower(word)) {
			cleanChirp = append(cleanChirp, "****")
			continue
		}
		cleanChirp = append(cleanChirp, word)
	}

	return strings.Join(cleanChirp, " ")
}
