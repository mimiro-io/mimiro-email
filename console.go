package email

import (
	"go.uber.org/zap"
)

type MailConsoleService struct {
	logger *zap.SugaredLogger
}

func NewMailConsoleService() *MailConsoleService {
	logger := zap.NewNop().Sugar()
	return &MailConsoleService{logger: logger.Named("email")}
}

func (s *MailConsoleService) SendTemplate(email Mail) error {
	s.logger.Debugf("mail sent to console with response %v", email)

	return nil
}
