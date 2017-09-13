package codeship

type BuildTestFixture struct {
	Build                 Build
	RestartBuildRef       string
	RestartBuildCommitSha string
	OrgUUID               string
	ProjectUUID           string
}

type BuildTestFixtures map[string]BuildTestFixture

func getBuildFixtures() BuildTestFixtures {
	return map[string]BuildTestFixture{
		"create": BuildTestFixture{
			Build: Build{
				OrganizationUUID: "",
				ProjectUUID:      "",
				Ref:              "heads/master",
				CommitSha:        "",
			},
			RestartBuildRef:       "",
			RestartBuildCommitSha: "",
			OrgUUID:               "",
		},
		"restart": BuildTestFixture{
			Build: Build{
				OrganizationUUID: "",
				ProjectUUID:      "",
				Ref:              "heads/master",
				CommitSha:        "",
				UUID:             "",
			},
			RestartBuildRef:       "",
			RestartBuildCommitSha: "",
			OrgUUID:               "",
		},
		"buildservices": BuildTestFixture{
			Build: Build{
				OrganizationUUID: "",
				ProjectUUID:      "",
				Ref:              "heads/master",
				CommitSha:        "",
				UUID:             "",
			},
			RestartBuildRef:       "",
			RestartBuildCommitSha: "",
			OrgUUID:               "",
		},
	}
}

type ProjectTestFixture struct {
	Project               Project
	RestartBuildRef       string
	RestartBuildCommitSha string
	OrgUUID               string
}

func getCreateProjectFixtures() []ProjectTestFixture {

	return []ProjectTestFixture{
		{
			Project: Project{
				RepositoryURL: "",
				TestPipelines: []struct {
					Commands []string `json:"commands,omitempty"`
					Name     string   `json:"name,omitempty"`
				}{
					{
						Commands: []string{""},
						Name:     "test pass",
					},
				},
				Type: TypeBasic,
			},
			RestartBuildRef:       "heads/master",
			RestartBuildCommitSha: "",
			OrgUUID:               "",
		},
		{
			Project: Project{
				RepositoryURL: "",
				Type:          TypePro,
			},
			RestartBuildRef:       "heads/master",
			RestartBuildCommitSha: "",
			OrgUUID:               "",
		},
	}
}
