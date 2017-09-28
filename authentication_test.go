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
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		status  int
		err     optionalError
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
			name: "unauthorized auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprint(w, fixture("auth/unauthorized.json"))
			},
			status: http.StatusUnauthorized,
			err:    optionalError(errors.New("authentication failed: invalid credentials")),
		},
		{
			name: "rate limit exceeded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
			},
			status: http.StatusForbidden,
			err:    optionalError(errors.New("authentication failed: rate limit exceeded")),
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			status: http.StatusInternalServerError,
			err:    optionalError(errors.New("authentication failed: HTTP status: 500")),
		},
		{
			name: "other status code",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
			},
			status: http.StatusTeapot,
			err:    optionalError(errors.New("authentication failed: HTTP status: 418")),
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
			err:    optionalError(errors.New("authentication failed: HTTP status: 418; content \"I'm a teapot\"")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux = http.NewServeMux()
			server = httptest.NewServer(mux)

			mux.HandleFunc("/auth", tt.handler)

			client, _ = codeship.New("username", "password", codeship.BaseURL(server.URL))
			org, _ = client.Scope(context.Background(), "codeship")

			defer func() {
				server.Close()
			}()

			assert := assert.New(t)

			resp, err := client.Authenticate(context.Background())
			assert.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if err != nil {
				if tt.err == nil {
					assert.Fail("Unexpected error: %s", err.Error())
				} else {
					assert.Equal(tt.err.Error(), err.Error())
				}
				return
			}
		})
	}
}
