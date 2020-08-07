package ses

import (
	"encoding/base64"
	"flag"
	"fmt"
	"testing"
)

var to, from string

func init() {
	flag.StringVar(&to, "to", "success@simulator.amazonses.com", "email recipient")
	flag.StringVar(&from, "from", "", "email sender")
}

func checkFlags(t *testing.T) {
	if len(from) == 0 {
		t.Fatal("must specify sender via -from flag.")
	}
}

func TestSendEmail(t *testing.T) {
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES text test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendEmailHTML(t *testing.T) {
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES text test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendRawEmail(t *testing.T) {
	checkFlags(t)
	attachment := base64.StdEncoding.EncodeToString([]byte(textBody))
	raw := `To: %s
From: %s
Subject: amazon SES raw test
Content-Type: multipart/mixed; boundary="_003_97DCB304C5294779BEBCFC8357FCC4D2"
MIME-Version: 1.0

--_003_97DCB304C5294779BEBCFC8357FCC4D2
Content-Type: text/plain;

%s

--_003_97DCB304C5294779BEBCFC8357FCC4D2
Content-Type: text/plain; name="test.txt"
Content-Description: test.txt
Content-Disposition: attachment; filename="test.txt"; size=%d;
Content-Transfer-Encoding: base64

%s

--_003_97DCB304C5294779BEBCFC8357FCC4D2
`

	_, err := EnvConfig.SendRawEmail([]byte(fmt.Sprintf(raw, to, from, textBody, len(attachment), attachment)))
	if err != nil {
		t.Fatal(err)
	}
}

var textBody = `This is an example email body for the amazon SES go package.`
