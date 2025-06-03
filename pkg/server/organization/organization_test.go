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

func TestListProjectsInOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.ListProjectsInOrganizationRequest
		want         *organization.ListProjectsInOrganizationResponse
		wantErr      bool
		mockResponse []*model.Project
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization.ListProjectsInOrganizationRequest{OrganizationId: 1},
			want: &organization.ListProjectsInOrganizationResponse{
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
			input:        &organization.ListProjectsInOrganizationRequest{OrganizationId: 999},
			want:         &organization.ListProjectsInOrganizationResponse{},
			mockResponse: []*model.Project{},
		},
		{
			name:    "NG Invalid params",
			input:   &organization.ListProjectsInOrganizationRequest{OrganizationId: 0},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization.ListProjectsInOrganizationRequest{OrganizationId: 1},
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
				mockDB.On("ListProjectsInOrganization", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListProjectsInOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestRemoveProjectsInOrganization(t *testing.T) {
	cases := []struct {
		name                             string
		input                            *organization.RemoveProjectsInOrganizationRequest
		wantErr                          bool
		mockError                        error
		callRemoveProjectsInOrganization bool
	}{
		{
			name:                             "OK",
			input:                            &organization.RemoveProjectsInOrganizationRequest{OrganizationId: 1, ProjectId: 1},
			callRemoveProjectsInOrganization: true,
		},
		{
			name:    "NG Invalid params - organization_id is zero",
			input:   &organization.RemoveProjectsInOrganizationRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:                             "Invalid DB error",
			input:                            &organization.RemoveProjectsInOrganizationRequest{OrganizationId: 1, ProjectId: 1},
			mockError:                        gorm.ErrInvalidDB,
			wantErr:                          true,
			callRemoveProjectsInOrganization: true,
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
			if c.callRemoveProjectsInOrganization {
				mockDB.On("RemoveProjectsInOrganization", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
			}
			_, err := svc.RemoveProjectsInOrganization(ctx, c.input)
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
					{OrganizationId: 1, ProjectId: 1, Status: organization.OrganizationInvitationStatus_PENDING, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
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
			input:   &organization.ListOrganizationInvitationRequest{OrganizationId: 0, ProjectId: 0},
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

func TestPutOrganizationInvitation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization.PutOrganizationInvitationRequest
		want         *organization.PutOrganizationInvitationResponse
		wantErr      bool
		mockResponse *model.OrganizationInvitation
		mockError    error
	}{
		{
			name: "OK",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
			want: &organization.PutOrganizationInvitationResponse{
				OrganizationInvitation: &organization.OrganizationInvitation{
					OrganizationId: 1,
					ProjectId:      1,
					Status:         organization.OrganizationInvitationStatus_PENDING,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         "PENDING",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "OK Update Status",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			want: &organization.PutOrganizationInvitationResponse{
				OrganizationInvitation: &organization.OrganizationInvitation{
					OrganizationId: 1,
					ProjectId:      1,
					Status:         organization.OrganizationInvitationStatus_ACCEPTED,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         "ACCEPTED",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "NG Invalid params - status is invalid",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_UNKNOWN,
			},
			wantErr: true,
		},
		{
			name: "NG Invalid params - organization_id is zero",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 0,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
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
				mockDB.On("PutOrganizationInvitation", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.PutOrganizationInvitation(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestDeleteOrganizationInvitation(t *testing.T) {
	cases := []struct {
		name                             string
		input                            *organization.DeleteOrganizationInvitationRequest
		wantErr                          bool
		mockError                        error
		callDeleteOrganizationInvitation bool
	}{
		{
			name: "OK",
			input: &organization.DeleteOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
			},
			callDeleteOrganizationInvitation: true,
		},
		{
			name:    "NG Invalid params - organization_id is zero",
			input:   &organization.DeleteOrganizationInvitationRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:                             "Invalid DB error",
			input:                            &organization.DeleteOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			mockError:                        gorm.ErrInvalidDB,
			wantErr:                          true,
			callDeleteOrganizationInvitation: true,
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
			if c.callDeleteOrganizationInvitation {
				mockDB.On("DeleteOrganizationInvitation", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
			}
			_, err := svc.DeleteOrganizationInvitation(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestReplyOrganizationInvitation(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                string
		input               *organization.ReplyOrganizationInvitationRequest
		want                *organization.ReplyOrganizationInvitationResponse
		wantErr             bool
		mockResponse        *model.OrganizationInvitation
		mockError           error
		mockOrgProjResponse *model.OrganizationProject
		mockOrgProjError    error
	}{
		{
			name: "OK ACCEPTED",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			want: &organization.ReplyOrganizationInvitationResponse{
				OrganizationProject: &organization.OrganizationProject{
					OrganizationId: 1,
					ProjectId:      1,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			mockOrgProjResponse: &model.OrganizationProject{
				OrganizationID: 1,
				ProjectID:      1,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "OK REJECTED",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_REJECTED,
			},
			want: &organization.ReplyOrganizationInvitationResponse{},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         organization.OrganizationInvitationStatus_REJECTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "NG Invalid params - status is invalid",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_UNKNOWN,
			},
			wantErr: true,
		},
		{
			name: "Invalid DB error - UpdateOrganizationInvitationStatus fails",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
		{
			name: "Invalid DB error - CreateOrganizationProject fails",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			mockOrgProjError: gorm.ErrInvalidDB,
			wantErr:          true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			svc := OrganizationService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("PutOrganizationInvitation", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
				if c.mockResponse != nil && c.mockResponse.Status == organization.OrganizationInvitationStatus_ACCEPTED.String() {
					mockDB.On("PutOrganizationProject", test.RepeatMockAnything(3)...).Return(c.mockOrgProjResponse, c.mockOrgProjError).Once()
				}
			}
			result, err := svc.ReplyOrganizationInvitation(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
