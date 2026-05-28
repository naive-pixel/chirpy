package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
)

type request struct {
	Body string `json:"body"`
}
type errResponse struct {
	Error string `json:"error"`
}
type validResponse struct {
	Valid bool `json:"valid"`
}

type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error while reading []byte stream: %v", err)
		w.WriteHeader(500)
		return

	}

	r := request{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		respondWithError(w, 500, "error while converting request from []byte to struct", nil)
		return

	}
	if len(r.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	cleanedChirp := censorChirp(r.Body)
	cleanJson := returnVals{
		CleanedBody: cleanedChirp,
	}
	respondWithJSON(w, 200,
		cleanJson)

}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respBody := errResponse{
		Error: msg,
	}
	respondWithJSON(w, code, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)

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
