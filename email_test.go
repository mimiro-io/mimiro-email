package email

import (
	"github.com/franela/goblin"
	"testing"
)

func TestNewEmail(t *testing.T) {
	g := goblin.Goblin(t)
	var srv Email
	g.Describe("test email interface", func() {
		g.Before(func() {
			consoleCfg := Configuration{
				Service: "Console",
			}
			srv, _ = NewEmail(consoleCfg)
		})

		g.It("test console ", func() {
			m := Mail{
				Sender:   "OpenFarm Dev <noreply@openfarm-dev.io>",
				To:       []string{"test1@test.com", "test2@test.com"},
				Cc:       nil,
				Bcc:      nil,
				Subject:  "Subject",
				BodyHtml: []byte("<h1>Heading</h1>"),
				BodyText: []byte("Heading"),
			}
			err := srv.SendTemplate(m)

			g.Assert(err).IsNil()

		})

		g.It("test undefined service ", func() {
			consoleCfg := Configuration{
				Service: "TEST",
			}

			_, err := NewEmail(consoleCfg)
			g.Assert(err).IsNotNil()
		})

		g.It("test configuration  is ok ", func() {
			consoleCfg := Configuration{
				Service: "AWS",
				Properties: map[string]interface{}{
					"QueueName":    "email",                                  // number => string
					"Sender":       "OpenFarm Dev <noreply@openfarm-dev.io>", // string => number
					"DelaySeconds": 10,
					"Region":       "east-1",
					"Url":          "https://domain.com/whatever",
					"ClientId":     "ClientId",
					"Secret":       "Secret",
				},
			}

			config, err := consoleCfg.Properties.ValidAWSConfig()

			g.Assert(err).IsNil()
			g.Assert(config).IsNotNil()
		})

		g.It("test configuration  is not ok ", func() {
			consoleCfg := Configuration{
				Service: "AWS",
				Properties: map[string]interface{}{
					"QueueName":    "email",                                  // number => string
					"Senderz":      "OpenFarm Dev <noreply@openfarm-dev.io>", // string => number
					"DelaySeconds": 10,
					"Url":          "https://domain.com/whatever",
					"ClientId":     "ClientId",
					"Secret":       "Secret",
				},
			}

			config, err := consoleCfg.Properties.ValidAWSConfig()

			g.Assert(err).IsNotNil()
			g.Assert(config).IsZero()
		})

	})

}
