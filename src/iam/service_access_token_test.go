package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"gorm.io/gorm"
)

func TestListAccessToken(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
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
				{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUesrId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
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
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListAccessToken").Return(c.mockResponce, c.mockError).Once()
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
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
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
			want:         &iam.AuthenticateAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUesrId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
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
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetActiveAccessTokenHash").Return(c.mockResponce, c.mockError).Once()
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
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
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
			input:       &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUesrId: 1}},
			want:        &iam.PutAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUesrId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{AccessTokenId: 1, PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUesrId: 1}},
			want:        &iam.PutAccessTokenResponse{AccessToken: &iam.AccessToken{AccessTokenId: 1, Name: "nm", Description: "desc", ProjectId: 1, ExpiredAt: now.Unix(), LastUpdatedUesrId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.AccessToken{AccessTokenID: 1, TokenHash: "xxx", Name: "nm", Description: "desc", ProjectID: 1, ExpiredAt: now, LastUpdatedUserID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &iam.PutAccessTokenRequest{ProjectId: 999, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUesrId: 1}},
			wantErr: true,
		},
		{
			name:       "NG DB error(Get)",
			input:      &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUesrId: 1}},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(Put)",
			input:      &iam.PutAccessTokenRequest{ProjectId: 1, AccessToken: &iam.AccessTokenForUpsert{PlainTextToken: "xxx", Name: "nm", Description: "desc", ProjectId: 1, LastUpdatedUesrId: 1}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetAccessTokenByUniqueKey").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutAccessToken").Return(c.mockUpdResp, c.mockUpdErr).Once()
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
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DeleteAccessTokenRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &iam.DeleteAccessTokenRequest{ProjectId: 1, AccessTokenId: 1},
			wantErr: false,
		},
		{
			name:    "NG Invalid parameters",
			input:   &iam.DeleteAccessTokenRequest{AccessTokenId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid DB error",
			input:   &iam.DeleteAccessTokenRequest{ProjectId: 1, AccessTokenId: 1},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DeleteAccessToken").Return(c.mockErr).Once()
			_, err := svc.DeleteAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachAccessTokenRole(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
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
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("AttachAccessTokenRole").Return(c.mockResp, c.mockErr).Once()
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
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DetachAccessTokenRoleRequest
		wantErr bool
		mockErr error
	}{
		{
			name:  "OK",
			input: &iam.DetachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.DetachAccessTokenRoleRequest{RoleId: 2, AccessTokenId: 3},
			wantErr: true,
		},
		{
			name:    "NG Invalid DB error",
			input:   &iam.DetachAccessTokenRoleRequest{ProjectId: 1, RoleId: 2, AccessTokenId: 3},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DetachAccessTokenRole").Return(c.mockErr).Once()
			_, err := svc.DetachAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
