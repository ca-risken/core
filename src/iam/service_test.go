package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.GetUserRequest
		want         *iam.GetUserResponse
		wantErr      bool
		mockResponce *model.User
		mockError    error
	}{
		{
			name:         "OK",
			input:        &iam.GetUserRequest{UserId: 111, Sub: "sub"},
			want:         &iam.GetUserResponse{User: &iam.User{UserId: 111, Sub: "sub", Name: "nm", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.User{UserID: 111, Sub: "sub", Name: "nm", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &iam.GetUserRequest{UserId: 111, Sub: "sub"},
			want:      &iam.GetUserResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &iam.GetUserRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &iam.GetUserRequest{UserId: 111, Sub: "sub"},
			wantErr:   true,
			mockError: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetUser").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetUser(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

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
				mock.On("GetUserPoicy").Return(c.mockResponce, c.mockError).Once()
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

/*
 * Mock Repository
 */
type mockIAMRepository struct {
	mock.Mock
}

func (m *mockIAMRepository) GetUser(uint32, string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockIAMRepository) GetUserPoicy(uint32) (*[]model.Policy, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Policy), args.Error(1)
}
