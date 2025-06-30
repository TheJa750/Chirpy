package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

const adminMetrics = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		a.fileserverHits.Add(1)

		next.ServeHTTP(w, req)
	})

}

func (a *apiConfig) getFileserverHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(adminMetrics, a.fileserverHits.Load())
	w.Write([]byte(body))
}

func (a *apiConfig) resetFileserverHitsHandler(w http.ResponseWriter, req *http.Request) {
	a.fileserverHits = atomic.Int32{}
	w.WriteHeader(http.StatusOK)
}

type JsonError struct {
	Message string `json:"error"`
}

type CleanedChirpBody struct {
	Body string `json:"cleaned_body"`
}
