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
				OrganizationId: 1,
				AccessToken: &organization_iam.OrganizationAccessTokenForUpsert{
					PlainTextToken:    "plain",
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					LastUpdatedUserId: 200,
				},
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
				OrganizationId: 1,
				AccessToken: &organization_iam.OrganizationAccessTokenForUpsert{
					AccessTokenId:     2,
					PlainTextToken:    "plain",
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					LastUpdatedUserId: 200,
				},
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
				OrganizationId: 1,
				AccessToken: &organization_iam.OrganizationAccessTokenForUpsert{
					PlainTextToken:    "plain",
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					LastUpdatedUserId: 200,
				},
			},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name: "NG put error",
			input: &organization_iam.PutOrganizationAccessTokenRequest{
				OrganizationId: 1,
				AccessToken: &organization_iam.OrganizationAccessTokenForUpsert{
					PlainTextToken:    "plain",
					Name:              "token",
					Description:       "desc",
					OrganizationId:    1,
					LastUpdatedUserId: 200,
				},
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
