package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheJa750/Chirpy/internal/auth"
)

func (a *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshToken, err := a.dbQueries.GetUserByToken(req.Context(), token)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Refresh token has been revoked", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, a.JWTSecret, 3600*time.Second)
	if err != nil {
		log.Printf("Error creating new access token: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": newAccessToken})
}

func (a *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshToken, err := a.dbQueries.GetUserByToken(req.Context(), token)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Refresh token has already been revoked", http.StatusBadRequest)
		return
	}

	err = a.dbQueries.RevokeUserToken(req.Context(), refreshToken.Token)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
