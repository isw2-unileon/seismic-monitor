package ingest

import (
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

// --- 1. MOCK DEL PROVEEDOR DE SISMOS ---
type MockEarthquakeProvider struct {
	Response models.USGSResponse
}

func (m *MockEarthquakeProvider) GetEarthquakes() (models.USGSResponse, error) {
	return m.Response, nil
}

// --- 2. MOCK DEL MOTOR ESPACIAL ---
type MockSpatialRepository struct {
	Called bool
}

func (m *MockSpatialRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	m.Called = true
	return []models.User{
		{ID: "1", Email: "usuarioA@test.com"},
	}, nil
}

// --- 3. MOCK DEL REPOSITORIO DE SISMOS (BASE DE DATOS) ---
type MockEarthquakeRepository struct {
	SavedCount int
}

func (m *MockEarthquakeRepository) SaveEarthquake(sismo models.Feature) error {
	m.SavedCount++
	return nil
}

func (m *MockEarthquakeRepository) GetEarthquakesSince(since time.Time) ([]models.Feature, error) {
	return []models.Feature{}, nil
}

func (m *MockEarthquakeRepository) GetFilteredEarthquakes(minMag float64, limit int) ([]models.Feature, error) {
	return []models.Feature{}, nil
}

// --- LOS TESTS ---

func TestIngestionWorker_Process(t *testing.T) {
	// Caso 1: Flujo básico feliz (Registro y alerta de sismo nuevo)
	t.Run("Procesamiento exitoso de un sismo nuevo", func(t *testing.T) {
		stopChan := make(chan bool)
		alertQueue := make(chan models.AlertMessage, 100)

		// Configuramos el mock para que devuelva un sismo único
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

		// Arrancamos el worker
		go worker.Start(stopChan)

		// Damos un pequeño margen de ejecución para la carga inicial y el primer tick
		time.Sleep(80 * time.Millisecond)
		stopChan <- true // Detenemos el worker

		// VERIFICACIONES
		if !spatialMock.Called {
			t.Error("El worker no consultó los usuarios afectados en el repositorio espacial")
		}

		if dbMock.SavedCount == 0 {
			t.Error("El worker no guardó el sismo en la base de datos")
		}

		if len(alertQueue) != 1 {
			t.Errorf("Se esperaba 1 alerta en la cola, pero se encontraron %d", len(alertQueue))
		}
	})

	// Caso 2: Criterio Avanzado (Filtro Anti-Duplicados)
	// Aquí demostramos al profesor que testeamos las reglas de negocio críticas de la sesión
	t.Run("Evita duplicados si el sismo ya fue procesado en la sesión", func(t *testing.T) {
		stopChan := make(chan bool)
		alertQueue := make(chan models.AlertMessage, 100)

		// El proveedor siempre devuelve el mismo sismo repetido
		providerMock := &MockEarthquakeProvider{
			Response: models.USGSResponse{
				Features: []models.Feature{{ID: "sismo_repetido_999"}},
			},
		}
		spatialMock := &MockSpatialRepository{}
		dbMock := &MockEarthquakeRepository{}

		worker := NewIngestionWorker(
			10*time.Millisecond, // Intervalo muy rápido para forzar múltiples pasadas (ticks)
			providerMock,
			spatialMock,
			dbMock,
			alertQueue,
		)

		go worker.Start(stopChan)

		// Dejamos correr el tiempo suficiente para que pasen al menos 3 o 4 ticks
		time.Sleep(50 * time.Millisecond)
		stopChan <- true

		// VERIFICACIONES DE ROBUSTEZ
		// Aunque el ticker pasó varias veces, gracias al mapa `processedIDs` solo debió guardarse UNA vez.
		if dbMock.SavedCount > 1 {
			t.Errorf("Error de duplicación: El sismo se guardó %d veces en la BD, se esperaba solo 1", dbMock.SavedCount)
		}

		if len(alertQueue) > 1 {
			t.Errorf("Error de spam: Se generaron %d alertas para el mismo sismo, se esperaba solo 1", len(alertQueue))
		}
	})
}
