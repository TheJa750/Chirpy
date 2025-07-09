package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/TheJa750/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Importing pq for PostgreSQL driver
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	mux := http.NewServeMux()

	svr := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
	}

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	//metrics handlers
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))
	mux.HandleFunc("GET /admin/metrics", cfg.getFileserverHitsHandler)
	//mux.HandleFunc("POST /admin/reset", cfg.resetFileserverHitsHandler)

	//api handlers
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", cfg.postChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpByIDHandler)
	mux.HandleFunc("POST /api/login", cfg.loginUserHandler)

	//dev handlers
	mux.HandleFunc("POST /admin/reset", cfg.resetUsersHandler)

	svr.ListenAndServe()

}
