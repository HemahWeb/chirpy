package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"

	// Test that hashing works
	hash1, err := HashPassword(password1)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}
	if hash1 == "" {
		t.Error("HashPassword() returned empty hash")
	}

	hash2, err := HashPassword(password2)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}
	if hash2 == "" {
		t.Error("HashPassword() returned empty hash")
	}

	// Test that different passwords produce different hashes
	if hash1 == hash2 {
		t.Error("Different passwords should produce different hashes")
	}

	// Test that hashing the same password multiple times produces different hashes (due to salt)
	hash1Again, err := HashPassword(password1)
	if err != nil {
		t.Fatalf("HashPassword() failed on second call: %v", err)
	}
	if hash1 == hash1Again {
		t.Error("Hashing the same password should produce different hashes due to salt")
	}
}

func TestHashPasswordWithEmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() should work with empty password: %v", err)
	}
	if hash == "" {
		t.Error("HashPassword() returned empty hash for empty password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPasswordWithSpecialCharacters(t *testing.T) {
	passwords := []string{
		"!@#$%^&*()",
		"password with spaces",
		"password\nwith\nnewlines",
		"password\twith\ttabs",
		"password with unicode: 🚀🎉",
		"this is a moderately long password that should still work correctly",
	}

	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Errorf("HashPassword() failed for password '%s': %v", password, err)
		}
		if hash == "" {
			t.Errorf("HashPassword() returned empty hash for password '%s'", password)
		}

		// Verify the hash can be validated
		err = CheckPasswordHash(password, hash)
		if err != nil {
			t.Errorf("Generated hash could not be validated for password '%s': %v", password, err)
		}
	}
}
