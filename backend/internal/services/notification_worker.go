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
	ctx context.Context,
	alertQueue <-chan models.AlertMessage,
	notifier ports.NotificationService,
	aiProvider ports.AIProvider,
) {
	log.Println("[Notification Worker] Iniciado y esperando alertas...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Notification Worker] Cerrando de forma limpia (Graceful Shutdown)...")
			return
		case msg, ok := <-alertQueue:
			if !ok {
				log.Println("[Notification Worker] Cola de alertas cerrada.")
				return
			}

			go func(m models.AlertMessage) {
				aiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				advice, err := aiProvider.GenerateSafetyAdvice(aiCtx, msg.Sismo)
				cancel()

				if err != nil {
					log.Printf("[Notification Worker] Error en IA para sismo %s: %v. Usando copia de seguridad.", msg.Sismo.ID, err)
					advice = "Atención: Siga las indicaciones de los equipos de emergencia de su localidad."
				}

				msg.Sismo.AIAdvice = advice

				err = notifier.SendAlert(msg.User, msg.Sismo)
				if err != nil {
					log.Printf("[Notification Worker] Error crítico de envío a %s: %v", msg.User.Email, err)
				}
			}(msg)
		}
	}
}
