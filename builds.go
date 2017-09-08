package codeship

import (
	"encoding/json"
	"fmt"
	"time"
)

type Build struct {
	AllocatedAt   time.Time `json:"allocated_at"`
	CommitMessage string    `json:"commit_message"`
	CommitSha     string    `json:"commit_sha"`
	FinishedAt    time.Time `json:"finished_at"`
	Links         struct {
		Pipelines []string `json:"pipelines"`
		Services  []string `json:"services"`
		Steps     []string `json:"steps"`
	} `json:"links"`
	OrganizationUUID string    `json:"organization_uuid"`
	ProjectUUID      string    `json:"project_uuid"`
	QueuedAt         time.Time `json:"queued_at"`
	Ref              string    `json:"ref"`
	Status           string    `json:"status"`
	Username         string    `json:"username"`
	UUID             string    `json:"uuid"`
}

type BuildList struct {
	Builds []Build
}

type buildResponse struct {
	Build Build
}

type BuildPipelines struct {
	Pipelines []struct {
		UUID       string    `json:"uuid"`
		BuildUUID  string    `json:"build_uuid"`
		Type       string    `json:"type"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		FinishedAt time.Time `json:"finished_at"`
		Metrics    struct {
			AmiID                 string `json:"ami_id,omitempty"`
			Queries               string `json:"queries,omitempty"`
			CPUUser               string `json:"cpu_user,omitempty"`
			Duration              string `json:"duration,omitempty"`
			CPUSystem             string `json:"cpu_system,omitempty"`
			InstanceID            string `json:"instance_id,omitempty"`
			Architecture          string `json:"architecture,omitempty"`
			InstanceType          string `json:"instance_type,omitempty"`
			CPUPerSecond          string `json:"cpu_per_second,omitempty"`
			DiskFreeBytes         string `json:"disk_free_bytes,omitempty"`
			DiskUsedBytes         string `json:"disk_used_bytes,omitempty"`
			NetworkRxBytes        string `json:"network_rx_bytes,omitempty"`
			NetworkTxBytes        string `json:"network_tx_bytes,omitempty"`
			MaxUsedConnections    string `json:"max_used_connections,omitempty"`
			MemoryMaxUsageInBytes string `json:"memory_max_usage_in_bytes,omitempty"`
		} `json:"metrics,omitempty"`
	} `json:"pipelines"`
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}

type BuildStep struct {
	BuildUUID   string    `json:"build_uuid,omitempty"`
	BuildingAt  time.Time `json:"building_at,omitempty"`
	Command     string    `json:"command,omitempty"`
	FinishedAt  time.Time `json:"finished_at,omitempty"`
	ImageName   string    `json:"image_name,omitempty"`
	Name        string    `json:"name,omitempty"`
	Registry    string    `json:"registry,omitempty"`
	ServiceUUID string    `json:"service_uuid,omitempty"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	Status      string    `json:"status,omitempty"`
	Tag         string    `json:"tag,omitempty"`
	Type        string    `json:"type,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	UUID        string    `json:"uuid,omitempty"`
}

type BuildSteps struct {
	Steps []BuildStep
}

type BuildService struct {
	BuildUUID  string    `json:"build_uuid,omitempty"`
	BuildingAt time.Time `json:"building_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	Name       string    `json:"name,omitempty"`
	PullingAt  time.Time `json:"pulling_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	UUID       string    `json:"uuid,omitempty"`
	Status     string    `json:"status,omitempty"`
}

type BuildServices struct {
	Services []BuildService
}

type BuildRequest struct {
	CommitSha string `json:"commit_sha"`
	Ref       string `json:"ref"`
}

// CreateBuild Create a new build
func (api *API) CreateBuild(build Build) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", build.OrganizationUUID, build.ProjectUUID)

	buildRequest := BuildRequest{
		Ref:       build.Ref,
		CommitSha: build.CommitSha,
	}

	_, err := api.makeRequest("POST", path, buildRequest)
	if err != nil {
		return false, fmt.Errorf("Unable to create build: %s", err)
	}

	return true, nil
}

// GetBuild Fetch a build
func (api *API) GetBuild(orgID, projectID, buildID string) (Build, error) {
	build := Build{}
	orgID = api.getOrgUUID(orgID)
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s", orgID, projectID, buildID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return build, fmt.Errorf("Unable to get build: %s", err)
	}

	buildResp := buildResponse{}
	json.Unmarshal(resp, &buildResp)

	return buildResp.Build, nil
}

// ListBuilds Fetch a list of builds for the given organization
func (api *API) ListBuilds(orgID string, projectID string) (BuildList, error) {
	buildList := BuildList{}
	orgID = api.getOrgUUID(orgID)
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", orgID, projectID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildList, fmt.Errorf("Unable to list builds: %s", err)
	}

	json.Unmarshal(resp, &buildList)

	return buildList, nil
}

// GetBuildPipelines Basic projects only
func (api *API) GetBuildPipelines(build Build) (BuildPipelines, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/pipelines", build.OrganizationUUID, build.ProjectUUID, build.UUID)

	buildPipelines := BuildPipelines{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildPipelines, fmt.Errorf("Unable to get build pipelines: %s", err)
	}

	json.Unmarshal(resp, &buildPipelines)

	return buildPipelines, nil
}

// StopBuild Stop a running build
func (api *API) StopBuild(build Build) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/stop", build.OrganizationUUID, build.ProjectUUID, build.UUID)

	_, err := api.makeRequest("POST", path, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to stop build, error: %s", err)
	}

	return true, nil
}

// RestartBuild Restart a previous build
func (api *API) RestartBuild(build Build) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/restart", build.OrganizationUUID, build.ProjectUUID, build.UUID)

	_, err := api.makeRequest("POST", path, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to restart build, error: %s", err)
	}

	return true, nil
}

// GetBuildServices Pro projects only
func (api *API) GetBuildServices(build Build) (BuildServices, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/services", build.OrganizationUUID, build.ProjectUUID, build.UUID)

	buildServices := BuildServices{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildServices, fmt.Errorf("Unable to get build services: %s", err)
	}

	json.Unmarshal(resp, &buildServices)

	return buildServices, nil
}

// GetBuildSteps Pro projects only
func (api *API) GetBuildSteps(build Build) (BuildSteps, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/steps", build.OrganizationUUID, build.ProjectUUID, build.UUID)

	buildSteps := BuildSteps{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildSteps, fmt.Errorf("Unable to get build steps: %s", err)
	}

	json.Unmarshal(resp, &buildSteps)

	return buildSteps, nil
}
