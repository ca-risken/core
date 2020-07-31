package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

func TestListProject(t *testing.T) {
	now := time.Now()
	var ctx context.Context
	mockDB := mockRepository{}
	svc := projectService{
		repository: &mockDB,
	}
	cases := []struct {
		name         string
		input        *project.ListProjectRequest
		want         *project.ListProjectResponse
		wantErr      bool
		mockResponce *[]model.Project
		mockError    error
	}{
		{
			name:  "OK",
			input: &project.ListProjectRequest{UserId: 1},
			want: &project.ListProjectResponse{
				Project: []*project.Project{
					{ProjectId: 111, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{ProjectId: 222, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponce: &[]model.Project{
				{ProjectID: 111, CreatedAt: now, UpdatedAt: now},
				{ProjectID: 222, CreatedAt: now, UpdatedAt: now},
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
			name:      "NG DB error",
			input:     &project.ListProjectRequest{UserId: 1},
			mockError: gorm.ErrCantStartTransaction,
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
	mockDB := mockRepository{}
	mockIAM := mockClient{}
	svc := projectService{
		repository: &mockDB,
		iamClient:  &mockIAM,
	}
	cases := []struct {
		name         string
		input        *project.CreateProjectRequest
		want         *project.CreateProjectResponse
		wantErr      bool
		mockResponce *model.Project
		mockError    error
		mockIAMError error
	}{
		{
			name:         "OK",
			input:        &project.CreateProjectRequest{UserId: 1, Name: "nm"},
			want:         &project.CreateProjectResponse{Project: &project.Project{ProjectId: 1, Name: "nm", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Project{ProjectID: 1, Name: "nm", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &project.CreateProjectRequest{UserId: 1},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &project.CreateProjectRequest{UserId: 1},
			mockError: gorm.ErrCantStartTransaction,
			wantErr:   true,
		},
		{
			name:         "NG IAM service error",
			input:        &project.CreateProjectRequest{UserId: 1},
			mockResponce: &model.Project{ProjectID: 1, Name: "nm", CreatedAt: now, UpdatedAt: now},
			mockIAMError: errors.New("Something error occured"),
			wantErr:      true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("CreateProject").Return(c.mockResponce, c.mockError).Once()
			}
			mockIAM.On("CreateDefaultRole").Return(c.mockIAMError).Once()
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
	mockDB := mockRepository{}
	svc := projectService{
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
			name:      "NG DB error",
			input:     &project.UpdateProjectRequest{ProjectId: 1, Name: "fix-name"},
			mockError: gorm.ErrCantStartTransaction,
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
	mockIAM := mockClient{}
	svc := projectService{
		iamClient: &mockIAM,
	}
	cases := []struct {
		name         string
		input        *project.DeleteProjectRequest
		wantErr      bool
		mockIAMError error
	}{
		{
			name:  "OK",
			input: &project.DeleteProjectRequest{ProjectId: 1},
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
			mockIAM.On("DeleteAllProjectRole").Return(c.mockIAMError).Once()
			_, err := svc.DeleteProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

/**
 * Mock Repository
**/
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) ListProject(uint32, uint32, string) (*[]model.Project, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Project), args.Error(1)
}
func (m *mockRepository) CreateProject(string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}
func (m *mockRepository) UpdateProject(uint32, string) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}

/**
 * Mock GRPC Client
**/
type mockClient struct {
	mock.Mock
}

func (m *mockClient) CreateDefaultRole(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockClient) DeleteAllProjectRole(context.Context, uint32) error {
	args := m.Called()
	return args.Error(0)
}
