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

// Basic mock to fulfill the report repository interface
type MockReportRepo struct{}

func (m *MockReportRepo) RegisterReport(r models.UserReport) (int, error) {
	return 1, nil
}

func TestReportAPI_Integration(t *testing.T) {
	// Configure Gin in test mode to avoid polluting the log console
	gin.SetMode(gin.TestMode)

	// Pass nil to UserRepo and AlertQueue because returning count=1 means the handler will never call them.
	repo := &MockReportRepo{}
	handler := handlers.NewReportHandler(repo, nil, nil)

	router := gin.New()
	router.POST("/api/report-feeling", handler.HandleReport)


	t.Run("Petición válida registra reporte exitosamente", func(t *testing.T) {
		payload := models.UserReport{
			Longitude: -70.64827, // Santiago de Chile
			Latitude:  -33.45694,
		}
		jsonBytes, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "192.168.1.10:1234" // Simulate user 1's IP

		responseRecorder := httptest.NewRecorder()

		router.ServeHTTP(responseRecorder, req)

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("Se esperaba código 200, se obtuvo %d", responseRecorder.Code)
		}
	})


	t.Run("Coordenadas fuera de los límites de la Tierra devuelve 400", func(t *testing.T) {
		payload := models.UserReport{
			Longitude: 250.0, // Impossible longitude (> 180)
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


	t.Run("Peticiones duplicadas desde la misma IP activan el Anti-Spam", func(t *testing.T) {
		payload := models.UserReport{Longitude: 10.0, Latitude: 10.0}
		jsonBytes, _ := json.Marshal(payload)
		targetIP := "10.0.0.5"

		// First request from IP 10.0.0.5 -> Should pass (200)
		req1, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req1.Header.Set("Content-Type", "application/json")
		req1.RemoteAddr = targetIP + ":4422" // IP + ephemeral network port
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		if w1.Code != http.StatusOK {
			t.Fatalf("La primera petición falló inesperadamente con %d", w1.Code)
		}

		// IMMEDIATE second request from the same IP 10.0.0.5 -> Must be blocked (429)
		req2, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
		req2.Header.Set("Content-Type", "application/json")
		req2.RemoteAddr = targetIP + ":4423" // Same IP, different port
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusTooManyRequests {
			t.Errorf("Anti-Spam falló: Se esperaba código 429, se obtuvo %d", w2.Code)
		}
	})
}
