package email

import (
	"fmt"
	"net/smtp"
	"seismic-monitor/backend/internal/models"
	"strings" // Necesario para detectar el prefijo del ID
)

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

// SendAlert implementa la interfaz NotificationService
func (s *SMTPSender) SendAlert(user models.User, sismo models.Feature) error {
	from := "alertas@seismicmonitor.com"
	to := []string{user.Email}

	var subject string
	var body string

	// Lógica para diferenciar el tipo de mensaje según el ID del sismo
	if strings.HasPrefix(sismo.ID, "COMUNIDAD-") {
		// --- FORMATO AVISO COMUNITARIO ---
		subject = "¡AVISO COMUNITARIO! Posible temblor detectado"
		body = fmt.Sprintf(
			"Hola %s,\n\nVarios usuarios en tu zona han reportado sentir un temblor justo ahora. "+
				"Esta es una alerta temprana basada en reportes ciudadanos. ¡Mantente alerta y ten precaución!",
			user.Email,
		)
	} else {
		// --- FORMATO SISMO OFICIAL (USGS) ---
		subject = fmt.Sprintf("¡PELIGRO! Sismo oficial %s detectado cerca de ti", sismo.ID)
		body = fmt.Sprintf(
			"Hola %s,\n\nSe ha detectado un sismo oficial confirmado de magnitud %.1f en %s.\n\n¡Mantente a salvo!",
			user.Email, sismo.Info.Mag, sismo.Info.Place,
		)
	}

	// Dentro de SendAlert, cuando construyes el body:
	body += "\n\n--- 💡 ANÁLISIS DE SEGURIDAD (IA) ---\n"
	body += sismo.AIAdvice
	body += "\n--------------------------------------\n"

	// Construimos el mensaje siguiendo el estándar de red SMTP
	// Importante: El formato debe ser "Subject: ... \n\n Cuerpo"
	messageStr := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s",
		from, user.Email, subject, body)

	message := []byte(messageStr)

	// Autenticación con Mailtrap o servidor SMTP
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Envío del correo
	err := smtp.SendMail(addr, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("error enviando email via SMTP: %w", err)
	}

	return nil
}
