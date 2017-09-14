package codeship

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Authentication object holds access token and scope information
type Authentication struct {
	AccessToken   string `json:"access_token"`
	Organizations []struct {
		Name   string   `json:"name"`
		UUID   string   `json:"uuid"`
		Scopes []string `json:"scopes"`
	} `json:"organizations"`
	ExpiresAt int64 `json:"expires_at"`
}

// Authenticate swaps username/password for an authentication token
func (c *Client) Authenticate() error {
	var err error
	c.authentication, err = c.authenticate()
	if err != nil {
		return errors.Wrap(err, "unable to exchange username/password for auth token")
	}

	return nil
}

// Exchange username and password for an authentication object.
func (c *Client) authenticate() (Authentication, error) {
	path := "/auth"
	req, _ := http.NewRequest("POST", c.BaseURL+path, nil)
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Authentication{}, errors.Wrap(err, fmt.Sprintf("Unable to call %s%s", c.BaseURL, path))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Authentication{}, errors.Wrap(err, "unable to read API response")
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		break
	case http.StatusUnauthorized:
		return Authentication{}, fmt.Errorf("HTTP status %d: invalid credentials", resp.StatusCode)
	case http.StatusForbidden:
		return Authentication{}, fmt.Errorf("HTTP status %d: insufficient permissions", resp.StatusCode)
	default:
		if resp.StatusCode >= 500 {
			return Authentication{}, fmt.Errorf("HTTP status %d: service failure", resp.StatusCode)
		}

		var s string
		if body != nil {
			s = string(body)
		}
		return Authentication{}, fmt.Errorf("HTTP status %d: content %q", resp.StatusCode, s)
	}

	authentication := Authentication{}
	err = json.Unmarshal(body, &authentication)
	if err != nil {
		return authentication, errors.Wrap(err, "unable to unmarshal JSON into Authentication")
	}

	return authentication, nil
}
