package email

import (
	"errors"
	"fmt"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/wneessen/go-mail"
)

type EmailService struct {
	client *mail.Client
	from   string
}

func NewEmailService(config *config.SmtpConfig) (*EmailService, error) {
	client, err := mail.NewClient(
		config.Host,
		mail.WithPort(config.Port),
		mail.WithTLSPolicy(map[bool]mail.TLSPolicy{
			true:  mail.TLSMandatory,
			false: mail.NoTLS,
		}[*config.Tls]),
		mail.WithTimeout(10*time.Second),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(config.User),
		mail.WithPassword(config.Password),
	)

	if err != nil {
		return nil, errors.New("Failed to create SMTP client")
	}

	return &EmailService{
		client: client,
		from:   config.From,
	}, nil
}

func (s *EmailService) Close() error {
	return s.client.Close()
}

func (s *EmailService) Send(message EmailMessage) error {
	msg := mail.NewMsg()

	if err := msg.From(s.from); err != nil {
		return domain.InternalError("Invalid sender email address")
	}

	if err := msg.To(message.To.Value()); err != nil {
		return domain.InternalError("Invalid recipient email address")
	}

	msg.Subject(message.Subject.Value())
	msg.SetBodyString(mail.TypeTextPlain, message.PlainText.Value())

	if message.HtmlText != nil {
		msg.AddAlternativeString(
			mail.TypeTextHTML,
			message.HtmlText.Value(),
		)
	}

	if err := s.client.DialAndSend(msg); err != nil {
		return domain.InternalError(fmt.Sprintf("Failed to send email: %v", err))
	}

	return nil
}
