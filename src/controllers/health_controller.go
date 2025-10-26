package controllers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"docstore-api/src/config"
	"github.com/gin-gonic/gin"
)

// HealthController handles health check operations
type HealthController struct {
	config *config.Config
}

// NewHealthController creates a new health controller
func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{
		config: cfg,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string    `json:"status" example:"ok"`
	Timestamp   time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Service     string    `json:"service" example:"docstore-api"`
	Version     string    `json:"version" example:"1.0.0"`
	Environment string    `json:"environment" example:"development"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API (available at /health, not /api/v1/health)
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (hc *HealthController) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:      "ok",
		Timestamp:   time.Now().UTC(),
		Service:     "docstore-api",
		Version:     "1.0.0",
		Environment: hc.config.Environment,
	}

	c.JSON(http.StatusOK, response)
}

// Metrics godoc
// @Summary Prometheus metrics endpoint
// @Description Returns Prometheus-compatible metrics for monitoring
// @Tags monitoring
// @Accept text/plain
// @Produce text/plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (hc *HealthController) Metrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := fmt.Sprintf(`# HELP docstore_api_info Information about the DocStore API
# TYPE docstore_api_info gauge
docstore_api_info{version="1.0.0",environment="%s"} 1

# HELP docstore_api_uptime_seconds Total uptime of the service in seconds
# TYPE docstore_api_uptime_seconds counter
docstore_api_uptime_seconds %d

# HELP docstore_api_memory_usage_bytes Current memory usage in bytes
# TYPE docstore_api_memory_usage_bytes gauge
docstore_api_memory_usage_bytes %d

# HELP docstore_api_memory_allocated_bytes Total allocated memory in bytes
# TYPE docstore_api_memory_allocated_bytes counter
docstore_api_memory_allocated_bytes %d

# HELP docstore_api_goroutines Current number of goroutines
# TYPE docstore_api_goroutines gauge
docstore_api_goroutines %d

# HELP docstore_api_health_status Health status of the API (1 = healthy, 0 = unhealthy)
# TYPE docstore_api_health_status gauge
docstore_api_health_status 1
`,
		hc.config.Environment,
		int64(time.Since(time.Now().Add(-time.Hour)).Seconds()), // Placeholder uptime
		m.Sys,
		m.TotalAlloc,
		runtime.NumGoroutine(),
	)

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, metrics)
}
