package codeship_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	codeship "github.com/codeship/codeship-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *codeship.Client
	org    *codeship.Organization
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("auth/success.json"))
	})

	client, _ = codeship.New(codeship.NewBasicAuth("test", "pass"), codeship.BaseURL(server.URL))
	org, _ = client.Organization(context.Background(), "codeship")

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestNew(t *testing.T) {
	type args struct {
		auth codeship.Authenticator
		opts []codeship.Option
	}
	tests := []struct {
		name string
		args args
		err  string
	}{
		{
			name: "basic auth happy path",
			args: args{
				auth: codeship.NewBasicAuth("foo", "bar"),
			},
		},
		{
			name: "requires authenticator",
			args: args{
				auth: nil,
			},
			err: "no authenticator provided",
		},
		{
			name: "handles error option func",
			args: args{
				auth: codeship.NewBasicAuth("foo", "bar"),
				opts: []codeship.Option{
					func(*codeship.Client) error {
						return errors.New("boom")
					},
				},
			},
			err: "options parsing failed: boom",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := codeship.New(tt.args.auth, tt.args.opts...)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.NotNil(got)
		})
	}
}

func TestOrganization(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		handler http.HandlerFunc
		args    args
		want    *codeship.Organization
		err     string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, fixture("auth/success.json"))
			},
			args: args{
				name: "codeship",
			},
			want: &codeship.Organization{
				Name: "codeship",
				UUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				Scopes: []string{
					"project.read",
					"project.write",
					"build.read",
					"build.write",
				},
			},
		},
		{
			name: "unauthorized",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprint(w, fixture("auth/unauthorized.json"))
			},
			args: args{
				name: "codeship",
			},
			err: "authentication failed: invalid credentials",
		},
		{
			name: "wrong organization",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				fmt.Fprint(w, fixture("auth/success.json"))
			},
			args: args{
				name: "foo",
			},
			err: "organization 'foo' not authorized. Authorized organizations: [{codeship 28123f10-e33d-5533-b53f-111ef8d7b14f [project.read project.write build.read build.write]}]",
		},
		{
			name: "empty organization",
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			args: args{
				name: "",
			},
			err: "no organization name provided",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux = http.NewServeMux()
			server = httptest.NewServer(mux)

			defer func() {
				server.Close()
			}()

			mux.HandleFunc("/auth", tt.handler)

			c, _ := codeship.New(codeship.NewBasicAuth("username", "password"), codeship.BaseURL(server.URL))
			got, err := c.Organization(context.Background(), tt.args.name)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.NotNil(got)
			assert.Equal(tt.want.UUID, got.UUID)
			assert.Equal(tt.want.Name, got.Name)
			assert.Equal(tt.want.Scopes, got.Scopes)
		})
	}
}

func TestVerboseLogger(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	defer func() {
		server.Close()
	}()

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, fixture("auth/success.json"))
	})

	var (
		buf bytes.Buffer
		err error
	)

	logger := log.New(&buf, "INFO: ", log.Lshortfile)

	c, _ := codeship.New(
		codeship.NewBasicAuth("username", "password"),
		codeship.BaseURL(server.URL),
		codeship.Verbose(true),
		codeship.Logger(logger),
	)

	org, err = c.Organization(context.Background(), "codeship")

	assert := assert.New(t)
	require := require.New(t)

	require.NoError(err)
	require.NotNil(org)
	assert.True(buf.Len() > 0)
}
