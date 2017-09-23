package codeship_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/codeship/codeship-go"

	"github.com/stretchr/testify/assert"
)

func TestCreateBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("POST", r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w)
	})

	_, resp, err := org.CreateBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "heads/master", "185ab4c7dc4eda2a027c284f7a669cac3f50a5ed")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusAccepted, resp.StatusCode)
}

func TestStopBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/stop", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("POST", r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w)
	})

	success, resp, err := org.StopBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert := assert.New(t)
	assert.NoError(err)
	assert.True(success)
	assert.NotNil(resp)
	assert.Equal(http.StatusAccepted, resp.StatusCode)
}

func TestRestartBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/restart", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("POST", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w)
	})

	success, resp, err := org.RestartBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert := assert.New(t)
	assert.NoError(err)
	assert.True(success)
	assert.NotNil(resp)
	assert.Equal(http.StatusAccepted, resp.StatusCode)
}

func TestGetBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/get.json"))
	})

	build, resp, err := org.GetBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusOK, resp.StatusCode)

	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:55.193+00:00")
	allocatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.967+00:00")
	queuedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:39.314+00:00")

	expected := codeship.Build{
		UUID:             "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
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
		Links: codeship.BuildLinks{
			Services: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/services",
			Steps:    "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/steps",
		},
	}

	assert.Equal(expected, build)
}

func TestListBuilds(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/list.json"))
	})

	builds, resp, err := org.ListBuilds("28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(2, len(builds.Builds))

	build := builds.Builds[0]

	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:55.193+00:00")
	allocatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.967+00:00")
	queuedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:39.314+00:00")

	expected := codeship.Build{
		UUID:             "25a3dd8c-eb3e-4e75-1298-8cbcbe621342",
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
		Links: codeship.BuildLinks{
			Services: "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/services",
			Steps:    "https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/steps",
		},
	}

	assert.Equal(expected, build)
	assert.Equal(2, builds.Total)
	assert.Equal(1, builds.Page)
	assert.Equal(30, builds.PerPage)
}

func TestListBuildPipelines(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/9ec4b230-76f8-0135-86b9-2ee351ae25fe/pipelines", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/pipelines.json"))
	})

	pipelines, resp, err := org.ListBuildPipelines("28123f10-e33d-5533-b53f-111ef8d7b14f", "9ec4b230-76f8-0135-86b9-2ee351ae25fe")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, len(pipelines.Pipelines))

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
}

func TestListBuildServices(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/28123f10-e33d-5533-b53f-111ef8d7b14f/services", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/services.json"))
	})

	buildServices, resp, err := org.ListBuildServices("28123f10-e33d-5533-b53f-111ef8d7b14f", "28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, len(buildServices.Services))

	service := buildServices.Services[0]

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
	assert.Equal(1, buildServices.Total)
	assert.Equal(1, buildServices.Page)
	assert.Equal(30, buildServices.PerPage)
}

func TestListBuildSteps(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/28123f10-e33d-5533-b53f-111ef8d7b14f/steps", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/steps.json"))
	})

	buildSteps, resp, err := org.ListBuildSteps("28123f10-e33d-5533-b53f-111ef8d7b14f", "28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(1, len(buildSteps.Steps))

	step := buildSteps.Steps[0]

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
	assert.Equal(1, buildSteps.Total)
	assert.Equal(1, buildSteps.Page)
	assert.Equal(30, buildSteps.PerPage)
}
