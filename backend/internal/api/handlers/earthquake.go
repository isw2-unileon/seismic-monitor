package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type EarthquakeHandler struct {
	Service *services.EarthquakeService
}

func NewEarthquakeHandler(service *services.EarthquakeService) *EarthquakeHandler {
	return &EarthquakeHandler{Service: service}
}

func (h *EarthquakeHandler) GetEarthquakes(c *gin.Context) {
	limit := 50
	minMag := 0.0

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	if m, err := strconv.ParseFloat(c.Query("min_magnitude"), 64); err == nil {
		minMag = m
	}

	// Delegate to the repository via a service if there was a need,
	// for now we can assume the service also has this method
	// or we can create it. To maintain the original interface, ideally
	// the service should expose GetFilteredEarthquakes.
	// We assume the Service was updated to delegate or Repo was used before.
	earthquakes, err := h.Service.GetFilteredEarthquakes(minMag, limit)
	if err != nil {
		fmt.Printf("Error en GetFilteredEarthquakes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los sismos"})
		return
	}

	features := make([]models.Feature, len(earthquakes))
	for i := range earthquakes {
		earthquakes[i].Type = "Feature"
		// GeoJSON FeatureCollection compatible structure
		features[i] = earthquakes[i]
	}

	response := models.USGSResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	if response.Features == nil {
		response.Features = []models.Feature{} // Avoid null in JSON
	}

	c.JSON(http.StatusOK, response)
}

func (h *EarthquakeHandler) GetHistory(c *gin.Context) {
	earthquakes, err := h.Service.GetRecentEarthquakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los sismos de la última hora"})
		return
	}

	features := make([]models.Feature, len(earthquakes))
	for i := range earthquakes {
		earthquakes[i].Type = "Feature"
		// Estructura compatible con FeatureCollection de GeoJSON
		features[i] = earthquakes[i]
	}

	response := models.USGSResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	if response.Features == nil {
		response.Features = []models.Feature{} // Evitar null en JSON
	}

	c.JSON(http.StatusOK, response)
}
