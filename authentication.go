package codeship

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Organization object holds organization information from Authentication.
type Organization struct {
	Name   string   `json:"name"`
	UUID   string   `json:"uuid"`
	Scopes []string `json:"scopes"`
}

// Authentication object holds access token and scope information.
type Authentication struct {
	AccessToken   string         `json:"access_token"`
	Organizations []Organization `json:"organizations"`
	ExpiresAt     int64          `json:"expires_at"`
}

// Authenticate swaps username/password for an authentication token and sets
// it in the API object for future requests.
func (api *API) Authenticate() error {
	var err error

	// Swap username/password for temporary auth token
	api.Authentication, err = api.authenticate()
	if err != nil {
		return errors.Wrap(err, "unable to exchange username/password for auth token")
	}

	// Get all organizations the user is authenticated with
	orgs := api.Authentication.GetOrganizations()
	var ok bool

	// Set current organization to the one they requested by name
	if api.Organization, ok = orgs[strings.ToLower(api.Organization.Name)]; !ok {
		validOrgs := ""
		for org := range orgs {
			validOrgs += " " + org
		}
		return fmt.Errorf("unable to find organization named %s. Valid options are: %s", api.Organization.Name, validOrgs)
	}
	return nil
}

// Exchange username and password for an authentication object.
func (api *API) authenticate() (Authentication, error) {
	path := "/auth"
	req, _ := http.NewRequest("POST", api.BaseURL+path, nil)
	req.SetBasicAuth(api.Username, api.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return Authentication{}, errors.Wrap(err, fmt.Sprintf("Unable to call %s%s", api.BaseURL, path))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	authentication := Authentication{}
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

	err = json.Unmarshal(body, &authentication)
	if err != nil {
		return Authentication{}, errors.Wrap(err, "unable to unmarshal JSON into Authentication")
	}

	return authentication, nil
}

// GetOrganizations returns a map of orgs with the org name being the key and Organization as the value.
func (auth *Authentication) GetOrganizations() map[string]Organization {
	orgs := map[string]Organization{}
	for _, org := range auth.Organizations {
		orgs[org.Name] = org
	}
	return orgs
}
