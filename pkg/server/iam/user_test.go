package iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestListUser(t *testing.T) {
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
			input:     &iam.ListUserRequest{ProjectId: 1, Activated: true, Name: "nm"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock}

			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListUser", test.RepeatMockAnything(7)...).Return(c.mockResponce, c.mockError).Once()
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
	now := time.Now()
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
			want:         &iam.GetUserResponse{User: &iam.User{UserId: 111, Sub: "sub", Name: "nm", Activated: true, IsAdmin: false, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.User{UserID: 111, Sub: "sub", Name: "nm", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
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
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetUser", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
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
	now := time.Now()
	cases := []struct {
		name                     string
		input                    *iam.PutUserRequest
		want                     *iam.PutUserResponse
		wantErr                  bool
		mockGetResp              *model.User
		mockGetErr               error
		mockInsertResp           *model.User
		mockInsertErr            error
		mockUpdResp              *model.User
		mockUpdErr               error
		mockListUserReservedResp *[]db.UserReservedWithProjectID
		mockListUserReservedErr  error

		callGetActiveUserCount bool
	}{
		{
			name:                     "OK Insert",
			input:                    &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			want:                     &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:               gorm.ErrRecordNotFound,
			mockInsertResp:           &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			mockListUserReservedResp: &[]db.UserReservedWithProjectID{},
			callGetActiveUserCount:   true,
		},
		{
			name:                     "OK Insert First User (Auto Admin)",
			input:                    &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			want:                     &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:               gorm.ErrRecordNotFound,
			mockInsertResp:           &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: true, CreatedAt: now, UpdatedAt: now},
			mockListUserReservedResp: &[]db.UserReservedWithProjectID{},
			callGetActiveUserCount:   true,
		},
		{
			name:                   "OK Update",
			input:                  &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			want:                   &iam.PutUserResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp:            &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			mockUpdResp:            &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			callGetActiveUserCount: false,
		},
		{
			name:                   "NG Invalid param",
			input:                  &iam.PutUserRequest{User: &iam.UserForUpsert{Name: "nm", Activated: true}},
			wantErr:                true,
			callGetActiveUserCount: false,
		},
		{
			name:                   "NG DB error(GetUserBySub)",
			input:                  &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			mockGetErr:             gorm.ErrInvalidTransaction,
			wantErr:                true,
			callGetActiveUserCount: false,
		},
		{
			name:                   "NG DB error(CreateUser)",
			input:                  &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			mockGetErr:             gorm.ErrRecordNotFound,
			mockInsertErr:          gorm.ErrInvalidTransaction,
			wantErr:                true,
			callGetActiveUserCount: true,
		},
		{
			name:                    "NG DB error(ListUserReserved)",
			input:                   &iam.PutUserRequest{User: &iam.UserForUpsert{Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true}},
			wantErr:                 true,
			mockGetErr:              gorm.ErrRecordNotFound,
			mockInsertResp:          &model.User{UserID: 1, Sub: "sub", Name: "nm", UserIdpKey: "uik", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			mockListUserReservedErr: gorm.ErrInvalidTransaction,
			callGetActiveUserCount:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.callGetActiveUserCount {
				if c.name == "OK Insert First User (Auto Admin)" {
					mock.On("GetActiveUserCount", test.RepeatMockAnything(2)...).Return(test.Int(0), nil)
				} else {
					mock.On("GetActiveUserCount", test.RepeatMockAnything(2)...).Return(test.Int(3), nil)
				}
			}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetUserBySub", test.RepeatMockAnything(2)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutUser", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			if c.mockInsertResp != nil || c.mockInsertErr != nil {
				mock.On("CreateUser", test.RepeatMockAnything(2)...).Return(c.mockInsertResp, c.mockInsertErr).Once()
			}
			if c.mockListUserReservedResp != nil || c.mockListUserReservedErr != nil {
				mock.On("ListUserReservedWithProjectID", test.RepeatMockAnything(2)...).Return(c.mockListUserReservedResp, c.mockListUserReservedErr).Once()
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

func TestUpdateUserAdmin(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *iam.UpdateUserAdminRequest
		want     *iam.UpdateUserAdminResponse
		wantErr  bool
		mockResp *model.User
		mockErr  error
	}{
		{
			name:     "OK Set Admin True",
			input:    &iam.UpdateUserAdminRequest{UserId: 1, IsAdmin: true},
			want:     &iam.UpdateUserAdminResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", Activated: true, IsAdmin: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", Activated: true, IsAdmin: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:     "OK Set Admin False",
			input:    &iam.UpdateUserAdminRequest{UserId: 1, IsAdmin: false},
			want:     &iam.UpdateUserAdminResponse{User: &iam.User{UserId: 1, Sub: "sub", Name: "nm", Activated: true, IsAdmin: false, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.User{UserID: 1, Sub: "sub", Name: "nm", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.UpdateUserAdminRequest{UserId: 0, IsAdmin: true},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &iam.UpdateUserAdminRequest{UserId: 1, IsAdmin: true},
			mockErr: gorm.ErrInvalidTransaction,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock, logger: logging.NewLogger()}

			if c.mockResp != nil || c.mockErr != nil {
				mock.On("UpdateUserAdmin", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.UpdateUserAdmin(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
