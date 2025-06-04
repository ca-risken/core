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
	"github.com/ca-risken/core/proto/organization_iam"
	organization_iammock "github.com/ca-risken/core/proto/organization_iam/mocks"
	"github.com/ca-risken/core/proto/project"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
		putPolicyResponse          *organization_iam.PutOrganizationPolicyResponse
		putRoleResponse            *organization_iam.PutOrganizationRoleResponse
		attachPolicyResponse       *organization_iam.AttachOrganizationPolicyResponse
		attachRoleResponse         *organization_iam.AttachOrganizationRoleResponse
		mockOrganizationIAMError   error
	}{
		{
			name:                       "OK",
			input:                      &organization.CreateOrganizationRequest{Name: "nm", Description: "desc", UserId: 1},
			want:                       &organization.CreateOrganizationResponse{Organization: &organization.Organization{OrganizationId: 1, Name: "nm", Description: "desc", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			createOrganizationResponse: &model.Organization{OrganizationID: 1, Name: "nm", Description: "desc", CreatedAt: now, UpdatedAt: now},
			putPolicyResponse:          &organization_iam.PutOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			putRoleResponse:            &organization_iam.PutOrganizationRoleResponse{Role: &organization_iam.OrganizationRole{RoleId: 1, OrganizationId: 1, Name: "role", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachPolicyResponse:       &organization_iam.AttachOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			attachRoleResponse:         &organization_iam.AttachOrganizationRoleResponse{Role: &organization_iam.OrganizationRole{RoleId: 1, OrganizationId: 1, Name: "role", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
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
			putPolicyResponse:          &organization_iam.PutOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 1, Name: "policy", ActionPtn: ".*", OrganizationId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockOrganizationIAMError:   errors.New("Something error occurred"),
			wantErr:                    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrganizationRepository(t)
			mockOrganizationIAM := organization_iammock.NewOrganizationIAMServiceClient(t)
			svc := OrganizationService{
				repository:            mockDB,
				organizationIamClient: mockOrganizationIAM,
				logger:                logging.NewLogger(),
			}
			if c.createOrganizationResponse != nil || c.createOrganizationError != nil {
				mockDB.On("CreateOrganization", test.RepeatMockAnything(3)...).Return(c.createOrganizationResponse, c.createOrganizationError).Once()
			}
			if c.putPolicyResponse != nil {
				if c.wantErr && c.mockOrganizationIAMError != nil {
					mockOrganizationIAM.On("PutOrganizationPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockOrganizationIAMError).Once()
				} else {
					mockOrganizationIAM.On("PutOrganizationPolicy", test.RepeatMockAnything(2)...).Return(c.putPolicyResponse, c.mockOrganizationIAMError).Times(2)
				}
			}
			if c.putRoleResponse != nil {
				mockOrganizationIAM.On("PutOrganizationRole", test.RepeatMockAnything(2)...).Return(c.putRoleResponse, c.mockOrganizationIAMError).Times(2)
			}
			if c.attachPolicyResponse != nil {
				mockOrganizationIAM.On("AttachOrganizationPolicy", test.RepeatMockAnything(2)...).Return(c.attachPolicyResponse, c.mockOrganizationIAMError).Times(2)
			}
			if c.attachRoleResponse != nil {
				mockOrganizationIAM.On("AttachOrganizationRole", test.RepeatMockAnything(2)...).Return(c.attachRoleResponse, c.mockOrganizationIAMError).Once()
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

// MockOrganizationIAMServiceClient is a manual mock for the OrganizationIAMServiceClient
type MockOrganizationIAMServiceClient struct {
	mock.Mock
}

func (m *MockOrganizationIAMServiceClient) ListOrganizationRole(ctx context.Context, in *organization_iam.ListOrganizationRoleRequest, opts ...grpc.CallOption) (*organization_iam.ListOrganizationRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.ListOrganizationRoleResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) GetOrganizationRole(ctx context.Context, in *organization_iam.GetOrganizationRoleRequest, opts ...grpc.CallOption) (*organization_iam.GetOrganizationRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.GetOrganizationRoleResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) PutOrganizationRole(ctx context.Context, in *organization_iam.PutOrganizationRoleRequest, opts ...grpc.CallOption) (*organization_iam.PutOrganizationRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.PutOrganizationRoleResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) DeleteOrganizationRole(ctx context.Context, in *organization_iam.DeleteOrganizationRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) AttachOrganizationRole(ctx context.Context, in *organization_iam.AttachOrganizationRoleRequest, opts ...grpc.CallOption) (*organization_iam.AttachOrganizationRoleResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.AttachOrganizationRoleResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) DetachOrganizationRole(ctx context.Context, in *organization_iam.DetachOrganizationRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) ListOrganizationPolicy(ctx context.Context, in *organization_iam.ListOrganizationPolicyRequest, opts ...grpc.CallOption) (*organization_iam.ListOrganizationPolicyResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.ListOrganizationPolicyResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) GetOrganizationPolicy(ctx context.Context, in *organization_iam.GetOrganizationPolicyRequest, opts ...grpc.CallOption) (*organization_iam.GetOrganizationPolicyResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.GetOrganizationPolicyResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) PutOrganizationPolicy(ctx context.Context, in *organization_iam.PutOrganizationPolicyRequest, opts ...grpc.CallOption) (*organization_iam.PutOrganizationPolicyResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.PutOrganizationPolicyResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) DeleteOrganizationPolicy(ctx context.Context, in *organization_iam.DeleteOrganizationPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) AttachOrganizationPolicy(ctx context.Context, in *organization_iam.AttachOrganizationPolicyRequest, opts ...grpc.CallOption) (*organization_iam.AttachOrganizationPolicyResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.AttachOrganizationPolicyResponse), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) DetachOrganizationPolicy(ctx context.Context, in *organization_iam.DetachOrganizationPolicyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockOrganizationIAMServiceClient) IsAuthorizedOrganization(ctx context.Context, in *organization_iam.IsAuthorizedOrganizationRequest, opts ...grpc.CallOption) (*organization_iam.IsAuthorizedOrganizationResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*organization_iam.IsAuthorizedOrganizationResponse), args.Error(1)
}
