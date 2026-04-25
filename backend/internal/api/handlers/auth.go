package handlers

import (
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

// Register maneja el registro de nuevos usuarios
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de registro inválidos"})
		return
	}

	// Verificar si el usuario ya existe
	existing, _ := h.Repo.FindUserByEmail(req.Email)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "el correo ya está registrado"})
		return
	}

	// Hashear contraseña
	hashed, err := database.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar la contraseña"})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashed,
		AlertRadius:  100, // Valor por defecto en km
	}

	if err := h.Repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "usuario registrado con éxito"})
}

// Login maneja la autenticación y devuelve un JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := h.ShouldBind(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email o contraseña requeridos"})
		return
	}

	user, err := h.Repo.FindUserByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if !database.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	token, err := h.JWTService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email": user.Email,
			"id":    user.ID,
		},
	})
}

// Helper para binding (evita repetir código)
func (h *AuthHandler) ShouldBind(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}
