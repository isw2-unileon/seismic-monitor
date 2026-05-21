package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"seismic-monitor/backend/internal/api/handlers"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Mock básico para cumplir con la interfaz del repositorio de reportes
type MockReportRepo struct{}

func (m *MockReportRepo) RegisterReport(r models.UserReport) (int, error) {
	return 1, nil
}

func TestReportAPI_Integration(t *testing.T) {
	// Configurar Gin en modo test para no ensuciar la consola de logs
	gin.SetMode(gin.TestMode)

	// 1. Inicializar dependencias.
	// Pasamos nil a UserRepo y AlertQueue porque al devolver conteo=1, el handler jamás los llamará.
	repo := &MockReportRepo{}
	handler := handlers.NewReportHandler(repo, nil, nil)

	// 2. Configurar el servidor/enrutador de Gin simulando la app real
	router := gin.New()
	router.POST("/api/report-feeling", handler.HandleReport)

	// --- CASO 1: Petición Válida Exitosa (200 OK) ---
	t.Run("Petición válida registra reporte exitosamente", func(t *testing.T) {
		payload := models.UserReport{
			Longitude: -70.64827, // Santiago de Chile
			Latitude:  -33.45694,
		}
		jsonBytes, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.10:1234" // Simulamos la IP del usuario 1

		// Creamos un grabador de respuestas HTTP
		responseRecorder := httptest.NewRecorder()

		// Lanzamos la petición al enrutador
		router.ServeHTTP(responseRecorder, req)

		// Verificaciones
		if responseRecorder.Code != http.StatusOK {
			t.Errorf("Se esperaba código 200, se obtuvo %d", responseRecorder.Code)
		}
	})

	// --- CASO 2: Validación de Coordenadas Erróneas (400 Bad Request) ---
	t.Run("Coordenadas fuera de los límites de la Tierra devuelve 400", func(t *testing.T) {
		payload := models.UserReport{
			Longitude: 250.0, // Longitud imposible (> 180)
			Latitude:  45.0,
		}
		jsonBytes, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.11:1234"

		responseRecorder := httptest.NewRecorder()
		router.ServeHTTP(responseRecorder, req)

		if responseRecorder.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba código 400 por coordenadas corruptas, se obtuvo %d", responseRecorder.Code)
		}
	})

	// --- CASO 3: Ataque/Abuso de Spam (429 Too Many Requests) ---
	t.Run("Peticiones duplicadas desde la misma IP activan el Anti-Spam", func(t *testing.T) {
		payload := models.UserReport{Longitude: 10.0, Latitude: 10.0}
		jsonBytes, _ := json.Marshal(payload)
		targetIP := "10.0.0.5"

		// Primera petición de la IP 10.0.0.5 -> Debería dejarle pasar (200)
		req1, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req1.Header.Set("Content-Type", "application/json")
		req1.RemoteAddr = targetIP + ":4422" // IP + puerto efímero de red
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		if w1.Code != http.StatusOK {
			t.Fatalf("La primera petición falló inesperadamente con %d", w1.Code)
		}

		// Segunda petición INMEDIATA desde la misma IP 10.0.0.5 -> Debe bloquearse (429)
		req2, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req2.Header.Set("Content-Type", "application/json")
		req2.RemoteAddr = targetIP + ":4423" // Misma IP, otro puerto
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusTooManyRequests {
			t.Errorf("Anti-Spam falló: Se esperaba código 429, se obtuvo %d", w2.Code)
		}
	})
}
