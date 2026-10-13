package util_test

import (
	"testing"
	"time"

	"tourism-booking/internal/util"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "my-secret-key-12345"
	userID := int64(42)
	role := "admin"
	expiryHours := 2

	// Test GenerateToken
	tokenString, err := util.GenerateToken(secret, userID, role, expiryHours)
	if err != nil {
		t.Fatalf("GenerateToken return error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("GenerateToken returned empty token string")
	}

	// Test ParseToken with correct secret
	claims, err := util.ParseToken(secret, tokenString)
	if err != nil {
		t.Fatalf("ParseToken return error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected Role %s, got %s", role, claims.Role)
	}

	// Test ParseToken with wrong secret
	_, err = util.ParseToken("wrong-secret", tokenString)
	if err == nil {
		t.Error("expected error when parsing token with wrong secret, got nil")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	testCases := []struct {
		name        string
		secret      string
		tokenString string
	}{
		{
			name:        "garbage token string",
			secret:      "secret",
			tokenString: "invalid.token.string",
		},
		{
			name:        "empty token string",
			secret:      "secret",
			tokenString: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := util.ParseToken(tc.secret, tc.tokenString)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
			if claims != nil {
				t.Errorf("expected nil claims for %s, got %+v", tc.name, claims)
			}
		})
	}
}

func TestGenerateToken_Expired(t *testing.T) {
	secret := "my-secret-key"
	// Generate token with negative expiry hours (already expired)
	tokenString, err := util.GenerateToken(secret, 1, "user", -1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Give a tiny pause to ensure timestamp has passed
	time.Sleep(10 * time.Millisecond)

	_, err = util.ParseToken(secret, tokenString)
	if err == nil {
		t.Error("expected error when parsing expired token, got nil")
	}
}
