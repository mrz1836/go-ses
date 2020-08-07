package ses

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awssigner "github.com/aws/aws-sdk-go/aws/signer/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Config specifies configuration options and credentials for accessing Amazon SES.
type Config struct {
	// Endpoint is the AWS endpoint to use for requests.
	Endpoint string

	// The region
	Region string

	// AccessKeyID is your Amazon AWS access key ID.
	AccessKeyID string

	// SecretAccessKey is your Amazon AWS secret key.
	SecretAccessKey string
}

// EnvConfig takes the access key ID and secret access key values from the environment variables
// $AWS_ACCESS_KEY_ID and $AWS_SECRET_KEY, respectively.
var EnvConfig = Config{
	Endpoint:        os.Getenv("AWS_SES_ENDPOINT"),
	Region:          os.Getenv("AWS_REGION"),
	AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
	SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
}

func (c *Config) fillRecipients(from string, to, cc, bcc []string, data url.Values) {
	data.Add("Source", from)

	if len(to) > 0 {
		for i := 0; i < len(to); i++ {
			data.Add(fmt.Sprintf("Destination.ToAddresses.member.%d", i+1), to[i])
		}
	}
	if len(cc) > 0 {
		for i := 0; i < len(cc); i++ {
			data.Add(fmt.Sprintf("Destination.CcAddresses.member.%d", i+1), cc[i])
		}
	}
	if len(bcc) > 0 {
		for i := 0; i < len(bcc); i++ {
			data.Add(fmt.Sprintf("Destination.BccAddresses.member.%d", i+1), bcc[i])
		}
	}
}

// SendEmail sends a plain text email. Note that from must be a verified
// address in the AWS control panel.
func (c *Config) SendEmail(from string, to, cc, bcc []string, subject, body string) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendEmail")
	c.fillRecipients(from, to, cc, bcc, data)
	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Text.Data", body)
	data.Add("AWSAccessKeyId", c.AccessKeyID)
	return c.sesPost(data)
}

// SendEmailHTML sends a HTML email. Note that from must be a verified address
// in the AWS control panel.
func (c *Config) SendEmailHTML(from string, to, cc, bcc []string, subject, bodyText, bodyHTML string) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendEmail")
	c.fillRecipients(from, to, cc, bcc, data)
	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Text.Data", bodyText)
	data.Add("Message.Body.Html.Data", bodyHTML)
	data.Add("AWSAccessKeyId", c.AccessKeyID)
	return c.sesPost(data)
}

// SendRawEmail sends a raw email. Note that from must be a verified address
// in the AWS control panel.
func (c *Config) SendRawEmail(raw []byte) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendRawEmail")
	data.Add("RawMessage.Data", base64.StdEncoding.EncodeToString(raw))
	data.Add("AWSAccessKeyId", c.AccessKeyID)
	return c.sesPost(data)
}

func (c *Config) sigv4(req *http.Request, body string, tm time.Time) {
	creds := credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey}})
	signer := awssigner.NewSigner(creds)
	signer.Sign(req, strings.NewReader(body), "email", c.Region, tm)
}

func (c *Config) sesPost(data url.Values) (string, error) {
	req, err := http.NewRequest("POST", c.Endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	now := time.Now().UTC()
	date := now.Format("Mon, 02 Jan 2006 15:04:05 -0700")
	req.Header.Set("Date", date)
	c.sigv4(req, data.Encode(), now)

	var r *http.Response
	if r, err = http.DefaultClient.Do(req); err != nil {
		return "", err
	}

	resultBody, _ := ioutil.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error code %d. response: %s", r.StatusCode, resultBody)
	}

	return string(resultBody), nil
}
