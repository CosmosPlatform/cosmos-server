package mail

import (
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"fmt"

	"github.com/wneessen/go-mail"
)

type Service interface {
	SendMail(to string, subject string, body string) error
}

type mailService struct {
	client     *mail.Client
	senderMail string
	logger     log.Logger
}

func NewMailService(mailConfig config.MailConfig, logger log.Logger) (Service, error) {
	client, err := mail.NewClient(
		mailConfig.SMTPHost,
		mail.WithPort(mailConfig.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(mailConfig.SenderEmail),
		mail.WithPassword(mailConfig.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %v", err)
	}

	return &mailService{
		client:     client,
		senderMail: mailConfig.SenderEmail,
		logger:     logger,
	}, nil
}

func (ms *mailService) SendMail(to string, subject string, body string) error {
	m := mail.NewMsg()
	err := m.From(ms.senderMail)
	if err != nil {
		return fmt.Errorf("failed to set sender email: %v", err)
	}
	err = m.To(to)
	if err != nil {
		return fmt.Errorf("failed to set recipient email: %v", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)

	err = ms.client.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
