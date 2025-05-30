package middleware

import (
	"diploma/modules/auth/jwt"
	contextkeys "diploma/pkg/context-keys"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockJWT is a mock implementation of the JWT service
type mockJWT struct {
	mock.Mock
}

func (m *mockJWT) VerifyToken(token string) (*jwt.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func (m *mockJWT) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	args := m.Called(id, username, role)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockJWT) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *mockJWT) GetJWTKey() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var ErrInvalidToken = errors.New("invalid token")

	tests := []struct {
		name           string
		setupAuth      func(*http.Request)
		setupMock      func(*mockJWT)
		expectedStatus int
		checkContext   bool
	}{
		{
			name: "Valid Token",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
			},
			setupMock: func(m *mockJWT) {
				claims := &jwt.Claims{
					UserID:   1,
					Username: "testuser",
					Role:     1,
				}
				m.On("VerifyToken", "valid-token").Return(claims, nil)
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name: "Invalid Token",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupMock: func(m *mockJWT) {
				m.On("VerifyToken", "invalid-token").Return(nil, ErrInvalidToken)
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "Missing Authorization Header",
			setupAuth: func(req *http.Request) {
				// Don't set any header
			},
			setupMock: func(m *mockJWT) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name: "Invalid Authorization Format",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "InvalidFormat")
			},
			setupMock: func(m *mockJWT) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockJWT := new(mockJWT)
			tt.setupMock(mockJWT)

			middleware := &AuthMiddleware{
				jwt: mockJWT,
			}

			router := gin.New()
			router.Use(middleware.AuthMiddleware())

			var contextUser *jwt.Claims
			router.GET("/test", func(c *gin.Context) {
				if tt.checkContext {
					if user, exists := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims); exists {
						contextUser = user
					}
				}
				c.Status(http.StatusOK)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			tt.setupAuth(req)

			// Perform request
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkContext && w.Code == http.StatusOK {
				assert.NotNil(t, contextUser)
				if contextUser != nil {
					assert.Equal(t, int64(1), contextUser.UserID)
					assert.Equal(t, "testuser", contextUser.Username)
					assert.Equal(t, 1, contextUser.Role)
				}
			}

			mockJWT.AssertExpectations(t)
		})
	}
}
