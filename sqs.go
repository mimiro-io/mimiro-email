package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	DelaySeconds int32
	Region       string
	Url          string
	ClientId     string
	Secret       string
	Auth         string
}

func NewMailSQSService(cfg Properties) *MailSQSService {
	logger := NewLogger()
	c, err := cfg.ValidAWSConfig()
	if err != nil {
		logger.Warnf(err.Error())
		return nil
	}
	if strings.ToUpper(c.Auth) == "CREDENTIALS" {
		return &MailSQSService{
			logger: logger.Named("email"),
			cfg:    c,
			client: sqs.NewFromConfig(aws.Config{
				Region:      c.Region,
				Credentials: credentials.NewStaticCredentialsProvider(c.ClientId, c.Secret, ""),
				EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service string, region string, Options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           c.Url,
						SigningRegion: c.Region,
					}, nil
				}),
			}),
		}
	} else {
		awsCfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic("configuration error, " + err.Error())
		}

		client := sqs.NewFromConfig(awsCfg)

		return &MailSQSService{
			logger: logger.Named("email"),
			cfg:    c,
			client: client,
		}

	}

}

func (p Properties) ValidAWSConfig() (Config, error) {
	required := []string{"QueueName", "DelaySeconds", "Region", "Url", "ClientId", "Secret"}
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
			DelaySeconds: cast.ToInt32(p[required[1]]),
			Region:       p[required[2]].(string),
			Url:          p[required[3]].(string),
			ClientId:     p[required[4]].(string),
			Secret:       p[required[5]].(string),
		}, nil
	} else {
		return Config{}, eris.New(errorMsg)
	}
}

func (s *MailSQSService) Send(email Mail) error {
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
	s.logger.Infof("mail sent to %s with subject %s to SQS", email.To, email.Subject)
	s.logger.Debugf("SQS response %v", sendMessageOutput)
	s.logger.Debug(string(email.BodyText))
	return nil
}
