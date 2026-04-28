package email

import (
	"log"
	"time"

	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
)

// MockSender es un adaptador que simula el envío de correos
type MockSender struct{}

func (m *MockSender) SendAlert(user models.User, sismo models.Feature) error {
	// Simulamos el tiempo que tarda en enviarse un correo por internet
	time.Sleep(500 * time.Millisecond)

	log.Printf("[EMAIL ENVIADO] Para: %s | Asunto: ¡PELIGRO! Sismo %s detectado cerca de ti.", user.Email, sismo.ID)
	return nil
}
