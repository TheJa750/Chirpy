package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/Chirpy/internal/auth"
	"github.com/TheJa750/Chirpy/internal/database"
)

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
		ChirpyRed: user.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonUser)
}

func (a *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	var userReq UserRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		log.Printf("Error decoding user request: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if userReq.Email == "" || userReq.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := a.dbQueries.GetUserByEmail(req.Context(), userReq.Email)
	if err != nil {
		log.Printf("Error getting user by email: %s", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, userReq.Password)
	if err != nil {
		log.Printf("Password check failed: %s", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	access_token, err := auth.MakeJWT(user.ID, a.JWTSecret, 3600*time.Second)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rtString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = a.dbQueries.CreateUserToken(req.Context(), database.CreateUserTokenParams{
		Token:  rtString,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error creating refresh token in database: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email,
		AccessToken:  access_token,
		RefreshToken: rtString,
		ChirpyRed:    user.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonUser)
}

func (a *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
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

	var userReq UserRequest
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&userReq)
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

	user, err := a.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          userReq.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
		ChirpyRed: user.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonUser)
}

func (a *apiConfig) upgradeUserHandler(w http.ResponseWriter, req *http.Request) {
	var event PolkaEvent
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&event)
	if err != nil {
		log.Printf("Error decoding Polka event: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if event.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = a.dbQueries.UpgradeUserToChirpyRed(req.Context(), event.Data.UserID)
	if err != nil {
		log.Printf("Error upgrading user to Chirpy Red: %s", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
