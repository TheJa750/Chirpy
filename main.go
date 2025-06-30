package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	svr := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.HandleFunc("GET /api/healthz", healthzHandler)

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /admin/metrics", cfg.getFileserverHitsHandler)

	mux.HandleFunc("POST /admin/reset", cfg.resetFileserverHitsHandler)

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	svr.ListenAndServe()

}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		msg := JsonError{
			Message: "Something went wrong",
		}
		dat, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		msg := JsonError{
			Message: "Chirp is too long",
		}
		dat, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	msg := cleanChirpBody(params.Body)
	dat, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func cleanChirpBody(s string) CleanedChirpBody {
	words := strings.Fields(s)

	for i, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}

	cleanBody := strings.Join(words, " ")

	return CleanedChirpBody{cleanBody}
}
