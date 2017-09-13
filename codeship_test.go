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
	type setup struct {
		usernameEnv optionalString
		passwordEnv optionalString
	}
	tests := []struct {
		name  string
		args  args
		setup setup
		want  *codeship.API
		err   optionalError
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
			setup: setup{
				usernameEnv: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "prefers password param",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "codeship",
			},
			setup: setup{
				passwordEnv: optionalString{want: true, value: "baz"},
			},
		},
		{
			name: "requires organization name",
			args: args{
				username: "foo",
				password: "bar",
				orgName:  "",
			},
			err: optionalError{want: true, value: errors.New("missing username or password")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup.usernameEnv.want {
				os.Setenv("CODESHIP_USERNAME", tt.setup.usernameEnv.value)
			}
			if tt.setup.passwordEnv.want {
				os.Setenv("CODESHIP_PASSWORD", tt.setup.passwordEnv.value)
			}

			got, err := codeship.New(tt.args.username, tt.args.password, tt.args.orgName, tt.args.opts...)

			if err != nil {
				if !tt.err.want {
					assert.Fail(t, "Unexpected error: %v", err)
				}
				assert.Error(t, err, tt.err.value)
				return
			}

			assert.NotNil(t, got)

			if tt.setup.usernameEnv.want && tt.args.username == "" {
				assert.Equal(t, got.Username, tt.setup.usernameEnv.value)
			} else {
				assert.Equal(t, got.Username, tt.args.username)
			}

			if tt.setup.passwordEnv.want && tt.args.password == "" {
				assert.Equal(t, got.Password, tt.setup.passwordEnv.value)
			} else {
				assert.Equal(t, got.Password, tt.args.password)
			}
		})
	}
}
