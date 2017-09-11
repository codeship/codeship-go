package codeship

import (
	"encoding/json"
	"os"
	"testing"
)

func TestListProjects(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	orgName := os.Getenv("CODESHIP_ORGNAME")
	apiClient, err := New(username, password, orgName)
	if err != nil {
		t.Error("New returned error:", err)
	}

	orgID := apiClient.Authentication.Organizations[0].UUID

	projectList, err := apiClient.ListProjects()
	if err != nil {
		t.Errorf("Unable to list projects. Org ID: %s, Error: %s", orgID, err)
	}

	if len(projectList.Projects) == 0 {
		t.Error("Zero projects returned")
	}

}

func TestGetProject(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	orgName := os.Getenv("CODESHIP_ORGNAME")
	apiClient, err := New(username, password, orgName)
	if err != nil {
		t.Error("New returned error:", err)
	}

	projectList, err := apiClient.ListProjects()
	if err != nil {
		t.Error("Unable to list projects:", err)
	}

	if len(projectList.Projects) == 0 {
		t.Error("Zero projects returned")
	}

	projectID := projectList.Projects[0].UUID

	project, err := apiClient.GetProject(projectID)
	if err != nil {
		t.Errorf("Unable to get project %s, error: %s ", projectID, err)
		t.Fail()
	}

	if project.UUID != projectID {
		t.Errorf("The returned project's UUID (%s) does not match expected (%s)", project.UUID, projectID)
		t.Fail()
	}

}

func TestCreateProject(t *testing.T) {
	t.SkipNow()
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	orgName := os.Getenv("CODESHIP_ORGNAME")
	apiClient, err := New(username, password, orgName)
	if err != nil {
		t.Error("New returned error:", err)
	}

	createProjectFixtures := getCreateProjectFixtures()

	for _, projectFixture := range createProjectFixtures {
		project, err := apiClient.CreateProject(projectFixture.Project)
		if err != nil {
			t.Errorf("Unable to create project, error: %s", err)
			t.Fail()
		}

		if project.UUID == "" {
			projectJSON, _ := json.Marshal(project)
			t.Errorf("Project not created properly, missing UUID. object json: %s", projectJSON)
			t.Fail()
		}
	}

}

// func TestListProjectsForOrg(t *testing.T) {
// 	testSetup()
// 	username := os.Getenv("CODESHIP_USERNAME")
// 	password := os.Getenv("CODESHIP_PASSWORD")
//  orgName := os.Getenv("CODESHIP_ORGNAME")
// 	apiClient, err := New(username, password, orgName)
// 	if err != nil {
// 		t.Error("New returned error:", err)
// 	}
//
// 	projects, _ := apiClient.ListProjects()
// 	projectStr := ""
// 	for _, p := range projects.Projects {
// 		projectStr += p.Name + "=" + p.UUID + ", "
// 	}
// 	t.Errorf("%s", projectStr)
// }
