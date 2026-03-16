package org_iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	iam_mocks "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestListOrgUserReserved(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.ListOrgUserReservedRequest
		want         *org_iam.ListOrgUserReservedResponse
		wantErr      bool
		mockResponse []*model.OrganizationUserReserved
		mockError    error
	}{
		{
			name:  "OK",
			input: &org_iam.ListOrgUserReservedRequest{OrganizationId: 1, UserIdpKey: "key1"},
			want: &org_iam.ListOrgUserReservedResponse{
				UserReserved: []*org_iam.OrgUserReserved{{ReservedId: 1, UserIdpKey: "key1", RoleId: 10, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			},
			mockResponse: []*model.OrganizationUserReserved{{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG DB error",
			input:     &org_iam.ListOrgUserReservedRequest{OrganizationId: 1, UserIdpKey: "key1"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}
			if len(c.mockResponse) > 0 || c.mockError != nil {
				mock.On("ListOrgUserReserved", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListOrgUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrgUserReserved(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name            string
		input           *org_iam.PutOrgUserReservedRequest
		want            *org_iam.PutOrgUserReservedResponse
		wantErr         bool
		mockGetResp     *iam.GetUserResponse
		mockGetErr      error
		mockGetRoleResp *model.OrganizationRole
		mockGetRoleErr  error
		mockResponse    *model.OrganizationUserReserved
		mockError       error
	}{
		{
			name:  "OK",
			input: &org_iam.PutOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			want: &org_iam.PutOrgUserReservedResponse{
				UserReserved: &org_iam.OrgUserReserved{
					ReservedId: 1,
					UserIdpKey: "key1",
					RoleId:     10,
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
			},
			mockGetResp: &iam.GetUserResponse{},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockGetRoleResp: &model.OrganizationRole{
				RoleID: 10,
			},
			mockGetRoleErr: nil,
			mockResponse: &model.OrganizationUserReserved{
				ReservedID: 1,
				UserIdpKey: "key1",
				RoleID:     10,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		{
			name:        "NG user already exists",
			input:       &org_iam.PutOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			wantErr:     true,
			mockGetResp: &iam.GetUserResponse{User: &iam.User{UserId: 123}},
		},
		{
			name:           "NG role not found",
			input:          &org_iam.PutOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			wantErr:        true,
			mockGetResp:    &iam.GetUserResponse{},
			mockGetErr:     gorm.ErrRecordNotFound,
			mockGetRoleErr: gorm.ErrRecordNotFound,
		},
		{
			name:        "NG DB error",
			input:       &org_iam.PutOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			wantErr:     true,
			mockGetResp: &iam.GetUserResponse{},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockGetRoleResp: &model.OrganizationRole{
				RoleID: 10,
			},
			mockResponse: nil,
			mockError:    gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repoMock := mocks.NewOrgIAMRepository(t)
			iamMock := iam_mocks.NewIAMServiceClient(t)
			svc := OrgIAMService{
				repository: repoMock,
				iamClient:  iamMock,
			}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				iamMock.On("GetUserByUserIdpKey", mock.Anything, &iam.GetUserByUserIdpKeyRequest{UserIdpKey: c.input.UserIdpKey}).Return(&iam.GetUserByUserIdpKeyResponse{User: c.mockGetResp.User}, c.mockGetErr).Once()
			}
			if c.mockGetRoleResp != nil || c.mockGetRoleErr != nil {
				repoMock.On("GetOrgRole", mock.Anything, c.input.OrganizationId, c.input.RoleId).Return(c.mockGetRoleResp, c.mockGetRoleErr).Once()
			}
			if c.mockResponse != nil || c.mockError != nil {
				expectedInput := &model.OrganizationUserReserved{
					ReservedID: c.input.ReservedId,
					UserIdpKey: c.input.UserIdpKey,
					RoleID:     c.input.RoleId,
				}
				repoMock.On("PutOrgUserReserved", mock.Anything, expectedInput).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutOrgUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrgUserReserved(t *testing.T) {
	cases := []struct {
		name      string
		input     *org_iam.DeleteOrgUserReservedRequest
		wantErr   bool
		mockError error
		mockCall  bool
	}{
		{
			name:     "OK",
			input:    &org_iam.DeleteOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1},
			mockCall: true,
		},
		{
			name:      "NG DB error",
			input:     &org_iam.DeleteOrgUserReservedRequest{OrganizationId: 1, ReservedId: 1},
			wantErr:   true,
			mockCall:  true,
			mockError: gorm.ErrInvalidDB,
		},
		{
			name:    "NG validation error",
			input:   &org_iam.DeleteOrgUserReservedRequest{OrganizationId: 1, ReservedId: 0},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}
			if c.mockCall {
				mock.On("DeleteOrgUserReserved", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
			}
			_, err := svc.DeleteOrgUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
