package handlers

import (
	"fmt"
	"net/http"

	"seismic-monitor/backend/internal/auth"
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Repo       *database.UserRepository
	JWTService *auth.JWTService
}

func NewAuthHandler(repo *database.UserRepository, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		Repo:       repo,
		JWTService: jwtService,
	}
}

// Register handles the registration of new users
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid registration data"})
		return
	}

	existing, _ := h.Repo.FindUserByEmail(req.Email)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashed, err := database.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing password"})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashed,
		AlertRadius:  100,
		MinMagnitude: 3.0,
	}

	if err := h.Repo.CreateUser(user); err != nil {
		fmt.Printf("Error during CreateUser: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login handles authentication and returns a JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := h.ShouldBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	user, err := h.Repo.FindUserByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !database.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email":         user.Email,
			"id":            user.ID,
			"alert_radius":  user.AlertRadius,
			"min_magnitude": user.MinMagnitude,
		},
	})
}

func (h *AuthHandler) ShouldBind(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}
