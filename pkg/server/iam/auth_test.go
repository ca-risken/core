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
	"gorm.io/gorm"
)

func TestIsAuthorized(t *testing.T) {
	cases := []struct {
		name          string
		input         *iam.IsAuthorizedRequest
		want          *iam.IsAuthorizedResponse
		wantErr       bool
		mockUser      *model.User
		mockUserErr   error
		mockPolicies  *[]model.Policy
		mockPolicyErr error
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
			name:  "OK Non-admin user without valid policies",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
			mockUser: &model.User{
				UserID: 111, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
			mockPolicyErr: gorm.ErrRecordNotFound,
		},
		{
			name:          "OK User not found",
			input:         &iam.IsAuthorizedRequest{UserId: 999, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:          &iam.IsAuthorizedResponse{Ok: false},
			mockUserErr:   gorm.ErrRecordNotFound,
			mockPolicyErr: gorm.ErrRecordNotFound,
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
			var ctx context.Context
			mockRepo := mocks.NewIAMRepository(t)

			// Create service with nil clients for basic testing
			svc := IAMService{
				repository: mockRepo,
				logger:     logging.NewLogger(),
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
			input: &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedAdminResponse{Ok: true},
			mockUser: &model.User{
				UserID: 1, Sub: "admin", Name: "Admin User", Activated: true, IsAdmin: true,
			},
		},
		{
			name:  "OK Non-admin user",
			input: &iam.IsAuthorizedAdminRequest{UserId: 2, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedAdminResponse{Ok: false},
			mockUser: &model.User{
				UserID: 2, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
		},
		{
			name:        "OK User not found",
			input:       &iam.IsAuthorizedAdminRequest{UserId: 999, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:        &iam.IsAuthorizedAdminResponse{Ok: false},
			mockUserErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:        "NG Invalid DB error",
			input:       &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
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
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

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

func TestCheckOrganizationAuthorization(t *testing.T) {
	// This test is focused on the new organization authorization logic
	// We'll test it indirectly through IsAuthorized with organization clients set to nil
	// to verify the error handling when organization services are not available

	cases := []struct {
		name          string
		input         *iam.IsAuthorizedRequest
		want          *iam.IsAuthorizedResponse
		mockUser      *model.User
		mockPolicies  *[]model.Policy
		mockPolicyErr error
	}{
		{
			name:  "OK Non-admin user falls back to false when no project policies and no org clients",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
			mockUser: &model.User{
				UserID: 111, Sub: "user", Name: "Regular User", Activated: true, IsAdmin: false,
			},
			mockPolicyErr: gorm.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockRepo := mocks.NewIAMRepository(t)

			// Create service with nil organization clients to test fallback behavior
			svc := IAMService{
				repository:            mockRepo,
				organizationClient:    nil,
				organizationIamClient: nil,
				logger:                logging.NewLogger(),
			}

			if c.mockUser != nil {
				mockRepo.On("GetUser", test.RepeatMockAnything(4)...).Return(c.mockUser, nil).Once()
			}
			if c.mockPolicyErr != nil {
				mockRepo.On("GetUserPolicy", test.RepeatMockAnything(2)...).Return(c.mockPolicies, c.mockPolicyErr).Once()
			}

			got, err := svc.IsAuthorized(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
