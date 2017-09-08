package codeship

import (
	"os"
	"testing"
)

func TestCreateBuild(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Error("New returned error:", err)
		t.FailNow()
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["create"]

	ok, err := apiClient.CreateBuild(fixture.Build)
	if !ok || err != nil {
		t.Errorf("Unable to create new build. Org ID: %s, Project ID: %s, From Build ID: %s, Ref: %s, Commit: %s, Error: %s", fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID, fixture.Build.UUID, fixture.Build.Ref, fixture.Build.CommitSha, err)
		t.FailNow()
	}
}

func TestGetBuild(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Error("New returned error:", err)
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["restart"]

	build, err := apiClient.GetBuild(fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID, fixture.Build.UUID)
	if err != nil {
		t.Errorf("Unable to GetBuild. Org ID: %s, Project ID: %s, Build ID: %s, Error: %s", fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID, fixture.Build.UUID, err)
	}

	if build.UUID != fixture.Build.UUID {
		t.Errorf("Build returned from GetBuild (%s) does not match expected (%s). ", build.UUID, fixture.Build.UUID)
	}
}

func TestListBuilds(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Error("New returned error:", err)
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["restart"]

	buildList, err := apiClient.ListBuilds(fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID)
	if err != nil {
		t.Errorf("Unable to list builds. Org ID: %s, Project ID: %s, Error: %s", fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID, err)
	}

	if len(buildList.Builds) == 0 {
		t.Error("Zero builds returned")
	}

}

func TestGetBuildPipelines(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Error("New returned error:", err)
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["restart"]

	buildPipelines, err := apiClient.GetBuildPipelines(fixture.Build)
	if err != nil {
		t.Errorf("Unable to get build pipelines. Org ID: %s, Project ID: %s, Error: %s", fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID, err)
	}

	if len(buildPipelines.Pipelines) == 0 {
		t.Error("Zero pipelines returned")
	}
}

func TestStopBuild(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Errorf("New returned error: %s", err)
		t.FailNow()
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["restart"]

	_, err = apiClient.RestartBuild(fixture.Build)
	if err != nil {
		t.Errorf("Unable to restart build: %s", err)
		t.FailNow()
	}

	_, err = apiClient.StopBuild(fixture.Build)
	if err != nil {
		t.Errorf("Unable to stop build: %s", err)
		t.FailNow()
	}
}

func TestRestartBuild(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Errorf("New returned error: %s", err)
		t.FailNow()
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["restart"]

	_, err = apiClient.RestartBuild(fixture.Build)
	if err != nil {
		t.Errorf("Unable to restart build: %s", err)
		t.FailNow()
	}
}

func TestGetBuildServices(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	apiClient, err := New(username, password, "")
	if err != nil {
		t.Errorf("New returned error: %s", err)
		t.FailNow()
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["buildservices"]

	buildServices, err := apiClient.GetBuildServices(fixture.Build)
	if err != nil {
		t.Errorf("Unable to get build services: %s", err)
		t.FailNow()
	}

	if len(buildServices.Services) == 0 {
		t.Error("Build did not have any services")
	}
}

func TestGetBuildSteps(t *testing.T) {
	testSetup()
	apiClient, err := New("", "", "")
	if err != nil {
		t.Errorf("New returned error: %s", err)
		t.FailNow()
	}

	buildTestFixtures := getBuildFixtures()
	fixture := buildTestFixtures["buildservices"]

	buildSteps, err := apiClient.GetBuildSteps(fixture.Build)
	if err != nil {
		t.Errorf("Unable to get build steps: %s", err)
		t.FailNow()
	}

	if len(buildSteps.Steps) == 0 {
		t.Error("Build did not have any steps")
	}
}

// func TestListBuildsForProject(t *testing.T) {
// 	testSetup()
// 	username := os.Getenv("CODESHIP_USERNAME")
// 	password := os.Getenv("CODESHIP_PASSWORD")
// 	apiClient, err := New(username, password, "")
// 	if err != nil {
// 		t.Error("New returned error:", err)
// 	}
//
// 	buildTestFixtures := getBuildFixtures()
// 	fixture := buildTestFixtures["restart"]
//
// 	build, _ := apiClient.ListBuilds(fixture.Build.OrganizationUUID, fixture.Build.ProjectUUID)
// 	buildStr := ""
// 	for _, b := range build.Builds {
// 		buildStr += b.UUID + " , "
// 	}
// 	t.Errorf("%s", buildStr)
// }
