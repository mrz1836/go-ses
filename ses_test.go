package ses

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
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
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES text test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendEmailHTML(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES text test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendRawEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	attachment := base64.StdEncoding.EncodeToString([]byte(textBody))
	_, err := EnvConfig.SendRawEmail([]byte(fmt.Sprintf(rawBody, to, from, textBody, len(attachment), attachment)))
	if err != nil {
		t.Fatal(err)
	}
}

var rawBody = `To: %s
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

var textBody = `This is an example email body for the amazon SES go package.`
var htmlBody = `<p>This is an example email body for the amazon SES go package.</p>`

func TestSendEmailLocal(t *testing.T) {
	var values url.Values
	var auth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := ioutil.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s"}
	subject := "amazon SES text test"
	_, err := cfg.SendEmail("from", []string{"to"}, nil, nil, subject, textBody)
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=a/%s/region/email/aws4_request, SignedHeaders=content-type;date;host;x-amz-date, Signature=", time.Now().UTC().Format("20060102"))
	if !strings.HasPrefix(auth, expected) {
		t.Errorf("Wrong signature: expected: %s got %s", expected, auth)
	}
	if values.Get("Action") != "SendEmail" {
		t.Errorf("Missing Action")
	}
	if values.Get("Message.Subject.Data") != subject {
		t.Errorf("Wrong subject")
	}
	if values.Get("Message.Body.Text.Data") != textBody {
		t.Errorf("Wrong body")
	}
	if values.Get("AWSAccessKeyId") != cfg.AccessKeyID {
		t.Errorf("Wrong key")
	}
}

func TestSendEmailHTMLLocal(t *testing.T) {
	var values url.Values
	var auth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := ioutil.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s"}
	subject := "amazon SES text test"
	_, err := cfg.SendEmailHTML("from", []string{"to"}, nil, nil, subject, textBody, htmlBody)
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=a/%s/region/email/aws4_request, SignedHeaders=content-type;date;host;x-amz-date, Signature=", time.Now().UTC().Format("20060102"))
	if !strings.HasPrefix(auth, expected) {
		t.Errorf("Wrong signature: expected: %s got %s", expected, auth)
	}
	if values.Get("Action") != "SendEmail" {
		t.Errorf("Missing Action")
	}
	if values.Get("Message.Subject.Data") != subject {
		t.Errorf("Wrong subject")
	}
	if values.Get("Message.Body.Text.Data") != textBody {
		t.Errorf("Wrong body")
	}
	if values.Get("Message.Body.Html.Data") != htmlBody {
		t.Errorf("Wrong body")
	}
	if values.Get("AWSAccessKeyId") != cfg.AccessKeyID {
		t.Errorf("Wrong key")
	}
}

func TestSendEmailRawLocal(t *testing.T) {
	var values url.Values
	var auth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := ioutil.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	attachment := base64.StdEncoding.EncodeToString([]byte(textBody))
	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s"}
	body := []byte(fmt.Sprintf(rawBody, "to", "from", textBody, len(attachment), attachment))
	_, err := cfg.SendRawEmail(body)
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=a/%s/region/email/aws4_request, SignedHeaders=content-type;date;host;x-amz-date, Signature=", time.Now().UTC().Format("20060102"))
	if !strings.HasPrefix(auth, expected) {
		t.Errorf("Wrong signature: expected: %s got %s", expected, auth)
	}
	if values.Get("Action") != "SendRawEmail" {
		t.Errorf("Missing Action")
	}
	if values.Get("AWSAccessKeyId") != cfg.AccessKeyID {
		t.Errorf("Wrong key")
	}
	if values.Get("RawMessage.Data") != base64.StdEncoding.EncodeToString(body) {
		t.Errorf("Wrong data")
	}
}
