package codeship_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	codeship "github.com/codeship/codeship-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBuild(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w)
			},
			status: http.StatusAccepted,
		},
		{
			name: "project not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "project"))
			},
			status: http.StatusNotFound,
			err:    "unable to create build: project not found",
		},
		{
			name: "bad request",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, fmt.Sprintf(fixture("errors.json"), "ref is required"))
			},
			status: http.StatusBadRequest,
			err:    "unable to create build: ref is required",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds",
				tt.args.organizationUUID,
				tt.args.projectUUID),
				tt.handler)

			success, resp, err := org.CreateBuild(context.Background(), tt.args.projectUUID, "heads/master", "12345")

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.True(success)
		})
	}
}

func TestStopBuild(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w)
			},
			status: http.StatusAccepted,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to stop build: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/stop",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			success, resp, err := org.StopBuild(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.True(success)
		})
	}
}

func TestRestartBuild(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w)
			},
			status: http.StatusAccepted,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to restart build: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/restart",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			success, resp, err := org.RestartBuild(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.True(success)
		})
	}
}

func TestGetBuild(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("builds/get.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to get build: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			build, resp, err := org.GetBuild(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)

			finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:55.193+00:00")
			allocatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.967+00:00")
			queuedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:39.314+00:00")

			expected := codeship.Build{
				UUID:             "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
				ProjectID:        1,
				ProjectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				OrganizationUUID: "28123g10-e33d-5533-b57f-111ef8d7b14f",
				Ref:              "heads/master",
				CommitSha:        "185ab4c7dc4eda2a027c284f7a669cac3f50a5ed",
				Status:           "success",
				Username:         "fillup",
				CommitMessage:    "implemented interface for handling tests",
				FinishedAt:       finishedAt,
				AllocatedAt:      allocatedAt,
				QueuedAt:         queuedAt,
				Branch:           "test-branch",
				Links: codeship.BuildLinks{
					Services: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/services",
					Steps:    "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/steps",
				},
			}

			assert.Equal(expected, build)
		})
	}
}

func TestListBuilds(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("builds/list.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "project not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "project"))
			},
			status: http.StatusNotFound,
			err:    "unable to list builds: project not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds",
				tt.args.organizationUUID,
				tt.args.projectUUID),
				tt.handler)

			builds, resp, err := org.ListBuilds(context.Background(), tt.args.projectUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.Equal(2, len(builds.Builds))

			build := builds.Builds[0]

			finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:55.193+00:00")
			allocatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.967+00:00")
			queuedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:39.314+00:00")

			expected := codeship.Build{
				UUID:             "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
				ProjectID:        1,
				ProjectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				OrganizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				Ref:              "heads/master",
				CommitSha:        "185ab4c7dc4eda2a027c284f7a669cac3f50a5ed",
				Status:           "success",
				Username:         "fillup",
				CommitMessage:    "implemented interface for handling tests",
				FinishedAt:       finishedAt,
				AllocatedAt:      allocatedAt,
				QueuedAt:         queuedAt,
				Branch:           "test-branch",
				Links: codeship.BuildLinks{
					Services: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/services",
					Steps:    "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/steps",
				},
			}

			assert.Equal(expected, build)
			assert.Equal(2, builds.Total)
			assert.Equal(1, builds.Page)
			assert.Equal(30, builds.PerPage)
		})
	}
}

func TestListBuildPipelines(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("builds/pipelines.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to list pipelines: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/pipelines",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			pipelines, resp, err := org.ListBuildPipelines(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.Equal(1, len(pipelines.Pipelines))

			pipeline := pipelines.Pipelines[0]

			createdAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:16.556Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:38.394Z")
			finishedAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:38.391Z")

			expected := codeship.BuildPipeline{
				UUID:       "0a341890-a899-4492-9c94-86ef24527f05",
				BuildUUID:  "9ec4b230-76f8-0135-86b9-2ee351ae25fe",
				Type:       "build",
				Status:     "success",
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				FinishedAt: finishedAt,
				Metrics: codeship.BuildPipelineMetrics{
					AmiID:                 "ami-02322b79",
					Queries:               "112",
					CPUUser:               "1142",
					Duration:              "11",
					CPUSystem:             "499",
					InstanceID:            "i-0cfcd05a46d4cdb12",
					Architecture:          "trusty_64",
					InstanceType:          "i3.8xlarge",
					CPUPerSecond:          "136",
					DiskFreeBytes:         "128536784896",
					DiskUsedBytes:         "362098688",
					NetworkRxBytes:        "32221720",
					NetworkTxBytes:        "310269",
					MaxUsedConnections:    "1",
					MemoryMaxUsageInBytes: "665427968",
				},
			}

			assert.Equal(expected, pipeline)
			assert.Equal(1, pipelines.Total)
			assert.Equal(1, pipelines.Page)
			assert.Equal(30, pipelines.PerPage)
		})
	}
}

func TestListBuildServices(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("builds/services.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to list build services: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/services",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			services, resp, err := org.ListBuildServices(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.Equal(1, len(services.Services))

			service := services.Services[0]

			updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:44+00:00")
			finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")

			expected := codeship.BuildService{
				UUID:       "b46c6c6c-1bdb-4413-8e55-a9a8b1b27526",
				BuildUUID:  "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
				Name:       "test",
				Status:     "finished",
				UpdatedAt:  updatedAt,
				FinishedAt: finishedAt,
			}

			assert.Equal(expected, service)
			assert.Equal(1, services.Total)
			assert.Equal(1, services.Page)
			assert.Equal(30, services.PerPage)
		})
	}
}

func TestListBuildSteps(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
		buildUUID        string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("builds/steps.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "build not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "28123f10-e33d-5533-b53f-111ef8d7b14f",
				buildUUID:        "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf(fixture("not_found.json"), "build"))
			},
			status: http.StatusNotFound,
			err:    "unable to list build steps: build not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/steps",
				tt.args.organizationUUID,
				tt.args.projectUUID,
				tt.args.buildUUID),
				tt.handler)

			steps, resp, err := org.ListBuildSteps(context.Background(), tt.args.projectUUID, tt.args.buildUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			require.Equal(1, len(steps.Steps))

			step := steps.Steps[0]

			updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:44+00:00")
			startedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")
			buildingAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")
			finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:42+00:00")

			expected := codeship.BuildStep{
				UUID:        "21adb0ec-5139-4547-b1cf-7c40160b6e9d",
				ServiceUUID: "b46c6c6c-1bdb-4413-8e55-a9a8b1b27526",
				BuildUUID:   "28123f10-e33d-5533-b53f-111ef8d7b14f",
				Name:        "test",
				Type:        "run",
				Status:      "success",
				Command:     "./run-tests.sh",
				UpdatedAt:   updatedAt,
				StartedAt:   startedAt,
				BuildingAt:  buildingAt,
				FinishedAt:  finishedAt,
				Steps:       []codeship.BuildStep{},
			}

			assert.Equal(expected, step)
			assert.Equal(1, steps.Total)
			assert.Equal(1, steps.Page)
			assert.Equal(30, steps.PerPage)
		})
	}
}
