package ses

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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

	// AccessKeyID is your Amazon AWS access key ID.
	AccessKeyID string

	// SecretAccessKey is your Amazon AWS secret key.
	SecretAccessKey string
}

// EnvConfig takes the access key ID and secret access key values from the environment variables
// $AWS_ACCESS_KEY_ID and $AWS_SECRET_KEY, respectively.
var EnvConfig = Config{
	Endpoint:        os.Getenv("AWS_SES_ENDPOINT"),
	AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
	SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
}

// SendEmail sends a plain text email. Note that from must be a verified
// address in the AWS control panel.
func (c *Config) SendEmail(from string, to, cc, bcc []string, subject, body string) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendEmail")
	data.Add("Source", from)

	if to != nil {
		for i := 0; i < len(to); i++ {
			data.Add(fmt.Sprintf("Destination.ToAddresses.member.%d", i+11), to[i])
		}
	}
	if cc != nil {
		for i := 0; i < len(cc); i++ {
			data.Add(fmt.Sprintf("Destination.CcAddresses.member.%d", i+1), cc[i])
		}
	}
	if bcc != nil {
		for i := 0; i < len(bcc); i++ {
			data.Add(fmt.Sprintf("Destination.BccAddresses.member.%d", i+1), bcc[i])
		}
	}

	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Text.Data", body)
	data.Add("AWSAccessKeyId", c.AccessKeyID)

	return sesPost(data, c.Endpoint, c.AccessKeyID, c.SecretAccessKey)
}

// SendEmailHTML sends a HTML email. Note that from must be a verified address
// in the AWS control panel.
func (c *Config) SendEmailHTML(from string, to, cc, bcc []string, subject, bodyText, bodyHTML string) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendEmail")
	data.Add("Source", from)

	if to != nil {
		for i := 0; i < len(to); i++ {
			data.Add(fmt.Sprintf("Destination.ToAddresses.member.%d", i+1), to[i])
		}
	}
	if cc != nil {
		for i := 0; i < len(cc); i++ {
			data.Add(fmt.Sprintf("Destination.CcAddresses.member.%d", i+1), cc[i])
		}
	}
	if bcc != nil {
		for i := 0; i < len(bcc); i++ {
			data.Add(fmt.Sprintf("Destination.BccAddresses.member.%d", i+1), bcc[i])
		}
	}

	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Text.Data", bodyText)
	data.Add("Message.Body.Html.Data", bodyHTML)
	data.Add("AWSAccessKeyId", c.AccessKeyID)

	return sesPost(data, c.Endpoint, c.AccessKeyID, c.SecretAccessKey)
}

// SendRawEmail sends a raw email. Note that from must be a verified address
// in the AWS control panel.
func (c *Config) SendRawEmail(raw []byte) (string, error) {
	data := make(url.Values)
	data.Add("Action", "SendRawEmail")
	data.Add("RawMessage.Data", base64.StdEncoding.EncodeToString(raw))
	data.Add("AWSAccessKeyId", c.AccessKeyID)

	return sesPost(data, c.Endpoint, c.AccessKeyID, c.SecretAccessKey)
}

func authorizationHeader(date, accessKeyID, secretAccessKey string) []string {
	h := hmac.New(sha256.New, []uint8(secretAccessKey))
	h.Write([]uint8(date))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	auth := fmt.Sprintf("AWS3-HTTPS AWSAccessKeyId=%s, Algorithm=HmacSHA256, Signature=%s", accessKeyID, signature)
	return []string{auth}
}

func sesGet(data url.Values, endpoint, accessKeyID, secretAccessKey string) (string, error) {
	urlString := fmt.Sprintf("%s?%s", endpoint, data.Encode())
	endpointURL, _ := url.Parse(urlString)
	headers := map[string][]string{}

	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 -0700")
	headers["Date"] = []string{date}

	h := hmac.New(sha256.New, []uint8(secretAccessKey))
	h.Write([]uint8(date))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	auth := fmt.Sprintf("AWS3-HTTPS AWSAccessKeyId=%s, Algorithm=HmacSHA256, Signature=%s", accessKeyID, signature)
	headers["X-Amzn-Authorization"] = []string{auth}

	req := http.Request{
		Close:      true,
		Header:     headers,
		Method:     "GET",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        endpointURL,
	}

	r, err := http.DefaultClient.Do(&req)
	if err != nil {
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

func sesPost(data url.Values, endpoint, accessKeyID, secretAccessKey string) (string, error) {
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 -0700")
	req.Header.Set("Date", date)

	h := hmac.New(sha256.New, []uint8(secretAccessKey))
	h.Write([]uint8(date))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	auth := fmt.Sprintf("AWS3-HTTPS AWSAccessKeyId=%s, Algorithm=HmacSHA256, Signature=%s", accessKeyID, signature)
	req.Header.Set("X-Amzn-Authorization", auth)

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
