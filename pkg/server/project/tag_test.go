package project

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/project"
	"gorm.io/gorm"
)

func TestTagProject(t *testing.T) {
	now := time.Now()
	var ctx context.Context
	mockDB := mocks.MockProjectRepository{}
	svc := ProjectService{
		repository: &mockDB,
	}
	cases := []struct {
		name         string
		input        *project.TagProjectRequest
		want         *project.TagProjectResponse
		wantErr      bool
		mockResponce *model.ProjectTag
		mockError    error
	}{
		{
			name:         "OK",
			input:        &project.TagProjectRequest{ProjectId: 1, Tag: "tag", Color: "blue"},
			want:         &project.TagProjectResponse{ProjectTag: &project.ProjectTag{ProjectId: 1, Tag: "tag", Color: "blue", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.ProjectTag{ProjectID: 1, Tag: "tag", Color: "blue", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid params",
			input:   &project.TagProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &project.TagProjectRequest{ProjectId: 1, Tag: "tag"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("TagProject").Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.TagProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestUntagProject(t *testing.T) {
	var ctx context.Context
	mockDB := mocks.MockProjectRepository{}
	svc := ProjectService{
		repository: &mockDB,
	}
	cases := []struct {
		name      string
		input     *project.UntagProjectRequest
		wantErr   bool
		mockError error
	}{
		{
			name:  "OK",
			input: &project.UntagProjectRequest{ProjectId: 1, Tag: "tag"},
		},
		{
			name:    "NG Invalid params",
			input:   &project.UntagProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &project.UntagProjectRequest{},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("UntagProject").Return(c.mockError).Once()
			_, err := svc.UntagProject(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
