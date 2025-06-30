package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		a.fileserverHits.Add(1)

		next.ServeHTTP(w, req)
	})

}

func (a *apiConfig) getFileserverHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("Hits: %v", a.fileserverHits.Load())
	w.Write([]byte(body))
}

func (a *apiConfig) resetFileserverHitsHandler(w http.ResponseWriter, req *http.Request) {
	a.fileserverHits = atomic.Int32{}
	w.WriteHeader(http.StatusOK)
}
