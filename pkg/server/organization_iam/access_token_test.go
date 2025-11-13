package organization_iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/organization_iam"
	"gorm.io/gorm"
)

func TestListOrganizationAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.ListOrganizationAccessTokenRequest
		want         *organization_iam.ListOrganizationAccessTokenResponse
		wantErr      bool
		mockResponse *[]model.OrgAccessToken
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization_iam.ListOrganizationAccessTokenRequest{OrganizationId: 1, Name: "token", AccessTokenId: 10},
			want: &organization_iam.ListOrganizationAccessTokenResponse{AccessToken: []*organization_iam.OrganizationAccessToken{
				{AccessTokenId: 10, Name: "token", Description: "desc", OrganizationId: 1, ExpiredAt: now.Unix(), LastUpdatedUserId: 100, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &[]model.OrgAccessToken{
				{AccessTokenID: 10, TokenHash: "hash", Name: "token", Description: "desc", OrgID: 1, ExpiredAt: now, LastUpdatedUserID: 100, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty response",
			input:     &organization_iam.ListOrganizationAccessTokenRequest{OrganizationId: 1, Name: "token", AccessTokenId: 10},
			want:      &organization_iam.ListOrganizationAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid request",
			input:   &organization_iam.ListOrganizationAccessTokenRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.ListOrganizationAccessTokenRequest{OrganizationId: 1},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockResponse != nil || c.mockError != nil {
				mock.On("ListOrgAccessToken", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListOrganizationAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrganizationAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *organization_iam.PutOrganizationAccessTokenRequest
		want        *organization_iam.PutOrganizationAccessTokenResponse
		wantErr     bool
		mockGetResp *model.OrgAccessToken
		mockGetErr  error
		mockPutResp *model.OrgAccessToken
		mockPutErr  error
	}{
		{
			name: "OK insert",
			input: &organization_iam.PutOrganizationAccessTokenRequest{
				OrganizationId:    1,
				PlainTextToken:    "plain",
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			want: &organization_iam.PutOrganizationAccessTokenResponse{
				AccessToken: &organization_iam.OrganizationAccessToken{
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
			mockPutResp: &model.OrgAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrgID:             1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name: "OK update",
			input: &organization_iam.PutOrganizationAccessTokenRequest{
				OrganizationId:    1,
				AccessTokenId:     2,
				Name:              "token",
				Description:       "desc",
				LastUpdatedUserId: 200,
			},
			want: &organization_iam.PutOrganizationAccessTokenResponse{
				AccessToken: &organization_iam.OrganizationAccessToken{
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
			mockGetResp: &model.OrgAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrgID:             1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			mockPutResp: &model.OrgAccessToken{
				AccessTokenID:     2,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrgID:             1,
				ExpiredAt:         now,
				LastUpdatedUserID: 200,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name: "NG validation error",
			input: &organization_iam.PutOrganizationAccessTokenRequest{
				OrganizationId: 1,
			},
			wantErr: true,
		},
		{
			name: "NG get error",
			input: &organization_iam.PutOrganizationAccessTokenRequest{
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
			input: &organization_iam.PutOrganizationAccessTokenRequest{
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
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetOrgAccessTokenByUniqueKey", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockPutResp != nil || c.mockPutErr != nil {
				mock.On("PutOrgAccessToken", test.RepeatMockAnything(2)...).Return(c.mockPutResp, c.mockPutErr).Once()
			}
			got, err := svc.PutOrganizationAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrganizationAccessToken(t *testing.T) {
	cases := []struct {
		name    string
		input   *organization_iam.DeleteOrganizationAccessTokenRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &organization_iam.DeleteOrganizationAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10},
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.DeleteOrganizationAccessTokenRequest{},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if !c.wantErr {
				mock.On("DeleteOrgAccessToken", test.RepeatMockAnything(3)...).Return(nil).Once()
			}
			_, err := svc.DeleteOrganizationAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAuthenticateOrganizationAccessToken(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name      string
		input     *organization_iam.AuthenticateOrganizationAccessTokenRequest
		want      *organization_iam.AuthenticateOrganizationAccessTokenResponse
		wantErr   bool
		mockResp  *model.OrgAccessToken
		mockError error
	}{
		{
			name:  "OK",
			input: &organization_iam.AuthenticateOrganizationAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			want: &organization_iam.AuthenticateOrganizationAccessTokenResponse{
				AccessToken: &organization_iam.OrganizationAccessToken{
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
			mockResp: &model.OrgAccessToken{
				AccessTokenID:     10,
				TokenHash:         "hash",
				Name:              "token",
				Description:       "desc",
				OrgID:             1,
				ExpiredAt:         now,
				LastUpdatedUserID: 100,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name:      "OK record not found",
			input:     &organization_iam.AuthenticateOrganizationAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			want:      &organization_iam.AuthenticateOrganizationAccessTokenResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.AuthenticateOrganizationAccessTokenRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.AuthenticateOrganizationAccessTokenRequest{OrganizationId: 1, AccessTokenId: 10, PlainTextToken: "plain"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockResp != nil || c.mockError != nil {
				mock.On("GetActiveOrgAccessTokenHash", test.RepeatMockAnything(4)...).Return(c.mockResp, c.mockError).Once()
			}
			got, err := svc.AuthenticateOrganizationAccessToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestAttachOrganizationAccessTokenRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name      string
		input     *organization_iam.AttachOrganizationAccessTokenRoleRequest
		want      *organization_iam.AttachOrganizationAccessTokenRoleResponse
		wantErr   bool
		mockResp  *model.OrgAccessTokenRole
		mockError error
	}{
		{
			name:  "OK",
			input: &organization_iam.AttachOrganizationAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			want: &organization_iam.AttachOrganizationAccessTokenRoleResponse{
				AccessTokenRole: &organization_iam.OrganizationAccessTokenRole{
					AccessTokenId: 2,
					RoleId:        3,
					CreatedAt:     now.Unix(),
					UpdatedAt:     now.Unix(),
				},
			},
			mockResp: &model.OrgAccessTokenRole{
				AccessTokenID: 2,
				RoleID:        3,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.AttachOrganizationAccessTokenRoleRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.AttachOrganizationAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockResp != nil || c.mockError != nil {
				mock.On("AttachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockResp, c.mockError).Once()
			}
			got, err := svc.AttachOrganizationAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrganizationAccessTokenRole(t *testing.T) {
	cases := []struct {
		name      string
		input     *organization_iam.DetachOrganizationAccessTokenRoleRequest
		wantErr   bool
		mockError error
	}{
		{
			name:  "OK",
			input: &organization_iam.DetachOrganizationAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.DetachOrganizationAccessTokenRoleRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.DetachOrganizationAccessTokenRoleRequest{OrganizationId: 1, AccessTokenId: 2, RoleId: 3},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockError == nil && !c.wantErr {
				mock.On("DetachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(nil).Once()
			}
			if c.mockError != nil {
				mock.On("DetachOrgAccessTokenRole", test.RepeatMockAnything(4)...).Return(c.mockError).Once()
			}
			_, err := svc.DetachOrganizationAccessTokenRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
