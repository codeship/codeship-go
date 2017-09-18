package codeship_test

import (
	"os"
	"testing"

	codeship "github.com/codeship/codeship-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type optionalError struct {
	want  bool
	value error
}

type optionalString struct {
	want  bool
	value string
}

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func teardown() {
	os.Unsetenv("CODESHIP_USERNAME")
	os.Unsetenv("CODESHIP_PASSWORD")
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env.username.want {
				os.Setenv("CODESHIP_USERNAME", tt.env.username.value)
			}
			if tt.env.password.want {
				os.Setenv("CODESHIP_PASSWORD", tt.env.password.value)
			}

			got, err := codeship.New(tt.args.username, tt.args.password, tt.args.opts...)

			if err != nil && !tt.err.want {
				assert.Fail(t, "Unexpected error: %s", err.Error())
				return
			} else if err != nil {
				assert.Equal(t, tt.err.value.Error(), err.Error())
				return
			}

			assert.NotNil(t, got)

			if tt.env.username.want && tt.args.username == "" {
				assert.Equal(t, tt.env.username.value, got.Username)
			} else {
				assert.Equal(t, tt.args.username, got.Username)
			}

			if tt.env.password.want && tt.args.password == "" {
				assert.Equal(t, tt.env.password.value, got.Password)
			} else {
				assert.Equal(t, tt.args.password, got.Password)
			}
		})
	}
}
