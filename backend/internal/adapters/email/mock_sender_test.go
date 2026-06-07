package email

import (
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

func TestMockSender_SendAlert(t *testing.T) {
	sender := &MockSender{}

	user := models.User{ID: "1", Email: "usuario_test@ejemplo.com"}
	sismo := models.Feature{ID: "sismo_test_123"}

	// Measure time to ensure it simulates the delay (optional but good practice)
	start := time.Now()

	err := sender.SendAlert(user, sismo)

	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("No se esperaba error al simular el envío, se obtuvo: %v", err)
	}

	// Verify that the mock actually simulated a delay (we set 500ms)
	if elapsed < 500*time.Millisecond {
		t.Errorf("El mock fue demasiado rápido, se esperaba un retraso simulado. Tiempo: %v", elapsed)
	}
}
