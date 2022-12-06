package iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestListUser(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.ListUserRequest
		want         *iam.ListUserResponse
		wantErr      bool
		mockResponce *[]model.User
		mockError    error
	}{
		{
			name:  "OK",
			input: &iam.ListUserRequest{ProjectId: 1, Activated: true, Name: "nm"},
			want:  &iam.ListUserResponse{UserId: []uint32{1, 2, 3}},
			mockResponce: &[]model.User{
				{UserID: 1, Sub: "sub", Name: "nm", Activated: true},
				{UserID: 2, Sub: "sub", Name: "nm", Activated: true},
				{UserID: 3, Sub: "sub", Name: "nm", Activated: true},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &iam.ListUserRequest{ProjectId: 1, Activated: true, Name: "nm"},
			want:      &iam.ListUserResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param (Name)",
			input:   &iam.ListUserRequest{ProjectId: 1, Activated: true, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
			wantErr: true,
		},
		{
			name:      "Invalid SQL error",
			input:     &iam.ListUserRequest{ProjectId: 1, Activated: true, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListUser").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListUser(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock, logger: logging.NewLogger()}
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
			name:      "Invalid DB error",
			input:     &iam.GetUserRequest{UserId: 111, Sub: "sub"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
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

func TestPutUser(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock, logger: logging.NewLogger()}
	cases := []struct {
		name           string
		input          *iam.PutUserRequest
		want           *iam.PutUserResponse
		wantErr        bool
		mockGetResp    *model.User
		mockGetErr     error
		mockInsertResp *model.User
		mockInsertErr  error
		mockUpdResp    *model.User
		mockUpdErr     error
	}{
		{
			name:           "OK Insert",
			input:          &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			want:           &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:     gorm.ErrRecordNotFound,
			mockInsertResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			want:        &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update (UserIdpKey in request is empty)",
			input:       &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", Activated: true}},
			want:        &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "saved_uik", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "saved_uik", Activated: true, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "saved_uik", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &iam.PutUserRequest{User: &iam.UserForUpsert{Name: "nm", Activated: true}},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetUserBySub)",
			input:      &iam.PutUserRequest{User: &iam.UserForUpsert{Name: "nm", Activated: true}},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:          "NG DB error(CreateUser)",
			input:         &iam.PutUserRequest{User: &iam.UserForUpsert{Name: "nm", Activated: true}},
			mockGetErr:    gorm.ErrRecordNotFound,
			mockInsertErr: gorm.ErrInvalidTransaction,
			wantErr:       true,
		},
	}
	testUsers := 3
	mock.On("GetActiveUserCount").Return(&testUsers, nil)
	mock.On("AttachAdminRole").Return(nil)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetUserBySub").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutUser").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			if c.mockInsertResp != nil || c.mockInsertErr != nil {
				mock.On("CreateUser").Return(c.mockInsertResp, c.mockInsertErr).Once()
			}
			got, err := svc.PutUser(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
