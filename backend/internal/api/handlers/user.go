package handlers

import (
	"fmt"
	"log/slog"
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

func (h *UserHandler) UpdateLocation(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo identificar al usuario"})
		return
	}

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Error binding JSON in UpdateLocation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de ubicación inválidos"})
		return
	}

	slog.Info("Updating user profile", "userID", userID, "name", req.Name, "lat", *req.Latitude, "lng", *req.Longitude)

	err := h.Repo.UpdateUserLocation(userID.(string), req.Name, *req.Latitude, *req.Longitude, *req.AlertRadius, *req.MinMagnitude)
	if err != nil {
		slog.Error("Error updating user profile in DB", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "perfil actualizado correctamente",
	})
}

func (h *UserHandler) AddLocation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo identificar al usuario"})
		return
	}

	var req models.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Error binding JSON in AddLocation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de ubicación inválidos"})
		return
	}

	slog.Info("Adding alert center", "userID", userID, "lat", *req.Latitude, "lng", *req.Longitude)

	// Use the requested min_magnitude if it exists, otherwise default to 3.0
	minMag := 3.0
	if req.MinMagnitude != nil {
		minMag = *req.MinMagnitude
	}

	newID, err := h.Repo.AddUserAlertCenter(userID.(string), *req.Latitude, *req.Longitude, *req.AlertRadius, minMag)
	if err != nil {
		slog.Error("Error adding alert center in DB", "error", err, "userID", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al añadir la ubicación"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ubicación añadida correctamente",
		"id":      newID,
	})
}

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
		// Log for debugging
		fmt.Printf("Error al borrar ubicación %s: %v\n", centerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar la ubicación de la base de datos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ubicación eliminada correctamente"})
}
