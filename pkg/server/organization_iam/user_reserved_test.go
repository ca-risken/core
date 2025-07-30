package organization_iam

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
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestListOrganizationUserReserved(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.ListOrganizationUserReservedRequest
		want         *organization_iam.ListOrganizationUserReservedResponse
		wantErr      bool
		mockResponse []*model.OrganizationUserReserved
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization_iam.ListOrganizationUserReservedRequest{OrganizationId: 1, UserIdpKey: "key1"},
			want: &organization_iam.ListOrganizationUserReservedResponse{
				UserReserved: []*organization_iam.OrganizationUserReserved{{ReservedId: 1, UserIdpKey: "key1", RoleId: 10, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			},
			mockResponse: []*model.OrganizationUserReserved{{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.ListOrganizationUserReservedRequest{OrganizationId: 1, UserIdpKey: "key1"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}
			if len(c.mockResponse) > 0 || c.mockError != nil {
				mock.On("ListOrganizationUserReserved", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListOrganizationUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrganizationUserReserved(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name            string
		input           *organization_iam.PutOrganizationUserReservedRequest
		want            *organization_iam.PutOrganizationUserReservedResponse
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
			input: &organization_iam.PutOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			want: &organization_iam.PutOrganizationUserReservedResponse{
				UserReserved: &organization_iam.OrganizationUserReserved{
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
			input:       &organization_iam.PutOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			wantErr:     true,
			mockGetResp: &iam.GetUserResponse{User: &iam.User{UserId: 123}},
		},
		{
			name:           "NG role not found",
			input:          &organization_iam.PutOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
			wantErr:        true,
			mockGetResp:    &iam.GetUserResponse{},
			mockGetErr:     gorm.ErrRecordNotFound,
			mockGetRoleErr: gorm.ErrRecordNotFound,
		},
		{
			name:        "NG DB error",
			input:       &organization_iam.PutOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1, UserIdpKey: "key1", RoleId: 10},
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
			repoMock := mocks.NewOrganizationIAMRepository(t)
			iamMock := iam_mocks.NewIAMServiceClient(t)
			svc := OrganizationIAMService{
				repository: repoMock,
				iamClient:  iamMock,
			}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				iamMock.On("GetUser", mock.Anything, &iam.GetUserRequest{UserIdpKey: c.input.UserIdpKey}).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockGetRoleResp != nil || c.mockGetRoleErr != nil {
				repoMock.On("GetOrganizationRole", mock.Anything, c.input.OrganizationId, c.input.RoleId).Return(c.mockGetRoleResp, c.mockGetRoleErr).Once()
			}
			if c.mockResponse != nil || c.mockError != nil {
				expectedInput := &model.OrganizationUserReserved{
					ReservedID: c.input.ReservedId,
					UserIdpKey: c.input.UserIdpKey,
					RoleID:     c.input.RoleId,
				}
				repoMock.On("PutOrganizationUserReserved", mock.Anything, expectedInput).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutOrganizationUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrganizationUserReserved(t *testing.T) {
	cases := []struct {
		name      string
		input     *organization_iam.DeleteOrganizationUserReservedRequest
		wantErr   bool
		mockError error
		mockCall  bool
	}{
		{
			name:     "OK",
			input:    &organization_iam.DeleteOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1},
			mockCall: true,
		},
		{
			name:      "NG DB error",
			input:     &organization_iam.DeleteOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 1},
			wantErr:   true,
			mockCall:  true,
			mockError: gorm.ErrInvalidDB,
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.DeleteOrganizationUserReservedRequest{OrganizationId: 1, ReservedId: 0},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}
			if c.mockCall {
				mock.On("DeleteOrganizationUserReserved", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
			}
			_, err := svc.DeleteOrganizationUserReserved(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
