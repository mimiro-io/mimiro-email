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
	"testing"

	"github.com/franela/goblin"
)

func TestNewEmail(t *testing.T) {
	g := goblin.Goblin(t)
	var srv Email
	g.Describe("test email interface", func() {
		g.Before(func() {
			consoleCfg := Configuration{
				Service: "Console",
				Sender:  "OpenFarm Dev <noreply@openfarm-dev.io>",
			}
			srv, _ = NewEmail(consoleCfg)
		})

		g.It("test console ", func() {
			m := Mail{
				To:       []string{"test1@test.com", "test2@test.com"},
				Cc:       nil,
				Bcc:      nil,
				Subject:  "Subject",
				BodyHtml: []byte("<h1>Heading</h1>"),
				BodyText: []byte("Heading"),
			}
			err := srv.Send(m)

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
