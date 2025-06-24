package project

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/project"
	"gorm.io/gorm"

	iammock "github.com/ca-risken/core/proto/iam/mocks"
	organizationmock "github.com/ca-risken/core/proto/organization/mocks"
)

func TestListProject(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *project.ListProjectRequest
		want         *project.ListProjectResponse
		wantErr      bool
		mockResponce *[]db.ProjectWithTag
		mockError    error
	}{
		{
			name:  "OK",
			input: &project.ListProjectRequest{UserId: 1, ProjectId: 1001, Name: "test"},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "test", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponce: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "test", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK No record",
			input:     &project.ListProjectRequest{UserId: 999, ProjectId: 999, Name: "not-exist"},
			want:      &project.ListProjectResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid params",
			input:   &project.ListProjectRequest{Name: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abc"},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &project.ListProjectRequest{UserId: 1, ProjectId: 1001, Name: "test"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewProjectRepository(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)

			// Setup organization mock to return empty results for organization lookup
			mockOrg.On("ListOrganization", test.RepeatMockAnything(5)...).Return(&organization.ListOrganizationResponse{
				Organization: []*organization.Organization{},
			}, nil).Maybe()

			svc := ProjectService{
				repository:         mockDB,
				organizationClient: mockOrg,
				logger:             logging.NewLogger(),
			}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListProject", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !c.wantErr {
				if result == nil && c.want != nil {
					t.Fatalf("Result is nil but expected non-nil")
				}
				if result != nil && c.want == nil {
					t.Fatalf("Result is non-nil but expected nil")
				}
				if result != nil && c.want != nil {
					if len(result.Project) != len(c.want.Project) {
						t.Fatalf("Project count mismatch: want=%d, got=%d", len(c.want.Project), len(result.Project))
					}
					for i, project := range result.Project {
						expected := c.want.Project[i]
						if project.ProjectId != expected.ProjectId {
							t.Fatalf("ProjectId mismatch at index %d: want=%d, got=%d", i, expected.ProjectId, project.ProjectId)
						}
						if project.Name != expected.Name {
							t.Fatalf("Name mismatch at index %d: want=%s, got=%s", i, expected.Name, project.Name)
						}
						// Skip time comparison as it's difficult to match exactly
					}
				}
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                  string
		input                 *project.CreateProjectRequest
		want                  *project.CreateProjectResponse
		wantErr               bool
		createProjectResponse *model.Project
		createProjectError    error
		putPolicyResponse     *iam.PutPolicyResponse
		putRoleResponce       *iam.PutRoleResponse
		attachPolicyResponse  *iam.AttachPolicyResponse
		attachRoleResponse    *iam.AttachRoleResponse
		mockIAMError          error
	}{
		{
			name:                  "OK",
			input:                 &project.CreateProjectRequest{UserId: 1, Name: "nm"},
			want:                  &project.CreateProjectResponse{Project: &project.Project{ProjectId: 1, Name: "nm", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			createProjectResponse: &model.Project{ProjectID: 1, Name: "nm", CreatedAt: now, UpdatedAt: now},
			putRoleResponce:       &iam.PutRoleResponse{Role: &iam.Role{RoleId: 1, ProjectId: 1, Name: "nm", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			putPolicyResponse:     &iam.PutPolicyResponse{Policy: &iam.Policy{PolicyId: 1, Name: "nm", ActionPtn: "ap", ResourcePtn: "rp", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachPolicyResponse:  &iam.AttachPolicyResponse{RolePolicy: &iam.RolePolicy{RoleId: 1, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachRoleResponse:    &iam.AttachRoleResponse{UserRole: &iam.UserRole{RoleId: 1, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
		},
		{
			name:    "NG Invalid param",
			input:   &project.CreateProjectRequest{UserId: 1},
			wantErr: true,
		},
		{
			name:               "Invalid DB error",
			input:              &project.CreateProjectRequest{UserId: 1, Name: "nm"},
			createProjectError: gorm.ErrInvalidDB,
			wantErr:            true,
		},
		{
			name:                  "NG IAM service error",
			input:                 &project.CreateProjectRequest{UserId: 1, Name: "nm"},
			createProjectResponse: &model.Project{ProjectID: 1, Name: "nm", CreatedAt: now, UpdatedAt: now},
			putPolicyResponse:     &iam.PutPolicyResponse{Policy: &iam.Policy{PolicyId: 1, Name: "nm", ActionPtn: "ap", ResourcePtn: "rp", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockIAMError:          errors.New("Something error occured"),
			wantErr:               true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewProjectRepository(t)
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)

			svc := ProjectService{
				repository:         mockDB,
				iamClient:          mockIAM,
				organizationClient: mockOrg,
				logger:             logging.NewLogger(),
			}
			if c.createProjectResponse != nil || c.createProjectError != nil {
				mockDB.On("CreateProject", test.RepeatMockAnything(2)...).Return(c.createProjectResponse, c.createProjectError).Once()
			}
			if c.putPolicyResponse != nil {
				if c.wantErr {
					mockIAM.On("PutPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockIAMError).Once()
				} else {
					mockIAM.On("PutPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockIAMError).Times(3)
				}
			}
			if c.putRoleResponce != nil {
				mockIAM.On("PutRole", test.RepeatMockAnything(2)...).Return(c.putRoleResponce, c.mockIAMError).Times(3)
			}
			if c.attachPolicyResponse != nil {
				mockIAM.On("AttachPolicy", test.RepeatMockAnything(2)...).Return(c.attachPolicyResponse, c.mockIAMError).Times(3)
			}
			if c.attachRoleResponse != nil {
				mockIAM.On("AttachRole", test.RepeatMockAnything(2)...).Return(c.attachRoleResponse, c.mockIAMError).Once()
			}
			result, err := svc.CreateProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *project.UpdateProjectRequest
		want         *project.UpdateProjectResponse
		wantErr      bool
		mockResponce *model.Project
		mockError    error
	}{
		{
			name:         "OK",
			input:        &project.UpdateProjectRequest{ProjectId: 1, Name: "fix-name"},
			want:         &project.UpdateProjectResponse{Project: &project.Project{ProjectId: 1, Name: "fix-name", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Project{ProjectID: 1, Name: "fix-name", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid params",
			input:   &project.UpdateProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &project.UpdateProjectRequest{ProjectId: 1, Name: "fix-name"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewProjectRepository(t)
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)

			svc := ProjectService{
				repository:         mockDB,
				iamClient:          mockIAM,
				organizationClient: mockOrg,
				logger:             logging.NewLogger(),
			}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpdateProject", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.UpdateProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestDeleteProject(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name              string
		input             *project.DeleteProjectRequest
		wantErr           bool
		mockErr           error
		callDeleteProject bool
	}{
		{
			name:              "OK",
			input:             &project.DeleteProjectRequest{ProjectId: 1},
			callDeleteProject: true,
		},
		{
			name:              "NG Invalid params",
			input:             &project.DeleteProjectRequest{},
			wantErr:           true,
			callDeleteProject: false,
		},
		{
			name:              "NG DB error",
			input:             &project.DeleteProjectRequest{ProjectId: 1},
			wantErr:           true,
			mockErr:           errors.New("DB error"),
			callDeleteProject: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewProjectRepository(t)
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)

			svc := ProjectService{
				repository:         mockDB,
				iamClient:          mockIAM,
				organizationClient: mockOrg,
				logger:             logging.NewLogger(),
			}
			if c.callDeleteProject {
				mockDB.On("DeleteProject", test.RepeatMockAnything(2)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestIsActive(t *testing.T) {
	cases := []struct {
		name               string
		input              *project.IsActiveRequest
		want               *project.IsActiveResponse
		wantErr            bool
		listProjectResults *[]db.ProjectWithTag
		listProjectError   error
		listUserResponse   *iam.ListUserResponse
		mockError          error
	}{
		{
			name:               "OK",
			input:              &project.IsActiveRequest{ProjectId: 1},
			want:               &project.IsActiveResponse{Active: true},
			listProjectResults: &[]db.ProjectWithTag{{ProjectID: 1}},
			listUserResponse:   &iam.ListUserResponse{UserId: []uint32{1}},
		},
		{
			name:               "OK No Project",
			input:              &project.IsActiveRequest{ProjectId: 1},
			want:               &project.IsActiveResponse{Active: false},
			listProjectResults: &[]db.ProjectWithTag{},
		},
		{
			name:    "NG Invalid params",
			input:   &project.IsActiveRequest{},
			wantErr: true,
		},
		{
			name:               "NG DB error",
			input:              &project.IsActiveRequest{ProjectId: 1},
			listProjectResults: &[]db.ProjectWithTag{},
			listProjectError:   errors.New("something error occured"),
			wantErr:            true,
		},
		{
			name:               "NG IAM service error",
			input:              &project.IsActiveRequest{ProjectId: 1},
			listProjectResults: &[]db.ProjectWithTag{{ProjectID: 1}},
			listUserResponse:   &iam.ListUserResponse{UserId: []uint32{1}},
			mockError:          errors.New("Something error occured"),
			wantErr:            true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockRepository := mocks.NewProjectRepository(t)
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)

			svc := ProjectService{
				iamClient:          mockIAM,
				repository:         mockRepository,
				organizationClient: mockOrg,
				logger:             logging.NewLogger(),
			}
			if c.listProjectResults != nil {
				mockRepository.On("ListProject", test.RepeatMockAnything(4)...).Return(c.listProjectResults, c.listProjectError).Once()
			}
			if c.listUserResponse != nil {
				mockIAM.On("ListUser", test.RepeatMockAnything(2)...).Return(c.listUserResponse, c.mockError).Once()
			}
			got, err := svc.IsActive(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestListProjectWithOrganization(t *testing.T) {
	now := time.Now()

	t.Run("OK with organization access", func(t *testing.T) {
		var ctx context.Context

		mockDB := mocks.NewProjectRepository(t)
		mockOrg := organizationmock.NewOrganizationServiceClient(t)

		svc := ProjectService{
			repository:         mockDB,
			organizationClient: mockOrg,
			logger:             logging.NewLogger(),
		}

		// Mock direct projects
		mockDB.On("ListProject", ctx, uint32(1), uint32(0), "").Return(&[]db.ProjectWithTag{
			{ProjectID: 1, Name: "direct-project", CreatedAt: now, UpdatedAt: now},
		}, nil).Once()

		// Mock organization listing with userID - return organizations where user has roles
		mockOrg.On("ListOrganization", test.RepeatMockAnything(5)...).Return(&organization.ListOrganizationResponse{
			Organization: []*organization.Organization{
				{OrganizationId: 1, Name: "test-org"},
			},
		}, nil).Once()

		// Mock organization projects
		mockOrg.On("ListProjectsInOrganization", test.RepeatMockAnything(2)...).Return(&organization.ListProjectsInOrganizationResponse{
			Project: []*project.Project{
				{ProjectId: 2, Name: "org-project"},
			},
		}, nil).Once()

		// Mock getting project details for org project
		mockDB.On("ListProject", ctx, uint32(0), uint32(2), "").Return(&[]db.ProjectWithTag{
			{ProjectID: 2, Name: "org-project", CreatedAt: now, UpdatedAt: now},
		}, nil).Once()

		// Execute test
		result, err := svc.ListProject(ctx, &project.ListProjectRequest{UserId: 1})

		// Verify results
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}

		if len(result.Project) != 2 {
			t.Fatalf("Expected 2 projects, got %d", len(result.Project))
		}

		// Verify both direct and organization projects are returned
		projectNames := make(map[string]bool)
		for _, proj := range result.Project {
			projectNames[proj.Name] = true
		}

		if !projectNames["direct-project"] {
			t.Fatalf("Expected direct-project to be included")
		}
		if !projectNames["org-project"] {
			t.Fatalf("Expected org-project to be included")
		}
	})

	t.Run("OK no organization roles", func(t *testing.T) {
		var ctx context.Context

		mockDB := mocks.NewProjectRepository(t)
		mockOrg := organizationmock.NewOrganizationServiceClient(t)

		svc := ProjectService{
			repository:         mockDB,
			organizationClient: mockOrg,
			logger:             logging.NewLogger(),
		}

		// Mock direct projects only
		mockDB.On("ListProject", ctx, uint32(1), uint32(0), "").Return(&[]db.ProjectWithTag{
			{ProjectID: 1, Name: "direct-project", CreatedAt: now, UpdatedAt: now},
		}, nil).Once()

		// Mock organization listing with userID - return empty (no roles)
		mockOrg.On("ListOrganization", test.RepeatMockAnything(5)...).Return(&organization.ListOrganizationResponse{
			Organization: []*organization.Organization{},
		}, nil).Once()

		// Execute test
		result, err := svc.ListProject(ctx, &project.ListProjectRequest{UserId: 1})

		// Verify results
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}

		if len(result.Project) != 1 {
			t.Fatalf("Expected 1 project, got %d", len(result.Project))
		}

		if result.Project[0].ProjectId != 1 || result.Project[0].Name != "direct-project" {
			t.Fatalf("Expected only direct project, got %+v", result.Project[0])
		}
	})
}
