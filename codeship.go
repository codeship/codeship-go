package codeship

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const apiURL = "https://api.codeship.com/v2"

// Client holds information necessary to make a request to the Codeship API
type Client struct {
	Username       string
	Password       string
	BaseURL        string
	authentication Authentication
	headers        http.Header
	httpClient     *http.Client
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
		BaseURL:  apiURL,
		headers:  make(http.Header),
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, errors.Wrap(err, "options parsing failed")
	}

	// Fall back to http.DefaultClient if the package user does not provide
	// their own
	if client.httpClient == nil {
		client.httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return client, nil
}

// Scope scopes a client to a single Organization, allowing the user to make calls to the API
func (c *Client) Scope(name string) (*Organization, error) {
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
	return nil, fmt.Errorf("organization %s not authorized. Authorized organizations: %v", name, c.authentication.Organizations)
}
