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
	"github.com/ca-risken/core/proto/org_iam"
	org_iammock "github.com/ca-risken/core/proto/org_iam/mocks"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
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
			name:  "OK with ProjectId",
			input: &organization.ListOrganizationRequest{ProjectId: 123},
			want: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1, Name: "org1", Description: "org1 desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{OrganizationId: 2, Name: "org2", Description: "org2 desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResponce: []*model.Organization{
				{OrganizationID: 1, Name: "org1", Description: "org1 desc", CreatedAt: now, UpdatedAt: now},
				{OrganizationID: 2, Name: "org2", Description: "org2 desc", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:  "OK with all filters",
			input: &organization.ListOrganizationRequest{OrganizationId: 1, Name: "test", ProjectId: 123},
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
				mockDB.On("ListOrganization", test.RepeatMockAnything(5)...).Return(c.mockResponce, c.mockError).Once()
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
		putPolicyResponse          *org_iam.PutOrgPolicyResponse
		putRoleResponse            *org_iam.PutOrgRoleResponse
		attachPolicyResponse       *org_iam.AttachOrgPolicyResponse
		attachRoleResponse         *org_iam.AttachOrgRoleResponse
		mockOrgIAMError   error
	}{
		{
			name:                       "OK",
			input:                      &organization.CreateOrganizationRequest{Name: "nm", Description: "desc", UserId: 1},
			want:                       &organization.CreateOrganizationResponse{Organization: &organization.Organization{OrganizationId: 1, Name: "nm", Description: "desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			createOrganizationResponse: &model.Organization{OrganizationID: 1, Name: "nm", Description: "desc", CreatedAt: now, UpdatedAt: now},
			putPolicyResponse:          &org_iam.PutOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			putRoleResponse:            &org_iam.PutOrgRoleResponse{Role: &org_iam.OrgRole{RoleId: 1, OrganizationId: 1, Name: "role", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachPolicyResponse:       &org_iam.AttachOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachRoleResponse:         &org_iam.AttachOrgRoleResponse{Role: &org_iam.OrgRole{RoleId: 1, OrganizationId: 1, Name: "role", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
		},
		{
			name:    "NG Invalid param",
			input:   &organization.CreateOrganizationRequest{Name: "", UserId: 0},
			wantErr: true,
		},
		{
			name:                    "Invalid DB error",
			input:                   &organization.CreateOrganizationRequest{Name: "nm", Description: "desc", UserId: 1},
			createOrganizationError: gorm.ErrInvalidDB,
			wantErr:                 true,
		},
		{
			name:                       "NG Organization IAM service error",
			input:                      &organization.CreateOrganizationRequest{Name: "nm", Description: "desc", UserId: 1},
			createOrganizationResponse: &model.Organization{OrganizationID: 1, Name: "nm", Description: "desc", CreatedAt: now, UpdatedAt: now},
			putPolicyResponse:          &org_iam.PutOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockOrgIAMError:   errors.New("Something error occurred"),
			wantErr:                    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			mockOrgIAM := org_iammock.NewOrgIAMServiceClient(t)
			svc := OrganizationService{
				repository:            mockDB,
				orgIamClient: mockOrgIAM,
				logger:                logging.NewLogger(),
			}
			if c.createOrganizationResponse != nil || c.createOrganizationError != nil {
				mockDB.On("CreateOrganization", test.RepeatMockAnything(3)...).Return(c.createOrganizationResponse, c.createOrganizationError).Once()
			}
			if c.putPolicyResponse != nil {
				if c.wantErr && c.mockOrgIAMError != nil {
					mockOrgIAM.On("PutOrgPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockOrgIAMError).Once()
				} else {
					mockOrgIAM.On("PutOrgPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockOrgIAMError).Times(3)
				}
			}
			if c.putRoleResponse != nil {
				mockOrgIAM.On("PutOrgRole", test.RepeatMockAnything(2)...).Return(c.putRoleResponse, c.mockOrgIAMError).Times(3)
			}
			if c.attachPolicyResponse != nil {
				mockOrgIAM.On("AttachOrgPolicy", test.RepeatMockAnything(2)...).Return(c.attachPolicyResponse, c.mockOrgIAMError).Times(3)
			}
			if c.attachRoleResponse != nil {
				mockOrgIAM.On("AttachOrgRole", test.RepeatMockAnything(2)...).Return(c.attachRoleResponse, c.mockOrgIAMError).Once()
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
		name                      string
		input                     *organization.PutOrganizationInvitationRequest
		want                      *organization.PutOrganizationInvitationResponse
		wantErr                   bool
		mockResponse              *model.OrganizationInvitation
		mockError                 error
		mockOrgProjectExists      bool
		mockOrgProjectExistsError error
	}{
		{
			name: "OK PENDING - no existing project",
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
				Status:         organization.OrganizationInvitationStatus_PENDING.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			mockOrgProjectExists: false,
		},
		{
			name: "OK ACCEPTED - no existing project",
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
				Status:         organization.OrganizationInvitationStatus_ACCEPTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			mockOrgProjectExists: false,
		},
		{
			name: "OK ACCEPTED - with existing project",
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
				Status:         organization.OrganizationInvitationStatus_ACCEPTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			mockOrgProjectExists: true,
		},
		{
			name: "NG PENDING - existing project",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
			mockOrgProjectExists: true,
			wantErr:              true,
		},
		{
			name: "NG REJECTED - existing project",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_REJECTED,
			},
			mockOrgProjectExists: true,
			wantErr:              true,
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
			name: "NG ExistsOrganizationProject error",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
			mockOrgProjectExistsError: gorm.ErrInvalidDB,
			wantErr:                   true,
		},
		{
			name: "NG PutOrganizationInvitation DB error",
			input: &organization.PutOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_PENDING,
			},
			mockOrgProjectExists: false,
			mockError:            gorm.ErrInvalidDB,
			wantErr:              true,
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
			if c.input.OrganizationId != 0 && c.input.ProjectId != 0 {
				mockDB.On("ExistsOrganizationProject", test.RepeatMockAnything(3)...).Return(c.mockOrgProjectExists, c.mockOrgProjectExistsError).Once()
				if c.mockOrgProjectExistsError == nil && (!c.mockOrgProjectExists || c.input.Status == organization.OrganizationInvitationStatus_ACCEPTED) {
					if c.mockResponse != nil || c.mockError != nil {
						mockDB.On("PutOrganizationInvitation", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
					}
				}
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
		mockOrgProjectExists             bool
		mockOrgProjectExistsError        error
		mockRemoveProjectsError          error
	}{
		{
			name: "OK - no existing project",
			input: &organization.DeleteOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
			},
			callDeleteOrganizationInvitation: true,
			mockOrgProjectExists:             false,
		},
		{
			name: "OK - removes existing project",
			input: &organization.DeleteOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
			},
			callDeleteOrganizationInvitation: true,
			mockOrgProjectExists:             true,
		},
		{
			name:    "NG Invalid params - organization_id is zero",
			input:   &organization.DeleteOrganizationInvitationRequest{OrganizationId: 0, ProjectId: 1},
			wantErr: true,
		},
		{
			name:                             "NG DeleteOrganizationInvitation DB error",
			input:                            &organization.DeleteOrganizationInvitationRequest{OrganizationId: 1, ProjectId: 1},
			mockError:                        gorm.ErrInvalidDB,
			wantErr:                          true,
			callDeleteOrganizationInvitation: true,
		},
		{
			name: "NG ExistsOrganizationProject error",
			input: &organization.DeleteOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
			},
			callDeleteOrganizationInvitation: true,
			mockOrgProjectExistsError:        gorm.ErrInvalidDB,
			wantErr:                          true,
		},
		{
			name: "NG RemoveProjectsInOrganization error",
			input: &organization.DeleteOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
			},
			callDeleteOrganizationInvitation: true,
			mockOrgProjectExists:             true,
			mockRemoveProjectsError:          gorm.ErrInvalidDB,
			wantErr:                          true,
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
				if c.mockError == nil {
					mockDB.On("ExistsOrganizationProject", test.RepeatMockAnything(3)...).Return(c.mockOrgProjectExists, c.mockOrgProjectExistsError).Once()
					if c.mockOrgProjectExists {
						mockDB.On("RemoveProjectsInOrganization", test.RepeatMockAnything(3)...).Return(c.mockRemoveProjectsError).Once()
					}
				}
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
		name                          string
		input                         *organization.ReplyOrganizationInvitationRequest
		want                          *organization.ReplyOrganizationInvitationResponse
		wantErr                       bool
		mockResponse                  *model.OrganizationInvitation
		mockError                     error
		mockOrgProjResponse           *model.OrganizationProject
		mockOrgProjError              error
		callExistsOrganizationProject bool
		mockOrgProjectExists          bool
		mockOrgProjectExistsError     error
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
			name: "OK REJECTED - no existing project",
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
			callExistsOrganizationProject: true,
			mockOrgProjectExists:          false,
		},
		{
			name: "OK REJECTED - existing project",
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
			callExistsOrganizationProject: true,
			mockOrgProjectExists:          true,
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
			name: "NG Invalid params - organization_id is zero",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 0,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			wantErr: true,
		},
		{
			name: "NG PutOrganizationInvitation fails",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_ACCEPTED,
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
		{
			name: "NG PutOrganizationProject fails",
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
		{
			name: "NG ExistsOrganizationProject error",
			input: &organization.ReplyOrganizationInvitationRequest{
				OrganizationId: 1,
				ProjectId:      1,
				Status:         organization.OrganizationInvitationStatus_REJECTED,
			},
			mockResponse: &model.OrganizationInvitation{
				OrganizationID: 1,
				ProjectID:      1,
				Status:         organization.OrganizationInvitationStatus_REJECTED.String(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			callExistsOrganizationProject: true,
			mockOrgProjectExistsError:     gorm.ErrInvalidDB,
			wantErr:                       true,
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
				mockDB.On("PutOrganizationInvitation", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
				if c.mockOrgProjResponse != nil || c.mockOrgProjError != nil {
					mockDB.On("PutOrganizationProject", test.RepeatMockAnything(3)...).Return(c.mockOrgProjResponse, c.mockOrgProjError).Once()
				}
				if c.callExistsOrganizationProject {
					mockDB.On("ExistsOrganizationProject", test.RepeatMockAnything(3)...).Return(c.mockOrgProjectExists, c.mockOrgProjectExistsError).Once()
					if c.mockOrgProjectExists && c.mockOrgProjectExistsError == nil {
						mockDB.On("RemoveProjectsInOrganization", test.RepeatMockAnything(3)...).Return(nil).Once()
					}
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

func TestCreateProjectWithOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                       string
		input                      *organization.CreateProjectWithOrganizationRequest
		want                       *organization.CreateProjectWithOrganizationResponse
		wantErr                    bool
		createProjectResponse      *project.CreateProjectResponse
		createProjectError         error
		putInvitationResponse      *model.OrganizationInvitation
		putInvitationError         error
		putOrgProjectResponse      *model.OrganizationProject
		putOrgProjectError         error
	}{
		{
			name:  "OK",
			input: &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "test-project", OrganizationId: 1},
			want: &organization.CreateProjectWithOrganizationResponse{
				Project: &project.Project{ProjectId: 100, Name: "test-project", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			createProjectResponse: &project.CreateProjectResponse{
				Project: &project.Project{ProjectId: 100, Name: "test-project", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			putInvitationResponse: &model.OrganizationInvitation{OrganizationID: 1, ProjectID: 100, Status: "ACCEPTED", CreatedAt: now, UpdatedAt: now},
			putOrgProjectResponse: &model.OrganizationProject{OrganizationID: 1, ProjectID: 100, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param (no name)",
			input:   &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "", OrganizationId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid param (no user_id)",
			input:   &organization.CreateProjectWithOrganizationRequest{UserId: 0, Name: "test", OrganizationId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid param (no organization_id)",
			input:   &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "test", OrganizationId: 0},
			wantErr: true,
		},
		{
			name:               "NG CreateProject error",
			input:              &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "test", OrganizationId: 1},
			createProjectError: errors.New("create project error"),
			wantErr:            true,
		},
		{
			name:  "NG PutOrganizationInvitation error",
			input: &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "test", OrganizationId: 1},
			createProjectResponse: &project.CreateProjectResponse{
				Project: &project.Project{ProjectId: 100, Name: "test", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			putInvitationError: errors.New("put invitation error"),
			wantErr:            true,
		},
		{
			name:  "NG PutOrganizationProject error",
			input: &organization.CreateProjectWithOrganizationRequest{UserId: 1, Name: "test", OrganizationId: 1},
			createProjectResponse: &project.CreateProjectResponse{
				Project: &project.Project{ProjectId: 100, Name: "test", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			putInvitationResponse: &model.OrganizationInvitation{OrganizationID: 1, ProjectID: 100, Status: "ACCEPTED", CreatedAt: now, UpdatedAt: now},
			putOrgProjectError:    errors.New("put org project error"),
			wantErr:               true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			mockProjectClient := projectmock.NewProjectServiceClient(t)
			svc := OrganizationService{
				repository:    mockDB,
				projectClient: mockProjectClient,
				logger:        logging.NewLogger(),
			}
			if c.createProjectResponse != nil || c.createProjectError != nil {
				mockProjectClient.On("CreateProject", test.RepeatMockAnything(2)...).Return(c.createProjectResponse, c.createProjectError).Once()
			}
			if c.putInvitationResponse != nil || c.putInvitationError != nil {
				mockDB.On("PutOrganizationInvitation", test.RepeatMockAnything(4)...).Return(c.putInvitationResponse, c.putInvitationError).Once()
			}
			if c.putOrgProjectResponse != nil || c.putOrgProjectError != nil {
				mockDB.On("PutOrganizationProject", test.RepeatMockAnything(3)...).Return(c.putOrgProjectResponse, c.putOrgProjectError).Once()
			}
			result, err := svc.CreateProjectWithOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
