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
			consoleCfg := EmailConfig{
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
			consoleCfg := EmailConfig{
				Service: "TEST",
			}

			_, err := NewEmail(consoleCfg)
			g.Assert(err).IsNotNil()
		})

	})

}
