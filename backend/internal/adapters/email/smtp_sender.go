package email

import (
	"fmt"
	"net/smtp"
	"seismic-monitor/backend/internal/models"
)

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

// SendAlert implementa la interfaz NotificationService
func (s *SMTPSender) SendAlert(user models.User, sismo models.Feature) error {
	// Configuramos la autenticación con Mailtrap
	/*
		auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

		// Montamos el correo (Cabeceras + Cuerpo)
		to := []string{user.Email}
		msg := []byte(fmt.Sprintf("To: %s\r\n"+
			"Subject: ¡PELIGRO! Sismo detectado en %s\r\n"+
			"\r\n"+
			"Hola %s,\r\n\r\n"+
			"Un sismo de magnitud %.1f ha ocurrido cerca de ti.\r\n"+
			"Por favor, toma precauciones.\r\n",
			user.Email, sismo.Info.Place, user.Username, sismo.Info.Mag))

		// Enviamos el correo físicamente
		address := fmt.Sprintf("%s:%s", s.Host, s.Port)
		err := smtp.SendMail(address, auth, s.Username, to, msg)
		if err != nil {
			return fmt.Errorf("fallo al enviar correo real: %w", err)
		}

		return nil
	*/
	from := "alertas@seismicmonitor.com"
	to := []string{user.Email}
	subject := fmt.Sprintf("Subject: ¡PELIGRO! Sismo %s detectado cerca de ti\n", sismo.ID)
	body := fmt.Sprintf("Hola %s,\n\nSe ha detectado un sismo de magnitud %.1f en %s.\n\n¡Mantente a salvo!",
		user.Email, sismo.Info.Mag, sismo.Info.Place)

	message := []byte(subject + "\n" + body)

	// 2. Autenticación
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// 3. Enviar
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	err := smtp.SendMail(addr, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("error enviando email via SMTP: %w", err)
	}

	return nil
}
