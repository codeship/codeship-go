package codeship

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Project structure for Project object
type Project struct {
	AesKey              string    `json:"aes_key,omitempty"`
	AuthenticationUser  string    `json:"authentication_user"`
	CreatedAt           time.Time `json:"created_at"`
	DeploymentPipelines []struct {
		Branch struct {
			BranchName string `json:"branch_name"`
			MatchNode  string `json:"match_node"`
		} `json:"branch"`
		Config   []string `json:"config"`
		Position int      `json:"position"`
	} `json:"deployment_pipelines,omitempty"`
	EnvironmentVariables []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environment_variables,omitempty"`
	Name              string `json:"name"`
	NotificationRules []struct {
		Branch      string `json:"branch"`
		BranchMatch string `json:"branch_match"`
		Notifier    string `json:"notifier"`
		Options     struct {
			Key  string `json:"key"`
			URL  string `json:"url"`
			Room string `json:"room"`
		} `json:"options"`
		BuildStatuses []string `json:"build_statuses"`
		EmailTarget   string   `json:"email_target"`
	} `json:"notification_rules"`
	OrganizationUUID   string   `json:"organization_uuid"`
	RepositoryProvider string   `json:"repository_provider"`
	RepositoryURL      string   `json:"repository_url"`
	SetupCommands      []string `json:"setup_commands,omitempty"`
	SSHKey             string   `json:"ssh_key"`
	TeamIds            []int    `json:"team_ids"`
	TestPipelines      []struct {
		Commands []string `json:"commands"`
		Name     string   `json:"name"`
	} `json:"test_pipelines,omitempty"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	UUID      string    `json:"uuid"`
}

// ProjectList holds a list of Project objects
type ProjectList struct {
	Projects []Project
	pagination
}

type projectResponse struct {
	Project Project
}

// ListProjects fetches a list of projects
func (o *Organization) ListProjects() (ProjectList, error) {
	path := fmt.Sprintf("/organizations/%s/projects", o.UUID)

	projectList := ProjectList{}
	resp, err := o.request("GET", path, nil)
	if err != nil {
		return projectList, errors.Wrap(err, "unable to list projects")
	}

	err = json.Unmarshal(resp, &projectList)
	if err != nil {
		return projectList, errors.Wrap(err, "unable to unmarshal response into ProjectList")
	}

	return projectList, nil
}

// GetProject fetches a project by ID
func (o *Organization) GetProject(projectID string) (Project, error) {
	path := fmt.Sprintf("/organizations/%s/projects/%s", o.UUID, projectID)

	project := projectResponse{}
	resp, err := o.request("GET", path, nil)
	if err != nil {
		return project.Project, errors.Wrap(err, "unable to get project")
	}

	err = json.Unmarshal(resp, &project)
	if err != nil {
		return project.Project, errors.Wrap(err, "unable to unmarshal response into Project")
	}

	return project.Project, nil
}

// CreateProject creates a new project
func (o *Organization) CreateProject(project Project) (Project, error) {
	path := fmt.Sprintf("/organizations/%s/projects", o.UUID)

	resp, err := o.request("POST", path, project)
	if err != nil {
		return project, errors.Wrap(err, "unable to create project")
	}

	projResponse := projectResponse{}
	err = json.Unmarshal(resp, &projResponse)
	if err != nil {
		return project, errors.Wrap(err, "unable to unmarshal response into Project")
	}

	return projResponse.Project, nil
}
