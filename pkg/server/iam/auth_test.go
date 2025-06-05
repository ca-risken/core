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
	"github.com/ca-risken/core/proto/organization_iam"
	oimocks "github.com/ca-risken/core/proto/organization_iam/mocks"
	"gorm.io/gorm"
)

func TestIsAuthorized(t *testing.T) {
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
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetUserPolicy", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
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
		name            string
		input           *iam.IsAuthorizedAdminRequest
		want            *iam.IsAuthorizedAdminResponse
		wantErr         bool
		mockIAMResponce *[]model.Policy
		mockIAMError    error
		mockOrgIAMResp  *organization_iam.GetSystemAdminOrganizationPolicyResponse
		mockOrgIAMError error
	}{
		{
			name:  "OK Authorized (both IAM and Org IAM policies)",
			input: &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedAdminResponse{Ok: true},
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 2, Name: "put for aws", ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
			mockOrgIAMResp: &organization_iam.GetSystemAdminOrganizationPolicyResponse{
				OrganizationPolicies: []*organization_iam.OrganizationPolicy{
					{PolicyId: 1, Name: "system-admin", OrganizationId: 0, ActionPtn: ".*"},
				},
			},
		},
		{
			name:         "OK Not Admin (no IAM policies)",
			input:        &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:         &iam.IsAuthorizedAdminResponse{Ok: false},
			mockIAMError: gorm.ErrRecordNotFound,
		},
		{
			name:  "OK Not Admin (no Org IAM policies)",
			input: &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedAdminResponse{Ok: false},
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 2, Name: "put for aws", ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
			mockOrgIAMResp: &organization_iam.GetSystemAdminOrganizationPolicyResponse{
				OrganizationPolicies: []*organization_iam.OrganizationPolicy{},
			},
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:         "NG Invalid IAM DB error",
			input:        &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr:      true,
			mockIAMError: gorm.ErrInvalidDB,
		},
		{
			name:    "NG Invalid Org IAM service error",
			input:   &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 2, Name: "put for aws", ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
			mockOrgIAMError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			orgIamMock := oimocks.NewOrganizationIAMServiceClient(t)
			svc := IAMService{repository: mock, organizationIamClient: orgIamMock, logger: logging.NewLogger()}

			if c.mockIAMResponce != nil || c.mockIAMError != nil {
				mock.On("GetAdminPolicy", test.RepeatMockAnything(2)...).Return(c.mockIAMResponce, c.mockIAMError).Once()
			}
			if c.mockOrgIAMResp != nil || c.mockOrgIAMError != nil {
				orgIamMock.On("GetSystemAdminOrganizationPolicy", test.RepeatMockAnything(3)...).Return(c.mockOrgIAMResp, c.mockOrgIAMError).Once()
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
		name            string
		input           *iam.IsAdminRequest
		want            *iam.IsAdminResponse
		wantErr         bool
		mockIAMResponce *[]model.Policy
		mockIAMError    error
		mockOrgIAMResp  *organization_iam.GetSystemAdminOrganizationPolicyResponse
		mockOrgIAMError error
	}{
		{
			name:  "OK Admin (both IAM and Org IAM policies)",
			input: &iam.IsAdminRequest{UserId: 1},
			want:  &iam.IsAdminResponse{Ok: true},
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "no-project-policy", ProjectID: 0, ActionPtn: ".*", ResourcePtn: ".*"},
			},
			mockOrgIAMResp: &organization_iam.GetSystemAdminOrganizationPolicyResponse{
				OrganizationPolicies: []*organization_iam.OrganizationPolicy{
					{PolicyId: 1, Name: "system-admin", OrganizationId: 0, ActionPtn: ".*"},
				},
			},
		},
		{
			name:         "OK Not Admin (no IAM policies)",
			input:        &iam.IsAdminRequest{UserId: 1},
			want:         &iam.IsAdminResponse{Ok: false},
			mockIAMError: gorm.ErrRecordNotFound,
		},
		{
			name:  "OK Not Admin (no Org IAM policies)",
			input: &iam.IsAdminRequest{UserId: 1},
			want:  &iam.IsAdminResponse{Ok: false},
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "no-project-policy", ProjectID: 0, ActionPtn: ".*", ResourcePtn: ".*"},
			},
			mockOrgIAMResp: &organization_iam.GetSystemAdminOrganizationPolicyResponse{
				OrganizationPolicies: []*organization_iam.OrganizationPolicy{},
			},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.IsAdminRequest{UserId: 0},
			wantErr: true,
		},
		{
			name:         "NG Invalid IAM DB error",
			input:        &iam.IsAdminRequest{UserId: 1},
			wantErr:      true,
			mockIAMError: gorm.ErrInvalidDB,
		},
		{
			name:    "NG Invalid Org IAM service error",
			input:   &iam.IsAdminRequest{UserId: 1},
			wantErr: true,
			mockIAMResponce: &[]model.Policy{
				{PolicyID: 1, Name: "no-project-policy", ProjectID: 0, ActionPtn: ".*", ResourcePtn: ".*"},
			},
			mockOrgIAMError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			orgIamMock := oimocks.NewOrganizationIAMServiceClient(t)
			svc := IAMService{repository: mock, organizationIamClient: orgIamMock, logger: logging.NewLogger()}

			if c.mockIAMResponce != nil || c.mockIAMError != nil {
				mock.On("GetAdminPolicy", test.RepeatMockAnything(2)...).Return(c.mockIAMResponce, c.mockIAMError).Once()
			}
			if c.mockOrgIAMResp != nil || c.mockOrgIAMError != nil {
				orgIamMock.On("GetSystemAdminOrganizationPolicy", test.RepeatMockAnything(3)...).Return(c.mockOrgIAMResp, c.mockOrgIAMError).Once()
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
