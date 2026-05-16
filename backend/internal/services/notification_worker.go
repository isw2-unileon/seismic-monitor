package services

import (
	"context"
	"log"
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
	"time"
)

// StartNotificationWorker escucha la cola y envía los correos
func StartNotificationWorker(
	alertQueue <-chan models.AlertMessage,
	notifier ports.NotificationService,
	aiProvider ports.AIProvider,
) {
	log.Println("[Notification Worker] Iniciado y esperando alertas...")

	for msg := range alertQueue {
		// 1. Generar consejo de IA (usamos un contexto con timeout por seguridad)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		advice, err := aiProvider.GenerateSafetyAdvice(ctx, msg.Sismo)
		cancel()

		if err != nil {
			log.Printf("[Notification Worker] Error IA: %v. Usando fallback.", err)
			advice = "Mantente informado por canales oficiales y sigue los protocolos de seguridad."
		}

		// 2. Inyectar el consejo en el sismo antes de enviarlo
		msg.Sismo.AIAdvice = advice

		// 3. Enviar la alerta final (vía SMTP, Mock, etc.)
		err = notifier.SendAlert(msg.User, msg.Sismo)
		if err != nil {
			log.Printf("[Notification Worker] Error enviando alerta a %s: %v", msg.User.Email, err)
		}
	}
}
