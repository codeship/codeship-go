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

func TestListProjects(t *testing.T) {
	setup()

	p, resp, err := org.ListProjects(context.Background())
	require.NoError(t, err)
	require.NotZero(t, p)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NotEmpty(t, p.Projects)
	projects := p.Projects

	for {
		for _, project := range projects {
			// found our project
			if project.UUID == proProjectUUID {
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
		p, resp, _ = org.ListProjects(context.Background(), codeship.Page(next))
		projects = p.Projects
	}
}

func TestGetProject(t *testing.T) {
	setup()

	project, resp, err := org.GetProject(context.Background(), proProjectUUID)
	require.NoError(t, err)
	require.NotZero(t, project)
	require.NotZero(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, organizationUUID, project.OrganizationUUID)
	assert.Equal(t, proProjectUUID, project.UUID)
	assert.Equal(t, proProjectName, project.Name)
	assert.NotZero(t, project.CreatedAt)
	assert.NotZero(t, project.UpdatedAt)
}
