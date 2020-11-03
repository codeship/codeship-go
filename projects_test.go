package codeship_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	codeship "github.com/codeship/codeship-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProjects(t *testing.T) {
	type args struct {
		organizationUUID string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Link", "<https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2>; rel=\"last\", <https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2>; rel=\"next\"")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("projects/list.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "organization not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, fixture("not_found.json"), "organization")
			},
			status: http.StatusNotFound,
			err:    "unable to list projects: organization not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects",
				tt.args.organizationUUID),
				tt.handler)

			projects, resp, err := org.ListProjects(context.Background())

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)

			current, _ := resp.CurrentPage()
			assert.Equal(1, current)
			assert.Equal("https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2", resp.Links.Last)
			assert.Equal("https://api.codeship.com/v2/organizations/28123f10-e33d-5533-b53f-111ef8d7b14f/projects/?page=2", resp.Links.Next)

			require.Equal(2, len(projects.Projects))

			project := projects.Projects[1]

			createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:09.556Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2017-09-08T19:19:55.252Z")

			expected := codeship.Project{
				ID:                 2,
				UUID:               "83605ef0-76f8-0135-8810-6e5f001a2e3c",
				OrganizationUUID:   "28123f10-e33d-5533-b53f-111ef8d7b14f",
				Name:               "org/another-project",
				Type:               codeship.ProjectTypeBasic,
				RepositoryURL:      "https://github.com/org/another-project",
				RepositoryProvider: "github",
				AuthenticationUser: "Test User",
				NotificationRules: []codeship.NotificationRule{
					{
						Notifier:      "github",
						BranchMatch:   "exact",
						BuildStatuses: []string{"failed", "started", "recovered", "success"},
						Target:        "all",
					},
					{
						Notifier:      "email",
						BranchMatch:   "exact",
						Options:       codeship.NotificationOptions{},
						BuildStatuses: []string{"failed", "recovered"},
						Target:        "all",
					},
				},
				SSHKey:        "ssh-rsa key",
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				TeamIDs:       []int{1007, 1009},
				SetupCommands: []string{},
				TestPipelines: []codeship.TestPipeline{
					{
						ID:       5,
						Name:     "Test Commands",
						Commands: []string{"./run-tests.sh"},
					},
				},
				DeploymentPipelines: []codeship.DeploymentPipeline{
					{
						ID: 4,
						Branch: codeship.DeploymentBranch{
							BranchName: "master",
							MatchMode:  "*",
						},
						Position: 1,
					},
				},
				EnvironmentVariables: []codeship.EnvironmentVariable{},
			}

			assert.Equal(expected, project)
			assert.Equal(1, projects.Page)
			assert.Equal(2, projects.Total)
			assert.Equal(30, projects.PerPage)
		})
	}
}

func TestGetProject(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "0059df30-7701-0135-8810-6e5f001a2e3c",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("projects/get.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "project not found",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "0059df30-7701-0135-8810-6e5f001a2e3c",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("GET", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, fixture("not_found.json"), "project")
			},
			status: http.StatusNotFound,
			err:    "unable to get project: project not found",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s",
				tt.args.organizationUUID,
				tt.args.projectUUID),
				tt.handler)

			project, resp, err := org.GetProject(context.Background(), tt.args.projectUUID)

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)

			createdAt, _ := time.Parse(time.RFC3339, "2017-09-08T20:19:55.199Z")
			updatedAt, _ := time.Parse(time.RFC3339, "2017-09-13T17:13:36.336Z")

			expected := codeship.Project{
				ID:                 1,
				UUID:               "0059df30-7701-0135-8810-6e5f001a2e3c",
				OrganizationUUID:   "28123f10-e33d-5533-b53f-111ef8d7b14f",
				Name:               "org/test-project",
				Type:               codeship.ProjectTypePro,
				RepositoryURL:      "https://github.com/org/test-project",
				RepositoryProvider: "github",
				AuthenticationUser: "Test User",
				NotificationRules: []codeship.NotificationRule{
					{
						Notifier:      "github",
						BranchMatch:   "exact",
						BuildStatuses: []string{"failed", "started", "recovered", "success"},
						Target:        "all",
						Options: codeship.NotificationOptions{
							Key:  "foo",
							Room: "devs",
							URL:  "https://google.com",
						},
					},
					{
						Notifier:      "email",
						BranchMatch:   "exact",
						Options:       codeship.NotificationOptions{},
						BuildStatuses: []string{"failed", "recovered"},
						Target:        "all",
					},
				},
				AesKey:    "aeskey",
				SSHKey:    "ssh-rsa key",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				TeamIDs:   []int{1007},
			}

			assert.Equal(expected, project)
		})
	}
}

