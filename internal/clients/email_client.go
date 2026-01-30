package clients

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
)

type EmailClient struct {
	addr string
	auth smtp.Auth
	from string
}

func (c *EmailClient) Send(subject, textBody, htmlBody string, to []string) error {
	boundary := "LingoStructBoundary12345"

	var builder strings.Builder

	fmt.Fprintf(&builder, "From: %s\r\n", c.from)
	fmt.Fprintf(&builder, "To: %s\r\n", strings.Join(to, ","))
	fmt.Fprintf(&builder, "Subject: %s\r\n", subject)
	builder.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(&builder, "Content-Type: multipart/alternative; boundary=%s\r\n", boundary)
	builder.WriteString("\r\n")

	fmt.Fprintf(&builder, "--%s\r\n", boundary)
	builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(textBody)
	builder.WriteString("\r\n")

	fmt.Fprintf(&builder, "--%s\r\n", boundary)
	builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(htmlBody)
	builder.WriteString("\r\n")

	fmt.Fprintf(&builder, "--%s--\r\n", boundary)

	return smtp.SendMail(c.addr, c.auth, c.from, to, []byte(builder.String()))
}

func NewEmailClient(cfg *config.SmtpConfig) (*EmailClient, error) {
	from := cfg.From
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("smtp dial failed: %w", err)
	}
	defer conn.Quit()

	err = conn.StartTLS(&tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.Host,
	})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("smtp STARTTLS failed: %w", err)
	}

	if err := conn.Auth(auth); err != nil {
		return nil, fmt.Errorf("smtp auth failed: %w", err)
	}

	return &EmailClient{addr, auth, from}, nil
}
