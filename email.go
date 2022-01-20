package email

import (
	"github.com/rotisserie/eris"
	"strings"
)

type Email interface {
	SendTemplate(email Mail) error
}

type Mail struct {
	Sender   string
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	BodyHtml []byte
	BodyText []byte
}

type Properties map[string]interface{}

type EmailConfig struct {
	Service    string // CONSOLE MIMIRO_EMAIL_SQS(AWS)
	Properties Properties
}

var ErrUndefinedService = eris.New("undefined service")

func NewEmail(cfg EmailConfig) (Email, error) {
	switch strings.ToUpper(cfg.Service) {
	case "AWS":
		return NewMailSQSService(cfg.Properties), nil
	case "CONSOLE":
		return NewMailConsoleService(), nil
	default:
		return nil, ErrUndefinedService
	}
}
