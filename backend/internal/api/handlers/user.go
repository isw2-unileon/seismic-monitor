package handlers

import (
	"net/http"

	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo *database.UserRepository
}

func NewUserHandler(repo *database.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// UpdateLocation permite al usuario actualizar su posición y radio de alerta
func (h *UserHandler) UpdateLocation(c *gin.Context) {
	// Obtener userID del contexto (puesto por el middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo identificar al usuario"})
		return
	}

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de ubicación inválidos"})
		return
	}

	if err := h.Repo.UpdateUserLocation(userID.(int), req.Latitude, req.Longitude, req.AlertRadius); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar la ubicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ubicación actualizada correctamente"})
}
