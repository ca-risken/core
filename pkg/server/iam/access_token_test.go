package iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestListAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *iam.ListAccessTokenRequest
		want         *iam.ListAccessTokenResponse
		wantErr      bool
		mockResponce *[]model.AccessToken
		mockError    error
	}{
		{
			name:  "OK",
			input: &iam.ListAccessTokenRequest{ProjectId: 1, Name: "nm", AccessTokenId: 1},
			want: &iam.ListAccessTokenResponse{AccessToken: []*iam.AccessToken{
				{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.AccessToken{
				{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &iam.ListAccessTokenRequest{ProjectId: 1, Name: "nm", AccessTokenId: 1},
			want:      &iam.ListAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &iam.ListAccessTokenRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:      "NG Invalid DB error",
			input:     &iam.ListAccessTokenRequest{ProjectId: 1, Name: "nm", AccessTokenId: 1},
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
				mock.On("ListAccessToken", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestAuthenticateAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *iam.AuthenticateAccessTokenRequest
		want         *iam.AuthenticateAccessTokenResponse
		wantErr      bool
		mockResponce *model.AccessToken
		mockError    error
	}{
		{
			name:         "OK",
			input:        &iam.AuthenticateAccessTokenRequest{ProjectId: 1, AccessTokenId: 1, PlainTextToken: "xxx"},
			want:         &iam.AuthenticateAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AccessToken{AccessTokenID: 1, Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &iam.AuthenticateAccessTokenRequest{ProjectId: 1, AccessTokenId: 1, PlainTextToken: "xxx"},
			want:      &iam.AuthenticateAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &iam.AuthenticateAccessTokenRequest{},
			wantErr: true,
		},
		{
			name:      "NG Invalid DB error",
			input:     &iam.AuthenticateAccessTokenRequest{ProjectId: 1, AccessTokenId: 1, PlainTextToken: "xxx"},
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
				mock.On("GetActiveAccessTokenHash", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.AuthenticateAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *iam.PutAccessTokenRequest
		want        *iam.PutAccessTokenResponse
		wantErr     bool
		mockGetResp *model.AccessToken
		mockGetErr  error
		mockUpdResp *model.AccessToken
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUserId: 1}},
			want:        &iam.PutAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{AccessTokenId: 1, PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUserId: 1}},
			want:        &iam.PutAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &iam.PutAccessTokenRequest{ProjectId: 999, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUserId: 1}},
			wantErr: true,
		},
		{
			name:       "NG DB error(Get)",
			input:      &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUserId: 1}},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(Put)",
			input:      &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUserId: 1}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetAccessTokenByUniqueKey", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutAccessToken", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAccessToken(t *testing.T) {
	cases := []struct {
		name     string
		input    *iam.DeleteAccessTokenRequest
		wantErr  bool
		callMock bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &iam.DeleteAccessTokenRequest{ProjectId: 1, AccessTokenId: 1},
			wantErr:  false,
			callMock: true,
		},
		{
			name:     "NG Invalid parameters",
			input:    &iam.DeleteAccessTokenRequest{AccessTokenId: 1},
			wantErr:  true,
			callMock: false,
		},
		{
			name:     "NG Invalid DB error",
			input:    &iam.DeleteAccessTokenRequest{ProjectId: 1, AccessTokenId: 1},
			wantErr:  true,
			callMock: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock}

			if c.callMock {
				mock.On("DeleteAccessToken", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachAccessTokenRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *iam.AttachAccessTokenRoleRequest
		want     *iam.AttachAccessTokenRoleResponse
		wantErr  bool
		mockResp *model.AccessTokenRole
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &iam.AttachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
			want:     &iam.AttachAccessTokenRoleResponse{AccessTokenRole: &iam.AccessTokenRole{RoleId: 2, AccessTokenId: 3, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.AccessTokenRole{RoleID: 2, AccessTokenID: 3, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.AttachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2},
			wantErr: true,
		},
		{
			name:    "NG Invalid DB error",
			input:   &iam.AttachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock}

			if c.mockResp != nil || c.mockErr != nil {
				mock.On("AttachAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.AttachAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachAccessTokenRole(t *testing.T) {
	cases := []struct {
		name     string
		input    *iam.DetachAccessTokenRoleRequest
		wantErr  bool
		mockCall bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &iam.DetachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
			mockCall: true,
		},
		{
			name:     "NG Invalid parameter",
			input:    &iam.DetachAccessTokenRoleRequest{RoleId: 2, AccessTokenId: 3},
			mockCall: false,
			wantErr:  true,
		},
		{
			name:     "NG Invalid DB error",
			input:    &iam.DetachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
			wantErr:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewIAMRepository(t)
			svc := IAMService{repository: mock}

			if c.mockCall {
				mock.On("DetachAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DetachAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
