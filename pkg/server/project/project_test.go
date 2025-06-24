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
	organizationmock "github.com/ca-risken/core/proto/organization/mocks"
	"github.com/ca-risken/core/proto/project"
	"gorm.io/gorm"

	iammock "github.com/ca-risken/core/proto/iam/mocks"
)

func TestListProject(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                      string
		input                     *project.ListProjectRequest
		want                      *project.ListProjectResponse
		wantErr                   bool
		mockDirectProjectsResp    *[]db.ProjectWithTag
		mockDirectProjectsErr     error
		mockListOrganizationResp  *organization.ListOrganizationResponse
		mockListOrganizationErr   error
		mockOrgProjectsResp       []*organization.ListProjectsInOrganizationResponse
		mockOrgProjectsErr        []error
		mockOrgProjectDetailsResp []*[]db.ProjectWithTag
		mockOrgProjectDetailsErr  []error
	}{
		{
			name:  "OK - Direct projects only",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "a", Tag: []*project.ProjectTag{
						{ProjectId: 1, Tag: "tag1", Color: "red"},
						{ProjectId: 1, Tag: "tag2", Color: "pink"},
					}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{ProjectId: 2, Name: "b", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockDirectProjectsResp: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "a", Tag: &[]model.ProjectTag{
					{ProjectID: 1, Tag: "tag1", Color: "red", CreatedAt: now, UpdatedAt: now},
					{ProjectID: 1, Tag: "tag2", Color: "pink", CreatedAt: now, UpdatedAt: now},
				}, CreatedAt: now, UpdatedAt: now},
				{ProjectID: 2, Name: "b", CreatedAt: now, UpdatedAt: now},
			},
			mockListOrganizationResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{},
			},
		},
		{
			name:  "OK - Organization projects only",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 3, Name: "org-project", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockDirectProjectsErr: gorm.ErrRecordNotFound,
			mockListOrganizationResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001, Name: "test-org"},
				},
			},
			mockOrgProjectsResp: []*organization.ListProjectsInOrganizationResponse{
				{
					Project: []*project.Project{
						{ProjectId: 3, Name: "org-project"},
					},
				},
			},
			mockOrgProjectDetailsResp: []*[]db.ProjectWithTag{
				{
					{ProjectID: 3, Name: "org-project", CreatedAt: now, UpdatedAt: now},
				},
			},
		},
		{
			name:  "OK - Mixed direct and organization projects",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "direct-project", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{ProjectId: 3, Name: "org-project", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockDirectProjectsResp: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "direct-project", CreatedAt: now, UpdatedAt: now},
			},
			mockListOrganizationResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001, Name: "test-org"},
				},
			},
			mockOrgProjectsResp: []*organization.ListProjectsInOrganizationResponse{
				{
					Project: []*project.Project{
						{ProjectId: 3, Name: "org-project"},
					},
				},
			},
			mockOrgProjectDetailsResp: []*[]db.ProjectWithTag{
				{
					{ProjectID: 3, Name: "org-project", CreatedAt: now, UpdatedAt: now},
				},
			},
		},
		{
			name:  "OK - Duplicate projects (organization project already in direct)",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "shared-project", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockDirectProjectsResp: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "shared-project", CreatedAt: now, UpdatedAt: now},
			},
			mockListOrganizationResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001, Name: "test-org"},
				},
			},
			mockOrgProjectsResp: []*organization.ListProjectsInOrganizationResponse{
				{
					Project: []*project.Project{
						{ProjectId: 1, Name: "shared-project"},
					},
				},
			},
			mockOrgProjectDetailsResp: []*[]db.ProjectWithTag{
				{
					{ProjectID: 1, Name: "shared-project", CreatedAt: now, UpdatedAt: now},
				},
			},
		},
		{
			name:                  "OK - No direct projects, no organizations",
			input:                 &project.ListProjectRequest{UserId: 1},
			want:                  &project.ListProjectResponse{Project: []*project.Project{}},
			mockDirectProjectsErr: gorm.ErrRecordNotFound,
			mockListOrganizationResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{},
			},
		},
		{
			name:  "OK - Organization service error (fallback to direct projects)",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "direct-only", Tag: []*project.ProjectTag{}, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockDirectProjectsResp: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "direct-only", CreatedAt: now, UpdatedAt: now},
			},
			mockListOrganizationErr: errors.New("organization service unavailable"),
		},
		{
			name:    "NG Invalid params",
			input:   &project.ListProjectRequest{Name: "12345678901234567890123456789012345678901234567890123456789012345"},
			wantErr: true,
		},
		{
			name:                  "Invalid DB error",
			input:                 &project.ListProjectRequest{UserId: 1, ProjectId: 1001, Name: "test"},
			wantErr:               true,
			mockDirectProjectsErr: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewProjectRepository(t)
			mockOrgClient := organizationmock.NewOrganizationServiceClient(t)
			svc := ProjectService{
				repository:         mockDB,
				organizationClient: mockOrgClient,
				logger:             logging.NewLogger(),
			}
			if c.mockDirectProjectsResp != nil || c.mockDirectProjectsErr != nil {
				mockDB.On("ListProject", test.RepeatMockAnything(4)...).Return(c.mockDirectProjectsResp, c.mockDirectProjectsErr).Once()
			}
			if c.mockListOrganizationResp != nil || c.mockListOrganizationErr != nil {
				mockOrgClient.On("ListOrganization", test.RepeatMockAnything(2)...).Return(c.mockListOrganizationResp, c.mockListOrganizationErr).Once()
			}
			if c.mockOrgProjectsResp != nil {
				for i, resp := range c.mockOrgProjectsResp {
					var err error
					if i < len(c.mockOrgProjectsErr) {
						err = c.mockOrgProjectsErr[i]
					}
					mockOrgClient.On("ListProjectsInOrganization", test.RepeatMockAnything(2)...).Return(resp, err).Once()
				}
			}
			if c.mockOrgProjectDetailsResp != nil {
				for i, resp := range c.mockOrgProjectDetailsResp {
					var err error
					if i < len(c.mockOrgProjectDetailsErr) {
						err = c.mockOrgProjectDetailsErr[i]
					}
					mockDB.On("ListProject", test.RepeatMockAnything(4)...).Return(resp, err).Once()
				}
			}
			result, err := svc.ListProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
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
			svc := ProjectService{
				repository: mockDB,
				iamClient:  mockIAM,
				logger:     logging.NewLogger(),
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
			svc := ProjectService{repository: mockDB}
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
			svc := ProjectService{
				repository: mockDB,
				logger:     logging.NewLogger(),
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
			svc := ProjectService{
				iamClient:  mockIAM,
				repository: mockRepository,
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
