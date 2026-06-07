package handlers

import (
	"fmt"
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

	err := h.Repo.UpdateUserLocation(userID.(string), req.Name, *req.Latitude, *req.Longitude, *req.AlertRadius, *req.MinMagnitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "perfil actualizado correctamente",
	})
}

// AddLocation añade una nueva zona de interés
func (h *UserHandler) AddLocation(c *gin.Context) {
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

	// Usamos el min_magnitude del request si existe, si no 3.0 por defecto
	minMag := 3.0
	if req.MinMagnitude != nil {
		minMag = *req.MinMagnitude
	}

	newID, err := h.Repo.AddUserAlertCenter(userID.(string), *req.Latitude, *req.Longitude, *req.AlertRadius, minMag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al añadir la ubicación"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ubicación añadida correctamente",
		"id":      newID,
	})
}

// DeleteLocation elimina un centro de alerta
func (h *UserHandler) DeleteLocation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo identificar al usuario"})
		return
	}

	centerID := c.Param("id")
	if centerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de centro requerido"})
		return
	}

	if err := h.Repo.DeleteUserAlertCenter(userID.(string), centerID); err != nil {
		// Log para depuración
		fmt.Printf("Error al borrar ubicación %s: %v\n", centerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar la ubicación de la base de datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ubicación eliminada correctamente"})
}



