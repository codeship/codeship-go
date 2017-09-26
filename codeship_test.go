package codeship_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	codeship "github.com/codeship/codeship-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *codeship.Client
	org    *codeship.Organization
)

type optionalError struct {
	want  bool
	value error
}

type optionalString struct {
	want  bool
	value string
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

func assertHeaders(t *testing.T, headers http.Header) {
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, "application/json", headers.Get("Accept"))
}

func TestNew(t *testing.T) {
	type args struct {
		username string
		password string
		orgName  string
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
		err  optionalError
	}{
		{
			name: "requires username",
			args: args{
				username: "",
				password: "foo",
				orgName:  "codeship",
			},
			err: optionalError{want: true, value: errors.New("missing username or password")},
		},
		{
			name: "requires password",
			args: args{
				username: "foo",
				password: "",
				orgName:  "codeship",
			},
			err: optionalError{want: true, value: errors.New("missing username or password")},
		},
		{
			name: "prefers username param",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "codeship",
			},
			env: env{
				username: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "prefers password param",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "codeship",
			},
			env: env{
				password: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "uses env username if not passed in",
			args: args{
				username: "",
				password: "bar",
				orgName:  "codeship",
			},
			env: env{
				username: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "uses env password if not passed in",
			args: args{
				username: "foo",
				password: "",
				orgName:  "codeship",
			},
			env: env{
				password: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "requires organization name",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "",
			},
			err: optionalError{want: true, value: errors.New("organization name is required")},
		},
		{
			name: "handles error option func",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "codeship",
				opts: []codeship.Option{
					func(*codeship.Client) error {
						return errors.New("boom")
					},
				},
			},
			err: optionalError{want: true, value: errors.New("options parsing failed: boom")},
		},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				_ = os.Unsetenv("CODESHIP_USERNAME")
				_ = os.Unsetenv("CODESHIP_PASSWORD")
			}()

			if tt.env.username.want {
				_ = os.Setenv("CODESHIP_USERNAME", tt.env.username.value)
			}
			if tt.env.password.want {
				_ = os.Setenv("CODESHIP_PASSWORD", tt.env.password.value)
			}

			got, err := codeship.New(tt.args.username, tt.args.password, tt.args.opts...)

			if err != nil {
				if !tt.err.want {
					assert.Fail("Unexpected error: %s", err.Error())
				} else {
					assert.Equal(tt.err.value.Error(), err.Error())
				}
				return
			}

			assert.NotNil(got)

			if tt.env.username.want && tt.args.username == "" {
				assert.Equal(tt.env.username.value, got.Username)
			} else {
				assert.Equal(tt.args.username, got.Username)
			}

			if tt.env.password.want && tt.args.password == "" {
				assert.Equal(tt.env.password.value, got.Password)
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
		err     optionalError
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
			err: optionalError{want: true, value: errors.New("authentication failed: invalid credentials")},
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
			err: optionalError{want: true, value: errors.New("organization 'foo' not authorized. Authorized organizations: [{codeship 28123f10-e33d-5533-b53f-111ef8d7b14f [project.read project.write build.read build.write]}]")},
		},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		mux = http.NewServeMux()
		server = httptest.NewServer(mux)

		defer func() {
			server.Close()
		}()

		mux.HandleFunc("/auth", tt.handler)

		t.Run(tt.name, func(t *testing.T) {
			c, _ := codeship.New("username", "password", codeship.BaseURL(server.URL))
			got, err := c.Scope(context.Background(), tt.args.name)

			if err != nil {
				if !tt.err.want {
					assert.Fail("Unexpected error: %s", err.Error())
				} else {
					assert.Equal(tt.err.value.Error(), err.Error())
				}
				return
			}

			assert.NotNil(got)
			assert.Equal(tt.want.UUID, got.UUID)
			assert.Equal(tt.want.Name, got.Name)
			assert.Equal(tt.want.Scopes, got.Scopes)
		})
	}
}
