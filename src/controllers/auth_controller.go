package controllers

import (
	"docstore-api/src/config"
	"docstore-api/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	config *config.Config
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		config: cfg,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simple authentication - verify against configured credentials
	if req.Username == ctrl.config.AdminUser && req.Password == ctrl.config.AdminPass {
		token, err := middleware.GenerateToken(req.Username, ctrl.config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			User:  req.Username,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}
