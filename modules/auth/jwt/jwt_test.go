package jwt

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTOperations(t *testing.T) {
	// Initialize JWT with a test secret
	jwtSecret := "test-secret-key"
	jwtService := NewJSONWebToken(jwtSecret)

	t.Run("Generate and Verify Tokens", func(t *testing.T) {
		// Test data
		userID := int64(1)
		username := "testuser"
		role := 1

		// Generate tokens
		accessToken, refreshToken, err := jwtService.GenerateJSONWebTokens(userID, username, role)
		require.NoError(t, err)
		require.NotEmpty(t, accessToken)
		require.NotEmpty(t, refreshToken)

		// Verify access token
		claims, err := jwtService.VerifyToken(accessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)

		// Verify refresh token
		claims, err = jwtService.VerifyToken(refreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Refresh Access Token", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"
		role := 1

		_, refreshToken, err := jwtService.GenerateJSONWebTokens(userID, username, role)
		require.NoError(t, err)

		newAccessToken, err := jwtService.RefreshAccessToken(refreshToken)
		require.NoError(t, err)
		require.NotEmpty(t, newAccessToken)

		claims, err := jwtService.VerifyToken(newAccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := jwtService.VerifyToken("invalid-token")
		assert.Error(t, err)
	})

	t.Run("Token Expiration", func(t *testing.T) {
		// Create an expired token
		expiredToken, err := jwtService.generateJSONWebToken(1, 1, "testuser", time.Now().Add(-1*time.Hour))
		require.NoError(t, err)

		_, err = jwtService.VerifyToken(expiredToken)
		assert.Error(t, err)
	})
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name          string
		header        string
		expectedToken string
		expectedError bool
	}{
		{
			name:          "Valid Bearer Token",
			header:        "Bearer valid-token",
			expectedToken: "valid-token",
			expectedError: false,
		},
		{
			name:          "Missing Bearer Prefix",
			header:        "valid-token",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Empty Header",
			header:        "",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Invalid Format",
			header:        "Bearer token extra",
			expectedToken: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			token, err := ExtractTokenFromHeader(req)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}
