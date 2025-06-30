package main

import (
	"net/http"
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

	mux.HandleFunc("GET /healthz", healthzHandler)

	fileHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /metrics", cfg.getFileserverHitsHandler)

	mux.HandleFunc("POST /reset", cfg.resetFileserverHitsHandler)

	svr.ListenAndServe()

}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
