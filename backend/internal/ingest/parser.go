package ingest

import (
	"encoding/json"
	"fmt"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// ParseUSGSData toma los bytes crudos del JSON y los convierte en el struct USGSResponse
func ParseUSGSData(data []byte) (models.USGSResponse, error) {
	var response models.USGSResponse

	err := json.Unmarshal(data, &response)
	if err != nil {
		return response, fmt.Errorf("error al deserializar el JSON del USGS: %w", err)
	}

	return response, nil
}
