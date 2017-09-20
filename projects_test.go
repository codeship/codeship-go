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
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/list.json"))
	})

	projects, err := org.ListProjects()

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(2, len(projects.Projects))
	project := projects.Projects[1]
	assert.Equal("83605ef0-76f8-0135-8810-6e5f001a2e3c", project.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", project.OrganizationUUID)
	assert.Equal("org/another-project", project.Name)
	assert.Equal(codeship.ProjectTypeBasic, project.Type)
	assert.Equal("https://github.com/org/another-project", project.RepositoryURL)
	assert.Equal("github", project.RepositoryProvider)
	assert.Equal("Test User", project.AuthenticationUser)
	assert.Equal(2, len(project.NotificationRules))
	notificationRule := project.NotificationRules[0]
	assert.Equal("github", notificationRule.Notifier)
	assert.Equal("exact", notificationRule.BranchMatch)
	assert.NotEmpty(notificationRule.BuildStatuses)
	assert.Equal("ssh-rsa key", project.SSHKey)
	createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:09.556Z")
	assert.Equal(createdAt, project.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:55.252Z")
	assert.Equal(updatedAt, project.UpdatedAt)
	assert.Equal(2, len(project.TeamIDs))
	assert.Equal(1007, project.TeamIDs[0])
	assert.Equal(1009, project.TeamIDs[1])
	assert.Equal(1, len(project.TestPipelines))
	assert.Equal("Test Commands", project.TestPipelines[0].Name)
	assert.Equal(1, len(project.TestPipelines[0].Commands))
	assert.Equal("./run-tests.sh", project.TestPipelines[0].Commands[0])
	assert.Equal(2, projects.Total)
	assert.Equal(1, projects.Page)
	assert.Equal(30, projects.PerPage)
}

func TestGetProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/0059df30-7701-0135-8810-6e5f001a2e3c", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("GET", r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/get.json"))
	})

	project, err := org.GetProject("0059df30-7701-0135-8810-6e5f001a2e3c")

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal("0059df30-7701-0135-8810-6e5f001a2e3c", project.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", project.OrganizationUUID)
	assert.Equal("org/test-project", project.Name)
	assert.Equal(codeship.ProjectTypePro, project.Type)
	assert.Equal("https://github.com/org/test-project", project.RepositoryURL)
	assert.Equal("github", project.RepositoryProvider)
	assert.Equal("Test User", project.AuthenticationUser)
	assert.Equal(2, len(project.NotificationRules))
	notificationRule := project.NotificationRules[0]
	assert.Equal("github", notificationRule.Notifier)
	assert.Equal("exact", notificationRule.BranchMatch)
	assert.NotEmpty(notificationRule.BuildStatuses)
	assert.Equal("all", notificationRule.Target)
	assert.NotEmpty(notificationRule.Options)
	assert.Equal("foo", notificationRule.Options.Key)
	assert.Equal("devs", notificationRule.Options.Room)
	assert.Equal("https://google.com", notificationRule.Options.URL)
	assert.Equal("ssh-rsa key", project.SSHKey)
	assert.Equal("aeskey", project.AesKey)
	createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T20:19:55.199Z")
	assert.Equal(createdAt, project.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.336Z")
	assert.Equal(updatedAt, project.UpdatedAt)
	assert.Equal(1, len(project.TeamIDs))
	assert.Equal(1007, project.TeamIDs[0])
}

func TestCreateProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("POST", r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

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
	assert.Equal("7de09100-7aeb-0135-b8e4-76a42f3a0b26", project.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", project.OrganizationUUID)
	assert.Equal(codeship.ProjectTypeBasic, project.Type)
	assert.Equal("https://github.com/org/example-repo", project.RepositoryURL)
	assert.Equal("github", project.RepositoryProvider)
	assert.Equal("Test User", project.AuthenticationUser)
	assert.NotEmpty(project.NotificationRules)
	assert.Equal(1, len(project.TeamIDs))
	assert.Equal("ssh-rsa key", project.SSHKey)
	assert.Empty(project.SetupCommands)
	assert.Empty(project.DeploymentPipelines)
	assert.Empty(project.EnvironmentVariables)
	assert.Equal(1, len(project.TestPipelines))
	assert.Equal("run tests", project.TestPipelines[0].Name)
	assert.Equal(1, len(project.TestPipelines[0].Commands))
	assert.Equal("./run-tests.sh", project.TestPipelines[0].Commands[0])
}

func TestUpdateProject(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/7de09100-7aeb-0135-b8e4-76a42f3a0b26", func(w http.ResponseWriter, r *http.Request) {
		assert := assert.New(t)
		assert.Equal("PUT", r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))
		assert.Equal("application/json", r.Header.Get("Accept"))

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
	assert.Equal("7de09100-7aeb-0135-b8e4-76a42f3a0b26", project.UUID)
	assert.Equal("28123f10-e33d-5533-b53f-111ef8d7b14f", project.OrganizationUUID)
	assert.Equal(codeship.ProjectTypePro, project.Type)
	assert.Equal("https://github.com/org/example-repo", project.RepositoryURL)
	assert.Equal("github", project.RepositoryProvider)
	assert.Equal("Test User", project.AuthenticationUser)
	assert.NotEmpty(project.NotificationRules)
	assert.Equal(2, len(project.TeamIDs))
	assert.Equal("ssh-rsa key", project.SSHKey)
	assert.Empty(project.SetupCommands)
	assert.Empty(project.DeploymentPipelines)
	assert.Empty(project.EnvironmentVariables)
	assert.Empty(project.TestPipelines)
}
