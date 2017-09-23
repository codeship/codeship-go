package codeship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Organization holds the configuration for the current API client scoped to the Organization. Should not
// be modified concurrently
type Organization struct {
	UUID   string
	Name   string
	Scopes []string
	client *Client
}

type Links struct {
	Next     string
	Previous string
	Last     string
	First    string
}

// Response is a Codeship response. This wraps the standard http.Response returned from Codeship.
type Response struct {
	*http.Response

	// Links that were returned with the response. These are parsed from the Link header.
	Links Links
}

func (r *Response) links() Links {
	var links Links

	// if linkText, ok := r.Response.Header["Link"]; ok {
	// }

	return links
}

const apiURL = "https://api.codeship.com/v2"

// Client holds information necessary to make a request to the Codeship API
type Client struct {
	Username       string
	Password       string
	baseURL        string
	authentication Authentication
	headers        http.Header
	httpClient     *http.Client
	logger         *log.Logger
	verbose        bool
}

// New creates a new Codeship API client
func New(username, password string, opts ...Option) (*Client, error) {
	if username == "" {
		username = os.Getenv("CODESHIP_USERNAME")
	}

	if password == "" {
		password = os.Getenv("CODESHIP_PASSWORD")
	}

	if username == "" || password == "" {
		return nil, errors.New("missing username or password")
	}

	client := &Client{
		Username: username,
		Password: password,
		baseURL:  apiURL,
		headers:  make(http.Header),
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, errors.Wrap(err, "options parsing failed")
	}

	// Fall back to http.DefaultClient if the user does not provide
	// their own
	if client.httpClient == nil {
		client.httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	// Fall back to default log.Logger (STDOUT) if the user does not provide
	// their own
	if client.logger == nil {
		client.logger = &log.Logger{}
		client.logger.SetOutput(os.Stdout)
	}

	return client, nil
}

// Scope scopes a client to a single Organization, allowing the user to make calls to the API
func (c *Client) Scope(name string) (*Organization, error) {
	if c.AuthenticationRequired() {
		if _, err := c.Authenticate(); err != nil {
			return nil, err
		}
	}

	for _, org := range c.authentication.Organizations {
		if org.Name == strings.ToLower(name) {
			return &Organization{
				UUID:   org.UUID,
				Name:   org.Name,
				Scopes: org.Scopes,
				client: c,
			}, nil
		}
	}
	return nil, ErrUnauthorized(fmt.Sprintf("organization '%s' not authorized. Authorized organizations: %v", name, c.authentication.Organizations))
}

// Authentication returns the client's current Authentication object
func (c *Client) Authentication() Authentication {
	return c.authentication
}

// AuthenticationRequired determines if a client must authenticate before making a request
func (c *Client) AuthenticationRequired() bool {
	return c.authentication.AccessToken == "" || c.authentication.ExpiresAt <= time.Now().Unix()
}

func (c *Client) request(method, path string, params interface{}) ([]byte, *Response, error) {
	url := c.baseURL + path
	// Replace nil with a JSON object if needed
	var reqBody io.Reader
	if params != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(params); err != nil {
			return nil, nil, err
		}
		reqBody = buf
	}

	if c.AuthenticationRequired() {
		if _, err := c.Authenticate(); err != nil {
			return nil, nil, err
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, errors.Wrap(err, "HTTP request creation failed")
	}

	// Apply any user-defined headers first
	req.Header = cloneHeader(c.headers)
	req.Header.Set("Authorization", "Bearer "+c.authentication.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.do(req)
}

func (c *Client) do(req *http.Request) ([]byte, *Response, error) {
	if c.verbose {
		dumpReq, _ := httputil.DumpRequest(req, true)
		c.logger.Println(string(dumpReq))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "HTTP request failed")
	}

	if c.verbose {
		dumpResp, _ := httputil.DumpResponse(resp, true)
		c.logger.Println(string(dumpResp))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	response := &Response{Response: resp}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, response, errors.Wrap(err, "could not read response body")
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		break
	case http.StatusUnauthorized:
		return nil, response, ErrUnauthorized("invalid credentials")
	case http.StatusForbidden:
		return nil, response, ErrUnauthorized("insufficient permissions")
	default:
		if len(body) > 0 {
			return nil, response, fmt.Errorf("HTTP status: %d; content %q", resp.StatusCode, string(body))
		}
		return nil, response, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	return body, response, nil
}

// cloneHeader returns a shallow copy of the header.
// copied from https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func cloneHeader(header http.Header) http.Header {
	h := make(http.Header)
	for k, vs := range header {
		h[k] = vs
	}
	return h
}
