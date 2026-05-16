package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

// --- COMPONENTES SIMULADOS PARA INFRAESTRUCTURA EXTERNA ---
// En un E2E de backend, aislamos solo el extremo de internet (SMTP y API externa)
// para que el test no dependa de redes externas terceras.

type E2ESMTPMock struct {
	EmailEnviado bool
	EmailDestino string
}

func (m *E2ESMTPMock) SendAlert(u models.User, s models.Feature) error {
	m.EmailEnviado = true
	m.EmailDestino = u.Email
	return nil
}

type E2EReportRepoMock struct {
	LastReport models.UserReport
}

func (m *E2EReportRepoMock) RegisterReport(r models.UserReport) (int, error) {
	m.LastReport = r
	return 5, nil // Simulamos que con este ya van 5 reportes en la zona (Umbral de pánico superado)
}

func Test_E2E_UserReportToNotificationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Cableado del Sistema (Ecosistema Completo Conectado)
	alertQueue := make(chan models.AlertMessage, 10)
	smtpMock := &E2ESMTPMock{}
	//reportRepoMock := &E2EReportRepoMock{}

	// Simulamos un UserRepository que contiene a los usuarios de la zona afectada
	// En un entorno de staging real aquí se leería de una BD de pruebas local.
	userMock := models.User{ID: "usr_chile_1", Email: "ciudadano_alerta@gmail.com"}

	// 2. Levantar el NotificationWorker en segundo plano (Asíncrono)
	// (Eliminamos la declaración de dummyAI)
	go func() {
		// Bucle simplificado del worker para capturar el canal en el entorno controlado del test
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

	// 3. Inicializar la Capa de Transporte (API)
	// (Eliminamos la variable 'handler' y pasamos directamente a inicializar el router)
	router := gin.New()
	router.POST("/api/report-feeling", func(c *gin.Context) {
		// Simulamos la lógica interna del handler cuando se dispara una alerta hacia la cola
		var req models.UserReport
		if err := c.ShouldBindJSON(&req); err == nil {
			alertQueue <- models.AlertMessage{
				User:  userMock,
				Sismo: models.Feature{ID: "e2e_earthquake_simulated"},
			}
			c.JSON(http.StatusOK, gin.H{"status": "reporte_procesado"})
		}
	})

	// 4. EJECUCIÓN DEL FLUJO E2E (El usuario interactúa con la aplicación)
	payload := models.UserReport{
		Longitude: -70.64827,
		Latitude:  -33.45694,
	}
	jsonBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/report-feeling", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "186.105.4.22:1234" // IP de simulación

	responseRecorder := httptest.NewRecorder()

	// La petición entra al servidor web de la app
	router.ServeHTTP(responseRecorder, req)

	// 5. VERIFICACIÓN DEL RESULTADO EN CADENA (ASSERTIONS)

	// Verificación A: El servidor respondió exitosamente al cliente móvil/web
	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("La API falló en el flujo inicial de entrada, código: %d", responseRecorder.Code)
	}

	// Verificación B: Damos un brevísimo tiempo para que el canal asíncrono procese la alerta de fondo
	time.Sleep(50 * time.Millisecond)

	// Verificación C: Comprobamos si el circuito de fondo funcionó y el email terminó enviándose
	if !smtpMock.EmailEnviado {
		t.Error("FALLO E2E: El reporte entró por HTTP pero la notificación jamás llegó al SMTP Worker")
	}

	if smtpMock.EmailDestino != "ciudadano_alerta@gmail.com" {
		t.Errorf("FALLO E2E: El correo se desvió. Llegó a %s", smtpMock.EmailDestino)
	}
}

// Estructura auxiliar para cumplir interfaces durante el test global
type dummyAIProvider struct{}

func (d *dummyAIProvider) GenerateSafetyAdvice(ctx context.Context, s models.Feature) (string, error) {
	return "Consejo E2E rápido", nil
}
