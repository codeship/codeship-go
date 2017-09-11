package codeship

import (
	"encoding/json"
	"fmt"
	"time"

	errors "github.com/pkg/errors"
)

// TypePro constant for Pro project type value
const TypePro = "pro"

// TypeBasic constant for Basic project type value
const TypeBasic = "basic"

// Project structure for Project object
type Project struct {
	AesKey              string    `json:"aes_key"`
	AuthenticationUser  string    `json:"authentication_user"`
	CreatedAt           time.Time `json:"created_at"`
	DeploymentPipelines []struct {
		Branch struct {
			BranchName string `json:"branch_name"`
			MatchNode  string `json:"match_node"`
		} `json:"branch"`
		Config   []string `json:"config"`
		Position int      `json:"position,omitempty"`
	} `json:"deployment_pipelines"`
	EnvironmentVariables []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environment_variables"`
	Name              string `json:"name"`
	NotificationRules []struct {
		Branch      string `json:"branch"`
		BranchMatch string `json:"branch_match"`
		Notifier    string `json:"notifier"`
		Options     struct {
			Campfire struct {
				Room string `json:"room"`
			} `json:"campfire"`
			FlowdockKey string `json:"flowdock_key"`
			Hipchat     struct {
				Key string `json:"key"`
			} `json:"hipchat"`
			WebhookURL string `json:"webhook_url"`
		} `json:"options"`
		BuildOwner    string   `json:"build_owner,omitempty"`
		BuildStatuses []string `json:"build_statuses,omitempty"`
		EmailTarget   string   `json:"email_target,omitempty"`
	} `json:"notification_rules"`
	RepositoryProvider string   `json:"repository_provider"`
	RepositoryURL      string   `json:"repository_url"`
	SetupCommands      []string `json:"setup_commands"`
	SSHKey             string   `json:"ssh_key"`
	TeamIds            []int    `json:"team_ids"`
	TestPipelines      []struct {
		Commands []string `json:"commands,omitempty"`
		Name     string   `json:"name,omitempty"`
	} `json:"test_pipelines"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	UUID      string    `json:"uuid"`
}

// ProjectList holds a list of Project objects
type ProjectList struct {
	Projects []Project
}

type projectResponse struct {
	Project Project
}

// ListProjects Fetch a list of projects for the given organization
func (api *API) ListProjects(orgID string) (ProjectList, error) {
	projectList := ProjectList{}
	orgID = api.getOrgUUID(orgID)
	path := fmt.Sprintf("/organizations/%s/projects", orgID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return projectList, errors.Wrap(err, "Unable to list projects")
	}

	err = json.Unmarshal(resp, &projectList)
	if err != nil {
		return projectList, errors.Wrap(err, "Unable to unmarshal JSON into ProjectList")
	}

	return projectList, nil
}

// GetProject Fetch a project by ID
func (api *API) GetProject(orgID string, projectID string) (Project, error) {
	project := projectResponse{}
	path := fmt.Sprintf("/organizations/%s/projects/%s", orgID, projectID)

	resp, err := api.makeRequest("GET", path, nil)
	if err != nil {
		return project.Project, errors.Wrap(err, "Unable to get project")
	}

	err = json.Unmarshal(resp, &project)
	if err != nil {
		return project.Project, errors.Wrap(err, "Unable to unmarshal API response, error")
	}

	return project.Project, nil
}

// CreateProject Create a new project
func (api *API) CreateProject(orgID string, project Project) (Project, error) {
	path := fmt.Sprintf("/organizations/%s/projects", orgID)

	resp, err := api.makeRequest("POST", path, project)
	if err != nil {
		return project, errors.Wrap(err, "Unable to create project, error")
	}

	projResponse := projectResponse{}
	err = json.Unmarshal(resp, &projResponse)
	if err != nil {
		return project, errors.Wrap(err, "Unable to unmarshal response into Project, error")
	}

	return projResponse.Project, nil
}
