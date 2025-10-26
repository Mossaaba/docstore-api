package controllers

import (
	"docstore-api/src/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthController(t *testing.T) {
	cfg := &config.Config{
		Environment: "test",
	}

	controller := NewHealthController(cfg)

	assert.NotNil(t, controller)
	assert.Equal(t, cfg, controller.config)
}

func TestHealthController_HealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		environment string
	}{
		{
			name:        "health check in development environment",
			environment: "development",
		},
		{
			name:        "health check in production environment",
			environment: "production",
		},
		{
			name:        "health check in test environment",
			environment: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Environment: tt.environment,
			}

			controller := NewHealthController(cfg)
			router := gin.New()
			router.GET("/health", controller.HealthCheck)

			// Create request
			req, err := http.NewRequest("GET", "/health", nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Record the time before the request
			beforeRequest := time.Now().UTC()

			// Perform request
			router.ServeHTTP(w, req)

			// Record the time after the request
			afterRequest := time.Now().UTC()

			// Check status code
			assert.Equal(t, http.StatusOK, w.Code)

			// Check content type
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			// Parse response
			var response HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify response fields
			assert.Equal(t, "ok", response.Status)
			assert.Equal(t, "docstore-api", response.Service)
			assert.Equal(t, "1.0.0", response.Version)
			assert.Equal(t, tt.environment, response.Environment)

			// Verify timestamp is within reasonable range
			assert.True(t, response.Timestamp.After(beforeRequest) || response.Timestamp.Equal(beforeRequest))
			assert.True(t, response.Timestamp.Before(afterRequest) || response.Timestamp.Equal(afterRequest))

			// Verify timestamp is in UTC
			assert.Equal(t, time.UTC, response.Timestamp.Location())
		})
	}
}

func TestHealthController_HealthCheck_JSONStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Environment: "test",
	}

	controller := NewHealthController(cfg)
	router := gin.New()
	router.GET("/health", controller.HealthCheck)

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse as generic map to check structure
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check all required fields exist
	requiredFields := []string{"status", "timestamp", "service", "version", "environment"}
	for _, field := range requiredFields {
		assert.Contains(t, response, field, "Response should contain field: %s", field)
	}

	// Check field types
	assert.IsType(t, "", response["status"])
	assert.IsType(t, "", response["timestamp"])
	assert.IsType(t, "", response["service"])
	assert.IsType(t, "", response["version"])
	assert.IsType(t, "", response["environment"])
}

func TestHealthController_Metrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		environment string
	}{
		{
			name:        "metrics in development environment",
			environment: "development",
		},
		{
			name:        "metrics in production environment",
			environment: "production",
		},
		{
			name:        "metrics in test environment",
			environment: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Environment: tt.environment,
			}

			controller := NewHealthController(cfg)
			router := gin.New()
			router.GET("/metrics", controller.Metrics)

			// Create request
			req, err := http.NewRequest("GET", "/metrics", nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, http.StatusOK, w.Code)

			// Check content type
			assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

			// Get response body
			responseBody := w.Body.String()

			// Check that response is not empty
			assert.NotEmpty(t, responseBody)

			// Check for expected Prometheus metric patterns
			expectedMetrics := []string{
				"docstore_api_info",
				"docstore_api_uptime_seconds",
				"docstore_api_memory_usage_bytes",
				"docstore_api_memory_allocated_bytes",
				"docstore_api_goroutines",
				"docstore_api_health_status",
			}

			for _, metric := range expectedMetrics {
				assert.Contains(t, responseBody, metric, "Response should contain metric: %s", metric)
			}

			// Check for environment in the info metric
			expectedEnvMetric := `environment="` + tt.environment + `"`
			assert.Contains(t, responseBody, expectedEnvMetric, "Response should contain environment: %s", tt.environment)

			// Check for version in the info metric
			assert.Contains(t, responseBody, `version="1.0.0"`, "Response should contain version")

			// Check for HELP and TYPE comments
			assert.Contains(t, responseBody, "# HELP", "Response should contain HELP comments")
			assert.Contains(t, responseBody, "# TYPE", "Response should contain TYPE comments")

			// Check that health status is 1 (healthy)
			assert.Contains(t, responseBody, "docstore_api_health_status 1", "Health status should be 1")
		})
	}
}

func TestHealthController_Metrics_PrometheusFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Environment: "test",
	}

	controller := NewHealthController(cfg)
	router := gin.New()
	router.GET("/metrics", controller.Metrics)

	req, err := http.NewRequest("GET", "/metrics", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	responseBody := w.Body.String()
	lines := strings.Split(responseBody, "\n")

	// Check that we have multiple lines
	assert.Greater(t, len(lines), 10, "Metrics should have multiple lines")

	// Check for proper Prometheus format patterns
	helpLines := 0
	typeLines := 0
	metricLines := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# HELP") {
			helpLines++
		} else if strings.HasPrefix(line, "# TYPE") {
			typeLines++
		} else if !strings.HasPrefix(line, "#") {
			metricLines++
			// Check that metric lines have proper format (metric_name value or metric_name{labels} value)
			parts := strings.Fields(line)
			assert.GreaterOrEqual(t, len(parts), 2, "Metric line should have at least metric name and value: %s", line)
		}
	}

	// Should have HELP and TYPE comments
	assert.Greater(t, helpLines, 0, "Should have HELP comments")
	assert.Greater(t, typeLines, 0, "Should have TYPE comments")
	assert.Greater(t, metricLines, 0, "Should have metric lines")
}

func TestHealthController_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Environment: "integration-test",
	}

	controller := NewHealthController(cfg)
	router := gin.New()

	// Set up routes like in the actual application
	router.GET("/health", controller.HealthCheck)
	router.GET("/metrics", controller.Metrics)

	// Test health endpoint
	t.Run("health endpoint integration", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "integration-test", response.Environment)
	})

	// Test metrics endpoint
	t.Run("metrics endpoint integration", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/metrics", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "integration-test")
	})
}

func TestHealthResponse_JSONTags(t *testing.T) {
	// Test that the struct can be properly marshaled and unmarshaled
	originalResponse := HealthResponse{
		Status:      "ok",
		Timestamp:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Service:     "docstore-api",
		Version:     "1.0.0",
		Environment: "test",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(originalResponse)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaledResponse HealthResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, originalResponse.Status, unmarshaledResponse.Status)
	assert.Equal(t, originalResponse.Service, unmarshaledResponse.Service)
	assert.Equal(t, originalResponse.Version, unmarshaledResponse.Version)
	assert.Equal(t, originalResponse.Environment, unmarshaledResponse.Environment)
	assert.True(t, originalResponse.Timestamp.Equal(unmarshaledResponse.Timestamp))
}
