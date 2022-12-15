package iam

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestListUserReserved(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.ListUserReservedRequest
		want         *iam.ListUserReservedResponse
		wantErr      bool
		mockResponce *[]model.UserReserved
		mockError    error
	}{
		{
			name:  "OK",
			input: &iam.ListUserReservedRequest{ProjectId: 1},
			want: &iam.ListUserReservedResponse{UserReserved: []*iam.UserReserved{
				{ReservedId: 1, RoleId: 1, UserIdpKey: "uik1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{ReservedId: 2, RoleId: 2, UserIdpKey: "uik1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.UserReserved{
				{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1", CreatedAt: now, UpdatedAt: now},
				{ReservedID: 2, RoleID: 2, UserIdpKey: "uik1", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:         "OK empty reponse",
			input:        &iam.ListUserReservedRequest{ProjectId: 1},
			want:         &iam.ListUserReservedResponse{UserReserved: nil},
			mockResponce: &[]model.UserReserved{},
			mockError:    nil,
		},
		{
			name:    "NG Invalid param",
			input:   &iam.ListUserReservedRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid SQL error",
			input:     &iam.ListUserReservedRequest{ProjectId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListUserReserved").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListUserReserved(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutUserReserved(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name            string
		input           *iam.PutUserReservedRequest
		want            *iam.PutUserReservedResponse
		wantErr         bool
		mockGetUserResp *model.User
		mockGetUserErr  error
		isCalledGetRole bool
		mockGetRoleErr  error
		mockUpdResp     *model.UserReserved
		mockUpdErr      error
	}{
		{
			name:            "OK",
			input:           &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{RoleId: 1, UserIdpKey: "uik1"}},
			want:            &iam.PutUserReservedResponse{UserReserved: &iam.UserReserved{ReservedId: 1, RoleId: 1, UserIdpKey: "uik1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetUserErr:  gorm.ErrRecordNotFound,
			isCalledGetRole: true,
			mockGetRoleErr:  nil,
			mockUpdResp:     &model.UserReserved{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid Param",
			input:   &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{}},
			wantErr: true,
		},
		{
			name:           "NG GetUser Error",
			input:          &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{RoleId: 1, UserIdpKey: "uik1"}},
			wantErr:        true,
			mockGetUserErr: errors.New("something error"),
		},
		{
			name:            "NG User is found",
			input:           &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{RoleId: 1, UserIdpKey: "uik1"}},
			wantErr:         true,
			mockGetUserResp: &model.User{UserID: 1, Name: "nm"},
		},
		{
			name:            "NG GetRole Error",
			input:           &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{RoleId: 1, UserIdpKey: "uik1"}},
			wantErr:         true,
			mockGetUserErr:  gorm.ErrRecordNotFound,
			isCalledGetRole: true,
			mockGetRoleErr:  gorm.ErrRecordNotFound,
		},
		{
			name:            "NG PutUserReserved Error",
			input:           &iam.PutUserReservedRequest{UserReserved: &iam.UserReservedForUpsert{RoleId: 1, UserIdpKey: "uik1"}},
			wantErr:         true,
			mockGetUserErr:  gorm.ErrRecordNotFound,
			isCalledGetRole: true,
			mockUpdErr:      errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetUserResp != nil || c.mockGetUserErr != nil {
				mock.On("GetUserByUserIdpKey").Return(c.mockGetUserResp, c.mockGetUserErr).Once()
			}
			if c.isCalledGetRole {
				mock.On("GetRole").Return(&model.Role{}, c.mockGetRoleErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutUserReserved").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutUserReserved(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteUserReserved(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DeleteUserReservedRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &iam.DeleteUserReservedRequest{ReservedId: 1, ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG Invalid parameters",
			input:   &iam.DeleteUserReservedRequest{ReservedId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid DB error",
			input:   &iam.DeleteUserReservedRequest{ReservedId: 1, ProjectId: 1},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DeleteUserReserved").Return(c.mockErr).Once()
			_, err := svc.DeleteUserReserved(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachRoleByUserReserved(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	type args struct {
		userID     uint32
		userIdpKey string
	}
	cases := []struct {
		name               string
		args               args
		wantErr            bool
		mockListResp       *[]db.UserReservedWithProjectID
		mockListErr        error
		isCalledAttachRole bool
		mockAttachRoleErr  error
	}{
		{
			name: "OK",
			args: args{userID: 1, userIdpKey: "uik"},
			mockListResp: &[]db.UserReservedWithProjectID{
				{ProjectID: 1, RoleID: 1},
			},
			isCalledAttachRole: true,
			wantErr:            false,
		},
		{
			name:         "OK No Attach Role",
			args:         args{userID: 1, userIdpKey: "uik"},
			mockListResp: &[]db.UserReservedWithProjectID{},
			wantErr:      false,
		},
		{
			name:        "NG List UserReserved Error",
			args:        args{userID: 1, userIdpKey: "uik"},
			mockListErr: errors.New("something error"),
			wantErr:     true,
		},
		{
			name: "NG Attach Role Error",
			args: args{userID: 1, userIdpKey: "uik"},
			mockListResp: &[]db.UserReservedWithProjectID{
				{ProjectID: 1, RoleID: 1},
			},
			isCalledAttachRole: true,
			mockAttachRoleErr:  errors.New("something error"),
			wantErr:            true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockListResp != nil || c.mockListErr != nil {
				mock.On("ListUserReservedWithProjectID").Return(c.mockListResp, c.mockListErr).Once()
			}
			if c.isCalledAttachRole {
				mock.On("AttachRole").Return(&model.UserRole{}, c.mockAttachRoleErr).Once()

			}
			err := svc.AttachRoleByUserReserved(ctx, c.args.userID, c.args.userIdpKey)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
