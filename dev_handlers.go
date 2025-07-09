package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (a *apiConfig) resetUsersHandler(w http.ResponseWriter, req *http.Request) {
	// Resetting the users table is a destructive operation, so we should be careful.
	// This is just for development purposes.
	// In a production environment, you would likely want to implement a more secure
	// and controlled way to reset the users table, such as through an admin interface
	// with proper authentication and authorization.
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(w, "This operation is only allowed in development mode", http.StatusForbidden)
		return
	}

	err := a.dbQueries.ResetUsers(req.Context())
	if err != nil {
		http.Error(w, "Failed to reset users", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Users table reset successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
