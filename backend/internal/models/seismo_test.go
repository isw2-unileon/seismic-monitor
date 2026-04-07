package models

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalUSGS(t *testing.T) {
	// Simulación de un fragmento real de la API del USGS
	jsonRaw := []byte(`{
		"id": "us7000lz8x",
		"properties": {
			"mag": 4.5,
			"place": "20km S of Quito, Ecuador",
			"time": 1673456789000
		},
		"geometry": {
			"coordinates": [-78.5, -0.2, 10.0]
		}
	}`)

	var s Earthquake
	err := json.Unmarshal(jsonRaw, &s)

	if err != nil {
		t.Fatalf("Error al deserializar: %v", err)
	}

	if s.ID != "us7000lz8x" {
		t.Errorf("ID incorrecto: se esperaba us7000lz8x, se obtuvo %s", s.ID)
	}

	if s.Info.Mag != 4.5 {
		t.Errorf("Magnitud incorrecta: %.1f", s.Info.Mag)
	}
}
