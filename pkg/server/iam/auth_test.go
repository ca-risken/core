package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestIsAuthorized(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
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

func TestIsAuthorizedAdmin(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.IsAuthorizedAdminRequest
		want         *iam.IsAuthorizedAdminResponse
		wantErr      bool
		mockResponce *[]model.Policy
		mockError    error
	}{
		{
			name:  "OK Authorized",
			input: &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedAdminResponse{Ok: true},
			mockResponce: &[]model.Policy{
				{PolicyID: 1, Name: "viewer", ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 2, Name: "put for aws", ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:      "OK Record not found",
			input:     &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:      &iam.IsAuthorizedAdminResponse{Ok: false},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid parameter (invalid actionName format)",
			input:   &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding----PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr: true,
		},
		{
			name:      "NG Invalid DB error",
			input:     &iam.IsAuthorizedAdminRequest{UserId: 1, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetAdminPolicy").Return(c.mockResponce, c.mockError).Once()
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
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
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
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.IsAdminRequest
		want         *iam.IsAdminResponse
		wantErr      bool
		mockResponce *[]model.Policy
		mockError    error
	}{
		{
			name:  "OK Admin",
			input: &iam.IsAdminRequest{UserId: 1},
			want:  &iam.IsAdminResponse{Ok: true},
			mockResponce: &[]model.Policy{
				{PolicyID: 1, Name: "no-project-policy", ProjectID: 0, ActionPtn: ".*", ResourcePtn: ".*"},
			},
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
