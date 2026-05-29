package main

import (
	"fmt"
	"net/http"
)

func (acfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

/*
func (acfg *apiConfig) getNumberOfRequestsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", acfg.fileserverHits.Load())

}*/

func (acfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	acfg.fileserverHits.Store(0)
	acfg.db.Reset(req.Context())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset users table and hit counter."))
}

func (acfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	welcomePage := []byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		acfg.fileserverHits.Load()))
	w.Write(welcomePage)
}
