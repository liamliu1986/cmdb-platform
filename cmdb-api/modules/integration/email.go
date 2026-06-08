package integration

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailConfig struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort string `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type EmailService struct {
	cfg *EmailConfig
}

func NewEmailService(cfg *EmailConfig) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendAlert(to []string, subject, body string) error {
	if s.cfg == nil || s.cfg.SMTPHost == "" {
		return fmt.Errorf("email service not configured")
	}
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.SMTPHost)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, strings.Join(to, ","), subject, body,
	)
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.From, to, []byte(msg))
}

// SendCIChangeNotification sends email notification for CI changes
func (s *EmailService) SendCIChangeNotification(to []string, ciName string, action string, operator string, detail string) error {
	subject := fmt.Sprintf("[CMDB] CI %s: %s by %s", action, ciName, operator)
	body := fmt.Sprintf(
		"Configuration Item Change Notification\n\nItem: %s\nAction: %s\nOperator: %s\nDetails: %s\n",
		ciName, action, operator, detail,
	)
	return s.SendAlert(to, subject, body)
}
