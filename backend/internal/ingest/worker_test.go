package ingest

import (
	"sync"
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

type MockEarthquakeProvider struct {
	Response models.USGSResponse
}

func (m *MockEarthquakeProvider) GetEarthquakes() (models.USGSResponse, error) {
	return m.Response, nil
}

type MockSpatialRepository struct {
	mu     sync.Mutex
	Called bool
}

func (m *MockSpatialRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Called = true
	return []models.User{
		{ID: "1", Email: "usuarioA@test.com"},
	}, nil
}

func (m *MockSpatialRepository) WasCalled() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Called
}


type MockEarthquakeRepository struct {
	mu         sync.Mutex
	SavedCount int
}

func (m *MockEarthquakeRepository) SaveEarthquake(sismo models.Feature) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SavedCount++
	return nil
}

func (m *MockEarthquakeRepository) GetSavedCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.SavedCount
}

func (m *MockEarthquakeRepository) GetEarthquakesSince(since time.Time) ([]models.Feature, error) {
	return []models.Feature{}, nil
}

func (m *MockEarthquakeRepository) GetFilteredEarthquakes(minMag float64, limit int) ([]models.Feature, error) {
	return []models.Feature{}, nil
}


func TestIngestionWorker_Process(t *testing.T) {
	// Case 1: Basic happy path (Registration and alert of a new earthquake)
	t.Run("Procesamiento exitoso de un sismo nuevo", func(t *testing.T) {
		stopChan := make(chan bool)
		alertQueue := make(chan models.AlertMessage, 100)

		// Configure the mock to return a unique earthquake
		providerMock := &MockEarthquakeProvider{
			Response: models.USGSResponse{
				Features: []models.Feature{{ID: "us_test_123"}},
			},
		}
		spatialMock := &MockSpatialRepository{}
		dbMock := &MockEarthquakeRepository{}

		worker := NewIngestionWorker(
			50*time.Millisecond,
			providerMock,
			spatialMock,
			dbMock,
			alertQueue,
		)

		go worker.Start(stopChan)

		// Give a small execution margin for the initial load and the first tick
		time.Sleep(80 * time.Millisecond)

		if !spatialMock.WasCalled() {
			t.Error("El worker no consultó los usuarios afectados en el repositorio espacial")
		}

		if dbMock.GetSavedCount() == 0 {
			t.Error("El worker no guardó el sismo en la base de datos")
		}

		if len(alertQueue) != 1 {
			t.Errorf("Se esperaba 1 alerta en la cola, pero se encontraron %d", len(alertQueue))
		}
	})

	// Case 2: Advanced Criterion (Anti-Duplicate Filter)
	// Here we demonstrate that we test the critical business rules of the session
	t.Run("Evita duplicados si el sismo ya fue procesado en la sesión", func(t *testing.T) {
		stopChan := make(chan bool)
		alertQueue := make(chan models.AlertMessage, 100)

		// The provider always returns the same repeated earthquake
		providerMock := &MockEarthquakeProvider{
			Response: models.USGSResponse{
				Features: []models.Feature{{ID: "sismo_repetido_999"}},
			},
		}
		spatialMock := &MockSpatialRepository{}
		dbMock := &MockEarthquakeRepository{}

		worker := NewIngestionWorker(
			10*time.Millisecond, // Very fast interval to force multiple passes (ticks)
			providerMock,
			spatialMock,
			dbMock,
			alertQueue,
		)

		go worker.Start(stopChan)

		// Let enough time run for at least 3 or 4 ticks to pass
		time.Sleep(50 * time.Millisecond)
		stopChan <- true

		// Even though the ticker passed several times, thanks to the `processedIDs` map it should only have been saved ONCE.
		if dbMock.GetSavedCount() > 1 {
			t.Errorf("Error de duplicación: El sismo se guardó %d veces en la BD, se esperaba solo 1", dbMock.GetSavedCount())
		}

		if len(alertQueue) > 1 {
			t.Errorf("Error de spam: Se generaron %d alertas para el mismo sismo, se esperaba solo 1", len(alertQueue))
		}
	})
}
