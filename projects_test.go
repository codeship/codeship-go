package codeship_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

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
	assert.Equal("basic", project.Type)
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
	assert.Equal(2, len(project.TeamIds))
	assert.Equal(1007, project.TeamIds[0])
	assert.Equal(1009, project.TeamIds[1])
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
	assert.Equal("pro", project.Type)
	assert.Equal("https://github.com/org/test-project", project.RepositoryURL)
	assert.Equal("github", project.RepositoryProvider)
	assert.Equal("Test User", project.AuthenticationUser)
	assert.Equal(2, len(project.NotificationRules))
	notificationRule := project.NotificationRules[0]
	assert.Equal("github", notificationRule.Notifier)
	assert.Equal("exact", notificationRule.BranchMatch)
	assert.NotEmpty(notificationRule.BuildStatuses)
	assert.Equal("all", notificationRule.EmailTarget)
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
	assert.Equal(1, len(project.TeamIds))
	assert.Equal(1007, project.TeamIds[0])
}

// //
// func TestCreateProject(t *testing.T) {
// 	testSetup()

// 	mux.HandleFunc("/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects", func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(200)
// 		fmt.Fprintf(w, `{
//     "project": {
//         "uuid": "7de09100-7aeb-0135-b8e4-76a42f3a0b26",
//         "name": "org/example-repo",
//         "type": "basic",
//         "repository_url": "https://github.com/org/example-repo",
//         "repository_provider": "github",
//         "authentication_user": "Test User",
//         "organization_uuid": "28123f10-e33d-5533-b53f-111ef8d7b14f",
//         "notification_rules": [
//             {
//                 "notifier": "github",
//                 "branch": null,
//                 "branch_match": "exact",
//                 "build_statuses": [
//                     "failed",
//                     "started",
//                     "recovered",
//                     "success"
//                 ],
//                 "build_owner": "all",
//                 "options": {}
//             },
//             {
//                 "notifier": "email",
//                 "branch": null,
//                 "branch_match": "exact",
//                 "build_statuses": [
//                     "failed",
//                     "recovered"
//                 ],
//                 "build_owner": "all",
//                 "options": {}
//             }
//         ],
//         "ssh_key": "ssh-rsa key",
//         "aes_key": null,
//         "created_at": "2017-09-13T19:56:01.398Z",
//         "updated_at": "2017-09-13T19:56:02.804Z",
//         "team_ids": [
//             61593
//         ],
//         "setup_commands": [],
//         "deployment_pipelines": [],
//         "environment_variables": [],
//         "test_pipelines": [
//             {
//                 "name": "run tests",
//                 "commands": []
//             }
//         ]
//     }
// }`)
// 	})

// 	createProject := Project{
// 		RepositoryURL: "git@github.com/org/repo-name",
// 		TestPipelines: []struct {
// 			Commands []string `json:"commands,omitempty"`
// 			Name     string   `json:"name,omitempty"`
// 		}{
// 			{
// 				Commands: []string{"./run-tests.sh"},
// 				Name:     "test pass",
// 			},
// 		},
// 		Type: TypeBasic,
// 	}

// 	project, err := client.CreateProject(createProject)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "7de09100-7aeb-0135-b8e4-76a42f3a0b26", project.UUID)

// 	testTeardown()
// }
