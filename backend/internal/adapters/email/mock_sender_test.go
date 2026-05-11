package email

import (
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

func TestMockSender_SendAlert(t *testing.T) {
	// 1. Instanciamos el adaptador
	sender := &MockSender{}

	// 2. Preparamos datos falsos
	user := models.User{ID: "1", Email: "usuario_test@ejemplo.com"}
	sismo := models.Feature{ID: "sismo_test_123"}

	// 3. Medimos el tiempo para asegurarnos de que simula el retraso (opcional pero buena práctica)
	start := time.Now()

	// 4. Ejecutamos la función
	err := sender.SendAlert(user, sismo)

	elapsed := time.Since(start)

	// Verificaciones
	if err != nil {
		t.Fatalf("No se esperaba error al simular el envío, se obtuvo: %v", err)
	}

	// Comprobamos que el mock realmente haya simulado un retardo (habíamos puesto 500ms)
	if elapsed < 500*time.Millisecond {
		t.Errorf("El mock fue demasiado rápido, se esperaba un retraso simulado. Tiempo: %v", elapsed)
	}
}
