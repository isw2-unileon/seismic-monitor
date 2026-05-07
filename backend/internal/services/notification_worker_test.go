package services

import (
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"
)

// --- ESPÍA (SPY MOCK) ---
// TestSpyNotifier nos permite registrar qué correos intentó enviar el worker
type TestSpyNotifier struct {
	CallsCount int
	LastUser   models.User
	LastSismo  models.Feature
}

// SendAlert cumple con ports.NotificationService
func (m *TestSpyNotifier) SendAlert(user models.User, sismo models.Feature) error {
	m.CallsCount++
	m.LastUser = user
	m.LastSismo = sismo
	return nil
}

// --- EL TEST ---
func TestStartNotificationWorker(t *testing.T) {
	// 1. Preparamos la tubería (canal) y el espía
	alertQueue := make(chan models.AlertMessage, 5) // Buffer de 5
	spyNotifier := &TestSpyNotifier{}

	// 2. Arrancamos el worker en segundo plano
	go StartNotificationWorker(alertQueue, spyNotifier)

	// 3. Preparamos el mensaje de prueba
	testUser := models.User{ID: "99", Email: "peligro@test.com"}
	testSismo := models.Feature{ID: "earthquake_99"}

	// 4. Mandamos el trabajo por la cola
	alertQueue <- models.AlertMessage{
		User:  testUser,
		Sismo: testSismo,
	}

	// 5. Damos un pequeño respiro para que la goroutine lea el canal
	time.Sleep(50 * time.Millisecond)

	// 6. Verificamos los resultados
	if spyNotifier.CallsCount != 1 {
		t.Errorf("Se esperaba que el notificador se llamara 1 vez, se llamó %d veces", spyNotifier.CallsCount)
	}

	if spyNotifier.LastUser.Email != "peligro@test.com" {
		t.Errorf("El correo enviado es incorrecto. Se obtuvo: %s", spyNotifier.LastUser.Email)
	}

	if spyNotifier.LastSismo.ID != "earthquake_99" {
		t.Errorf("El sismo reportado es incorrecto. Se obtuvo: %v", spyNotifier.LastSismo.ID)
	}
}
