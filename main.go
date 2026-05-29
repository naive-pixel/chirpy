package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/naive-pixel/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL) // "postgress" is registered by the postgres driver "pq"

	if err != nil {
		fmt.Printf("Error while connecting to DB: %v\n", err)
		return
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChrip)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetSingleChirp)

	mux.HandleFunc("POST /api/login/", apiCfg.handlerLogin)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}
