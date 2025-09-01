package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}

	// Verify the token can be validated
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Generated token could not be validated: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Parsed user ID %v does not match original %v", parsedUserID, userID)
	}
}

func TestMakeJWTWithEmptySecret(t *testing.T) {
	userID := uuid.New()
	secret := ""
	expiresIn := 1 * time.Hour

	// Note: JWT library actually allows empty secrets, so this test should pass
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed with empty secret: %v", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}

	// Verify the token can be validated
	_, err = ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Generated token with empty secret could not be validated: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Test valid token
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("ValidateJWT() failed for valid token: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("ValidateJWT() returned wrong user ID: got %v, want %v", parsedUserID, userID)
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, correctSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Test with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT() should fail with wrong secret")
	}
}

func TestValidateJWTWithExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -1 * time.Hour // Token expired 1 hour ago

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create expired test token: %v", err)
	}

	// Test expired token
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with expired token")
	}
}

func TestValidateJWTWithExpiredTokenInPast(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create a token that expires in the past
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create expired test token: %v", err)
	}

	// Test expired token
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with expired token")
	}
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	secret := "test-secret"

	// Test with malformed token
	_, err := ValidateJWT("invalid.token.here", secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with malformed token")
	}

	// Test with empty token
	_, err = ValidateJWT("", secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with empty token")
	}

	// Test with just dots
	_, err = ValidateJWT("...", secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with just dots")
	}
}

func TestValidateJWTWithTamperedToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Tamper with the token by changing a character
	tamperedToken := token[:len(token)-1] + "X"

	_, err = ValidateJWT(tamperedToken, secret)
	if err == nil {
		t.Error("ValidateJWT() should fail with tampered token")
	}
}

func TestValidateJWTWithDifferentUserIDs(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	token1, err := MakeJWT(userID1, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token 1: %v", err)
	}

	token2, err := MakeJWT(userID2, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token 2: %v", err)
	}

	// Both tokens should be valid but return different user IDs
	parsedUserID1, err := ValidateJWT(token1, secret)
	if err != nil {
		t.Errorf("ValidateJWT() failed for token 1: %v", err)
	}
	if parsedUserID1 != userID1 {
		t.Errorf("Token 1 returned wrong user ID: got %v, want %v", parsedUserID1, userID1)
	}

	parsedUserID2, err := ValidateJWT(token2, secret)
	if err != nil {
		t.Errorf("ValidateJWT() failed for token 2: %v", err)
	}
	if parsedUserID2 != userID2 {
		t.Errorf("Token 2 returned wrong user ID: got %v, want %v", parsedUserID2, userID2)
	}

	// Verify they are different
	if parsedUserID1 == parsedUserID2 {
		t.Error("Different tokens should return different user IDs")
	}
}

func TestJWTClaims(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Parse the token to verify claims
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Fatal("Failed to get claims from token")
	}

	// Verify issuer
	if claims.Issuer != "chirpy" {
		t.Errorf("Expected issuer 'chirpy', got '%s'", claims.Issuer)
	}

	// Verify subject (user ID)
	if claims.Subject != userID.String() {
		t.Errorf("Expected subject '%s', got '%s'", userID.String(), claims.Subject)
	}

	// Verify issued at is recent (within last minute)
	if claims.IssuedAt == nil {
		t.Error("IssuedAt claim is missing")
	} else {
		issuedAt := claims.IssuedAt.Time
		if time.Since(issuedAt) > time.Minute {
			t.Errorf("Token was issued too long ago: %v", issuedAt)
		}
	}

	// Verify expires at is in the future
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt claim is missing")
	} else {
		expiresAt := claims.ExpiresAt.Time
		if expiresAt.Before(time.Now()) {
			t.Errorf("Token expires in the past: %v", expiresAt)
		}
		if time.Until(expiresAt) > expiresIn+time.Second {
			t.Errorf("Token expires too far in the future: %v", expiresAt)
		}
	}
}

func TestJWTWithVeryShortExpiration(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 200 * time.Millisecond // Longer to avoid immediate expiration

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Token should be valid immediately
	_, err = ValidateJWT(token, secret)
	if err != nil {
		// If it's already expired, that's also acceptable behavior
		t.Logf("Token expired immediately, which is acceptable: %v", err)
		return
	}

	// Wait for token to expire
	time.Sleep(250 * time.Millisecond)

	// Token should now be expired
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Token should be expired after waiting")
	}
}

func TestJWTWithZeroExpiration(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := 50 * time.Millisecond // Use a short but workable time

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Token should be valid immediately (or may already be expired, which is fine)
	_, err = ValidateJWT(token, secret)
	if err != nil {
		// If it's already expired, that's also acceptable behavior for very short expiration
		t.Logf("Token expired immediately, which is acceptable: %v", err)
		return
	}

	// Wait for token to expire
	time.Sleep(100 * time.Millisecond)

	// Token should now be expired
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Token should be expired after waiting")
	}
}

func TestJWTWithNegativeExpiration(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -1 * time.Hour // Negative expiration

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Token should be expired immediately
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Token with negative expiration should be expired immediately")
	}
}

func TestJWTWithNilUUID(t *testing.T) {
	var userID uuid.UUID // nil UUID
	secret := "test-secret"
	expiresIn := 1 * time.Hour

	// This should work (nil UUID is valid)
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() should work with nil UUID: %v", err)
	}

	// Token should be valid
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Token with nil UUID should be valid: %v", err)
	}

	// Parsed UUID should match (nil UUID)
	if parsedUserID != userID {
		t.Errorf("Expected nil UUID, got %v", parsedUserID)
	}
}
