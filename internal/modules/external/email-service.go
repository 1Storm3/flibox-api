package external

import (
	"github.com/1Storm3/flibox-api/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailServiceInterface interface {
	SendEmail(email, body, subject string) error
}

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

func (s *EmailService) SendEmail(email, body, subject string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.cfg.App.FromEmail)

	m.SetHeader("To", email)

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.cfg.SMTP.Host,
		s.cfg.SMTP.Port,
		s.cfg.SMTP.Username,
		s.cfg.SMTP.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
