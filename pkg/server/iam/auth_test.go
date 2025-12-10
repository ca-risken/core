package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	organizationmock "github.com/ca-risken/core/proto/organization/mocks"
	"github.com/ca-risken/core/proto/organization_iam"
	organizationiammock "github.com/ca-risken/core/proto/organization_iam/mocks"
	"gorm.io/gorm"
)

func TestIsAuthorized(t *testing.T) {
	cases := []struct {
		name                        string
		input                       *iam.IsAuthorizedRequest
		want                        *iam.IsAuthorizedResponse
		wantErr                     bool
		mockUser                    *model.User
		mockUserErr                 error
		mockPolicies                *[]model.Policy
		mockPolicyErr               error
		mockOrganizationListResp    *organization.ListOrganizationResponse
		mockOrganizationListErr     error
		mockOrgIAMAuthResp          *organization_iam.IsAuthorizedOrganizationResponse
		mockOrgIAMAuthErr           error
		expectOrganizationAuthCheck bool
	}{
		{
			name:  "OK Admin user (always authorized)",
			input: &iam.IsAuthorizedRequest{UserId: 1, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: true},
			mockUser: &model.User{
				UserID: 1, Sub: "admin", Name: "Admin User", Activated: true, IsAdmin: true,
			},
		},
		{
			name:  "OK Non-admin user with valid project policies",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: true},
			mockUser: &model.User{
				UserID: 111, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
			mockPolicies: &[]model.Policy{
				{PolicyID: 101, Name: "viewer", ProjectID: 1001, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 102, Name: "put for aws", ProjectID: 1001, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:  "OK Non-admin user authorized through organization",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: true},
			mockUser: &model.User{
				UserID: 111, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
			mockPolicyErr:               gorm.ErrRecordNotFound,
			expectOrganizationAuthCheck: true,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 3001, Name: "test-org"},
				},
			},
			mockOrgIAMAuthResp: &organization_iam.IsAuthorizedOrganizationResponse{Ok: true},
		},
		{
			name:  "OK Non-admin user without valid policies and no organization access",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
			mockUser: &model.User{
				UserID: 111, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
			mockPolicyErr: gorm.ErrRecordNotFound,

			expectOrganizationAuthCheck: true,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 3001, Name: "test-org"},
				},
			},
			mockOrgIAMAuthResp: &organization_iam.IsAuthorizedOrganizationResponse{Ok: false},
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:        "NG Invalid DB error",
			input:       &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr:     true,
			mockUserErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewIAMRepository(t)
			mockOrgClient := organizationmock.NewOrganizationServiceClient(t)
			mockOrgIAMClient := organizationiammock.NewOrganizationIAMServiceClient(t)
			svc := &IAMService{
				repository:            mockRepo,
				organizationClient:    mockOrgClient,
				organizationIamClient: mockOrgIAMClient,
				logger:                logging.NewLogger(),
			}
			if c.expectOrganizationAuthCheck {
				mockOrgClient.On("ListOrganization", test.RepeatMockAnything(2)...).Return(c.mockOrganizationListResp, c.mockOrganizationListErr).Once()
				if c.mockOrganizationListResp != nil && len(c.mockOrganizationListResp.Organization) > 0 {
					mockOrgIAMClient.On("IsAuthorizedOrganization", test.RepeatMockAnything(2)...).Return(c.mockOrgIAMAuthResp, c.mockOrgIAMAuthErr).Once()
				}
			}
			if c.mockUser != nil || c.mockUserErr != nil {
				mockRepo.On("GetUser", test.RepeatMockAnything(4)...).Return(c.mockUser, c.mockUserErr).Once()
			}
			if c.mockPolicies != nil || c.mockPolicyErr != nil {
				mockRepo.On("GetUserPolicy", test.RepeatMockAnything(2)...).Return(c.mockPolicies, c.mockPolicyErr).Once()
			}
			got, err := svc.IsAuthorized(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsAuthorizedAdmin(t *testing.T) {
	cases := []struct {
		name        string
		input       *iam.IsAuthorizedAdminRequest
		want        *iam.IsAuthorizedAdminResponse
		wantErr     bool
		mockUser    *model.User
		mockUserErr error
	}{
		{
			name:  "OK Admin user",
			input: &iam.IsAuthorizedAdminRequest{UserId: 1},
			want:  &iam.IsAuthorizedAdminResponse{Ok: true},
			mockUser: &model.User{
				UserID: 1, Sub: "admin", Name: "Admin User", Activated: true, IsAdmin: true,
			},
		},
		{
			name:  "OK Non-admin user",
			input: &iam.IsAuthorizedAdminRequest{UserId: 2},
			want:  &iam.IsAuthorizedAdminResponse{Ok: false},
			mockUser: &model.User{
				UserID: 2, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
		},
		{
			name:        "OK User not found",
			input:       &iam.IsAuthorizedAdminRequest{UserId: 999},
			want:        &iam.IsAuthorizedAdminResponse{Ok: false},
			mockUserErr: gorm.ErrRecordNotFound,
		},
		{
			name:        "NG Invalid DB error",
			input:       &iam.IsAuthorizedAdminRequest{UserId: 1},
			wantErr:     true,
			mockUserErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.mockUser != nil || c.mockUserErr != nil {
				mock.On("GetUser", test.RepeatMockAnything(4)...).Return(c.mockUser, c.mockUserErr).Once()
			}
			got, err := svc.IsAuthorizedAdmin(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsAuthorizedToken(t *testing.T) {
	cases := []struct {
		name    string
		input   *iam.IsAuthorizedTokenRequest
		want    *iam.IsAuthorizedTokenResponse
		wantErr bool

		// Mock setting
		callMockMaintainerCheck bool
		mockMaintainerCheckResp bool
		mockMaintainerCheckErr  error

		callMockGetTokenPolicy  bool
		mockGetTokenPolicyResp  *[]model.Policy
		mockGetTokenPolicyError error

		callMockListOrganization bool
		mockOrganizationListResp *organization.ListOrganizationResponse
		mockOrganizationListErr  error

		callMockOrgIAMToken bool
		mockOrgIAMTokenResp *organization_iam.IsAuthorizedOrganizationTokenResponse
		mockOrgIAMTokenErr  error

		setupOrgClients bool
	}{
		{
			name:  "OK Authorized",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: true},

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: true,
			callMockGetTokenPolicy:  true,
			mockGetTokenPolicyResp: &[]model.Policy{
				{PolicyID: 101, Name: "viewer", ProjectID: 1001, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 102, Name: "put for aws", ProjectID: 1001, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:  "OK Record not found",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: false},

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: true,
			callMockGetTokenPolicy:  true,
			mockGetTokenPolicyError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:    "NG Invalid DB error",
			input:   &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: true,
			callMockGetTokenPolicy:  true,
			mockGetTokenPolicyError: gorm.ErrInvalidDB,
		},
		{
			name:  "NG No maintainer token",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: false},

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: false,
		},
		{
			name:    "NG maintainer check error",
			input:   &iam.IsAuthorizedTokenRequest{AccessTokenId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: false,
			mockMaintainerCheckErr:  gorm.ErrInvalidDB,
		},
		{
			name:  "OK skip organization authorization without organization_id",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 444, ProjectId: 1004, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: false},

			callMockMaintainerCheck: true,
			mockMaintainerCheckResp: false,
			setupOrgClients:         true,
		},
		{
			name:  "OK Authorized by organization policy",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 222, ProjectId: 1002, OrganizationId: 3001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: true},

			callMockMaintainerCheck:  true,
			mockMaintainerCheckResp:  false,
			callMockListOrganization: true,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{{OrganizationId: 3001}},
			},
			callMockOrgIAMToken: true,
			mockOrgIAMTokenResp: &organization_iam.IsAuthorizedOrganizationTokenResponse{Ok: true},
		},
		{
			name:  "OK organization authorization error ignored",
			input: &iam.IsAuthorizedTokenRequest{AccessTokenId: 333, ProjectId: 1003, OrganizationId: 4001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedTokenResponse{Ok: false},

			callMockMaintainerCheck:  true,
			mockMaintainerCheckResp:  false,
			callMockListOrganization: true,
			mockOrganizationListErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.setupOrgClients || c.callMockListOrganization || c.callMockOrgIAMToken {
				mockOrgClient := organizationmock.NewOrganizationServiceClient(t)
				mockOrgIAMClient := organizationiammock.NewOrganizationIAMServiceClient(t)
				svc.organizationClient = mockOrgClient
				svc.organizationIamClient = mockOrgIAMClient

				if c.callMockListOrganization {
					mockOrgClient.On("ListOrganization", test.RepeatMockAnything(2)...).Return(c.mockOrganizationListResp, c.mockOrganizationListErr).Once()
				}
				if c.callMockOrgIAMToken {
					mockOrgIAMClient.On("IsAuthorizedOrganizationToken", test.RepeatMockAnything(2)...).Return(c.mockOrgIAMTokenResp, c.mockOrgIAMTokenErr).Once()
				}
			}

			if c.callMockMaintainerCheck {
				mock.On("ExistsAccessTokenMaintainer", test.RepeatMockAnything(3)...).Return(c.mockMaintainerCheckResp, c.mockMaintainerCheckErr).Once()
			}
			if c.callMockGetTokenPolicy {
				mock.On("GetTokenPolicy", test.RepeatMockAnything(2)...).Return(c.mockGetTokenPolicyResp, c.mockGetTokenPolicyError).Once()
			}
			got, err := svc.IsAuthorizedToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsAuthorizedByPolicy(t *testing.T) {
	// test data
	validPolicies := &[]model.Policy{
		{PolicyID: 1, Name: "viewer", ProjectID: 1, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
		{PolicyID: 2, Name: "put", ProjectID: 1, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
	}
	invalidPolicies := &[]model.Policy{
		{PolicyID: 1, Name: "viewer", ProjectID: 1, ActionPtn: "(", ResourcePtn: ")"},
	}
	// test cases
	cases := []struct {
		name      string
		projectID uint32
		action    string
		resource  string
		policy    *[]model.Policy
		want      bool
		wantErr   bool
	}{
		{
			name:      "OK Authorized 1",
			projectID: 1,
			action:    "finding/PutFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    validPolicies,
			want:      true,
		},
		{
			name:      "Authorized 2",
			projectID: 1,
			action:    "finding/ListFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    validPolicies,
			want:      true,
		},
		{
			name:      "Authorized 3",
			projectID: 1,
			action:    "finding/DescribeFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    validPolicies,
			want:      true,
		},
		{
			name:      "Unauthorized (Not allow project)",
			projectID: 999,
			action:    "finding/PutFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    validPolicies,
			want:      false,
		},
		{
			name:      "Unauthorized (Not allow action)",
			projectID: 1,
			action:    "finding/DeleteFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    validPolicies,
			want:      false,
		},
		{
			name:      "Unauthorized (Not allow resource)",
			projectID: 1,
			action:    "finding/PutFinding",
			resource:  "github:code-scan/repository-name",
			policy:    validPolicies,
			want:      false,
		},
		{
			name:      "Error",
			projectID: 1,
			action:    "finding/PutFinding",
			resource:  "aws:guardduty/ec2-instance-id",
			policy:    invalidPolicies,
			want:      false,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := isAuthorizedByPolicy(c.projectID, c.action, c.resource, c.policy)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsAdmin(t *testing.T) {
	cases := []struct {
		name      string
		input     *iam.IsAdminRequest
		want      *iam.IsAdminResponse
		wantErr   bool
		mockUser  *model.User
		mockError error
	}{
		{
			name:  "OK Admin",
			input: &iam.IsAdminRequest{UserId: 1},
			want:  &iam.IsAdminResponse{Ok: true},
			mockUser: &model.User{
				UserID: 1, Sub: "admin", Name: "Admin User", Activated: true, IsAdmin: true,
			},
		},
		{
			name:  "OK Not Admin",
			input: &iam.IsAdminRequest{UserId: 2},
			want:  &iam.IsAdminResponse{Ok: false},
			mockUser: &model.User{
				UserID: 2, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
		},
		{
			name:      "OK User not found",
			input:     &iam.IsAdminRequest{UserId: 999},
			want:      &iam.IsAdminResponse{Ok: false},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.IsAdminRequest{UserId: 0},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &iam.IsAdminRequest{UserId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.mockUser != nil || c.mockError != nil {
				mock.On("GetUser", test.RepeatMockAnything(4)...).Return(c.mockUser, c.mockError).Once()
			}
			got, err := svc.IsAdmin(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsAuthorizedByOrganizations(t *testing.T) {
	cases := []struct {
		name                     string
		userID                   uint32
		projectID                uint32
		actionName               string
		want                     bool
		wantErr                  bool
		mockOrganizationListResp *organization.ListOrganizationResponse
		mockOrganizationListErr  error
		mockOrgIAMAuthResp       []*organization_iam.IsAuthorizedOrganizationResponse
		mockOrgIAMAuthErr        []error
		expectOrgIAMAuthCalls    int
	}{
		{
			name:       "OK - User authorized through organization",
			userID:     1001,
			projectID:  2001,
			actionName: "finding/put-finding",
			want:       true,
			wantErr:    false,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 3001, Name: "test-org-1"},
					{OrganizationId: 3002, Name: "test-org-2"},
				},
			},
			expectOrgIAMAuthCalls: 2,
			mockOrgIAMAuthResp: []*organization_iam.IsAuthorizedOrganizationResponse{
				{Ok: false},
				{Ok: true},
			},
			mockOrgIAMAuthErr: []error{nil, nil},
		},
		{
			name:       "OK - User not authorized through any organization",
			userID:     1001,
			projectID:  2001,
			actionName: "finding/delete-finding",
			want:       false,
			wantErr:    false,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 3001, Name: "test-org-1"},
					{OrganizationId: 3002, Name: "test-org-2"},
				},
			},
			expectOrgIAMAuthCalls: 2,
			mockOrgIAMAuthResp: []*organization_iam.IsAuthorizedOrganizationResponse{
				{Ok: false},
				{Ok: false},
			},
			mockOrgIAMAuthErr: []error{nil, nil},
		},
		{
			name:       "OK - No organizations found for project",
			userID:     1001,
			projectID:  2001,
			actionName: "finding/put-finding",
			want:       false,
			wantErr:    false,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{},
			},
			expectOrgIAMAuthCalls: 0,
		},
		{
			name:                    "NG - Organization list error",
			userID:                  1001,
			projectID:               2001,
			actionName:              "finding/put-finding",
			want:                    false,
			wantErr:                 true,
			mockOrganizationListErr: gorm.ErrInvalidDB,
			expectOrgIAMAuthCalls:   0,
		},
		{
			name:       "OK - Organization IAM auth error (continues to next org)",
			userID:     1001,
			projectID:  2001,
			actionName: "finding/put-finding",
			want:       true,
			wantErr:    false,
			mockOrganizationListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 3001, Name: "test-org-1"},
					{OrganizationId: 3002, Name: "test-org-2"},
				},
			},
			expectOrgIAMAuthCalls: 2,
			mockOrgIAMAuthResp: []*organization_iam.IsAuthorizedOrganizationResponse{
				nil,
				{Ok: true},
			},
			mockOrgIAMAuthErr: []error{gorm.ErrInvalidDB, nil},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewIAMRepository(t)
			logger := logging.NewLogger()
			mockOrgClient := organizationmock.NewOrganizationServiceClient(t)
			mockOrgIAMClient := organizationiammock.NewOrganizationIAMServiceClient(t)
			svc := &IAMService{
				repository:            mockRepo,
				logger:                logger,
				organizationClient:    mockOrgClient,
				organizationIamClient: mockOrgIAMClient,
			}
			mockOrgClient.On("ListOrganization", test.RepeatMockAnything(2)...).Return(c.mockOrganizationListResp, c.mockOrganizationListErr).Once()
			for i := 0; i < c.expectOrgIAMAuthCalls; i++ {
				var resp *organization_iam.IsAuthorizedOrganizationResponse
				var err error
				if i < len(c.mockOrgIAMAuthResp) {
					resp = c.mockOrgIAMAuthResp[i]
				}
				if i < len(c.mockOrgIAMAuthErr) {
					err = c.mockOrgIAMAuthErr[i]
				}
				mockOrgIAMClient.On("IsAuthorizedOrganization", test.RepeatMockAnything(2)...).Return(resp, err).Once()
			}
			got, err := svc.isAuthorizedByOrganizations(ctx, c.userID, c.projectID, c.actionName)
			if err != nil && !c.wantErr {
				t.Errorf("isAuthorizedByOrganizations() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("isAuthorizedByOrganizations() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestIsAuthorizedByProject(t *testing.T) {
	cases := []struct {
		name           string
		userID         uint32
		projectID      uint32
		actionName     string
		resourceName   string
		want           bool
		wantErr        bool
		mockPolicyResp *[]model.Policy
		mockPolicyErr  error
	}{
		{
			name:         "OK - User authorized by project policy",
			userID:       1001,
			projectID:    1,
			actionName:   "finding/get-finding",
			resourceName: "aws:guardduty/ec2-instance-id",
			want:         true,
			wantErr:      false,
			mockPolicyResp: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ProjectID: 1, ActionPtn: "finding/(get|list|describe)", ResourcePtn: "aws:guardduty/ec2-instance-id"},
			},
		},
		{
			name:         "OK - User not authorized by project policy",
			userID:       1001,
			projectID:    1,
			actionName:   "finding/delete-finding",
			resourceName: "aws:guardduty/ec2-instance-id",
			want:         false,
			wantErr:      false,
			mockPolicyResp: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ProjectID: 1, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: "aws:guardduty/ec2-instance-id"},
			},
		},
		{
			name:           "OK - No policies found for project",
			userID:         1001,
			projectID:      1,
			actionName:     "finding/put-finding",
			resourceName:   "aws:guardduty/ec2-instance-id",
			want:           false,
			wantErr:        false,
			mockPolicyResp: &[]model.Policy{},
			mockPolicyErr:  gorm.ErrRecordNotFound,
		},
		{
			name:       "NG - Policy compile error",
			userID:     1001,
			projectID:  1,
			actionName: "finding/put-finding",
			want:       false,
			wantErr:    true,
			mockPolicyResp: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ProjectID: 1, ActionPtn: "[", ResourcePtn: ".*"},
			},
		},
		{
			name:          "NG - Policy list error",
			userID:        1001,
			projectID:     1,
			actionName:    "finding/put-finding",
			resourceName:  "aws:guardduty/ec2-instance-id",
			want:          false,
			wantErr:       true,
			mockPolicyErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewIAMRepository(t)
			logger := logging.NewLogger()
			svc := &IAMService{
				repository: mockRepo,
				logger:     logger,
			}
			mockRepo.On("GetUserPolicy", test.RepeatMockAnything(2)...).Return(c.mockPolicyResp, c.mockPolicyErr).Once()
			got, err := svc.isAuthorizedByProject(ctx, c.userID, c.projectID, c.actionName, c.resourceName)
			if err != nil && !c.wantErr {
				t.Errorf("isAuthorizedByOrganizations() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("isAuthorizedByOrganizations() = %v, want %v", got, c.want)
			}
		})
	}
}
