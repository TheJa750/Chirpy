package auth

import (
	"testing"
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
