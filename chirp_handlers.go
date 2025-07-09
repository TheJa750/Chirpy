package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/TheJa750/Chirpy/internal/auth"
	"github.com/TheJa750/Chirpy/internal/database"
	"github.com/google/uuid"
)

func validateChirp(body string) (CleanedChirpBody, JsonError) {
	if len(body) > 140 {
		msg := JsonError{
			Message: "Chirp is too long",
		}
		return CleanedChirpBody{}, msg
	}

	msg := cleanChirpBody(body)

	return msg, JsonError{}
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

func (a *apiConfig) postChirpHandler(w http.ResponseWriter, req *http.Request) {
	var chirpReq chirpRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		log.Printf("Error decoding chirp request: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cleanedChirp, errMsg := validateChirp(chirpReq.Body)
	if errMsg.Message != "" {
		log.Printf("Chirp validation error: %s", errMsg.Message)
		http.Error(w, errMsg.Message, http.StatusBadRequest)
		return
	}

	chirp, err := a.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp.Body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		UserID:    chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonChirp)
}

func (a *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	userIDStr := req.URL.Query().Get("author_id")
	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}
	var chirps []database.Chirp
	var err error
	if userIDStr == "" {
		chirps, err = a.dbQueries.GetChirps(req.Context())
		if err != nil {
			log.Printf("Error getting chirps: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("Invalid user ID: %s", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		chirps, err = a.dbQueries.GetChirpsByUserID(req.Context(), userID)
		if err != nil {
			log.Printf("Error getting chirps by user ID: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	jsonChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		jsonChirps[i] = Chirp{
			ID:        chirp.ID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			UserID:    chirp.UserID,
		}
	}

	if sortOrder == "desc" {
		sort.Slice(jsonChirps, func(i, j int) bool {
			return jsonChirps[i].CreatedAt.After(jsonChirps[j].CreatedAt)
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonChirps)
}

func (a *apiConfig) getChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	chirpIDstr := req.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		log.Printf("Invalid chirp ID: %s", err)
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}
	chirp, err := a.dbQueries.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp by ID: %s", err)
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		UserID:    chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonChirp)
}

func (a *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpIDstr := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		log.Printf("Invalid chirp ID: %s", err)
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirp, err := a.dbQueries.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp by ID: %s", err)
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	if chirp.UserID != userID {
		log.Printf("Unauthorized attempt to delete chirp by user %s", userID)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = a.dbQueries.DeleteChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
