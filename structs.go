package main

import (
	"sync/atomic"
	"time"

	"github.com/TheJa750/Chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	JWTSecret      string
	PolkaKey       string
}

const adminMetrics = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

type JsonError struct {
	Message string `json:"error"`
}

type CleanedChirpBody struct {
	Body string `json:"cleaned_body"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ChirpyRed    bool      `json:"is_chirpy_red"`
}

type chirpRequest struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PolkaEvent struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}
