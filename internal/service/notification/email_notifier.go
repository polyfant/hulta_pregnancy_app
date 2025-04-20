package notification

import (
	"context"
	
	"crypto/tls"
	"fmt"
	
	"net/smtp"
	"os"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPass   string
	TLSInsecure  bool
}

type EmailNotifier interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type EmailNotifierImpl struct {
	config EmailConfig
}

func NewEmailNotifier() *EmailNotifierImpl {
	return &EmailNotifierImpl{
		config: EmailConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SenderEmail:  os.Getenv("SMTP_SENDER_EMAIL"),
			SenderPass:   os.Getenv("SMTP_SENDER_PASS"),
			TLSInsecure:  os.Getenv("SMTP_TLS_INSECURE") == "true",
		},
	}
}

func (e *EmailNotifierImpl) SendEmail(ctx context.Context, to, subject, body string) error {
	// Construct email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", 
		to, subject, mime, body))

	// TLS Configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: e.config.TLSInsecure,
		ServerName:         e.config.SMTPHost,
	}

	// Connect to SMTP Server
	conn, err := tls.Dial("tcp", e.config.SMTPHost+":"+e.config.SMTPPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial error: %w", err)
	}

	client, err := smtp.NewClient(conn, e.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Close()

	// Authentication
	if err = client.Auth(smtp.PlainAuth("", e.config.SenderEmail, e.config.SenderPass, e.config.SMTPHost)); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(e.config.SenderEmail); err != nil {
		return fmt.Errorf("sender error: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("recipient error: %w", err)
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data error: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}

	return client.Quit()
}
