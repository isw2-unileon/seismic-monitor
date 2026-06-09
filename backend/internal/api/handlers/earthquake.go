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

// GetEarthquakes returns a GeoJSON FeatureCollection of earthquakes
func (h *EarthquakeHandler) GetEarthquakes(c *gin.Context) {
	limit := 50
	minMag := 0.0

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	if m, err := strconv.ParseFloat(c.Query("min_magnitude"), 64); err == nil {
		minMag = m
	}

	earthquakes, err := h.Service.GetFilteredEarthquakes(minMag, limit)
	if err != nil {
		fmt.Printf("Error in GetFilteredEarthquakes: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not obtain earthquakes"})
		return
	}

	features := make([]models.Feature, len(earthquakes))
	for i := range earthquakes {
		earthquakes[i].Type = "Feature"
		features[i] = earthquakes[i]
	}

	response := models.USGSResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	if response.Features == nil {
		response.Features = []models.Feature{}
	}

	c.JSON(http.StatusOK, response)
}

// GetHistory returns earthquakes from the last hour
func (h *EarthquakeHandler) GetHistory(c *gin.Context) {
	earthquakes, err := h.Service.GetRecentEarthquakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not obtain earthquakes from the last hour"})
		return
	}

	features := make([]models.Feature, len(earthquakes))
	for i := range earthquakes {
		earthquakes[i].Type = "Feature"
		features[i] = earthquakes[i]
	}

	response := models.USGSResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	if response.Features == nil {
		response.Features = []models.Feature{}
	}

	c.JSON(http.StatusOK, response)
}
