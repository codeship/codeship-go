package codeship

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	type args struct {
		client HTTPClient
	}
	tests := []struct {
		name string
		args args
		want HTTPClient
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
			codeship, err := New(NewBasicAuth("username", "password"), HttpClient(tt.args.client))

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

func TestLogger(t *testing.T) {
	type args struct {
		logger *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *log.Logger
	}{
		{
			name: "sets logger successfully",
			args: args{
				logger: log.New(os.Stderr, "DEBUG: ", log.LUTC),
			},
			want: log.New(os.Stderr, "DEBUG: ", log.LUTC),
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeship, err := New(NewBasicAuth("username", "password"), Logger(tt.args.logger))

			require.NoError(err)
			assert.Equal(codeship.logger, tt.want)
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
