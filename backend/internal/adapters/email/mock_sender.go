package email

import (
	"log"
	"time"

	"seismic-monitor/backend/internal/models"
)

// MockSender is an adapter that simulates sending emails
type MockSender struct{}

func (m *MockSender) SendAlert(user models.User, sismo models.Feature) error {
	time.Sleep(500 * time.Millisecond)

	log.Printf("[EMAIL SENT] To: %s | Subject: DANGER! Earthquake %s detected near you.", user.Email, sismo.ID)
	return nil
}
