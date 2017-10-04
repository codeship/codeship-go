package codeship_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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

type optionalString *string

func newOptionalString(value string) optionalString {
	return &value
}

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("auth/success.json"))
	})

	client, _ = codeship.New("test", "pass", codeship.BaseURL(server.URL))
	org, _ = client.Scope(context.Background(), "codeship")

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
		username string
		password string
		opts     []codeship.Option
	}
	type env struct {
		username optionalString
		password optionalString
	}
	tests := []struct {
		name string
		args args
		env  env
		want *codeship.Client
		err  error
	}{
		{
			name: "requires username",
			args: args{
				username: "",
				password: "foo",
			},
			err: errors.New("missing username or password"),
		},
		{
			name: "requires password",
			args: args{
				username: "foo",
				password: "",
			},
			err: errors.New("missing username or password"),
		},
		{
			name: "prefers username param",
			args: args{
				username: "foo",
				password: "bar",
			},
			env: env{
				username: newOptionalString("baz"),
			},
		},
		{
			name: "prefers password param",
			args: args{
				username: "foo",
				password: "bar",
			},
			env: env{
				password: newOptionalString("baz"),
			},
		},
		{
			name: "uses env username if not passed in",
			args: args{
				username: "",
				password: "bar",
			},
			env: env{
				username: newOptionalString("baz"),
			},
		},
		{
			name: "uses env password if not passed in",
			args: args{
				username: "foo",
				password: "",
			},
			env: env{
				password: newOptionalString("baz"),
			},
		},
		{
			name: "handles error option func",
			args: args{
				username: "foo",
				password: "bar",
				opts: []codeship.Option{
					func(*codeship.Client) error {
						return errors.New("boom")
					},
				},
			},
			err: errors.New("options parsing failed: boom"),
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				_ = os.Unsetenv("CODESHIP_USERNAME")
				_ = os.Unsetenv("CODESHIP_PASSWORD")
			}()

			if tt.env.username != nil {
				_ = os.Setenv("CODESHIP_USERNAME", *tt.env.username)
			}
			if tt.env.password != nil {
				_ = os.Setenv("CODESHIP_PASSWORD", *tt.env.password)
			}

			got, err := codeship.New(tt.args.username, tt.args.password, tt.args.opts...)

			if tt.err != nil {
				require.Error(err)
				assert.EqualError(tt.err, err.Error())
				return
			}

			require.NoError(err)
			assert.NotNil(got)

			if tt.env.username != nil && tt.args.username == "" {
				assert.Equal(*tt.env.username, got.Username)
			} else {
				assert.Equal(tt.args.username, got.Username)
			}

			if tt.env.password != nil && tt.args.password == "" {
				assert.Equal(*tt.env.password, got.Password)
			} else {
				assert.Equal(tt.args.password, got.Password)
			}
		})
	}
}

func TestScope(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		handler http.HandlerFunc
		args    args
		want    *codeship.Organization
		err     error
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
			err: errors.New("authentication failed: invalid credentials"),
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
			err: errors.New("organization 'foo' not authorized. Authorized organizations: [{codeship 28123f10-e33d-5533-b53f-111ef8d7b14f [project.read project.write build.read build.write]}]"),
		},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux = http.NewServeMux()
			server = httptest.NewServer(mux)

			defer func() {
				server.Close()
			}()

			mux.HandleFunc("/auth", tt.handler)

			c, _ := codeship.New("username", "password", codeship.BaseURL(server.URL))
			got, err := c.Scope(context.Background(), tt.args.name)

			if tt.err != nil {
				assert.Error(err)
				assert.Equal(tt.err.Error(), err.Error())
				return
			}

			assert.NoError(err)
			assert.NotNil(got)
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

	var buf bytes.Buffer
	logger := log.New(&buf, "INFO: ", log.Lshortfile)

	c, _ := codeship.New("username", "password",
		codeship.BaseURL(server.URL),
		codeship.Verbose(true),
		codeship.Logger(logger),
	)

	org, err := c.Scope(context.Background(), "codeship")

	assert := assert.New(t)

	assert.NoError(err)
	assert.NotNil(org)
	assert.True(buf.Len() > 0)
}
