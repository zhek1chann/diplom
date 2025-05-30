package jwt

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
	jwt *JSONWebToken
}

func TestJWT(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func (s *JWTTestSuite) SetupTest() {
	s.jwt = NewJSONWebToken("test-secret-key")
}

func (s *JWTTestSuite) TestGenerateJSONWebTokens_Success() {
	userID := int64(1)
	username := "testuser"
	role := 0 // CustomerRole

	accessToken, refreshToken, err := s.jwt.GenerateJSONWebTokens(userID, username, role)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken)
	assert.NotEmpty(s.T(), refreshToken)

	// Verify access token
	claims, err := s.jwt.VerifyToken(accessToken)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userID, claims.UserID)
	assert.Equal(s.T(), username, claims.Username)
	assert.Equal(s.T(), role, claims.Role)
}

func (s *JWTTestSuite) TestRefreshAccessToken_Success() {
	userID := int64(1)
	username := "testuser"
	role := 0 // CustomerRole

	// Generate initial tokens
	_, refreshToken, err := s.jwt.GenerateJSONWebTokens(userID, username, role)
	assert.NoError(s.T(), err)

	// Use refresh token to get new access token
	newAccessToken, err := s.jwt.RefreshAccessToken(refreshToken)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), newAccessToken)

	// Verify new access token
	claims, err := s.jwt.VerifyToken(newAccessToken)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userID, claims.UserID)
	assert.Equal(s.T(), username, claims.Username)
	assert.Equal(s.T(), role, claims.Role)
}

func (s *JWTTestSuite) TestVerifyToken_InvalidToken() {
	// Test with invalid token
	_, err := s.jwt.VerifyToken("invalid.token.here")
	assert.Error(s.T(), err)

	// Test with expired token
	expiredToken := s.generateExpiredToken()
	_, err = s.jwt.VerifyToken(expiredToken)
	assert.Error(s.T(), err)
}

func (s *JWTTestSuite) TestRefreshAccessToken_InvalidToken() {
	_, err := s.jwt.RefreshAccessToken("invalid.refresh.token")
	assert.Error(s.T(), err)
}

func (s *JWTTestSuite) generateExpiredToken() string {
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		Role:     0,
	}
	claims.ExpiresAt.Time = time.Now().Add(-time.Hour) // Expired 1 hour ago

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(s.jwt.jwtKey)
	return tokenString
}

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
