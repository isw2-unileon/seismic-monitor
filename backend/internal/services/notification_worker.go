package services

import (
	"log"
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
)

// StartNotificationWorker escucha la cola y envía los correos
func StartNotificationWorker(
	alertQueue <-chan models.AlertMessage,
	notifier ports.NotificationService,
) {
	log.Println("[Notification Worker] Iniciado y esperando alertas...")

	for msg := range alertQueue {
		// Aquí usamos la arquitectura hexagonal: no sabemos si es un Mock, SendGrid o SMTP
		err := notifier.SendAlert(msg.User, msg.Sismo)
		if err != nil {
			log.Printf("[Notification Worker] Error enviando alerta a %s: %v", msg.User.Email, err)
			// En un sistema real, aquí podrías reencolar el mensaje para reintentarlo luego
		}
	}
}
