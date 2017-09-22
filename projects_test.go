package codeship_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	codeship "github.com/codeship/codeship-go"
	"github.com/stretchr/testify/assert"
)

func TestListProjects(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/list.json"))
	})

	projects, err := org.ListProjects()

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(2, len(projects.Projects))

	project := projects.Projects[1]

	createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:09.556Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:55.252Z")

	expected := codeship.Project{
		UUID:               "83605ef0-76f8-0135-8810-6e5f001a2e3c",
		OrganizationUUID:   "28123f10-e33d-5533-b53f-111ef8d7b14f",
		Name:               "org/another-project",
		Type:               codeship.ProjectTypeBasic,
		RepositoryURL:      "https://github.com/org/another-project",
		RepositoryProvider: "github",
		AuthenticationUser: "Test User",
		NotificationRules: []codeship.NotificationRule{
			codeship.NotificationRule{
				Notifier:      "github",
				BranchMatch:   "exact",
				BuildStatuses: []string{"failed", "started", "recovered", "success"},
				Target:        "all",
			},
			codeship.NotificationRule{
				Notifier:      "email",
				BranchMatch:   "exact",
				Options:       codeship.NotificationOptions{},
				BuildStatuses: []string{"failed", "recovered"},
				Target:        "all",
			},
		},
		SSHKey:        "ssh-rsa key",
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		TeamIDs:       []int{1007, 1009},
		SetupCommands: []string{},
		TestPipelines: []codeship.TestPipeline{
			codeship.TestPipeline{
				Name:     "Test Commands",
				Commands: []string{"./run-tests.sh"},
			},
		},
		DeploymentPipelines:  []codeship.DeploymentPipeline{},
		EnvironmentVariables: []codeship.EnvironmentVariable{},
	}

	assert.Equal(expected, project)
	assert.Equal(1, projects.Page)
	assert.Equal(2, projects.Total)
	assert.Equal(30, projects.PerPage)
}

func TestGetProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/0059df30-7701-0135-8810-6e5f001a2e3c", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/get.json"))
	})

	project, err := org.GetProject("0059df30-7701-0135-8810-6e5f001a2e3c")

	assert := assert.New(t)
	assert.NoError(err)

	createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T20:19:55.199Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.336Z")

	expected := codeship.Project{
		UUID:               "0059df30-7701-0135-8810-6e5f001a2e3c",
		OrganizationUUID:   "28123f10-e33d-5533-b53f-111ef8d7b14f",
		Name:               "org/test-project",
		Type:               codeship.ProjectTypePro,
		RepositoryURL:      "https://github.com/org/test-project",
		RepositoryProvider: "github",
		AuthenticationUser: "Test User",
		NotificationRules: []codeship.NotificationRule{
			codeship.NotificationRule{
				Notifier:      "github",
				BranchMatch:   "exact",
				BuildStatuses: []string{"failed", "started", "recovered", "success"},
				Target:        "all",
				Options: codeship.NotificationOptions{
					Key:  "foo",
					Room: "devs",
					URL:  "https://google.com",
				},
			},
			codeship.NotificationRule{
				Notifier:      "email",
				BranchMatch:   "exact",
				Options:       codeship.NotificationOptions{},
				BuildStatuses: []string{"failed", "recovered"},
				Target:        "all",
			},
		},
		AesKey:    "aeskey",
		SSHKey:    "ssh-rsa key",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TeamIDs:   []int{1007},
	}

	assert.Equal(expected, project)
}

func TestCreateProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("POST", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, fixture("projects/create.json"))
	})

	project, err := org.CreateProject(codeship.ProjectCreateRequest{
		RepositoryURL: "git@github.com/org/repo-name",
		TestPipelines: []codeship.TestPipeline{
			{
				Commands: []string{"./run-tests.sh"},
				Name:     "run tests",
			},
		},
		Type: codeship.ProjectTypeBasic,
	})

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(project)
}

func TestUpdateProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/7de09100-7aeb-0135-b8e4-76a42f3a0b26", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("PUT", r.Method)
		assertHeaders(t, r.Header)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/update.json"))
	})

	project, err := org.UpdateProject("7de09100-7aeb-0135-b8e4-76a42f3a0b26", codeship.ProjectUpdateRequest{
		Type: codeship.ProjectTypePro,
		TeamIDs: []int{
			61593, 70000,
		},
	})

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(project)
}
