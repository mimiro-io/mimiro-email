package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.uber.org/zap"
	"strings"
)

type MailSQSService struct {
	logger *zap.SugaredLogger
	cfg    EmailConfig
	client *sqs.Client
}

func NewMailSQSService(cfg EmailConfig) *MailSQSService {
	logger := zap.NewNop().Sugar()
	return &MailSQSService{
		logger: logger.Named("email"),
		cfg:    cfg,
		client: sqs.NewFromConfig(aws.Config{
			Region:      cfg.Aws.Region,
			Credentials: credentials.NewStaticCredentialsProvider(cfg.Aws.ClientId, cfg.Aws.Secret, ""),
			EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service string, region string, Options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.Aws.Url,
					SigningRegion: cfg.Aws.Region,
				}, nil
			}),
		}),
	}
}

func (s *MailSQSService) SendTemplate(email Mail) error {
	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &s.cfg.Aws.QueueName,
	}
	queueUrl, err := s.client.GetQueueUrl(context.Background(), gQInput)
	if err != nil {
		return err
	}

	input := &sqs.SendMessageInput{
		DelaySeconds: s.cfg.Aws.DelaySeconds,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"SENDER_EMAIL": {
				DataType:    aws.String("String"),
				StringValue: aws.String(email.Sender),
			},
			"TO_ADDRESSES": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("[\"%s\"]", strings.Join(email.To, "\", \""))),
			},
			"CC_ADDRESSES": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("[\"%s\"]", strings.Join(email.Cc, "\", \""))),
			},
			"BCC_ADDRESSES": {
				DataType:    aws.String("String"),
				StringValue: aws.String(fmt.Sprintf("[\"%s\"]", strings.Join(email.Bcc, "\", \""))),
			},
			"SUBJECT": {
				DataType:    aws.String("String"),
				StringValue: aws.String(email.Subject),
			},
			"BODY_HTML": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(email.BodyHtml)),
			},
			"BODY_TEXT": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(email.BodyText)),
			},
		},
		MessageBody: nil,
		QueueUrl:    queueUrl.QueueUrl,
	}

	sendMessageOutput, err := s.client.SendMessage(context.Background(), input)
	if err != nil {
		return err
	}
	s.logger.Debugf("mail sent to sqs with response %v", sendMessageOutput)
	return nil
}
