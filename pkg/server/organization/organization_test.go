package organization

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/project"
	"gorm.io/gorm"
)

func TestListOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.ListOrganizationRequest
		want         *organization.ListOrganizationResponse
		wantErr      bool
		mockResponce []*model.Organization
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.ListOrganizationRequest{OrganizationId: 1, Name: "test"},
			want: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1, Name: "test", Description: "test desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponce: []*model.Organization{
				{OrganizationID: 1, Name: "test", Description: "test desc", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK No record",
			input:     &organization.ListOrganizationRequest{OrganizationId: 999, Name: "not-exist"},
			want:      &organization.ListOrganizationResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid params",
			input:   &organization.ListOrganizationRequest{Name: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abc"},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.ListOrganizationRequest{OrganizationId: 1, Name: "test"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOrganization", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestCreateOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                       string
		input                      *organization.CreateOrganizationRequest
		want                       *organization.CreateOrganizationResponse
		wantErr                    bool
		createOrganizationResponse *model.Organization
		createOrganizationError    error
	}{
		{
			name:                       "OK",
			input:                      &organization.CreateOrganizationRequest{Name: "nm", Description: "desc"},
			want:                       &organization.CreateOrganizationResponse{Organization: &organization.Organization{OrganizationId: 1, Name: "nm", Description: "desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			createOrganizationResponse: &model.Organization{OrganizationID: 1, Name: "nm", Description: "desc", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &organization.CreateOrganizationRequest{Name: ""},
			wantErr: true,
		},
		{
			name:                    "Invalid DB error",
			input:                   &organization.CreateOrganizationRequest{Name: "nm", Description: "desc"},
			createOrganizationError: gorm.ErrInvalidDB,
			wantErr:                 true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{
				repository: mockDB,
				logger:     logging.NewLogger(),
			}
			if c.createOrganizationResponse != nil || c.createOrganizationError != nil {
				mockDB.On("CreateOrganization", test.RepeatMockAnything(3)...).Return(c.createOrganizationResponse, c.createOrganizationError).Once()
			}
			result, err := svc.CreateOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}

}

func TestUpdateOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.UpdateOrganizationRequest
		want         *organization.UpdateOrganizationResponse
		wantErr      bool
		mockResponce *model.Organization
		mockError    error
	}{
		{
			name:         "OK",
			input:        &organization.UpdateOrganizationRequest{OrganizationId: 1, Name: "fix-name", Description: "fix-desc"},
			want:         &organization.UpdateOrganizationResponse{Organization: &organization.Organization{OrganizationId: 1, Name: "fix-name", Description: "fix-desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Organization{OrganizationID: 1, Name: "fix-name", Description: "fix-desc", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid params",
			input:   &organization.UpdateOrganizationRequest{OrganizationId: 1, Name: ""},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.UpdateOrganizationRequest{OrganizationId: 1, Name: "fix-name", Description: "fix-desc"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpdateOrganization", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.UpdateOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                   string
		input                  *organization.DeleteOrganizationRequest
		wantErr                bool
		mockErr                error
		callDeleteOrganization bool
	}{
		{
			name:                   "OK",
			input:                  &organization.DeleteOrganizationRequest{OrganizationId: 1},
			callDeleteOrganization: true,
		},
		{
			name:                   "NG Invalid params",
			input:                  &organization.DeleteOrganizationRequest{OrganizationId: 0},
			wantErr:                true,
			callDeleteOrganization: false,
		},
		{
			name:                   "NG DB error",
			input:                  &organization.DeleteOrganizationRequest{OrganizationId: 1},
			wantErr:                true,
			mockErr:                errors.New("DB error"),
			callDeleteOrganization: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{
				repository: mockDB,
				logger:     logging.NewLogger(),
			}
			if c.callDeleteOrganization {
				mockDB.On("DeleteOrganization", test.RepeatMockAnything(2)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestListOrganizationInvitation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.ListOrganizationInvitationRequest
		want         *organization.ListOrganizationInvitationResponse
		wantErr      bool
		mockResponse []*model.OrganizationInvitation
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.ListOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			want: &organization.ListOrganizationInvitationResponse{
				OrganizationInvitations: []*organization.OrganizationInvitation{
					{OrganizationId: 1, ProjectId: 1, Status: "PENDING", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponse: []*model.OrganizationInvitation{
				{OrganizationID: 1, ProjectID: 1, Status: "PENDING", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:         "OK No record",
			input:        &organization.ListOrganizationInvitationRequest{OrganizationId: 999, ProjectId: 999},
			want:         &organization.ListOrganizationInvitationResponse{},
			mockResponse: []*model.OrganizationInvitation{},
		},
		{
			name:    "NG Invalid params",
			input:   &organization.ListOrganizationInvitationRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.ListOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListOrganizationInvitation", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListOrganizationInvitation(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestCreateOrganizationInvitation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.CreateOrganizationInvitationRequest
		want         *organization.CreateOrganizationInvitationResponse
		wantErr      bool
		mockResponse *model.OrganizationInvitation
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.CreateOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			want: &organization.CreateOrganizationInvitationResponse{
				OrganizationInvitation: &organization.OrganizationInvitation{
					OrganizationId: 1, ProjectId: 1, Status: "PENDING", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				},
			},
			mockResponse: &model.OrganizationInvitation{OrganizationID: 1, ProjectID: 1, Status: "PENDING", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid params",
			input:   &organization.CreateOrganizationInvitationRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.CreateOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{
				repository: mockDB,
				logger:     logging.NewLogger(),
			}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("CreateOrganizationInvitation", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.CreateOrganizationInvitation(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestListProjectsByOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.ListProjectsByOrganizationRequest
		want         *organization.ListProjectsByOrganizationResponse
		wantErr      bool
		mockResponse []*model.Project
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.ListProjectsByOrganizationRequest{OrganizationId: 1},
			want: &organization.ListProjectsByOrganizationResponse{
				Project: []*project.Project{
					{ProjectId: 1, Name: "test-project", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponse: []*model.Project{
				{ProjectID: 1, Name: "test-project", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:         "OK No record",
			input:        &organization.ListProjectsByOrganizationRequest{OrganizationId: 999},
			want:         &organization.ListProjectsByOrganizationResponse{},
			mockResponse: []*model.Project{},
		},
		{
			name:    "NG Invalid params",
			input:   &organization.ListProjectsByOrganizationRequest{OrganizationId: 0},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.ListProjectsByOrganizationRequest{OrganizationId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListProjectsByOrganization", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListProjectsByOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestAddProjects(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.AddProjectsRequest
		want         *organization.AddProjectsResponse
		wantErr      bool
		mockResponse *model.OrganizationProject
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.AddProjectsRequest{OrganizationId: 1, ProjectId: 1},
			want: &organization.AddProjectsResponse{
				OrganizationProject: &organization.OrganizationProject{
					OrganizationId: 1,
					ProjectId:      1,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationProject{
				OrganizationID: 1,
				ProjectID:      1,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name:    "NG Invalid params - organization_id is zero",
			input:   &organization.AddProjectsRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid params - project_id is zero",
			input:   &organization.AddProjectsRequest{OrganizationId: 1, ProjectId: 0},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.AddProjectsRequest{OrganizationId: 1, ProjectId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("AddProjects", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.AddProjects(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
