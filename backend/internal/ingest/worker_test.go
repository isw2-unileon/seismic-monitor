package ingest

import (
	"errors"
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

// --- EL MOCK ---
// MockEarthquakeProvider simula ser la API de USGS
type MockEarthquakeProvider struct {
	CallsCount int
	ShouldFail bool
}

// GetEarthquakes es la implementación del mock
func (m *MockEarthquakeProvider) GetEarthquakes() (models.USGSResponse, error) {
	m.CallsCount++
	if m.ShouldFail {
		return models.USGSResponse{}, errors.New("error simulado")
	}
	return models.USGSResponse{}, nil
}

// --- EL TEST ---
func TestStartIngestionWorker_WithMock(t *testing.T) {
	stopChan := make(chan bool)
	mock := &MockEarthquakeProvider{ShouldFail: false}

	// Ejecutamos el worker con un intervalo muy corto para el test
	// Usamos una goroutine porque StartIngestionWorker tiene un bucle infinito
	go StartIngestionWorker(50*time.Millisecond, stopChan, mock)

	// Esperamos a que pasen un par de ticks (150ms debería dar para 2-3 llamadas)
	time.Sleep(150 * time.Millisecond)

	// Verificamos que se haya llamado al menos una vez
	if mock.CallsCount == 0 {
		t.Error("El worker no llamó al proveedor de sismos")
	}

	// Probamos el frenado: Enviamos señal de parada
	stopChan <- true

	// Guardamos las llamadas hechas hasta ahora
	callsBeforeStop := mock.CallsCount

	// Esperamos un poco más para ver si sigue llamando (no debería)
	time.Sleep(100 * time.Millisecond)

	if mock.CallsCount > callsBeforeStop {
		t.Errorf("El worker siguió ejecutándose después de la señal de parada. Llamadas: %d -> %d", callsBeforeStop, mock.CallsCount)
	}
}

func TestStartIngestionWorker_ErrorHandling(t *testing.T) {
	stopChan := make(chan bool)
	// Forzamos que el mock falle
	mock := &MockEarthquakeProvider{ShouldFail: true}

	go StartIngestionWorker(50*time.Millisecond, stopChan, mock)

	time.Sleep(100 * time.Millisecond)
	stopChan <- true

	if mock.CallsCount == 0 {
		t.Error("El worker debería haber llamado al proveedor incluso si falla")
	}
	// Si el test llega aquí sin entrar en pánico, es que el worker gestiona bien los errores
}
