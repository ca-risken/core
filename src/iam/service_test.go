package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestIsAuthorized(t *testing.T) {
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.IsAuthorizedRequest
		want         *iam.IsAuthorizedResponse
		wantErr      bool
		mockResponce *[]model.Policy
		mockError    error
	}{
		{
			name:  "OK Authorized",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: true},
			mockResponce: &[]model.Policy{
				{PolicyID: 101, Name: "viewer", ProjectID: 1001, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 102, Name: "put for aws", ProjectID: 1001, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:      "OK Record not found",
			input:     &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:      &iam.IsAuthorizedResponse{Ok: false},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:      "NG Invalid DB error",
			input:     &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetUserPolicy").Return(c.mockResponce, c.mockError).Once()
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

func TestIsAuthorizedToken(t *testing.T) {
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
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
			if c.callMockMaintainerCheck {
				mock.On("ExistsAccessTokenMaintainer").Return(c.mockMaintainerCheckResp, c.mockMaintainerCheckErr).Once()
			}
			if c.callMockGetTokenPolicy {
				mock.On("GetTokenPolicy").Return(c.mockGetTokenPolicyResp, c.mockGetTokenPolicyError).Once()
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
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.IsAdminRequest
		want         *iam.IsAdminResponse
		wantErr      bool
		mockResponce *model.Policy
		mockError    error
	}{
		{
			name:         "OK Admin",
			input:        &iam.IsAdminRequest{UserId: 1},
			want:         &iam.IsAdminResponse{Ok: true},
			mockResponce: &model.Policy{PolicyID: 1, Name: "no-project-policy", ProjectID: 0, ActionPtn: ".*", ResourcePtn: ".*"},
		},
		{
			name:      "OK Not Admin",
			input:     &iam.IsAdminRequest{UserId: 1},
			want:      &iam.IsAdminResponse{Ok: false},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAdminRequest{UserId: 0},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &iam.IsAdminRequest{UserId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetAdminPolicy").Return(c.mockResponce, c.mockError).Once()
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

/*
 * Mock Repository
 */
type mockIAMRepository struct {
	mock.Mock
}

func (m *mockIAMRepository) ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32) (*[]model.User, error) {
	args := m.Called()
	return args.Get(0).(*[]model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUser(context.Context, uint32, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUserBySub(context.Context, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUserPolicy(context.Context, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *mockIAMRepository) GetTokenPolicy(context.Context, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *mockIAMRepository) PutUser(context.Context, *model.User) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Role), args.Error(1)
}
func (m *mockIAMRepository) GetRole(context.Context, uint32, uint32) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) GetRoleByName(context.Context, uint32, string) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) PutRole(ctx context.Context, r *model.Role) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) DeleteRole(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) AttachRole(context.Context, uint32, uint32, uint32) (*model.UserRole, error) {
	args := m.Called()
	return args.Get(0).(*model.UserRole), args.Error(1)
}
func (m *mockIAMRepository) DetachRole(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) ListPolicy(context.Context, uint32, string, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *mockIAMRepository) GetPolicy(context.Context, uint32, uint32) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) GetPolicyByName(context.Context, uint32, string) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) PutPolicy(context.Context, *model.Policy) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) DeletePolicy(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) AttachPolicy(context.Context, uint32, uint32, uint32) (*model.RolePolicy, error) {
	args := m.Called()
	return args.Get(0).(*model.RolePolicy), args.Error(1)
}
func (m *mockIAMRepository) DetachPolicy(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) GetAdminPolicy(context.Context, uint32) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AccessToken), args.Error(1)
}
func (m *mockIAMRepository) GetAccessTokenByID(ctx context.Context, projectID, accessTokenID uint32) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *mockIAMRepository) GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *mockIAMRepository) GetActiveAccessTokenHash(ctx context.Context, projectID, accessTokenID uint32, tokenHash string) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *mockIAMRepository) PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessToken), args.Error(1)
}
func (m *mockIAMRepository) DeleteAccessToken(ctx context.Context, projectID, accessTokenID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) AttachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) (*model.AccessTokenRole, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessTokenRole), args.Error(1)
}
func (m *mockIAMRepository) GetAccessTokenRole(ctx context.Context, accessTokenID, roleID uint32) (*model.AccessTokenRole, error) {
	args := m.Called()
	return args.Get(0).(*model.AccessTokenRole), args.Error(1)
}
func (m *mockIAMRepository) DetachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) ExistsAccessTokenMaintainer(ctx context.Context, projectID, accessTokenID uint32) (bool, error) {
	args := m.Called()
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockIAMRepository) ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AccessToken), args.Error(1)
}
