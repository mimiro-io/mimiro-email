package email

import (
	"go.uber.org/zap"
)

type MailConsoleService struct {
	logger *zap.SugaredLogger
}

func NewMailConsoleService() *MailConsoleService {
	logger := NewLogger()

	return &MailConsoleService{logger: logger.Named("email")}
}

func (s *MailConsoleService) Send(email Mail) error {
	s.logger.Infof("mail sent to %s with subject %s to console", email.To, email.Subject)
	s.logger.Infof("body : %s", string(email.BodyText))

	return nil
}
