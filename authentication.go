package codeship

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// ErrUnauthorized occurs when Codeship returns a 401 Unauthorized response
type ErrUnauthorized string

func (e ErrUnauthorized) Error() string {
	return string(e)
}

// Authentication object holds access token and scope information
type Authentication struct {
	AccessToken   string `json:"access_token,omitempty"`
	Organizations []struct {
		Name   string   `json:"name,omitempty"`
		UUID   string   `json:"uuid,omitempty"`
		Scopes []string `json:"scopes,omitempty"`
	} `json:"organizations,omitempty"`
	ExpiresAt int64 `json:"expires_at,omitempty"`
}

// Authenticate swaps username/password for an authentication token
//
// Codeship API docs: https://apidocs.codeship.com/v2/authentication/authentication-endpoint
func (c *Client) Authenticate(ctx context.Context) (Response, error) {
	path := "/auth"
	req, _ := http.NewRequest("POST", c.baseURL+path, nil)
	c.authenticator.SetAuth(req)
	req.Header.Set("Content-Type", "application/json")

	c.authentication = Authentication{}

	body, resp, err := c.do(req.WithContext(ctx))
	if err != nil {
		return resp, errors.Wrap(err, "authentication failed")
	}
	if err = json.Unmarshal(body, &c.authentication); err != nil {
		return resp, errors.Wrap(err, "unable to unmarshal JSON into Authentication")
	}

	return resp, err
}
