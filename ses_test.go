package ses

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var to, cc, bcc, from string

const rawBody = `To: %s
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

const textBody = `This is an example email body for the amazon SES go package.`
const htmlBody = `<p>This is an example email body for the amazon SES go package.</p>`

// mockHTTPBadRequest for mocking requests
type mockHTTPBadRequest struct{}

// Do is a mock http request
func (m *mockHTTPBadRequest) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.Body = io.NopCloser(bytes.NewBufferString(`{"error":"message failed"}`))

	// Default is valid
	return resp, nil
}

func init() { //nolint:gochecknoinits // ignore
	flag.StringVar(&to, "to", "success@simulator.amazonses.com", "email recipient")
	flag.StringVar(&cc, "cc", "success@simulator.amazonses.com", "cc email recipient")
	flag.StringVar(&bcc, "bcc", "success@simulator.amazonses.com", "bcc email recipient")
	flag.StringVar(&from, "from", "", "email sender")
}

// checkFlags checks for from address
func checkFlags(t *testing.T) {
	if len(from) == 0 {
		t.Fatal("must specify sender via -from flag.")
	}
}

// TestConfig_SendEmail will test the method SendEmail()
func TestConfig_SendEmail(t *testing.T) {
	var values url.Values
	var auth string
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s", HTTPClient: http.DefaultClient}
	subject := "amazon SES text test"
	_, err := cfg.SendEmail("from", []string{to}, []string{cc}, []string{bcc}, subject, textBody)
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

// TestConfig_SendEmailHTML will test the method SendEmailHTML()
func TestConfig_SendEmailHTML(t *testing.T) {
	var values url.Values
	var auth string

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s", HTTPClient: http.DefaultClient}
	subject := "amazon SES text test"
	_, err := cfg.SendEmailHTML("from", []string{to}, []string{cc}, []string{bcc}, subject, textBody, htmlBody)
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

// TestConfig_SendRawEmail will test the method SendRawEmail()
func TestConfig_SendRawEmail(t *testing.T) {
	var values url.Values
	var auth string

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		values, _ = url.ParseQuery(string(body))
	}))
	defer server.Close()

	attachment := base64.StdEncoding.EncodeToString([]byte(textBody))
	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s", HTTPClient: http.DefaultClient}
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

// TestConfig_SendEmailError will test the method SendEmail()
func TestConfig_SendEmailError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer server.Close()

	attachment := base64.StdEncoding.EncodeToString([]byte(textBody))
	cfg := Config{Endpoint: server.URL, Region: "region", AccessKeyID: "a", SecretAccessKey: "s", HTTPClient: &mockHTTPBadRequest{}}
	body := []byte(fmt.Sprintf(rawBody, "to", "from", textBody, len(attachment), attachment))
	_, err := cfg.SendRawEmail(body)
	if err == nil {
		t.Fatalf("expected to get an error")
	}
}

//
// Live Integration Tests
//

// TestConfig_SendEmailLiveTo will test the method SendEmail()
func TestConfig_SendEmailLiveTo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES text test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

// TestConfig_SendEmailLiveCC will test the method SendEmail()
func TestConfig_SendEmailLiveCC(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, []string{cc}, nil, "amazon SES cc test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

// TestConfig_SendEmailLiveBCC will test the method SendEmail()
func TestConfig_SendEmailLiveBCC(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, []string{bcc}, "amazon SES bcc test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

// TestConfig_SendEmailHTMLLive will test the method SendEmailHTML()
func TestConfig_SendEmailHTMLLive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped")
	}
	checkFlags(t)
	_, err := EnvConfig.SendEmail(from, []string{to}, nil, nil, "amazon SES html test", textBody)
	if err != nil {
		t.Fatal(err)
	}
}

// TestConfig_SendRawEmailLive will test the method SendRawEmail()
func TestConfig_SendRawEmailLive(t *testing.T) {
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
