package codeship_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w)
	})

	_, err := org.CreateBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "heads/master", "185ab4c7dc4eda2a027c284f7a669cac3f50a5ed")

	assert.NoError(t, err)
}

func TestStopBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/stop", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w)
	})

	_, err := org.StopBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert.NoError(t, err)
}

func TestRestartBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/restart", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, ``)
	})

	_, err := org.RestartBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert.NoError(t, err)
}

func TestGetBuild(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/get.json"))
	})

	build, err := org.GetBuild("28123f10-e33d-5533-b53f-111ef8d7b14f", "25a3dd8c-eb3e-4e75-1298-8cbcbe621342")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal("25a3dd8c-eb3e-4e75-1298-8cbcbe621342", build.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", build.ProjectUUID)
	assert.Equal("28123g10-e33d-5533-b57f-111ef8d7b14f", build.OrganizationUUID)
	assert.Equal("heads/master", build.Ref)
	assert.Equal("185ab4c7dc4eda2a027c284f7a669cac3f50a5ed", build.CommitSha)
	assert.Equal("success", build.Status)
	assert.Equal("fillup", build.Username)
	assert.Equal("implemented interface for handling tests", build.CommitMessage)
	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:55.193+00:00")
	assert.Equal(finishedAt, build.FinishedAt)
	allocatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.967+00:00")
	assert.Equal(allocatedAt, build.AllocatedAt)
	queuedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:39.314+00:00")
	assert.Equal(queuedAt, build.QueuedAt)
	assert.NotEmpty(build.Links)
	assert.Equal("https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/services", build.Links.Services)
	assert.Equal("https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/25a3dd8c-eb3e-4e75-1298-8cbcbe621342/steps", build.Links.Steps)
	assert.Zero(build.Links.Pipelines)
}

func TestListBuilds(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/list.json"))
	})

	builds, err := org.ListBuilds("28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(2, len(builds.Builds))
	assert.Equal("25a3dd8c-eb3e-4e75-1298-8cbcbe621342", builds.Builds[0].UUID)
	assert.Equal("25a3dd8c-eb3e-4e75-1298-8cbcbe611111", builds.Builds[1].UUID)
	assert.Equal(2, builds.Total)
	assert.Equal(1, builds.Page)
	assert.Equal(30, builds.PerPage)
}

func TestGetBuildPipelines(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/9ec4b230-76f8-0135-86b9-2ee351ae25fe/pipelines", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/pipelines.json"))
	})

	buildPipelines, err := org.GetBuildPipelines("28123f10-e33d-5533-b53f-111ef8d7b14f", "9ec4b230-76f8-0135-86b9-2ee351ae25fe")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(buildPipelines.Pipelines))

	pipeline := buildPipelines.Pipelines[0]
	assert.Equal("0a341890-a899-4492-9c94-86ef24527f05", pipeline.UUID)
	assert.Equal("9ec4b230-76f8-0135-86b9-2ee351ae25fe", pipeline.BuildUUID)
	assert.Equal("build", pipeline.Type)
	assert.Equal("success", pipeline.Status)
	createdAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:16.556Z")
	assert.Equal(createdAt, pipeline.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:38.394Z")
	assert.Equal(updatedAt, pipeline.UpdatedAt)
	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-11T19:54:38.391Z")
	assert.Equal(finishedAt, pipeline.FinishedAt)
	assert.Equal(1, buildPipelines.Total)
	assert.Equal(1, buildPipelines.Page)
	assert.Equal(30, buildPipelines.PerPage)
}

func TestGetBuildServices(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/28123f10-e33d-5533-b53f-111ef8d7b14f/services", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/services.json"))
	})

	buildServices, err := org.GetBuildServices("28123f10-e33d-5533-b53f-111ef8d7b14f", "28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(buildServices.Services))
	service := buildServices.Services[0]
	assert.Equal("b46c6c6c-1bdb-4413-8e55-a9a8b1b27526", service.UUID)
	assert.Equal("25a3dd8c-eb3e-4e75-1298-8cbcbe621342", service.BuildUUID)
	assert.Equal("test", service.Name)
	assert.Equal("finished", service.Status)
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:44+00:00")
	assert.Equal(updatedAt, service.UpdatedAt)
	assert.Zero(service.PullingAt)
	assert.Zero(service.BuildingAt)
	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")
	assert.Equal(service.FinishedAt, finishedAt)
	assert.Equal(1, buildServices.Total)
	assert.Equal(1, buildServices.Page)
	assert.Equal(30, buildServices.PerPage)
}

func TestGetBuildSteps(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/28123f10-e33d-5533-b53f-111ef8d7b14f/builds/28123f10-e33d-5533-b53f-111ef8d7b14f/steps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("builds/steps.json"))
	})

	buildSteps, err := org.GetBuildSteps("28123f10-e33d-5533-b53f-111ef8d7b14f", "28123f10-e33d-5533-b53f-111ef8d7b14f")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(buildSteps.Steps))
	step := buildSteps.Steps[0]
	assert.Equal("21adb0ec-5139-4547-b1cf-7c40160b6e9d", step.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", step.BuildUUID)
	assert.Equal("b46c6c6c-1bdb-4413-8e55-a9a8b1b27526", step.ServiceUUID)
	assert.Equal("test", step.Name)
	assert.Equal("", step.Tag)
	assert.Equal("run", step.Type)
	assert.Equal("success", step.Status)
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:44+00:00")
	assert.Equal(updatedAt, step.UpdatedAt)
	startedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")
	assert.Equal(startedAt, step.StartedAt)
	buildingAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:41+00:00")
	assert.Equal(buildingAt, step.BuildingAt)
	finishedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:42+00:00")
	assert.Equal(finishedAt, step.FinishedAt)
	assert.Empty(step.Steps)
	assert.Equal("./run-tests.sh", step.Command)
	assert.Equal("", step.ImageName)
	assert.Equal("", step.Registry)
	assert.Equal(1, buildSteps.Total)
	assert.Equal(1, buildSteps.Page)
	assert.Equal(30, buildSteps.PerPage)
}
