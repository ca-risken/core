package organization_iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/organization_iam"
	"gorm.io/gorm"
)

func TestIsAuthorizedByOrganizationPolicy(t *testing.T) {
	validPolicies := &[]model.OrganizationPolicy{
		{PolicyID: 1, Name: "organization-admin", OrganizationID: 1, ActionPtn: "organization/.*"},
		{PolicyID: 2, Name: "organizatino-viewer", OrganizationID: 1, ActionPtn: "project/(get|list)"},
	}
	cases := []struct {
		name           string
		organizationID uint32
		action         string
		policy         *[]model.OrganizationPolicy
		want           bool
		wantErr        bool
	}{
		{
			name:           "OK Authorized organization get",
			organizationID: 1,
			action:         "organization/get-organization",
			policy:         validPolicies,
			want:           true,
		},
		{
			name:           "OK Unauthorized different organization",
			organizationID: 999,
			action:         "organization/create-organization",
			policy:         validPolicies,
			want:           false,
		},
		{
			name:           "OK Unauthorized action not allowed",
			organizationID: 1,
			action:         "organization/delete-organization",
			policy:         &[]model.OrganizationPolicy{{PolicyID: 2, Name: "organization-viewer", OrganizationID: 1, ActionPtn: "organization/(get|list)"}},
			want:           false,
		},
		{
			name:           "NG Error invalid regex pattern",
			organizationID: 1,
			action:         "organization/get",
			policy:         &[]model.OrganizationPolicy{{PolicyID: 1, Name: "invalid-pattern", OrganizationID: 1, ActionPtn: "[invalid regex"}},
			wantErr:        true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := isAuthorizedByOrganizationPolicy(c.organizationID, c.action, c.policy)
			if (err != nil) != c.wantErr {
				t.Errorf("isAuthorizedByOrganizationPolicy() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("isAuthorizedByOrganizationPolicy() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestIsAuthorizedOrganization(t *testing.T) {
	cases := []struct {
		name         string
		input        *organization_iam.IsAuthorizedOrganizationRequest
		want         *organization_iam.IsAuthorizedOrganizationResponse
		wantErr      bool
		mockResponse *[]model.OrganizationPolicy
		mockError    error
	}{
		{
			name: "OK Authorized",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			want: &organization_iam.IsAuthorizedOrganizationResponse{Ok: true},
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 101, Name: "organization-admin", OrganizationID: 1001, ActionPtn: "organization/.*"},
			},
		},
		{
			name: "OK Unauthorized - no matching policy",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/delete-organization",
			},
			want: &organization_iam.IsAuthorizedOrganizationResponse{Ok: false},
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 102, Name: "organization-viewer", OrganizationID: 1001, ActionPtn: "organization/(get|list)"},
			},
		},
		{
			name: "OK No policies found",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			want:      &organization_iam.IsAuthorizedOrganizationResponse{Ok: false},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name: "NG Invalid params - organization_id is zero",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 0,
				ActionName:     "organization/create-organization",
			},
			wantErr: true,
		},
		{
			name: "NG Invalid params - action_name is invalid",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "",
			},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewOrganizationIAMRepository(t)
			logger := logging.NewLogger()
			svc := NewOrganizationIAMService(mockRepo, logger)

			if c.mockResponse != nil || c.mockError != nil {
				mockRepo.On("GetOrganizationPolicyByUserID", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}

			result, err := svc.IsAuthorizedOrganization(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatal("Expected error but got nil")
			}
			if !c.wantErr && !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
