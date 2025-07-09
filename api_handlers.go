package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (a *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	var userReq UserRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		log.Printf("Error decoding user request: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if userReq.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if userReq.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user, err := a.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          userReq.Email,
		HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonUser)
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

	cleanedChirp, errMsg := validateChirp(chirpReq.Body)
	if errMsg.Message != "" {
		log.Printf("Chirp validation error: %s", errMsg.Message)
		http.Error(w, errMsg.Message, http.StatusBadRequest)
		return
	}

	chirp, err := a.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp.Body,
		UserID: chirpReq.UserID,
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
	chirps, err := a.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonChirps)
}

func (a *apiConfig) getChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	chirpIDstr := req.PathValue("chirpID")
	fmt.Println(chirpIDstr)
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
