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

type EmailConfig struct {
	Service string // CONSOLE MIMIRO_EMAIL_SQS(AWS)
	Aws     *Aws   //TODO:: make smarter and generic map[string]interface{}
	Console *Console
}

type Aws struct {
	QueueName    string
	SenderEmail  string
	DelaySeconds int32
	Region       string
	Url          string
	ClientId     string
	Secret       string
}

type Console struct {
}

var ErrUndefinedService = eris.New("undefined service")

func NewEmail(cfg EmailConfig) (Email, error) {
	switch strings.ToUpper(cfg.Service) {
	case "AWS":
		return NewMailSQSService(cfg), nil
	case "CONSOLE":
		return NewMailConsoleService(), nil
	default:
		return nil, ErrUndefinedService
	}
}

func (m Mail) Send() {

}

func main() {
	m := Mail{
		Sender:   "",
		To:       nil,
		Cc:       nil,
		Bcc:      nil,
		Subject:  "",
		BodyHtml: nil,
		BodyText: nil,
	}

	m.Send()
}
