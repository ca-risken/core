package organization_iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	iammock "github.com/ca-risken/core/proto/iam/mocks"
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
			name:   "OK Authorized organization get",
			action: "organization/get-organization",
			policy: validPolicies,
			want:   true,
		},
		{
			name:   "OK Unauthorized action not allowed",
			action: "organization/delete-organization",
			policy: &[]model.OrganizationPolicy{{PolicyID: 2, Name: "organization-viewer", OrganizationID: 1, ActionPtn: "organization/(get|list)"}},
			want:   false,
		},
		{
			name:    "NG Error invalid regex pattern",
			action:  "organization/get",
			policy:  &[]model.OrganizationPolicy{{PolicyID: 1, Name: "invalid-pattern", OrganizationID: 1, ActionPtn: "[invalid regex"}},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := isAuthorizedByOrganizationPolicy(c.action, c.policy)
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
		name              string
		input             *organization_iam.IsAuthorizedOrganizationRequest
		want              *organization_iam.IsAuthorizedOrganizationResponse
		wantErr           bool
		mockResponse      *[]model.OrganizationPolicy
		mockError         error
		mockIsAdminResp   bool
		mockIsAdminErr    error
		expectIsAdminCall bool
	}{
		{
			name: "OK Admin user - immediate authorization",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			want:              &organization_iam.IsAuthorizedOrganizationResponse{Ok: true},
			mockIsAdminResp:   true,
			expectIsAdminCall: true,
		},
		{
			name: "OK Authorized - non-admin user with matching policy",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			want:              &organization_iam.IsAuthorizedOrganizationResponse{Ok: true},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 101, Name: "organization-admin", OrganizationID: 1001, ActionPtn: "organization/.*"},
			},
		},
		{
			name: "OK Unauthorized - non-admin user with no matching policy",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/delete-organization",
			},
			want:              &organization_iam.IsAuthorizedOrganizationResponse{Ok: false},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 102, Name: "organization-viewer", OrganizationID: 1001, ActionPtn: "organization/(get|list)"},
			},
		},
		{
			name: "OK Unauthorized - non-admin user with no policies found",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			want:              &organization_iam.IsAuthorizedOrganizationResponse{Ok: false},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockError:         gorm.ErrRecordNotFound,
		},
		{
			name: "NG IsAdmin check error",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			mockIsAdminErr:    gorm.ErrInvalidDB,
			expectIsAdminCall: true,
			wantErr:           true,
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
			name: "NG Invalid DB error in policy check",
			input: &organization_iam.IsAuthorizedOrganizationRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockError:         gorm.ErrInvalidDB,
			wantErr:           true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewOrganizationIAMRepository(t)
			logger := logging.NewLogger()
			mockIAM := iammock.NewIAMServiceClient(t)
			svc := NewOrganizationIAMService(mockRepo, mockIAM, logger)

			if c.expectIsAdminCall {
				if c.mockIsAdminErr != nil {
					mockIAM.On("IsAdmin", test.RepeatMockAnything(2)...).Return(nil, c.mockIsAdminErr).Once()
				} else {
					mockIAM.On("IsAdmin", test.RepeatMockAnything(2)...).Return(&iam.IsAdminResponse{Ok: c.mockIsAdminResp}, nil).Once()
				}
			}

			if !c.mockIsAdminResp && c.mockIsAdminErr == nil && c.expectIsAdminCall && (c.mockResponse != nil || c.mockError != nil) {
				mockRepo.On("GetOrganizationPolicyByUserID", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
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

func TestIsAuthorizedOrganizationToken(t *testing.T) {
	cases := []struct {
		name          string
		input         *organization_iam.IsAuthorizedOrganizationTokenRequest
		want          *organization_iam.IsAuthorizedOrganizationTokenResponse
		wantErr       bool
		callExists    bool
		existsResp    bool
		existsErr     error
		callGetPolicy bool
		getPolicyResp *[]model.OrganizationPolicy
		getPolicyErr  error
	}{
		{
			name: "OK Authorized",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:          &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: true},
			callExists:    true,
			existsResp:    true,
			callGetPolicy: true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/.*"},
			},
		},
		{
			name: "OK Unauthorized policy not found",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/delete",
			},
			want:          &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: false},
			callExists:    true,
			existsResp:    true,
			callGetPolicy: true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/(get|list)"},
			},
		},
		{
			name: "OK Record not found",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:          &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: false},
			callExists:    true,
			existsResp:    true,
			callGetPolicy: true,
			getPolicyErr:  gorm.ErrRecordNotFound,
		},
		{
			name: "NG Invalid parameter - action",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "",
			},
			wantErr: true,
		},
		{
			name: "NG Exists check error",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			wantErr:    true,
			callExists: true,
			existsErr:  gorm.ErrInvalidDB,
		},
		{
			name: "NG Token not active",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:       &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: false},
			callExists: true,
			existsResp: false,
		},
		{
			name: "NG Get policy error",
			input: &organization_iam.IsAuthorizedOrganizationTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			wantErr:       true,
			callExists:    true,
			existsResp:    true,
			callGetPolicy: true,
			getPolicyErr:  gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewOrganizationIAMRepository(t)
			logger := logging.NewLogger()
			mockIAM := iammock.NewIAMServiceClient(t)
			svc := NewOrganizationIAMService(mockRepo, mockIAM, logger)

			if c.callExists {
				mockRepo.On("ExistsOrgActiveAccessToken", test.RepeatMockAnything(3)...).Return(c.existsResp, c.existsErr).Once()
			}
			if c.callGetPolicy {
				mockRepo.On("GetOrgTokenPolicy", test.RepeatMockAnything(3)...).Return(c.getPolicyResp, c.getPolicyErr).Once()
			}

			got, err := svc.IsAuthorizedOrganizationToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
