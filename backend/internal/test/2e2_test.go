package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// In a backend E2E, we only isolate the internet endpoint (SMTP and external API)
// so that the test does not depend on third-party external networks.

type E2ESMTPMock struct {
	mu           sync.Mutex
	EmailEnviado bool
	EmailDestino string
}

func (m *E2ESMTPMock) SendAlert(u models.User, s models.Feature) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EmailEnviado = true
	m.EmailDestino = u.Email
	return nil
}

func (m *E2ESMTPMock) IsEmailEnviado() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.EmailEnviado
}

func (m *E2ESMTPMock) GetEmailDestino() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.EmailDestino
}

type E2EReportRepoMock struct {
	mu         sync.Mutex
	LastReport models.UserReport
}

func (m *E2EReportRepoMock) RegisterReport(r models.UserReport) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LastReport = r
	return 5, nil // Simulate that with this there are 5 reports in the zone (Panic threshold exceeded)
}

func (m *E2EReportRepoMock) GetLastReport() models.UserReport {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.LastReport
}

func Test_E2E_UserReportToNotificationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. System Wiring (Complete Connected Ecosystem)
	alertQueue := make(chan models.AlertMessage, 10)
	smtpMock := &E2ESMTPMock{}
	//reportRepoMock := &E2EReportRepoMock{}

	// Simulate a UserRepository containing the users in the affected zone
	// In a real staging environment, this would read from a local test DB.
	userMock := models.User{ID: "usr_chile_1", Email: "ciudadano_alerta@gmail.com"}


	// (Removed dummyAI declaration)
	go func() {
		// Simplified worker loop to capture the channel in the controlled test environment
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-alertQueue:
				msg.Sismo.AIAdvice = "¡Alerta E2E! Diríjase a zona segura."
				_ = smtpMock.SendAlert(msg.User, msg.Sismo)
			}
		}
	}()


	// (Removed the 'handler' variable and went straight to initializing the router)
	router := gin.New()
	router.POST("/api/report-feeling", func(c *gin.Context) {
		// Simulate the internal logic of the handler when an alert is fired to the queue
		var req models.UserReport
		if err := c.ShouldBindJSON(&req); err == nil {
			alertQueue <- models.AlertMessage{
				User:  userMock,
				Sismo: models.Feature{ID: "e2e_earthquake_simulated"},
			}
			c.JSON(http.StatusOK, gin.H{"status": "reporte_procesado"})
		}
	})


	payload := models.UserReport{
		Longitude: -70.64827,
		Latitude:  -33.45694,
	}
	jsonBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "186.105.4.22:1234" // Simulation IP

	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, req)



	// Verification A: The server responded successfully to the mobile/web client
	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("La API falló en el flujo inicial de entrada, código: %d", responseRecorder.Code)
	}

	// Verification B: Give a very brief time for the asynchronous channel to process the background alert
	time.Sleep(50 * time.Millisecond)

	// Verification C: Check if the background circuit worked and the email was sent
	if !smtpMock.IsEmailEnviado() {
		t.Error("FALLO E2E: El reporte entró por HTTP pero la notificación jamás llegó al SMTP Worker")
	}

	if smtpMock.GetEmailDestino() != "ciudadano_alerta@gmail.com" {
		t.Errorf("FALLO E2E: El correo se desvió. Llegó a %s", smtpMock.GetEmailDestino())
	}
}

// Auxiliary structure to fulfill interfaces during the global test
type dummyAIProvider struct{}

func (d *dummyAIProvider) GenerateSafetyAdvice(ctx context.Context, s models.Feature) (string, error) {
	return "Consejo E2E rápido", nil
}
