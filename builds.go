package codeship

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
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
	Builds []Build `json:"builds"`
	Pagination
}

type buildResponse struct {
	Build Build
}

// BuildPipeline structure of BuildPipeline object
type BuildPipeline struct {
	UUID       string    `json:"uuid"`
	BuildUUID  string    `json:"build_uuid"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at"`
	Metrics    struct {
		AmiID                 string `json:"ami_id"`
		Queries               string `json:"queries"`
		CPUUser               string `json:"cpu_user"`
		Duration              string `json:"duration"`
		CPUSystem             string `json:"cpu_system"`
		InstanceID            string `json:"instance_id"`
		Architecture          string `json:"architecture"`
		InstanceType          string `json:"instance_type"`
		CPUPerSecond          string `json:"cpu_per_second"`
		DiskFreeBytes         string `json:"disk_free_bytes"`
		DiskUsedBytes         string `json:"disk_used_bytes"`
		NetworkRxBytes        string `json:"network_rx_bytes"`
		NetworkTxBytes        string `json:"network_tx_bytes"`
		MaxUsedConnections    string `json:"max_used_connections"`
		MemoryMaxUsageInBytes string `json:"memory_max_usage_in_bytes"`
	} `json:"metrics"`
}

// BuildPipelines holds a list of BuildPipeline objects
type BuildPipelines struct {
	Pipelines []BuildPipeline `json:"pipelines"`
	Pagination
}

// BuildStep structure of BuildStep object
type BuildStep struct {
	BuildUUID   string      `json:"build_uuid"`
	BuildingAt  time.Time   `json:"building_at"`
	Command     string      `json:"command"`
	FinishedAt  time.Time   `json:"finished_at"`
	ImageName   string      `json:"image_name"`
	Name        string      `json:"name"`
	Registry    string      `json:"registry"`
	ServiceUUID string      `json:"service_uuid"`
	StartedAt   time.Time   `json:"started_at"`
	Status      string      `json:"status"`
	Steps       []BuildStep `json:"steps"`
	Tag         string      `json:"tag"`
	Type        string      `json:"type"`
	UpdatedAt   time.Time   `json:"updated_at"`
	UUID        string      `json:"uuid"`
}

// BuildSteps holds a list of BuildStep objects
type BuildSteps struct {
	Steps []BuildStep `json:"steps"`
	Pagination
}

// BuildService structure of BuildService object
type BuildService struct {
	BuildUUID  string    `json:"build_uuid"`
	BuildingAt time.Time `json:"building_at"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Name       string    `json:"name"`
	PullingAt  time.Time `json:"pulling_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UUID       string    `json:"uuid"`
	Status     string    `json:"status"`
}

// BuildServices holds a list of BuildService objects
type BuildServices struct {
	Services []BuildService `json:"services"`
	Pagination
}

type buildRequest struct {
	CommitSha string `json:"commit_sha"`
	Ref       string `json:"ref"`
}

// CreateBuild creates a new build
func (o *Organization) CreateBuild(projectUUID, ref, commitSha string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", o.UUID, projectUUID)

	buildReq := buildRequest{
		Ref:       ref,
		CommitSha: commitSha,
	}

	_, err := o.makeRequest("POST", path, buildReq)
	if err != nil {
		return false, errors.Wrap(err, "unable to create build")
	}

	return true, nil
}

// GetBuild fetches a build
func (o *Organization) GetBuild(projectUUID, buildUUID string) (Build, error) {
	build := Build{}
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s", o.UUID, projectUUID, buildUUID)

	resp, err := o.makeRequest("GET", path, nil)
	if err != nil {
		return build, errors.Wrap(err, "unable to get build")
	}

	buildResp := buildResponse{}
	err = json.Unmarshal(resp, &buildResp)
	if err != nil {
		return build, errors.Wrap(err, "unable to unmarshal response into Build")
	}

	return buildResp.Build, nil
}

// ListBuilds fetches a list of builds for the given organization
func (o *Organization) ListBuilds(projectUUID string) (BuildList, error) {
	buildList := BuildList{}
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", o.UUID, projectUUID)

	resp, err := o.makeRequest("GET", path, nil)
	if err != nil {
		return buildList, errors.Wrap(err, "unable to list builds")
	}

	err = json.Unmarshal(resp, &buildList)
	if err != nil {
		return buildList, errors.Wrap(err, "unable to unmarshal response into BuildList")
	}

	return buildList, nil
}

// GetBuildPipelines gets Basic build pipelines
func (o *Organization) GetBuildPipelines(projectUUID, buildUUID string) (BuildPipelines, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/pipelines", o.UUID, projectUUID, buildUUID)

	buildPipelines := BuildPipelines{}
	resp, err := o.makeRequest("GET", path, nil)
	if err != nil {
		return buildPipelines, errors.Wrap(err, "unable to get build pipelines")
	}

	err = json.Unmarshal(resp, &buildPipelines)
	if err != nil {
		return buildPipelines, errors.Wrap(err, "unable to unmarshal response into BuildPipelines")
	}

	return buildPipelines, nil
}

// StopBuild stops a running build
func (o *Organization) StopBuild(projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/stop", o.UUID, projectUUID, buildUUID)

	_, err := o.makeRequest("POST", path, nil)
	if err != nil {
		return false, errors.Wrap(err, "unable to stop build")
	}

	return true, nil
}

// RestartBuild restarts a previous build
func (o *Organization) RestartBuild(projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/restart", o.UUID, projectUUID, buildUUID)

	_, err := o.makeRequest("POST", path, nil)
	if err != nil {
		return false, errors.Wrap(err, "unable to restart build")
	}

	return true, nil
}

// GetBuildServices gets Pro build services
func (o *Organization) GetBuildServices(projectUUID, buildUUID string) (BuildServices, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/services", o.UUID, projectUUID, buildUUID)

	buildServices := BuildServices{}
	resp, err := o.makeRequest("GET", path, nil)
	if err != nil {
		return buildServices, errors.Wrap(err, "unable to get build services")
	}

	err = json.Unmarshal(resp, &buildServices)
	if err != nil {
		return buildServices, errors.Wrap(err, "unable to unmarshal response into BuildServices")
	}

	return buildServices, nil
}

// GetBuildSteps gets Pro build steps
func (o *Organization) GetBuildSteps(projectUUID, buildUUID string) (BuildSteps, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/steps", o.UUID, projectUUID, buildUUID)

	buildSteps := BuildSteps{}
	resp, err := o.makeRequest("GET", path, nil)
	if err != nil {
		return buildSteps, errors.Wrap(err, "unable to get build steps")
	}

	err = json.Unmarshal(resp, &buildSteps)
	if err != nil {
		return buildSteps, errors.Wrap(err, "unable to unmarshal response into BuildSteps")
	}

	return buildSteps, nil
}
