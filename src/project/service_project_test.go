package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/src/project/model"
	"gorm.io/gorm"
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
		mockResponce *[]projectWithTag
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
			mockResponce: &[]projectWithTag{
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
			name:      "Invalid DB error",
			input:     &project.CreateProjectRequest{UserId: 1},
			mockError: gorm.ErrInvalidDB,
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
