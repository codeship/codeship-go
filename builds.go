package codeship

import (
	"encoding/json"
	"fmt"
	"time"

	errors "github.com/pkg/errors"
)

// Build structure of Build object
type Build struct {
	AllocatedAt   time.Time `json:"allocated_at"`
	CommitMessage string    `json:"commit_message"`
	CommitSha     string    `json:"commit_sha"`
	FinishedAt    time.Time `json:"finished_at"`
	Links         struct {
		Pipelines string `json:"pipelines"`
		Services  string `json:"services"`
		Steps     string `json:"steps"`
	} `json:"links"`
	OrganizationUUID string    `json:"organization_uuid"`
	ProjectUUID      string    `json:"project_uuid"`
	QueuedAt         time.Time `json:"queued_at"`
	Ref              string    `json:"ref"`
	Status           string    `json:"status"`
	Username         string    `json:"username"`
	UUID             string    `json:"uuid"`
}

// BuildList holds a list of Build objects
type BuildList struct {
	Builds []Build
}

type buildResponse struct {
	Build Build
}

// BuildPipelines holds a list of Pipeline objects
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

// BuildStep structure of BuildStep object
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

// BuildSteps holds a list of BuildStep objects
type BuildSteps struct {
	Steps []BuildStep
}

// BuildService structure of BuildService object
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

// BuildServices holds a list of BuildService objects
type BuildServices struct {
	Services []BuildService
}

type buildRequest struct {
	CommitSha string `json:"commit_sha"`
	Ref       string `json:"ref"`
}

// CreateBuild Create a new build
func (api *API) CreateBuild(orgUUID, projectUUID, ref, commitSha string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", orgUUID, projectUUID)

	buildReq := buildRequest{
		Ref:       ref,
		CommitSha: commitSha,
	}

	_, err := api.makeRequest("POST", path, buildReq)
	if err != nil {
		return false, errors.Wrap(err, "Unable to create build")
	}

	return true, nil
}

// GetBuild Fetch a build
func (api *API) GetBuild(orgUUID, projectUUID, buildUUID string) (Build, error) {
	build := Build{}
	orgUUID = api.getOrgUUID(orgUUID)
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s", orgUUID, projectUUID, buildUUID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return build, errors.Wrap(err, "Unable to get build")
	}

	buildResp := buildResponse{}
	err = json.Unmarshal(resp, &buildResp)
	if err != nil {
		return build, errors.Wrap(err, "Unable to unmarshal JSON into Build")
	}

	return buildResp.Build, nil
}

// ListBuilds Fetch a list of builds for the given organization
func (api *API) ListBuilds(orgUUID string, projectUUID string) (BuildList, error) {
	buildList := BuildList{}
	orgUUID = api.getOrgUUID(orgUUID)
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", orgUUID, projectUUID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildList, errors.Wrap(err, "Unable to list builds")
	}

	err = json.Unmarshal(resp, &buildList)
	if err != nil {
		return buildList, errors.Wrap(err, "Unable to unmarshal JSON into BuildList")
	}

	return buildList, nil
}

// GetBuildPipelines Basic projects only
func (api *API) GetBuildPipelines(orgUUID, projectUUID, buildUUID string) (BuildPipelines, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/pipelines", orgUUID, projectUUID, buildUUID)

	buildPipelines := BuildPipelines{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildPipelines, errors.Wrap(err, "Unable to get build pipelines")
	}

	err = json.Unmarshal(resp, &buildPipelines)
	if err != nil {
		return buildPipelines, errors.Wrap(err, "Unable to unmarshal JSON into BuildPipelines")
	}

	return buildPipelines, nil
}

// StopBuild Stop a running build
func (api *API) StopBuild(orgUUID, projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/stop", orgUUID, projectUUID, buildUUID)

	_, err := api.makeRequest("POST", path, nil)
	if err != nil {
		return false, errors.Wrap(err, "Unable to stop build")
	}

	return true, nil
}

// RestartBuild Restart a previous build
func (api *API) RestartBuild(orgUUID, projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/restart", orgUUID, projectUUID, buildUUID)

	_, err := api.makeRequest("POST", path, nil)
	if err != nil {
		return false, errors.Wrap(err, "Unable to restart build, error")
	}

	return true, nil
}

// GetBuildServices Pro projects only
func (api *API) GetBuildServices(orgUUID, projectUUID, buildUUID string) (BuildServices, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/services", orgUUID, projectUUID, buildUUID)

	buildServices := BuildServices{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildServices, errors.Wrap(err, "Unable to get build services")
	}

	err = json.Unmarshal(resp, &buildServices)
	if err != nil {
		return buildServices, errors.Wrap(err, "Unable to unmarshal JSON into BuildServices")
	}

	return buildServices, nil
}

// GetBuildSteps Pro projects only
func (api *API) GetBuildSteps(orgUUID, projectUUID, buildUUID string) (BuildSteps, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/steps", orgUUID, projectUUID, buildUUID)

	buildSteps := BuildSteps{}
	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return buildSteps, errors.Wrap(err, "Unable to get build steps")
	}

	err = json.Unmarshal(resp, &buildSteps)
	if err != nil {
		return buildSteps, errors.Wrap(err, "Unable to unmarshal JSON into BuildSteps")
	}

	return buildSteps, nil
}
