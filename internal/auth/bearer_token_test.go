package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       map[string]string
		expectedToken string
		expectedError string
	}{
		{
			name: "Valid Bearer Token",
			headers: map[string]string{
				"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			},
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectedError: "",
		},
		{
			name: "Valid Bearer Token with Short JWT",
			headers: map[string]string{
				"Authorization": "Bearer abc123",
			},
			expectedToken: "abc123",
			expectedError: "",
		},
		{
			name: "Valid Bearer Token with Long JWT",
			headers: map[string]string{
				"Authorization": "Bearer " + string(make([]byte, 1000)), // Very long token
			},
			expectedToken: string(make([]byte, 1000)),
			expectedError: "",
		},
		{
			name: "Valid Bearer Token with Special Characters",
			headers: map[string]string{
				"Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjJ9.4Adcj3UFYzPUVaVF43FmMo",
			},
			expectedToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyNDI2MjJ9.4Adcj3UFYzPUVaVF43FmMo",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create headers
			headers := http.Header{}
			for key, value := range tt.headers {
				headers.Set(key, value)
			}

			// Call function
			token, err := GetBearerToken(headers)

			// Check token
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token %q, got %q", tt.expectedToken, token)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
			}
		})
	}
}

func TestGetBearerTokenMissingHeader(t *testing.T) {
	headers := http.Header{}
	// Don't set Authorization header

	token, err := GetBearerToken(headers)

	if err == nil {
		t.Error("Expected error for missing Authorization header, got nil")
	}
	if err.Error() != "no authorization header" {
		t.Errorf("Expected error 'no authorization header', got %q", err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token, got %q", token)
	}
}

func TestGetBearerTokenEmptyHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "")

	token, err := GetBearerToken(headers)

	if err == nil {
		t.Error("Expected error for empty Authorization header, got nil")
	}
	if err.Error() != "no authorization header" {
		t.Errorf("Expected error 'no authorization header', got %q", err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token, got %q", token)
	}
}

func TestGetBearerTokenInvalidPrefix(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedError string
	}{
		{
			name:          "Wrong prefix - Bearer lowercase",
			authHeader:    "bearer token123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - BEARER uppercase",
			authHeader:    "BEARER token123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with no space",
			authHeader:    "Bearer",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with tab",
			authHeader:    "Bearer\ttoken123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with newline",
			authHeader:    "Bearer\ntoken123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with carriage return",
			authHeader:    "Bearer\rtoken123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with mixed case",
			authHeader:    "BeArEr token123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with numbers",
			authHeader:    "Bearer123 token123",
			expectedError: "invalid authorization header",
		},
		{
			name:          "Wrong prefix - Bearer with special chars",
			authHeader:    "Bearer! token123",
			expectedError: "invalid authorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tt.authHeader)

			token, err := GetBearerToken(headers)

			if err == nil {
				t.Errorf("Expected error for invalid header %q, got nil", tt.authHeader)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
			}
			if token != "" {
				t.Errorf("Expected empty token, got %q", token)
			}
		})
	}
}

func TestGetBearerTokenEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectedError string
	}{
		{
			name:          "Bearer with exactly 7 chars (minimum valid)",
			authHeader:    "Bearer ",
			expectedToken: "",
			expectedError: "",
		},
		{
			name:          "Bearer with exactly 8 chars (minimum valid with token)",
			authHeader:    "Bearer a",
			expectedToken: "a",
			expectedError: "",
		},
		{
			name:          "Bearer with space in token",
			authHeader:    "Bearer token with spaces",
			expectedToken: "token with spaces",
			expectedError: "",
		},
		{
			name:          "Bearer with special characters in token",
			authHeader:    "Bearer !@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectedToken: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectedError: "",
		},
		{
			name:          "Bearer with unicode in token",
			authHeader:    "Bearer 🚀🎉✨",
			expectedToken: "🚀🎉✨",
			expectedError: "",
		},
		{
			name:          "Bearer with newlines in token",
			authHeader:    "Bearer line1\nline2\r\nline3",
			expectedToken: "line1\nline2\r\nline3",
			expectedError: "",
		},
		{
			name:          "Bearer with tabs in token",
			authHeader:    "Bearer tab1\ttab2\t\t\ttab3",
			expectedToken: "tab1\ttab2\t\t\ttab3",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tt.authHeader)

			token, err := GetBearerToken(headers)

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token %q, got %q", tt.expectedToken, token)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %q, got %q", tt.expectedError, err.Error())
				}
			}
		})
	}
}

func TestGetBearerTokenMultipleHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer valid-token")
	headers.Add("Authorization", "Bearer another-token") // Multiple headers with same name

	token, err := GetBearerToken(headers)

	// Should use the first one set
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if token != "valid-token" {
		t.Errorf("Expected token 'valid-token', got %q", token)
	}
}

func TestGetBearerTokenCaseInsensitiveHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("authorization", "Bearer valid-token") // lowercase header name

	token, err := GetBearerToken(headers)

	// Go's http.Header.Get is case-insensitive
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if token != "valid-token" {
		t.Errorf("Expected token 'valid-token', got %q", token)
	}
}

func TestGetBearerTokenWithNilHeaders(t *testing.T) {
	// This test ensures the function handles nil headers gracefully
	// In practice, this shouldn't happen with real HTTP requests
	var headers http.Header

	token, err := GetBearerToken(headers)

	if err == nil {
		t.Error("Expected error for nil headers, got nil")
	}
	if err.Error() != "no authorization header" {
		t.Errorf("Expected error 'no authorization header', got %q", err.Error())
	}
	if token != "" {
		t.Errorf("Expected empty token, got %q", token)
	}
}
