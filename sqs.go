// Copyright 2021 MIMIRO AS
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rotisserie/eris"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type MailSQSService struct {
	logger *zap.SugaredLogger
	cfg    Config
	client *sqs.Client
}

type Config struct {
	QueueName    string
	DelaySeconds int32
}

func NewMailSQSService(cfg Properties) *MailSQSService {
	logger := NewLogger()
	c, err := cfg.ValidAWSConfig()
	if err != nil {
		logger.Warnf(err.Error())
		return nil
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())
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

func (p Properties) ValidAWSConfig() (Config, error) {
	required := []string{"QueueName", "DelaySeconds"}
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
		MessageBody: aws.String("Sender Service: mimiro-email"),
		QueueUrl:    queueUrl.QueueUrl,
	}

	sendMessageOutput, err := s.client.SendMessage(context.Background(), input)
	if err != nil {
		return err
	}
	s.logger.Debugf("log sqs response %v", sendMessageOutput)
	return nil
}
