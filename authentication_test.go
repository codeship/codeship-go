package codeship_test

import (
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
		err     optionalError
	}{
		{
			name: "successful auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, fixture("auth/success.json"))
			},
		},
		{
			name: "unauthorized auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprint(w, fixture("auth/unauthorized.json"))
			},
			err: optionalError{want: true, value: errors.New("authentication failed: invalid credentials")},
		},
		{
			name: "forbidden auth",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)

				fmt.Fprint(w, fixture("auth/unauthorized.json"))
			},
			err: optionalError{want: true, value: errors.New("authentication failed: insufficient permissions")},
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			err: optionalError{want: true, value: errors.New("authentication failed: HTTP status: 500")},
		},
		{
			name: "other status code",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
			},
			err: optionalError{want: true, value: errors.New("authentication failed: HTTP status: 418")},
		},
		{
			name: "other status code with body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTeapot)
				fmt.Fprint(w, "I'm a teapot")
			},
			err: optionalError{want: true, value: errors.New("authentication failed: HTTP status: 418; content \"I'm a teapot\"")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux = http.NewServeMux()
			server = httptest.NewServer(mux)

			mux.HandleFunc("/auth", tt.handler)

			client, _ = codeship.New("test", "pass", codeship.BaseURL(server.URL))
			org, _ = client.Scope("codeship")

			defer func() {
				server.Close()
			}()

			assert := assert.New(t)

			err := client.Authenticate()
			if err != nil {
				if !tt.err.want {
					assert.Fail("Unexpected error: %s", err.Error())
				} else {
					assert.Equal(tt.err.value.Error(), err.Error())
				}
				return
			}
		})
	}
}
