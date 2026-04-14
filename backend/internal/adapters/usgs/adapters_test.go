package usgs

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUSGSAdapter_GetEarthquakes(t *testing.T) {
	// 1. Usamos tu JSON mockeado del test del parser
	mockJSON := `{"features": [{"id": "us1000", "properties": {"mag": 5.5, "place": "10km Sur de León, España", "time": 1673456789000}}]}`

	// 2. Usamos tu servidor mockeado del test del fetcher
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	// 3. Instanciamos nuestro adaptador con la URL falsa
	adapter := &USGSAdapter{URL: server.URL}

	// 4. Ejecutamos la función de la interfaz
	response, err := adapter.GetEarthquakes()

	if err != nil {
		t.Fatalf("No se esperaba error, se obtuvo: %v", err)
	}

	if len(response.Features) != 1 {
		t.Fatalf("Se esperaba 1 sismo, se obtuvieron %d", len(response.Features))
	}

	if response.Features[0].ID != "us1000" {
		t.Errorf("ID incorrecto, se obtuvo: %s", response.Features[0].ID)
	}
}
