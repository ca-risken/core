package org_iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/org_iam"
	"gorm.io/gorm"
)

func TestListOrgAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.ListOrgAccessTokenRequest
		want         *org_iam.ListOrgAccessTokenResponse
		wantErr      bool
		mockResponse *[]model.OrganizationAccessToken
		mockError    error
	}{
		{
			name:  "OK",
			input: &org_iam.ListOrgAccessTokenRequest{OrganizationId: 1, Name: "token", AccessTokenId: 10},
			want: &org_iam.ListOrgAccessTokenResponse{AccessToken: []*org_iam.OrgAccessToken{
				{AccessTokenId: 10, Name: "token", Description: "desc", OrganizationId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 100, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &[]model.OrganizationAccessToken{
				{AccessTokenID: 10, TokenHash: "hash", Name: "token", Description: "desc", OrganizationID: 1, ExpiredAt: now, LastUpdatedUserID: 100, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty response",
			input:     &org_iam.ListOrgAccessTokenRequest{OrganizationId: 1, Name: "token", AccessTokenId: 10},
			want:      &org_iam.ListOrgAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid request",
			input:   &org_iam.ListOrgAccessTokenRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &org_iam.ListOrgAccessTokenRequest{OrganizationId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockResponse != nil || c.mockError != nil {
				mock.On("ListOrgAccessToken", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListOrgAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrgAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *org_iam.PutOrgAccessTokenRequest
		want        *org_iam.PutOrgAccessTokenResponse
		wantErr     bool
		mockGetResp *model.OrganizationAccessToken
		mockGetErr  error
		mockPutResp *model.OrganizationAccessToken
		mockPutErr  error
	}{
		{
			name: "OK insert",
			input: &org_iam.PutOrgAccessTokenRequest{
				OrganizationId:    1,
				PlainTextToken:    "plain",
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			want: &org_iam.PutOrgAccessTokenResponse{
				AccessToken: &org_iam.OrgAccessToken{
					AccessTokenId:     2,
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					ExpiredAt:         now.Unix(),
					LastUpdatedUserId: 200,
					CreatedAt:         now.Unix(),
					UpdatedAt:         now.Unix(),
				},
			},
			mockGetErr: gorm.ErrRecordNotFound,
			mockPutResp: &model.OrganizationAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrganizationID:    1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name: "OK update",
			input: &org_iam.PutOrgAccessTokenRequest{
				OrganizationId:    1,
				AccessTokenId:     2,
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			want: &org_iam.PutOrgAccessTokenResponse{
				AccessToken: &org_iam.OrgAccessToken{
					AccessTokenId:     2,
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					ExpiredAt:         now.Unix(),
					LastUpdatedUserId: 200,
					CreatedAt:         now.Unix(),
					UpdatedAt:         now.Unix(),
				},
			},
			mockGetResp: &model.OrganizationAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrganizationID:    1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			mockPutResp: &model.OrganizationAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrganizationID:    1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name: "NG validation error",
			input: &org_iam.PutOrgAccessTokenRequest{
				OrganizationId: 1,
			},
			wantErr: true,
		},
		{
			name: "NG get error",
			input: &org_iam.PutOrgAccessTokenRequest{
				OrganizationId:    1,
				PlainTextToken:    "plain",
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name: "NG put error",
			input: &org_iam.PutOrgAccessTokenRequest{
				OrganizationId:    1,
				PlainTextToken:    "plain",
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			wantErr:    true,
			mockGetErr: gorm.ErrRecordNotFound,
			mockPutErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetOrgAccessTokenByUniqueKey", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockPutResp != nil || c.mockPutErr != nil {
				mock.On("PutOrgAccessToken", test.RepeatMockAnything(2)...).Return(c.mockPutResp, c.mockPutErr).Once()
			}
			got, err := svc.PutOrgAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrgAccessToken(t *testing.T) {
	cases := []struct {
		name    string
		input   *org_iam.DeleteOrgAccessTokenRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &org_iam.DeleteOrgAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10},
		},
		{
			name:    "NG validation error",
			input:   &org_iam.DeleteOrgAccessTokenRequest{},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if !c.wantErr {
				mock.On("DeleteOrgAccessToken", test.RepeatMockAnything(3)...).Return(nil).Once()
			}
			_, err := svc.DeleteOrgAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAuthenticateOrgAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name      string
		input     *org_iam.AuthenticateOrgAccessTokenRequest
		want      *org_iam.AuthenticateOrgAccessTokenResponse
		wantErr   bool
		mockResp  *model.OrganizationAccessToken
		mockError error
	}{
		{
			name:  "OK",
			input: &org_iam.AuthenticateOrgAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			want: &org_iam.AuthenticateOrgAccessTokenResponse{
				AccessToken: &org_iam.OrgAccessToken{
					AccessTokenId:     10,
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					ExpiredAt:         now.Unix(),
					LastUpdatedUserId: 100,
					CreatedAt:         now.Unix(),
					UpdatedAt:         now.Unix(),
				},
			},
			mockResp: &model.OrganizationAccessToken{
				AccessTokenID:     10,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrganizationID:    1,
				ExpiredAt:         now,
				LastUpdatedUserID: 100,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name:      "OK record not found",
			input:     &org_iam.AuthenticateOrgAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			want:      &org_iam.AuthenticateOrgAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &org_iam.AuthenticateOrgAccessTokenRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &org_iam.AuthenticateOrgAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockResp != nil || c.mockError != nil {
				mock.On("GetActiveOrgAccessTokenHash", test.RepeatMockAnything(4)...).Return(c.mockResp, c.mockError).Once()
			}
			got, err := svc.AuthenticateOrgAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestAttachOrgAccessTokenRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name      string
		input     *org_iam.AttachOrgAccessTokenRoleRequest
		want      *org_iam.AttachOrgAccessTokenRoleResponse
		wantErr   bool
		mockResp  *model.OrganizationAccessTokenRole
		mockError error
	}{
		{
			name:  "OK",
			input: &org_iam.AttachOrgAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			want: &org_iam.AttachOrgAccessTokenRoleResponse{
				AccessTokenRole: &org_iam.OrgAccessTokenRole{
					AccessTokenId: 2,
					RoleId:        3,
					CreatedAt:     now.Unix(),
					UpdatedAt:     now.Unix(),
				},
			},
			mockResp: &model.OrganizationAccessTokenRole{
				AccessTokenID: 2,
				RoleID:        3,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
		{
			name:    "NG validation error",
			input:   &org_iam.AttachOrgAccessTokenRoleRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &org_iam.AttachOrgAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockResp != nil || c.mockError != nil {
				mock.On("AttachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockResp, c.mockError).Once()
			}
			got, err := svc.AttachOrgAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrgAccessTokenRole(t *testing.T) {
	cases := []struct {
		name      string
		input     *org_iam.DetachOrgAccessTokenRoleRequest
		wantErr   bool
		mockError error
	}{
		{
			name:  "OK",
			input: &org_iam.DetachOrgAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
		},
		{
			name:    "NG validation error",
			input:   &org_iam.DetachOrgAccessTokenRoleRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &org_iam.DetachOrgAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockError == nil && !c.wantErr {
				mock.On("DetachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(nil).Once()
			}
			if c.mockError != nil {
				mock.On("DetachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockError).Once()
			}
			_, err := svc.DetachOrgAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
