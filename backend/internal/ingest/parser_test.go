package ingest

import (
	"testing"
)

func TestParseUSGSData(t *testing.T) {
	mockJSON := []byte(`{
		"features": [
			{
				"id": "us1000",
				"properties": {
					"mag": 5.5,
					"place": "10km Sur de León, España",
					"time": 1673456789000
				}
			}
		]
	}`)

	response, err := ParseUSGSData(mockJSON)

	if err != nil {
		t.Fatalf("No se esperaba error, se obtuvo: %v", err)
	}

	if len(response.Features) != 1 {
		t.Fatalf("Se esperaba 1 sismo, se obtuvieron %d", len(response.Features))
	}

	if response.Features[0].ID != "us1000" {
		t.Errorf("ID incorrecto, se obtuvo: %s", response.Features[0].ID)
	}

	if response.Features[0].Info.Mag != 5.5 {
		t.Errorf("Magnitud incorrecta, se obtuvo: %f", response.Features[0].Info.Mag)
	}
}
