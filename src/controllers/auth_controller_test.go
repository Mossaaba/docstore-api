package controllers

import (
	"bytes"
	"docstore-api/src/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthController(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "test-secret",
		AdminUser: "admin",
		AdminPass: "password",
	}

	controller := NewAuthController(cfg)

	assert.NotNil(t, controller)
	assert.Equal(t, cfg, controller.config)
}

func TestAuthController_Login(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret: "test-secret-key",
		AdminUser: "admin",
		AdminPass: "password123",
	}

	controller := NewAuthController(cfg)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedFields []string
		checkToken     bool
	}{
		{
			name: "successful login with valid credentials",
			requestBody: LoginRequest{
				Username: "admin",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"token", "user"},
			checkToken:     true,
		},
		{
			name: "failed login with invalid username",
			requestBody: LoginRequest{
				Username: "wronguser",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			checkToken:     false,
		},
		{
			name: "failed login with invalid password",
			requestBody: LoginRequest{
				Username: "admin",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			checkToken:     false,
		},
		{
			name: "failed login with missing username",
			requestBody: LoginRequest{
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			checkToken:     false,
		},
		{
			name: "failed login with missing password",
			requestBody: LoginRequest{
				Username: "admin",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			checkToken:     false,
		},
		{
			name:           "failed login with invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			checkToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test
			router := gin.New()
			router.POST("/login", controller.Login)

			// Prepare request body
			var requestBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check expected fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Additional checks for successful login
			if tt.checkToken && tt.expectedStatus == http.StatusOK {
				token, exists := response["token"].(string)
				assert.True(t, exists, "Token should exist in response")
				assert.NotEmpty(t, token, "Token should not be empty")

				user, exists := response["user"].(string)
				assert.True(t, exists, "User should exist in response")
				assert.Equal(t, "admin", user, "User should match the logged in user")
			}

			// Additional checks for error responses
			if tt.expectedStatus != http.StatusOK {
				errorMsg, exists := response["error"].(string)
				assert.True(t, exists, "Error message should exist in response")
				assert.NotEmpty(t, errorMsg, "Error message should not be empty")
			}
		})
	}
}

func TestAuthController_Login_Integration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret: "integration-test-secret",
		AdminUser: "testadmin",
		AdminPass: "testpass123",
	}

	controller := NewAuthController(cfg)
	router := gin.New()
	router.POST("/api/v1/auth/login", controller.Login)

	// Test successful login flow
	loginReq := LoginRequest{
		Username: "testadmin",
		Password: "testpass123",
	}

	requestBody, err := json.Marshal(loginReq)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "testadmin", response.User)

	// Verify the token format (should be a JWT with 3 parts separated by dots)
	tokenParts := len(bytes.Split([]byte(response.Token), []byte(".")))
	assert.Equal(t, 3, tokenParts, "JWT token should have 3 parts separated by dots")
}

func TestLoginRequest_Validation(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret: "test-secret",
		AdminUser: "admin",
		AdminPass: "password",
	}

	controller := NewAuthController(cfg)
	router := gin.New()
	router.POST("/login", controller.Login)

	tests := []struct {
		name        string
		requestBody map[string]interface{}
		expectError bool
	}{
		{
			name: "valid request with all fields",
			requestBody: map[string]interface{}{
				"username": "admin",
				"password": "password",
			},
			expectError: false,
		},
		{
			name: "empty username",
			requestBody: map[string]interface{}{
				"username": "",
				"password": "password",
			},
			expectError: true,
		},
		{
			name: "empty password",
			requestBody: map[string]interface{}{
				"username": "admin",
				"password": "",
			},
			expectError: true,
		},
		{
			name: "missing username field",
			requestBody: map[string]interface{}{
				"password": "password",
			},
			expectError: true,
		},
		{
			name: "missing password field",
			requestBody: map[string]interface{}{
				"username": "admin",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				// Should either succeed (200) or fail with unauthorized (401) based on credentials
				assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized)
			}
		})
	}
}
