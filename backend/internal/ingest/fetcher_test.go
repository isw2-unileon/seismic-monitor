package ingest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchSeismicData(t *testing.T) {
	// Creamos un servidor de prueba para no llamar a la API real
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type": "FeatureCollection", "features": []}`))
	}))
	defer server.Close()

	data, err := FetchData(server.URL)

	if err != nil {
		t.Fatalf("Se esperaba error nil, se obtuvo: %v", err)
	}

	if len(data) == 0 {
		t.Error("Se esperaba contenido en la respuesta, se obtuvo vacío")
	}
}
