package handlers

import (
	"net/http"
	"strconv"

	"seismic-monitor/backend/internal/services"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type EarthquakeHandler struct {
	Service *services.EarthquakeService
}

func NewEarthquakeHandler(service *services.EarthquakeService) *EarthquakeHandler {
	return &EarthquakeHandler{Service: service}
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

	// Delegar al repositorio a través de un servicio si existiera la necesidad,
	// por ahora podemos asumir que el servicio también tiene este método 
	// o podemos crearlo. Para mantener la interfaz original, idealmente
	// el servicio debería exponer GetFilteredEarthquakes.
	// Asumiremos que el Service fue actualizado para delegar o se usaba Repo antes.
	earthquakes, err := h.Service.GetFilteredEarthquakes(minMag, limit)
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

// GetHistory devuelve los sismos de la última hora
func (h *EarthquakeHandler) GetHistory(c *gin.Context) {
	earthquakes, err := h.Service.GetRecentEarthquakes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los sismos de la última hora"})
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

