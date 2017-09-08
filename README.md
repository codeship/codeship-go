# Codeship API (v2) Client for Go

This is the start of an API client for the Codeship API written in Go.

> As a warning to all, I am really new to Go and this library may be crap.

## Usage
This library is intended to make integrating with Codeship fairly simple.

To start, you need to import the package:

```go
package main

import (
	codeship "github.com/fillup/codeship-go"
)
```

This library exposes the package `codeship`.

Getting a new API Client from it is done by calling `codeship.New()`:

```go
codeshipClient := codeship.New("username", "password")
```

With `codeshipClient` you can perform many actions on `Projects` and `Builds`:

### Projects

The `Project` type is defined as:

```go
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
		BuildOwner    string `json:"build_owner,omitempty"`
		BuildStatuses []string `json:"build_statuses,omitempty"`
		EmailTarget   string `json:"email_target,omitempty"`
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
```

#### List Projects

The `ProjectList` type is defined as:

```go
type ProjectList struct {
	Projects []Project
}
```

Get a list of all projects:
```go
projectList, err := codeshipClient.ListProjects(orgUUID)
// Loop through projects:
for project := range projectList.Projects {
  // do something with project
}
```

#### Get Project
Get a specific project

```go
project, err := codeshipClient.GetProject(orgUUID, projectUUID)
```

#### Create Project
Create a new project

```go
newProject := Project{
  Type:          codeship.TypePro,
  RepositoryURL: "git@github.com:my/repo.git",
  // ...
}

project, err := codeshipClient.CreateProject(orgUUID, newProject)
```

## Testing
Testing for this library is actually integration testing, not just unit testing.
It requires credentials to be able to talk to the Codeship API, so only run
tests if you know what you are doing and want to potentially change things
on your Codeship projects or builds.

If you really want to run tests:

 - Export ENV vars for `CODESHIP_USERNAME` and `CODESHIP_PASSWORD`, or copy
   `local.go.dist` to `local.go` and set username and password in that file.
 - Copy `test_fixtures.go.dist` to `test_fixtures.go` and fill it in with
   fixture information for real org/project/build data
 - Run `go test`

## Todo

- [ ] Iterate through pages of projects in List Projects to get full list or
add support for passing pagination parameters to the function
- [ ] Iterate through pages of builds in List Builds to get full list or
add support for passing pagination parameters to the function
- [ ] Decide on proper interface for methods, should each parameter be sent
separately or all together in the relevant struct?
