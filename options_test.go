package codeship

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		args args
		want *http.Client
	}{
		{
			name: "sets client successfully",
			args: args{
				client: &http.Client{
					Timeout: 5 * time.Second,
				},
			},
			want: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeship, err := New(NewBasicAuth("username", "password"), HTTPClient(tt.args.client))

			require.NoError(err)
			assert.Equal(codeship.httpClient, tt.want)
		})
	}
}

func TestHeaders(t *testing.T) {
	type args struct {
		headers http.Header
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			name: "sets headers successfully",
			args: args{
				headers: http.Header{
					"Content-Type": []string{"text/xml"},
					"Accept":       []string{"text/html"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"text/xml"},
				"Accept":       []string{"text/html"},
			},
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeship, err := New(NewBasicAuth("username", "password"), Headers(tt.args.headers))

			require.NoError(err)
			assert.Equal(codeship.headers, tt.want)
		})
	}
}

func TestBaseURL(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sets baseURL successfully",
			args: args{
				baseURL: "http://localhost:8080/api/v2",
			},
			want: "http://localhost:8080/api/v2",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeship, err := New(NewBasicAuth("username", "password"), BaseURL(tt.args.baseURL))

			require.NoError(err)
			assert.Equal(codeship.baseURL, tt.want)
		})
	}
}

// avocadoLogger is a StdLogger with more avocados
type avocadoLogger struct{}

func (a *avocadoLogger) Println(v ...interface{}) {
	fmt.Printf("\xF0\x9F\xA5\x91 %v \xF0\x9F\xA5\x91\n", v...)
}

func TestLogger(t *testing.T) {
	type args struct {
		logger StdLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sets logger successfully",
			args: args{
				logger: log.New(os.Stderr, "DEBUG: ", log.LUTC),
			},
		},
		{
			name: "sets custom logger successfully",
			args: args{
				logger: &avocadoLogger{},
			},
		},
		{
			name: "sets third-party logger successfully",
			args: args{
				logger: logrus.New(),
			},
		},
	}

	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(NewBasicAuth("username", "password"), Logger(tt.args.logger))
			require.NoError(err)
		})
	}
}

func TestVerbose(t *testing.T) {
	type args struct {
		verbose bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "sets verbose successfully",
			args: args{
				verbose: true,
			},
			want: true,
		},
		{
			name: "unsets verbose successfully",
			args: args{
				verbose: false,
			},
			want: false,
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeship, err := New(NewBasicAuth("username", "password"), Verbose(tt.args.verbose))

			require.NoError(err)
			assert.Equal(codeship.verbose, tt.want)
		})
	}
}
