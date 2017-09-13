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
	ExpiresAt int `json:"expires_at"`
}

func (api *API) Authenticate() error {
	var err error

	// Swap username/password for temporary auth token
	api.Authentication, err = api.authenticate()
	if err != nil {
		return errors.Wrap(err, "unable to exchange username/password for auth token")
	}

	// Get OrganizationUUID based on orgName
	orgs := api.Authentication.GetOrganizations()
	var ok bool

	if api.Organization.UUID, ok = orgs[api.Organization.Name]; !ok {
		validOrgs := ""
		for org := range orgs {
			validOrgs += " " + org
		}
		return fmt.Errorf("unable to find organization named %s. Valid options are: %s", api.Organization.Name, validOrgs)
	}
	return nil
}

// Exchange username and password for an authentication object
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

// GetOrganizations Return a map of orgs with the org name being the key and uuid as value
func (auth *Authentication) GetOrganizations() map[string]string {
	orgMap := map[string]string{}
	for _, org := range auth.Organizations {
		orgMap[org.Name] = org.UUID
	}
	return orgMap
}
