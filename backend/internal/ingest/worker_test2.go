package ingest

import (
	"testing"
	"time"

	"github.com/isw2-unileon/seismic-monitor/backend/internal/models"
)

// --- MOCK DEL PROVEEDOR DE SISMOS ---
type MockEarthquakeProvider2 struct{}

func (m *MockEarthquakeProvider2) GetEarthquakes() (models.USGSResponse, error) {
	// Devolvemos un sismo de prueba
	return models.USGSResponse{
		Features: []models.Feature{
			{ID: "us_test_123"},
		},
	}, nil
}

// --- MOCK DEL MOTOR ESPACIAL ---
type MockSpatialRepository struct {
	Called bool
}

func (m *MockSpatialRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	m.Called = true
	// Simulamos que el sismo afectó a 3 usuarios
	return []models.User{
		{ID: 1, Email: "usuarioA@test.com"},
		{ID: 2, Email: "usuarioB@test.com"},
		{ID: 3, Email: "usuarioC@test.com"},
	}, nil
}

// --- EL TEST ---
func TestStartIngestionWorker_CollisionDetection(t *testing.T) {
	stopChan := make(chan bool)

	providerMock := &MockEarthquakeProvider2{}
	spatialMock := &MockSpatialRepository{}

	// Arrancamos el worker con nuestros mocks
	go StartIngestionWorker(50*time.Millisecond, stopChan, providerMock, spatialMock)

	// Damos tiempo a que se ejecute al menos un tick
	time.Sleep(100 * time.Millisecond)
	stopChan <- true // Detenemos el worker

	// Verificamos el criterio de éxito: ¿Llamó el worker al motor espacial?
	if !spatialMock.Called {
		t.Error("El worker no llamó a la función GetAffectedUsers del repositorio espacial")
	}
}
