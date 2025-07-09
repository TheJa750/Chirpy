package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("Hashed password should not be empty")
	}

	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Fatalf("Password check failed: %v", err)
	}
}
func TestCheckPasswordHash(t *testing.T) {
	hashedPassword, err := HashPassword("testPassword")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, "testPassword")
	if err != nil {
		t.Fatalf("Password check failed: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, "wrongPassword")
	if err == nil {
		t.Fatal("Expected password check to fail with wrong password")
	}
}

func TestJWT(t *testing.T) {
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	token, err := MakeJWT(userID, "testSecret", 3600*time.Second)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	if token == "" {
		t.Fatal("JWT should not be empty")
	}

	parsedID, err := ValidateJWT(token, "testSecret")
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("Expected user ID %s, got %s", userID, parsedID)
	}
	_, err = ValidateJWT("invalidToken", "testSecret")
	if err == nil {
		t.Fatal("Expected validation to fail with invalid token")
	}
}

func TestGetBearerToken(t *testing.T) {
	emptyHeaders := http.Header{}
	validHeaders := http.Header{
		"Authorization": []string{"Bearer validToken"},
	}
	invalidHeaders := http.Header{
		"Authorization": []string{"InvalidToken"},
	}

	_, err := GetBearerToken(emptyHeaders)
	if err == nil {
		t.Fatal("Expected error for empty Authorization header")
	}

	token, err := GetBearerToken(validHeaders)
	if err != nil {
		t.Fatalf("Failed to get bearer token from valid headers: %v", err)
	}
	if token != "validToken" {
		t.Fatalf("Expected token 'validToken', got '%s'", token)
	}

	_, err = GetBearerToken(invalidHeaders)
	if err == nil {
		t.Fatal("Expected error for invalid Authorization header format")
	}

}
