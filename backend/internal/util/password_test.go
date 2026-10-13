package util_test

import (
	"testing"

	"tourism-booking/internal/util"
)

func TestHashAndPassword(t *testing.T) {
	plainPassword := "SecretP@ssw0rd!"

	// Test HashPassword
	hashed, err := util.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashed == plainPassword {
		t.Fatal("HashPassword returned plain text password without hashing")
	}

	// Test CheckPassword - Correct password
	if !util.CheckPassword(hashed, plainPassword) {
		t.Error("CheckPassword failed for correct password")
	}

	// Test CheckPassword - Wrong password
	if util.CheckPassword(hashed, "WrongPassword123") {
		t.Error("CheckPassword returned true for incorrect password")
	}

	// Test CheckPassword - Empty password
	if util.CheckPassword(hashed, "") {
		t.Error("CheckPassword returned true for empty password")
	}
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	invalidHash := "not-a-valid-bcrypt-hash"
	if util.CheckPassword(invalidHash, "password") {
		t.Error("CheckPassword returned true for invalid hash string")
	}
}
