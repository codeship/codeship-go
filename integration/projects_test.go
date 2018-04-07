// +build integration

package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProjects(t *testing.T) {
	projects, resp, err := org.ListProjects(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, projects)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	for _, project := range projects.Projects {
		assert.Equal(t, organizationUUID, project.OrganizationUUID)
		assert.NotZero(t, project.ID)
		assert.NotEmpty(t, project.UUID)
		assert.NotEmpty(t, project.Name)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	}
}

func TestGetProject(t *testing.T) {
	tests := []struct {
		name, projectName, projectUUID string
	}{
		{
			name:        "pro project",
			projectUUID: proProjectUUID,
			projectName: proProjectName,
		},
		{
			name:        "basic project",
			projectUUID: basicProjectUUID,
			projectName: basicProjectName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			project, resp, err := org.GetProject(context.Background(), tc.projectUUID)
			require.NoError(t, err)
			require.NotZero(t, project)
			require.NotZero(t, resp)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			assert.Equal(t, organizationUUID, project.OrganizationUUID)
			assert.NotZero(t, project.ID)
			assert.Equal(t, tc.projectUUID, project.UUID)
			assert.Equal(t, tc.projectName, project.Name)
			assert.NotZero(t, project.CreatedAt)
			assert.NotZero(t, project.UpdatedAt)
		})
	}
}
