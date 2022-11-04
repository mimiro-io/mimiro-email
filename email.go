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
	"github.com/rotisserie/eris"
	"strings"
)

type Email interface {
	Send(email Mail) error
}

type Mail struct {
	Sender   string
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	BodyHtml []byte
	BodyText []byte
}

type Properties map[string]interface{}

type Configuration struct {
	Service    string // CONSOLE MIMIRO_EMAIL_SQS(AWS)
	Sender     string
	Properties Properties
}

var ErrUndefinedService = eris.New("undefined service")

func NewEmail(cfg Configuration) (Email, error) {
	switch strings.ToUpper(cfg.Service) {
	case "AWS":
		return NewMailSQSService(cfg.Properties), nil
	case "CONSOLE":
		return NewMailConsoleService(), nil
	default:
		return nil, ErrUndefinedService
	}
}
