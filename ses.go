package ses

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	awssigner "github.com/aws/aws-sdk-go/aws/signer/v4"
)

// httpInterface is used for the http client (mocking)
type httpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

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

	// HTTPClient is a http client to use
	HTTPClient httpInterface
}

// EnvConfig takes the access key ID and secret access key values from the environment variables
// $AWS_ACCESS_KEY_ID and $AWS_SECRET_KEY, respectively.
var EnvConfig = Config{
	AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"), // Set from ENV using standard name
	Endpoint:        os.Getenv("AWS_SES_ENDPOINT"),  // Set from ENV using standard name
	HTTPClient:      http.DefaultClient,             // Use a default client unless overridden
	Region:          os.Getenv("AWS_REGION"),        // Set from ENV using standard name
	SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),    // Set from ENV using standard name
}

// fillRecipients will fill all recipients into the data.values
func (c *Config) fillRecipients(from string, to, cc, bcc []string, data url.Values) {
	data.Add("Source", from)

	// todo: remove IF cases, since if empty, for loop will skip anyway?
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

// SendEmailHTML sends an HTML email. Note that from must be a verified address
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

// sigv4 signs using the new V4 signature method
func (c *Config) sigv4(req *http.Request, body string, timestamp time.Time) error {
	awsCredentials := credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey}})
	_, err := awssigner.NewSigner(awsCredentials).Sign(req, strings.NewReader(body), "email", c.Region, timestamp)
	return err
}

// sesPost fires the actual HTTP post request with the email data
func (c *Config) sesPost(data url.Values) (string, error) {

	// Set the request with context
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.Endpoint, nil)
	if err != nil {
		return "", err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the date/time
	now := time.Now().UTC()
	req.Header.Set("Date", now.Format("Mon, 02 Jan 2006 15:04:05 -0700"))

	// Sign with AWS SigV4
	if err = c.sigv4(req, data.Encode(), now); err != nil {
		return "", err
	}

	// Fire the request
	var resp *http.Response
	if resp, err = c.HTTPClient.Do(req); err != nil {
		return "", err
	}

	// Read the body
	var resultBody []byte
	if resultBody, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	// Close the body reader
	defer func() {
		_ = resp.Body.Close()
	}()

	// Test the status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error code %d. response: %s", resp.StatusCode, resultBody)
	}

	// Return the body as a string
	return string(resultBody), nil
}
