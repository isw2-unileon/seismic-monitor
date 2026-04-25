package handlers

import (
	"net/http"
	"strconv"

	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type EarthquakeHandler struct {
	Repo *database.EarthquakeRepository
}

func NewEarthquakeHandler(repo *database.EarthquakeRepository) *EarthquakeHandler {
	return &EarthquakeHandler{Repo: repo}
}

// GetEarthquakes devuelve una FeatureCollection GeoJSON de sismos
func (h *EarthquakeHandler) GetEarthquakes(c *gin.Context) {
	// Filtros por defecto
	limit := 50
	minMag := 0.0

	// Leer parámetros de la URL
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	if m, err := strconv.ParseFloat(c.Query("min_magnitude"), 64); err == nil {
		minMag = m
	}

	earthquakes, err := h.Repo.GetFilteredEarthquakes(minMag, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los sismos"})
		return
	}

	// Asegurar que cada sismo tiene el tipo "Feature"
	for i := range earthquakes {
		earthquakes[i].Type = "Feature"
	}

	response := models.FeatureCollection{
		Type:     "FeatureCollection",
		Features: earthquakes,
	}

	if response.Features == nil {
		response.Features = []models.Earthquake{} // Evitar null en JSON
	}

	c.JSON(http.StatusOK, response)
}
