# mimiro-email

The purpose of this package is to have a service with an interface to send email with different backends. Supported now is Console and a Mimiro AWS flow using SQS and a lambda picking up a message from a sqs queue named "email" and sending the email with AWS SES.

Documentation regarding the SQS and the format used is found at https://github.com/mimiro-io/infra-aws-email

For local test purpose you can configure it to use console, then it will write the body text out in the log. 
In development configured with AWS the emails end up in Slack channel #infra-aws-email-dev, in prod it's sent to recipients.

## Configuration

### Common 
Service : *AWS or CONSOLE, stating backend technology*
- AWS - AWS implements a Mimiro flow with AWS SQS, lambda and AWS SES  
- CONSOLE - write body text in console

Sender : *Senders email address ex "OpenFarm Dev <noreply@openfarm-dev.io>"*
### Console config parameters
No params

### AWS config parameters
QueueName : *Name of sqs queue for Mimiro dev and prod it's named email*

DelaySeconds : *Delay queues let you postpone the delivery of new messages to consumers for a number of seconds, for example, when your consumer application needs additional time to process messages. If you create a delay queue, any messages that you send to the queue remain invisible to consumers for the duration of the delay period. The default (minimum) delay for a queue is 0 seconds. The maximum is 15 minutes*

## Use

Run
export GOPRIVATE=github.com/mimiro-io/mimiro-email

then 
go get -u github.com/mimiro-io/mimiro-email   


```
var srv mail.Email

awsCfg := mail.Configuration{
    Service: "AWS",
    Sender:  "noreply <noreply@noreply.io>",
    Properties: map[string]interface{}{
        "QueueName":    "email",
        "DelaySeconds": 10,
    },
}

consoleCfg := mail.Configuration{
    Service: "CONSOLE",
    Sender:  "noreply <noreply@noreply.io>",
    },
}

awsSrv, err = mail.NewEmail(awsCfg)
if err != nil {
    return nil
}

consoleSrv, err = mail.NewEmail(consoleCfg)
if err != nil {
    return nil
}

var mail = mail.Mail{
    Sender:   service.sender,
    To:       []string{"test@test.com"},
    Cc:       nil,
    Bcc:      nil,
    Subject:  "Subject",
    BodyHtml: bodyHtml.Bytes(),
    BodyText: bodyText.Bytes(),
}

err = service.awsSrv.Send(mail)
if err != nil {
    return err
}

err = service.consoleSrv.Send(mail)
if err != nil {
    return err
}

```

### Add to policies.json when configured with AWS
```
        {
            "Effect": "Allow",
            "Action": [
                "sqs:GetQueueUrl",
                "sqs:SendMessage"
            ],
            "Resource": [
                "arn:aws:sqs:eu-west-1:${data.aws_caller_identity.current.account_id}:email"
            ]
        }
```

### Github workflow
Application using package "github.com/mimiro-io/mimiro-email"  needs to pass build-args in GHA


          build-args: |
            GITHUB_PAT=${{ secrets.MIMIRO_BUILD_GITHUB_TOKEN }}

Se example below from application

```
name: CI
on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]
  release:
    types:
      - published
jobs:
  Test:
    runs-on: ops
    steps:
      - uses: actions/checkout@v2
      - name: Set up Docker Buildx (enable caching)
        uses: docker/setup-buildx-action@v1

      - name: Build Tests Runner
        uses: docker/build-push-action@v2
        id: test-image
        with:
          push: false
          load: true
          cache-from: type=gha
          cache-to: type=gha,mode=max
          target: builder
          build-args: |
            GITHUB_PAT=${{ secrets.MIMIRO_BUILD_GITHUB_TOKEN }}
          tags: |
            ${{ github.event.repository.name }}-tester:${{ github.sha }}
      - name: Run Tests
        id: run-tests
        run: |
          docker run -v /var/run/docker.sock:/var/run/docker.sock \
          -v $(pwd)/migration/sql:$(pwd)/migration/sql \
          -v $(pwd)/test-resources:/app/test-resources \
          -e MIGRATION_DIR=$(pwd)/migration/sql \
          ${{ github.event.repository.name }}-tester:${{ github.sha }} go test -v ./...
```