package email

import (
	"log"
	"time"

	"seismic-monitor/backend/internal/models"
)

type MockSender struct{}

func (m *MockSender) SendAlert(user models.User, sismo models.Feature) error {
	// Simulate the time it takes to send an email over the internet
	time.Sleep(500 * time.Millisecond)

	log.Printf("[EMAIL ENVIADO] Para: %s | Asunto: ¡PELIGRO! Sismo %s detectado cerca de ti.", user.Email, sismo.ID)
	return nil
}
