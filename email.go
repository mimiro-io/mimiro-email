package main

import (
	"github.com/rotisserie/eris"
	"strings"
)

type Email interface {
	Send(email Mail) error
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

type Configuration struct {
	Service    string // CONSOLE MIMIRO_EMAIL_SQS(AWS)
	Sender     string
	Properties Properties
}

var ErrUndefinedService = eris.New("undefined service")

func NewEmail(cfg Configuration) (Email, error) {
	switch strings.ToUpper(cfg.Service) {
	case "AWS":
		return NewMailSQSService(cfg.Properties), nil
	case "CONSOLE":
		return NewMailConsoleService(), nil
	default:
		return nil, ErrUndefinedService
	}
}

func main() {
	consoleCfg := Configuration{
		Service: "Console",
		Sender:  "OpenFarm Dev <noreply@openfarm-dev.io>",
	}
	srv, _ := NewEmail(consoleCfg)
	m := Mail{
		To:       []string{"test1@test.com", "test2@test.com"},
		Cc:       nil,
		Bcc:      nil,
		Subject:  "Subject",
		BodyHtml: []byte("<h1>Heading</h1>"),
		BodyText: []byte("Heading"),
	}
	err := srv.Send(m)
	if err != nil {
		return
	}

}
