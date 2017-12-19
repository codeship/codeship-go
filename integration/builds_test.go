package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBuilds(t *testing.T) {
	setup()

	builds, resp, err := org.ListBuilds(context.Background(), proProjectUUID)
	require.NoError(t, err)
	require.NotEmpty(t, builds)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	build := builds.Builds[0]
	assert.Equal(t, organizationUUID, build.OrganizationUUID)
	assert.Equal(t, proProjectUUID, build.ProjectUUID)
	assert.NotEmpty(t, build.UUID)
	assert.NotZero(t, build.AllocatedAt)
	assert.NotZero(t, build.QueuedAt)
	assert.NotZero(t, build.FinishedAt)
}

func TestGetBuild(t *testing.T) {
	setup()

	builds, resp, err := org.ListBuilds(context.Background(), proProjectUUID)
	require.NoError(t, err)
	require.NotEmpty(t, builds)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	build := builds.Builds[0]
	require.NotEmpty(t, build.UUID)

	got, resp, err := org.GetBuild(context.Background(), proProjectUUID, build.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, build, got)
}

func TestListBuildPipelines(t *testing.T) {
	setup()

	builds, resp, err := org.ListBuilds(context.Background(), basicProjectUUID)
	require.NoError(t, err)
	require.NotEmpty(t, builds)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	build := builds.Builds[0]
	require.NotEmpty(t, build.UUID)

	pipelines, resp, err := org.ListBuildPipelines(context.Background(), basicProjectUUID, build.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, pipelines)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	pipeline := pipelines.Pipelines[0]
	assert.NotEmpty(t, pipeline)
	assert.Equal(t, build.UUID, pipeline.BuildUUID)
	assert.NotEmpty(t, pipeline.UUID)
}

func TestListBuildServices(t *testing.T) {
	setup()

	builds, resp, err := org.ListBuilds(context.Background(), proProjectUUID)
	require.NoError(t, err)
	require.NotEmpty(t, builds)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	build := builds.Builds[0]
	require.NotEmpty(t, build.UUID)

	services, resp, err := org.ListBuildServices(context.Background(), proProjectUUID, build.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, services)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	service := services.Services[0]
	assert.NotEmpty(t, service)
	assert.Equal(t, build.UUID, service.BuildUUID)
	assert.NotEmpty(t, service.UUID)
	assert.NotEmpty(t, service.Name)
}

func TestListBuildSteps(t *testing.T) {
	setup()

	builds, resp, err := org.ListBuilds(context.Background(), proProjectUUID)
	require.NoError(t, err)
	require.NotEmpty(t, builds)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	build := builds.Builds[0]
	require.NotEmpty(t, build.UUID)

	steps, resp, err := org.ListBuildSteps(context.Background(), proProjectUUID, build.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, steps)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	step := steps.Steps[0]
	assert.NotEmpty(t, step)
	assert.Equal(t, build.UUID, step.BuildUUID)
	assert.NotEmpty(t, step.UUID)
	assert.NotEmpty(t, step.Name)
}
