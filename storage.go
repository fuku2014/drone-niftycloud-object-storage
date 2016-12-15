package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/smartystreets/go-aws-auth"
)

// ErrorResponse struct
type ErrorResponse struct {
	Error Error `xml:"Error"`
}

// Error struct
type Error struct {
	Code string `xml:"Code"`
}

// Error return ErrorMesssage
func (err Error) Error() string {
	return err.Code
}

// Client interface
type Client interface {
	PutObject(bucket, key string, file string, acl string) (err error)
}

// client struct
type client struct {
	endpoint    string
	credentials []awsauth.Credentials
}

// NewClient returns a new clauth  API client
func NewClient(endpoint string, credentials ...awsauth.Credentials) Client {
	return &client{
		endpoint:    endpoint,
		credentials: credentials,
	}
}

// PutObject create object to storage
func (c *client) PutObject(bucket, key string, file string, acl string) (err error) {
	resource := fmt.Sprintf("/%s/%s", bucket, key)

	body, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	req := c.makeRequest("PUT", resource, bytes.NewReader(body), nil)

	fileSize := int64(len(body))
	req.ContentLength = fileSize

	sum := c.getContentSum(req)
	req.Header.Add("Content-MD5", sum)

	content := contentType(file)
	req.Header.Add("Content-Type", content)

	req.Header.Add("x-amz-acl", acl)

	awsauth.SignS3(req, c.credentials...)
	err = c.doRequest(req, nil)
	return
}

// makeRequest return clauth  request
func (c *client) makeRequest(method, resource string, body io.Reader, query map[string]string) (req *http.Request) {
	endpoint, _ := url.Parse(c.endpoint)
	endpoint.Path = path.Join(endpoint.Path, resource)

	req, err := http.NewRequest(method, endpoint.String(), body)

	if err != nil {
		return
	}

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return
}

// doRequest sends an API request and returns the API response decode xml.
func (c *client) doRequest(req *http.Request, result interface{}) (err error) {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		var errResp = ErrorResponse{}
		err = xml.NewDecoder(res.Body).Decode(&errResp)
		if err != nil {
			return
		}
		return errResp.Error
	}

	if result != nil {
		err = xml.NewDecoder(res.Body).Decode(&result)
	}
	defer res.Body.Close()
	return
}

// getContentSum return base64 md5sum
func (c *client) getContentSum(req *http.Request) (sum string) {
	payload, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(payload))
	h := md5.New()
	h.Write(payload)
	sum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}

// contentType is a helper function that returns the content type for the file
// based on extension. If the file extension is unknown application/octet-stream
// is returned.
func contentType(file string) string {
	ext := filepath.Ext(file)
	typ := mime.TypeByExtension(ext)
	if typ == "" {
		typ = "application/octet-stream"
	}
	return typ
}
