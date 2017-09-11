package codeship

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	errors "github.com/pkg/errors"
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
	defer resp.Body.Close()

	authentication := Authentication{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Authentication{}, errors.Wrap(err, "Unable to read API response")
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		break
	case http.StatusUnauthorized:
		return Authentication{}, fmt.Errorf("HTTP status %d: invalid credentials", resp.StatusCode)
	case http.StatusForbidden:
		return Authentication{}, fmt.Errorf("HTTP status %d: insufficient permissions", resp.StatusCode)
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout,
		522, 523, 524:
		return Authentication{}, fmt.Errorf("HTTP status %d: service failure", resp.StatusCode)
	default:
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

// GetOrgMap Return a map of orgs with the org name being the key and uuid as value
func (auth *Authentication) GetOrgMap() map[string]string {
	orgMap := map[string]string{}
	for _, org := range auth.Organizations {
		orgMap[org.Name] = org.UUID
	}
	return orgMap
}
