package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
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
			name:  "OK Unauthorized (Not allow project)",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 9999, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
			mockResponce: &[]model.Policy{
				{PolicyID: 101, Name: "viewer", ProjectID: 1001, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 102, Name: "put for aws", ProjectID: 1001, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:  "OK Unauthorized (Not allow action)",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/DeleteFinding", ResourceName: "aws:guardduty/ec2-instance-id"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
			mockResponce: &[]model.Policy{
				{PolicyID: 101, Name: "viewer", ProjectID: 1001, ActionPtn: "finding/(Get|List|Describe)", ResourcePtn: ".*"},
				{PolicyID: 102, Name: "put for aws", ProjectID: 1001, ActionPtn: "finding/Put.*", ResourcePtn: "aws:.*"},
			},
		},
		{
			name:  "OK Unauthorized (Not allow resource)",
			input: &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			want:  &iam.IsAuthorizedResponse{Ok: false},
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
			name:      "NG DB error",
			input:     &iam.IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "github:code-scan/repository-name"},
			wantErr:   true,
			mockError: gorm.ErrInvalidSQL,
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
			name:      "NG DB error",
			input:     &iam.IsAdminRequest{UserId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidSQL,
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

func (m *mockIAMRepository) ListUser(activated bool, projectID uint32, name string, userID uint32) (*[]model.User, error) {
	args := m.Called()
	return args.Get(0).(*[]model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUser(uint32, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUserBySub(string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUserPolicy(uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *mockIAMRepository) PutUser(*model.User) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) ListRole(uint32, string, uint32) (*[]model.Role, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Role), args.Error(1)
}
func (m *mockIAMRepository) GetRole(uint32, uint32) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) GetRoleByName(uint32, string) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) PutRole(r *model.Role) (*model.Role, error) {
	args := m.Called()
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *mockIAMRepository) DeleteRole(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) AttachRole(uint32, uint32, uint32) (*model.UserRole, error) {
	args := m.Called()
	return args.Get(0).(*model.UserRole), args.Error(1)
}
func (m *mockIAMRepository) DetachRole(uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) ListPolicy(uint32, string, uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
func (m *mockIAMRepository) GetPolicy(uint32, uint32) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) GetPolicyByName(uint32, string) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) PutPolicy(*model.Policy) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
func (m *mockIAMRepository) DeletePolicy(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockIAMRepository) AttachPolicy(uint32, uint32, uint32) (*model.RolePolicy, error) {
	args := m.Called()
	return args.Get(0).(*model.RolePolicy), args.Error(1)
}
func (m *mockIAMRepository) DetachPolicy(uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockIAMRepository) GetAdminPolicy(uint32) (*model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*model.Policy), args.Error(1)
}
