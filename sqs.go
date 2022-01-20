package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rotisserie/eris"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
)

type MailSQSService struct {
	logger *zap.SugaredLogger
	cfg    Config
	client *sqs.Client
}

type Config struct {
	QueueName    string
	SenderEmail  string
	DelaySeconds int32
	Region       string
	Url          string
	ClientId     string
	Secret       string
}

func NewMailSQSService(cfg Properties) *MailSQSService {
	logger := zap.NewNop().Sugar()
	config, err := cfg.ValidAWSConfig()
	if err != nil {
		logger.Warnf(err.Error())
		return nil
	}
	return &MailSQSService{
		logger: logger.Named("email"),
		cfg:    config,
		client: sqs.NewFromConfig(aws.Config{
			Region:      config.Region,
			Credentials: credentials.NewStaticCredentialsProvider(config.ClientId, config.Secret, ""),
			EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service string, region string, Options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           config.Url,
					SigningRegion: config.Region,
				}, nil
			}),
		}),
	}
}

func (p Properties) ValidAWSConfig() (Config, error) {
	required := []string{"QueueName", "Sender", "DelaySeconds", "Region", "Url", "ClientId", "Secret"}
	var valid = true
	errorMsg := "Missing required properties :"
	for i, _ := range required {
		if p[required[i]] == nil {
			valid = false
			errorMsg = errorMsg + fmt.Sprintf(" %s", required[i])
		}
	}

	if valid {
		return Config{
			QueueName:    p[required[0]].(string),
			SenderEmail:  p[required[1]].(string),
			DelaySeconds: cast.ToInt32(p[required[2]]),
			Region:       p[required[3]].(string),
			Url:          p[required[4]].(string),
			ClientId:     p[required[5]].(string),
			Secret:       p[required[6]].(string),
		}, nil
	} else {
		return Config{}, eris.New(errorMsg)
	}
}

func (s *MailSQSService) SendTemplate(email Mail) error {
	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &s.cfg.QueueName,
	}
	queueUrl, err := s.client.GetQueueUrl(context.Background(), gQInput)
	if err != nil {
		return err
	}

	input := &sqs.SendMessageInput{
		DelaySeconds: s.cfg.DelaySeconds,
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
