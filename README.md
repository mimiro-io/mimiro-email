# mimiro-email

[![CI](https://github.com/mimiro-io/mimiro-email/actions/workflows/ci.yaml/badge.svg)](https://github.com/mimiro-io/mimiro-email/actions/workflows/ci.yaml)

The purpose of this package is to have a service with an interface to send email with different backends. Supported now is Console and a Mimiro AWS flow using SQS and a lambda picking up a message from a sqs queue named "email" and sending the email with AWS SES.

Documentation regarding the SQS and the format used is found at https://github.com/mimiro-io/infra-aws-email

For test purpose you can configure it to use console, then it will write the body text out in the log

## Configuration

//TODO::

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
