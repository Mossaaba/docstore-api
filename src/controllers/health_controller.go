package controllers

import (
	"net/http"
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