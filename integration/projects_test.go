package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	codeship "github.com/codeship/codeship-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	projectName = "codeship/codeship-go"
	projectUUID = "c38f3280-792b-0135-21bb-4e0cf8ff365b"
)

func TestListProjects(t *testing.T) {
	setup()

	projects, resp, err := org.ListProjects(context.Background())
	require.NoError(t, err)
	require.NotZero(t, projects)
	require.NotZero(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, projects.Projects)

	for {
		for _, p := range projects.Projects {
			// found our project
			if p.Name == projectName {
				assert.Equal(t, p.UUID, projectUUID)
				return
			}
		}

		if resp.IsLastPage() || resp.Next == "" {
			// we paged through all the results
			// and did not find our project
			t.FailNow()
		}

		next, _ := resp.NextPage()

		// so we don't hit our rate limit as fast
		time.Sleep(1 * time.Second)
		projects, resp, _ = org.ListProjects(context.Background(), codeship.Page(next))
	}
}

func TestGetProject(t *testing.T) {
	setup()

	project, resp, err := org.GetProject(context.Background(), projectUUID)
	require.NoError(t, err)
	require.NotZero(t, project)
	require.NotZero(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, project.UUID, projectUUID)
	assert.Equal(t, project.Name, projectName)
	assert.NotZero(t, project.CreatedAt)
	assert.NotZero(t, project.UpdatedAt)
}
