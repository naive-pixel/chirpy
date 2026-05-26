package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("error while reading []byte stream: %w", err)
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
	respondWithJSON(w, 200,
		validResponse{Valid: true})

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
