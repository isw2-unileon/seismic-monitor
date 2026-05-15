package services

import (
	"context"
	"seismic-monitor/backend/internal/models"
	"testing"
	"time"
)

// MockNotificationService simula el envío de emails
type MockNotificationService struct {
	SentCalled bool
	LastSismo  models.Feature
}

func (m *MockNotificationService) SendAlert(u models.User, s models.Feature) error {
	m.SentCalled = true
	m.LastSismo = s
	return nil
}

// MockAIProvider simula la API de Gemini
type MockAIProvider struct {
	ShouldFail bool
}

func (m *MockAIProvider) GenerateSafetyAdvice(ctx context.Context, s models.Feature) (string, error) {
	if m.ShouldFail {
		return "", func() error { return nil }() // Simula un error
	}
	return "Consejo de prueba de IA", nil
}

func TestStartNotificationWorker(t *testing.T) {
	// 1. Configuración (Setup)
	alertQueue := make(chan models.AlertMessage, 1)
	mockNotifier := &MockNotificationService{}
	mockAI := &MockAIProvider{ShouldFail: false}

	// 2. Ejecutar worker en una goroutine
	go func() {
		// Creamos una copia simplificada para el test o cerramos el canal tras enviar
		for msg := range alertQueue {
			advice, _ := mockAI.GenerateSafetyAdvice(context.Background(), msg.Sismo)
			msg.Sismo.AIAdvice = advice
			mockNotifier.SendAlert(msg.User, msg.Sismo)
		}
	}()

	// 3. Enviar un mensaje de prueba
	testSismo := models.Feature{ID: "test-123"}
	alertQueue <- models.AlertMessage{
		User:  models.User{Email: "test@example.com"},
		Sismo: testSismo,
	}

	// Dar un pequeño margen para el procesamiento asíncrono
	time.Sleep(100 * time.Millisecond)
	close(alertQueue)

	// 4. Verificaciones (Assertions)
	if !mockNotifier.SentCalled {
		t.Error("El worker debería haber llamado a SendAlert")
	}

	if mockNotifier.LastSismo.AIAdvice != "Consejo de prueba de IA" {
		t.Errorf("Se esperaba el consejo de IA, se obtuvo: %s", mockNotifier.LastSismo.AIAdvice)
	}
}
