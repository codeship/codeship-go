package codeship

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// BuildLinks structure of BuildLinks object for a Build
type BuildLinks struct {
	Pipelines string `json:"pipelines,omitempty"`
	Services  string `json:"services,omitempty"`
	Steps     string `json:"steps,omitempty"`
}

// Build structure of Build object
type Build struct {
	AllocatedAt      time.Time  `json:"allocated_at,omitempty"`
	CommitMessage    string     `json:"commit_message,omitempty"`
	CommitSha        string     `json:"commit_sha,omitempty"`
	FinishedAt       time.Time  `json:"finished_at,omitempty"`
	Links            BuildLinks `json:"links,omitempty"`
	OrganizationUUID string     `json:"organization_uuid,omitempty"`
	ProjectUUID      string     `json:"project_uuid,omitempty"`
	QueuedAt         time.Time  `json:"queued_at,omitempty"`
	Ref              string     `json:"ref,omitempty"`
	Status           string     `json:"status,omitempty"`
	Username         string     `json:"username,omitempty"`
	UUID             string     `json:"uuid,omitempty"`
}

// BuildList holds a list of Build objects
type BuildList struct {
	Builds []Build `json:"builds"`
	pagination
}

type buildResponse struct {
	Build Build `json:"build"`
}

// BuildPipelineMetrics structure of BuildPipelineMetrics object for a BuildPipeline
type BuildPipelineMetrics struct {
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
}

// BuildPipeline structure of BuildPipeline object for a Basic Project
type BuildPipeline struct {
	UUID       string               `json:"uuid,omitempty"`
	BuildUUID  string               `json:"build_uuid,omitempty"`
	Type       string               `json:"type,omitempty"`
	Status     string               `json:"status,omitempty"`
	CreatedAt  time.Time            `json:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty"`
	FinishedAt time.Time            `json:"finished_at,omitempty"`
	Metrics    BuildPipelineMetrics `json:"metrics,omitempty"`
}

// BuildPipelines holds a list of BuildPipeline objects for a Basic Project
type BuildPipelines struct {
	Pipelines []BuildPipeline `json:"pipelines"`
	pagination
}

// BuildStep structure of BuildStep object for a Pro Project
type BuildStep struct {
	BuildUUID   string      `json:"build_uuid,omitempty"`
	BuildingAt  time.Time   `json:"building_at,omitempty"`
	Command     string      `json:"command,omitempty"`
	FinishedAt  time.Time   `json:"finished_at,omitempty"`
	ImageName   string      `json:"image_name,omitempty"`
	Name        string      `json:"name,omitempty"`
	Registry    string      `json:"registry,omitempty"`
	ServiceUUID string      `json:"service_uuid,omitempty"`
	StartedAt   time.Time   `json:"started_at,omitempty"`
	Status      string      `json:"status,omitempty"`
	Steps       []BuildStep `json:"steps,omitempty"`
	Tag         string      `json:"tag,omitempty"`
	Type        string      `json:"type,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
}

// BuildSteps holds a list of BuildStep objects for a Pro Project
type BuildSteps struct {
	Steps []BuildStep `json:"steps"`
	pagination
}

// BuildService structure of BuildService object for a Pro Project
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

// BuildServices holds a list of BuildService objects for a Pro Project
type BuildServices struct {
	Services []BuildService `json:"services"`
	pagination
}

type buildRequest struct {
	CommitSha string `json:"commit_sha,omitempty"`
	Ref       string `json:"ref,omitempty"`
}

// CreateBuild creates a new build
func (o *Organization) CreateBuild(ctx context.Context, projectUUID, ref, commitSha string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", o.UUID, projectUUID)

	_, err := o.client.request(ctx, "POST", path, buildRequest{
		Ref:       ref,
		CommitSha: commitSha,
	})
	if err != nil {
		return false, errors.Wrap(err, "unable to create build")
	}

	return true, nil
}

// GetBuild fetches a build by UUID
func (o *Organization) GetBuild(ctx context.Context, projectUUID, buildUUID string) (Build, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s", o.UUID, projectUUID, buildUUID)

	resp, err := o.client.request(ctx, "GET", path, nil)
	if err != nil {
		return Build{}, errors.Wrap(err, "unable to get build")
	}

	var build buildResponse
	if err = json.Unmarshal(resp, &build); err != nil {
		return Build{}, errors.Wrap(err, "unable to unmarshal response into Build")
	}

	return build.Build, nil
}

// ListBuilds fetches a list of builds for the given organization
func (o *Organization) ListBuilds(ctx context.Context, projectUUID string) (BuildList, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds", o.UUID, projectUUID)

	resp, err := o.client.request(ctx, "GET", path, nil)
	if err != nil {
		return BuildList{}, errors.Wrap(err, "unable to list builds")
	}

	var builds BuildList
	if err = json.Unmarshal(resp, &builds); err != nil {
		return BuildList{}, errors.Wrap(err, "unable to unmarshal response into BuildList")
	}

	return builds, nil
}

// GetBuildPipelines gets Basic build pipelines
func (o *Organization) GetBuildPipelines(ctx context.Context, projectUUID, buildUUID string) (BuildPipelines, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/pipelines", o.UUID, projectUUID, buildUUID)

	resp, err := o.client.request(ctx, "GET", path, nil)
	if err != nil {
		return BuildPipelines{}, errors.Wrap(err, "unable to get build pipelines")
	}

	var pipelines BuildPipelines
	if err = json.Unmarshal(resp, &pipelines); err != nil {
		return BuildPipelines{}, errors.Wrap(err, "unable to unmarshal response into BuildPipelines")
	}

	return pipelines, nil
}

// StopBuild stops a running build
func (o *Organization) StopBuild(ctx context.Context, projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/stop", o.UUID, projectUUID, buildUUID)

	if _, err := o.client.request(ctx, "POST", path, nil); err != nil {
		return false, errors.Wrap(err, "unable to stop build")
	}

	return true, nil
}

// RestartBuild restarts a previous build
func (o *Organization) RestartBuild(ctx context.Context, projectUUID, buildUUID string) (bool, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/restart", o.UUID, projectUUID, buildUUID)

	if _, err := o.client.request(ctx, "POST", path, nil); err != nil {
		return false, errors.Wrap(err, "unable to restart build")
	}

	return true, nil
}

// GetBuildServices gets Pro build services
func (o *Organization) GetBuildServices(ctx context.Context, projectUUID, buildUUID string) (BuildServices, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/services", o.UUID, projectUUID, buildUUID)

	resp, err := o.client.request(ctx, "GET", path, nil)
	if err != nil {
		return BuildServices{}, errors.Wrap(err, "unable to get build services")
	}

	var services BuildServices
	if err = json.Unmarshal(resp, &services); err != nil {
		return BuildServices{}, errors.Wrap(err, "unable to unmarshal response into BuildServices")
	}

	return services, nil
}

// GetBuildSteps gets Pro build steps
func (o *Organization) GetBuildSteps(ctx context.Context, projectUUID, buildUUID string) (BuildSteps, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s/builds/%s/steps", o.UUID, projectUUID, buildUUID)

	resp, err := o.client.request(ctx, "GET", path, nil)
	if err != nil {
		return BuildSteps{}, errors.Wrap(err, "unable to get build steps")
	}

	var steps BuildSteps
	if err = json.Unmarshal(resp, &steps); err != nil {
		return BuildSteps{}, errors.Wrap(err, "unable to unmarshal response into BuildSteps")
	}

	return steps, nil
}
