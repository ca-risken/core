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
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	iammock "github.com/ca-risken/core/proto/iam/mocks"
)

func TestListProject(t *testing.T) {
	now := time.Now()
	var ctx context.Context
	mockDB := mocks.MockProjectRepository{}
	svc := ProjectService{
		repository: &mockDB,
	}
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
			mockResponce: &[]db.ProjectWithTag{
				{ProjectID: 1, Name: "a", Tag: &[]model.ProjectTag{
					{ProjectID: 1, Tag: "tag1", Color: "red", CreatedAt: now, UpdatedAt: now},
					{ProjectID: 1, Tag: "tag2", Color: "pink", CreatedAt: now, UpdatedAt: now},
				}, CreatedAt: now, UpdatedAt: now},
				{ProjectID: 2, Name: "b", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK No record",
			input:     &project.ListProjectRequest{UserId: 1},
			want:      &project.ListProjectResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid params",
			input:   &project.ListProjectRequest{Name: "12345678901234567890123456789012345678901234567890123456789012345"},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &project.ListProjectRequest{UserId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListProject").Return(c.mockResponce, c.mockError).Once()
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
	var ctx context.Context
	mockDB := mocks.MockProjectRepository{}
	mockIAM := iammock.IAMServiceClient{}
	svc := ProjectService{
		repository: &mockDB,
		iamClient:  &mockIAM,
		logger:     logging.NewLogger(),
	}
	cases := []struct {
		name                  string
		input                 *project.CreateProjectRequest
		want                  *project.CreateProjectResponse
		wantErr               bool
		createProjectResponse *model.Project
		putPolicyResponse     *iam.PutPolicyResponse
		putRoleResponce       *iam.PutRoleResponse
		attachPolicyResponse  *iam.AttachPolicyResponse
		attachRoleResponse    *iam.AttachRoleResponse
		createProjectError    error
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
			input:              &project.CreateProjectRequest{UserId: 1},
			createProjectError: gorm.ErrInvalidDB,
			wantErr:            true,
		},
		{
			name:                  "NG IAM service error",
			input:                 &project.CreateProjectRequest{UserId: 1},
			createProjectResponse: &model.Project{ProjectID: 1, Name: "nm", CreatedAt: now, UpdatedAt: now},
			mockIAMError:          errors.New("Something error occured"),
			wantErr:               true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.createProjectResponse != nil || c.createProjectError != nil {
				mockDB.On("CreateProject").Return(c.createProjectResponse, c.createProjectError).Once()
			}
			if c.putRoleResponce != nil {
				mockIAM.On("PutRole", mock.Anything, mock.Anything).Return(c.putRoleResponce, nil).Once()
			}
			if c.attachRoleResponse != nil {
				mockIAM.On("AttachRole", mock.Anything, mock.Anything).Return(c.attachRoleResponse, nil).Once()
			}
			if c.putPolicyResponse != nil {
				mockIAM.On("PutPolicy", mock.Anything, mock.Anything).Return(c.putPolicyResponse, nil).Once()
			}
			if c.attachPolicyResponse != nil {
				mockIAM.On("AttachPolicy", mock.Anything, mock.Anything).Return(c.attachPolicyResponse, nil).Once()
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
	var ctx context.Context
	mockDB := mocks.MockProjectRepository{}
	svc := ProjectService{
		repository: &mockDB,
	}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpdateProject").Return(c.mockResponce, c.mockError).Once()
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
	mockIAM := iammock.IAMServiceClient{}
	svc := ProjectService{
		iamClient: &mockIAM,
		logger:    logging.NewLogger(),
	}
	cases := []struct {
		name             string
		input            *project.DeleteProjectRequest
		wantErr          bool
		listRoleResponse *iam.ListRoleResponse
		mockIAMError     error
	}{
		{
			name:             "OK",
			input:            &project.DeleteProjectRequest{ProjectId: 1},
			listRoleResponse: &iam.ListRoleResponse{RoleId: []uint32{1}},
		},
		{
			name:    "NG Invalid params",
			input:   &project.DeleteProjectRequest{},
			wantErr: true,
		},
		{
			name:         "NG IAM service error",
			input:        &project.DeleteProjectRequest{},
			mockIAMError: errors.New("Something error occured"),
			wantErr:      true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.listRoleResponse != nil {
				mockIAM.On("ListRole", mock.Anything, mock.Anything).Return(c.listRoleResponse, nil).Once()
			}
			mockIAM.On("DeleteRole", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil).Once()
			_, err := svc.DeleteProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestIsActive(t *testing.T) {
	var ctx context.Context
	mockIAM := iammock.IAMServiceClient{}
	svc := ProjectService{
		iamClient: &mockIAM,
	}
	cases := []struct {
		name             string
		input            *project.IsActiveRequest
		want             *project.IsActiveResponse
		wantErr          bool
		listUserResponse *iam.ListUserResponse
		mockError        error
	}{
		{
			name:             "OK",
			input:            &project.IsActiveRequest{ProjectId: 1},
			want:             &project.IsActiveResponse{Active: true},
			listUserResponse: &iam.ListUserResponse{UserId: []uint32{1}},
		},
		{
			name:    "NG Invalid params",
			input:   &project.IsActiveRequest{},
			wantErr: true,
		},
		{
			name:      "NG IAM service error",
			input:     &project.IsActiveRequest{},
			mockError: errors.New("Something error occured"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.listUserResponse != nil {
				mockIAM.On("ListUser", mock.Anything, mock.Anything).Return(c.listUserResponse, nil).Once()
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
