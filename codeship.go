package codeship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	errors "github.com/pkg/errors"
)

const apiURL = "https://api.codeship.com/v2"

// HTTP Status Codes

// OK - Status 200 - OK
const OK = 200

// BadRequest - Status 400 - Bad Request
const BadRequest = 400

// Unauthorized - Status 401 - Unauthorized
const Unauthorized = 401

// Forbidden - Status 403 - Forbidden
const Forbidden = 403

// MethodNotSupported - Status 405 - Method not supported
const MethodNotSupported = 405

// TooManyRequests - Status 429 - Too many requests
const TooManyRequests = 429

// ServerError - Status 500 - Server Error
const ServerError = 500

// API holds the configuration for the current API client. A client should not
// be modified concurrently.
type API struct {
	Username       string
	Password       string
	Authentication Authentication
	BaseURL        string
	DefaultOrg     string
	headers        http.Header
	httpClient     *http.Client
}

// New creates a new Codeship API client.
func New(username, password string, orgName string, opts ...Option) (*API, error) {
	if username == "" {
		username = os.Getenv("CODESHIP_USERNAME")
	}

	if password == "" {
		password = os.Getenv("CODESHIP_PASSWORD")
	}

	if username == "" || password == "" {
		return nil, fmt.Errorf("Missing username or password")
	}

	api := &API{
		Username: username,
		Password: password,
		BaseURL:  apiURL,
		headers:  make(http.Header),
	}

	err := api.parseOptions(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "options parsing failed")
	}

	// Fall back to http.DefaultClient if the package user does not provide
	// their own.
	if api.httpClient == nil {
		api.httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	// Swap username/password for temporary auth token
	api.Authentication, err = api.authenticate()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to exchange username/password for auth token")
	}

	// If orgName provided, get UUID for it and store in api.DefaultOrg
	// if orgName != "" {
	// 	orgMap := api.Authentication.GetOrgMap()
	// 	ok := false
	// 	if api.DefaultOrg, ok = orgMap[orgName]; !ok {
	// 		validOrgs := ""
	// 		for org := range orgMap {
	// 			validOrgs += " " + org
	// 		}
	// 		return api, fmt.Errorf("API initialized successfully, but unable to find organization named %s. Valid options are: %s", orgName, validOrgs)
	// 	}
	// }

	return api, nil
}

func (api *API) makeRequest(method, path string, params interface{}) ([]byte, error) {
	url := api.BaseURL + path
	// Replace nil with a JSON object if needed
	var reqBody io.Reader
	if params != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(params); err != nil {
			return nil, err
		}
		reqBody = buf
	}

	resp, err := api.request(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		break
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("HTTP status %d: invalid credentials", resp.StatusCode)
	case http.StatusForbidden:
		return nil, fmt.Errorf("HTTP status %d: insufficient permissions", resp.StatusCode)
	default:
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("HTTP status %d: server error", resp.StatusCode)
		}

		var s string
		if body != nil {
			s = string(body)
		}
		return nil, fmt.Errorf("HTTP status %d: content %q", resp.StatusCode, s)
	}

	return body, nil
}

// request makes a HTTP request to the given API endpoint, returning the raw
// *http.Response, or an error if one occurred. The caller is responsible for
// closing the response body.
func (api *API) request(method, url string, reqBody io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request creation failed")
	}

	// Apply any user-defined headers first.
	req.Header = cloneHeader(api.headers)
	req.Header.Set("Authorization", "Bearer "+api.Authentication.AccessToken)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}

	return resp, nil
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

func (api *API) getOrgUUID(orgID string) string {
	if orgID != "" {
		return orgID
	}

	return api.DefaultOrg
}
