package codeship_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	codeship "github.com/codeship/codeship-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		status  int
		err     error
	}{
		{
			name: "successful auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, fixture("auth/success.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, "{ \"foo\": }")
			},
			status: http.StatusOK,
			err:    errors.New("unable to unmarshal JSON into Authentication: invalid character '}' looking for beginning of value"),
		},
		{
			name: "unauthorized auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprint(w, fixture("auth/unauthorized.json"))
			},
			status: http.StatusUnauthorized,
			err:    errors.New("authentication failed: invalid credentials"),
		},
		{
			name: "rate limit exceeded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
			},
			status: http.StatusForbidden,
			err:    errors.New("authentication failed: rate limit exceeded"),
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			status: http.StatusInternalServerError,
			err:    errors.New("authentication failed: HTTP status: 500"),
		},
		{
			name: "other status code",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
			},
			status: http.StatusTeapot,
			err:    errors.New("authentication failed: HTTP status: 418"),
		},
		{
			name: "other status code with body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
				fmt.Fprint(w, "I'm a teapot")
			},
			status: http.StatusTeapot,
			err:    errors.New("authentication failed: HTTP status: 418; content \"I'm a teapot\""),
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux = http.NewServeMux()
			server = httptest.NewServer(mux)

			mux.HandleFunc("/auth", tt.handler)

			client, _ = codeship.New("username", "password", codeship.BaseURL(server.URL))
			org, _ = client.Organization(context.Background(), "codeship")

			defer func() {
				server.Close()
			}()

			resp, err := client.Authenticate(context.Background())
			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err == nil {
				require.NoError(err)
			} else {
				require.Error(err)
				assert.EqualError(tt.err, err.Error())
			}
		})
	}
}
