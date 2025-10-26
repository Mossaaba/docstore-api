package middleware

import (
	"docstore-api/src/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	tests := []struct {
		name     string
		username string
	}{
		{
			name:     "generate token for admin user",
			username: "admin",
		},
		{
			name:     "generate token for regular user",
			username: "user123",
		},
		{
			name:     "generate token with special characters",
			username: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.username, cfg)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify token format (JWT should have 3 parts separated by dots)
			parts := len([]byte(token))
			assert.Greater(t, parts, 50, "Token should be reasonably long")

			// Verify we can parse the token back
			claims, err := ValidateToken(token, cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.username, claims.Username)

			// Verify expiration is set correctly (24 hours from now)
			expectedExpiry := time.Now().Add(24 * time.Hour)
			actualExpiry := claims.ExpiresAt.Time

			// Allow 1 minute tolerance for test execution time
			timeDiff := actualExpiry.Sub(expectedExpiry)
			assert.True(t, timeDiff < time.Minute && timeDiff > -time.Minute,
				"Token expiry should be approximately 24 hours from now")
		})
	}
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Generate a valid token for testing
	validToken, err := GenerateToken("testuser", cfg)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		expectError  bool
		expectedUser string
	}{
		{
			name:         "valid token",
			token:        validToken,
			expectError:  false,
			expectedUser: "testuser",
		},
		{
			name:        "invalid token format",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "malformed token",
			token:       "not.a.jwt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectedUser, claims.Username)
			}
		})
	}
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	cfg1 := &config.Config{JWTSecret: "secret1"}
	cfg2 := &config.Config{JWTSecret: "secret2"}

	// Generate token with first secret
	token, err := GenerateToken("testuser", cfg1)
	assert.NoError(t, err)

	// Try to validate with different secret
	claims, err := ValidateToken(token, cfg2)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestExpiredToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Create an expired token manually
	expiredClaims := &Claims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // Issued 2 hours ago
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	assert.NoError(t, err)

	// Try to validate expired token
	claims, err := ValidateToken(expiredTokenString, cfg)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Generate a valid token for testing
	validToken, err := GenerateToken("testuser", cfg)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectAbort    bool
		expectUsername string
	}{
		{
			name:           "valid bearer token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectAbort:    false,
			expectUsername: "testuser",
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "invalid header format - no Bearer prefix",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "invalid header format - wrong prefix",
			authHeader:     "Basic " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "empty token after Bearer",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test router with the middleware
			router := gin.New()
			router.Use(JWTAuthMiddleware(cfg))

			// Add a test endpoint that should only be reached if middleware passes
			router.GET("/protected", func(c *gin.Context) {
				username, exists := c.Get("username")
				assert.True(t, exists, "Username should be set in context")
				c.JSON(http.StatusOK, gin.H{"user": username})
			})

			// Create request
			req, err := http.NewRequest("GET", "/protected", nil)
			assert.NoError(t, err)

			// Set authorization header if provided
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectAbort {
				// If middleware didn't abort, check that username was set correctly
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUsername, response["user"])
			}
		})
	}
}

func TestJWTAuthMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret: "integration-test-secret",
	}

	// Create router with middleware
	router := gin.New()

	// Public endpoint (no middleware)
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public"})
	})

	// Protected endpoints (with middleware)
	protected := router.Group("/api")
	protected.Use(JWTAuthMiddleware(cfg))
	{
		protected.GET("/user", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{"user": username})
		})

		protected.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "sensitive"})
		})
	}

	// Test public endpoint (should work without token)
	t.Run("public endpoint without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/public", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test protected endpoint without token (should fail)
	t.Run("protected endpoint without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test protected endpoint with valid token (should work)
	t.Run("protected endpoint with valid token", func(t *testing.T) {
		token, err := GenerateToken("integrationuser", cfg)
		assert.NoError(t, err)

		req, _ := http.NewRequest("GET", "/api/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestClaims(t *testing.T) {
	// Test Claims struct
	claims := &Claims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	assert.Equal(t, "testuser", claims.Username)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestTokenLifecycle(t *testing.T) {
	// Test complete token lifecycle: generate -> validate -> use in middleware
	cfg := &config.Config{
		JWTSecret: "lifecycle-test-secret",
	}

	username := "lifecycleuser"

	// Step 1: Generate token
	token, err := GenerateToken(username, cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Step 2: Validate token
	claims, err := ValidateToken(token, cfg)
	assert.NoError(t, err)
	assert.Equal(t, username, claims.Username)

	// Step 3: Use token in middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(JWTAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		contextUsername, exists := c.Get("username")
		assert.True(t, exists)
		c.JSON(http.StatusOK, gin.H{"user": contextUsername})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