func TestCreateProject(t *testing.T) {
	type args struct {
		organizationUUID string
		projectType      codeship.ProjectType
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success (basic)",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectType:      codeship.ProjectTypeBasic,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				b, err := ioutil.ReadAll(r.Body)
				assert.NoError(err)
				defer r.Body.Close()
				assert.NotEmpty(b)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, fixture("projects/create_basic.json"))
			},
			status: http.StatusCreated,
		},
		{
			name: "success (pro)",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectType:      codeship.ProjectTypePro,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				b, err := ioutil.ReadAll(r.Body)
				assert.NoError(err)
				defer r.Body.Close()
				assert.NotEmpty(b)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, fixture("projects/create_pro.json"))
			},
			status: http.StatusCreated,
		},
		{
			name: "bad request",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectType:      codeship.ProjectTypeBasic,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("POST", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, fixture("errors.json"), "repository_url is required")
			},
			status: http.StatusBadRequest,
			err:    "unable to create project: repository_url is required",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects",
				tt.args.organizationUUID),
				tt.handler)

			project, resp, err := org.CreateProject(context.Background(), codeship.ProjectCreateRequest{
				RepositoryURL: "git@github.com/org/repo-name",
				TestPipelines: []codeship.TestPipeline{
					{
						Commands: []string{"./run-tests.sh"},
						Name:     "run tests",
					},
				},
				Type: tt.args.projectType,
			})

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.NotNil(project)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	type args struct {
		organizationUUID string
		projectUUID      string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		status  int
		err     string
	}{
		{
			name: "success",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "0059df30-7701-0135-8810-6e5f001a2e3c",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("PUT", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, fixture("projects/update.json"))
			},
			status: http.StatusOK,
		},
		{
			name: "bad request",
			args: args{
				organizationUUID: "28123f10-e33d-5533-b53f-111ef8d7b14f",
				projectUUID:      "0059df30-7701-0135-8810-6e5f001a2e3c",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert := assert.New(t)
				assert.Equal("PUT", r.Method)
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, fixture("errors.json"), "repository_url is required")
			},
			status: http.StatusBadRequest,
			err:    "unable to update project: repository_url is required",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/organizations/%s/projects/%s",
				tt.args.organizationUUID,
				tt.args.projectUUID),
				tt.handler)

			project, resp, err := org.UpdateProject(context.Background(), tt.args.projectUUID, codeship.ProjectUpdateRequest{
				Type: codeship.ProjectTypePro,
				TeamIDs: []int{
					61593, 70000,
				},
			})

			require.NotNil(resp)
			assert.Equal(tt.status, resp.StatusCode)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.NotNil(project)
		})
	}
}

func TestProjectType_String(t *testing.T) {
	tests := []struct {
		name        string
		projectType codeship.ProjectType
		want        string
	}{
		{
			name:        "pro",
			projectType: codeship.ProjectTypePro,
			want:        "pro",
		},
		{
			name:        "basic",
			projectType: codeship.ProjectTypeBasic,
			want:        "basic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.projectType.String())
		})
	}
}

func TestProjectType_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		projectType codeship.ProjectType
		args        args
		want        codeship.ProjectType
		err         string
	}{
		{
			name: "basic",
			args: args{
				data: []byte("\"basic\""),
			},
			want: codeship.ProjectTypeBasic,
		},
		{
			name: "pro",
			args: args{
				data: []byte("\"pro\""),
			},
			want: codeship.ProjectTypePro,
		},
		{
			name: "invalid",
			args: args{
				data: []byte("\"invalid\""),
			},
			err: "invalid ProjectType: invalid",
		},
		{
			name: "not string",
			args: args{
				data: []byte{},
			},
			err: "ProjectType should be a string, got []uint8",
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.projectType.UnmarshalJSON(tt.args.data)

			if tt.err != "" {
				require.Error(err)
				assert.EqualError(err, tt.err)
				return
			}

			require.NoError(err)
			assert.Equal(tt.want, tt.projectType)
		})
	}
}

func TestProjectType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		projectType codeship.ProjectType
		want        string
	}{
		{
			name:        "basic",
			projectType: codeship.ProjectTypeBasic,
			want:        `"basic"`,
		},
		{
			name:        "pro",
			projectType: codeship.ProjectTypePro,
			want:        `"pro"`,
		},
		{
			name:        "invalid",
			projectType: 2,
			want:        `""`,
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.projectType.MarshalJSON()
			require.NoError(err)
			assert.Equal(tt.want, string(b))
		})
	}
}
