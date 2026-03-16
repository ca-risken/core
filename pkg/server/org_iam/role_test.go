package org_iam

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const (
	length65string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListOrgRole(t *testing.T) {
	cases := []struct {
		name         string
		input        *org_iam.ListOrgRoleRequest
		want         *org_iam.ListOrgRoleResponse
		wantErr      bool
		mockResponce []*model.OrganizationRole
		mockError    error
	}{
		{
			name:  "OK",
			input: &org_iam.ListOrgRoleRequest{OrganizationId: 1, Name: "nm", UserId: 1},
			want:  &org_iam.ListOrgRoleResponse{RoleId: []uint32{1, 2, 3}},
			mockResponce: []*model.OrganizationRole{
				{RoleID: 1, Name: "nm"},
				{RoleID: 2, Name: "nm"},
				{RoleID: 3, Name: "nm"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &org_iam.ListOrgRoleRequest{OrganizationId: 1, Name: "nm", UserId: 1},
			want:      &org_iam.ListOrgRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &org_iam.ListOrgRoleRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid SQL error",
			input:     &org_iam.ListOrgRoleRequest{OrganizationId: 1, Name: "nm"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if len(c.mockResponce) > 0 || c.mockError != nil {
				mock.On("ListOrgRole", test.RepeatMockAnything(5)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrgRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.GetOrgRoleRequest
		want         *org_iam.GetOrgRoleResponse
		wantErr      bool
		mockResponce *model.OrganizationRole
		mockError    error
	}{
		{
			name:         "OK",
			input:        &org_iam.GetOrgRoleRequest{RoleId: 111, OrganizationId: 123},
			want:         &org_iam.GetOrgRoleResponse{Role: &org_iam.OrgRole{RoleId: 111, Name: "nm", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.OrganizationRole{RoleID: 111, Name: "nm", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &org_iam.GetOrgRoleRequest{RoleId: 111, OrganizationId: 123},
			want:      &org_iam.GetOrgRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &org_iam.GetOrgRoleRequest{},
			wantErr: true,
		},
		{
			name:      "invalid DB error",
			input:     &org_iam.GetOrgRoleRequest{RoleId: 111, OrganizationId: 123},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetOrgRole", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrgRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *org_iam.PutOrgRoleRequest
		want        *org_iam.PutOrgRoleResponse
		wantErr     bool
		mockGetResp *model.OrganizationRole
		mockGetErr  error
		mockUpdResp *model.OrganizationRole
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &org_iam.PutOrgRoleRequest{Name: "nm", OrganizationId: 123},
			want:        &org_iam.PutOrgRoleResponse{Role: &org_iam.OrgRole{RoleId: 1, Name: "nm", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.OrganizationRole{RoleID: 1, Name: "nm", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &org_iam.PutOrgRoleRequest{Name: "after", OrganizationId: 123},
			want:        &org_iam.PutOrgRoleResponse{Role: &org_iam.OrgRole{RoleId: 1, Name: "after", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.OrganizationRole{RoleID: 1, Name: "before", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.OrganizationRole{RoleID: 1, Name: "after", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &org_iam.PutOrgRoleRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetOrgRoleByName)",
			input:      &org_iam.PutOrgRoleRequest{Name: "nm", OrganizationId: 123},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutOrgRole)",
			input:      &org_iam.PutOrgRoleRequest{Name: "nm", OrganizationId: 123},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetOrgRoleByName", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutOrgRole", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrgRole(t *testing.T) {
	cases := []struct {
		name     string
		input    *org_iam.DeleteOrgRoleRequest
		wantErr  bool
		mockCall bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &org_iam.DeleteOrgRoleRequest{OrganizationId: 1, RoleId: 1},
			wantErr:  false,
			mockCall: true,
		},
		{
			name:     "NG Invalid parameters",
			input:    &org_iam.DeleteOrgRoleRequest{OrganizationId: 1},
			wantErr:  true,
			mockCall: false,
		},
		{
			name:     "Invalid DB error",
			input:    &org_iam.DeleteOrgRoleRequest{OrganizationId: 1, RoleId: 1},
			wantErr:  true,
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.mockCall {
				mock.On("DeleteOrgRole", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachOrgRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.AttachOrgRoleRequest
		want         *org_iam.AttachOrgRoleResponse
		mockResponse *model.OrganizationRole
		mockErr      error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &org_iam.AttachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			want: &org_iam.AttachOrgRoleResponse{
				Role: &org_iam.OrgRole{
					RoleId:         1,
					OrganizationId: 1,
					Name:           "test-role",
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationRole{
				RoleID:         1,
				OrganizationID: 1,
				Name:           "test-role",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "NG Invalid param",
			input: &org_iam.AttachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         0,
				UserId:         1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &org_iam.AttachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockDB := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mockDB}
			if c.mockErr != nil || c.mockResponse != nil {
				mockDB.On("AttachOrgRole", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockErr).Once()
			}
			got, err := svc.AttachOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrgRole(t *testing.T) {
	cases := []struct {
		name     string
		input    *org_iam.DetachOrgRoleRequest
		mockErr  error
		wantErr  bool
		mockCall bool
	}{
		{
			name: "OK",
			input: &org_iam.DetachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			mockCall: true,
		},
		{
			name: "NG Invalid param",
			input: &org_iam.DetachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         0,
				UserId:         1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &org_iam.DetachOrgRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
			wantErr:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockDB := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mockDB}
			if c.mockCall {
				mockDB.On("DetachOrgRole", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DetachOrgRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
		})
	}
}

func TestAttachOrgRoleByOrgUserReserved(t *testing.T) {
	testUserIdpKey := "uik"
	userID := uint32(100)

	cases := []struct {
		name               string
		mockListResp       *[]db.UserReservedWithOrganizationID
		mockListErr        error
		mockAttachErrIndex int // -1なら全て成功
		mockDeleteErrIndex int // -1なら全て成功
		wantErr            bool
	}{
		{
			name: "OK",
			mockListResp: &[]db.UserReservedWithOrganizationID{
				{OrganizationID: 1, ReservedID: 1, RoleID: 10},
				{OrganizationID: 2, ReservedID: 2, RoleID: 20},
			},
			mockListErr:        nil,
			mockAttachErrIndex: -1,
			mockDeleteErrIndex: -1,
			wantErr:            false,
		},
		{
			name:         "ListOrganizationUserReservedWithOrganizationID error",
			mockListResp: nil,
			mockListErr:  errors.New("list error"),
			wantErr:      true,
		},
		{
			name: "AttachOrgRole error",
			mockListResp: &[]db.UserReservedWithOrganizationID{
				{OrganizationID: 1, ReservedID: 1, RoleID: 10},
				{OrganizationID: 2, ReservedID: 2, RoleID: 20},
			},
			mockListErr:        nil,
			mockAttachErrIndex: 1,
			mockDeleteErrIndex: -1,
			wantErr:            true,
		},
		{
			name: "DeleteOrgUserReserved error",
			mockListResp: &[]db.UserReservedWithOrganizationID{
				{OrganizationID: 1, ReservedID: 1, RoleID: 10},
			},
			mockListErr:        nil,
			mockAttachErrIndex: -1,
			mockDeleteErrIndex: 0,
			wantErr:            true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repoMock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: repoMock}

			repoMock.On("ListOrgUserReservedWithOrganizationID", mock.Anything, testUserIdpKey).Return(c.mockListResp, c.mockListErr).Once()

			if c.mockListErr == nil && c.mockListResp != nil {
				for i, u := range *c.mockListResp {
					if c.mockAttachErrIndex == i {
						repoMock.On("AttachOrgRole", mock.Anything, u.OrganizationID, u.RoleID, userID).Return(nil, gorm.ErrInvalidDB).Once()
					} else {
						repoMock.On("AttachOrgRole", mock.Anything, u.OrganizationID, u.RoleID, userID).Return(&model.OrganizationRole{}, nil).Once()
					}

					if c.mockAttachErrIndex != i {
						if c.mockDeleteErrIndex == i {
							repoMock.On("DeleteOrgUserReserved", mock.Anything, u.OrganizationID, u.ReservedID).Return(errors.New("delete error")).Once()
						} else {
							repoMock.On("DeleteOrgUserReserved", mock.Anything, u.OrganizationID, u.ReservedID).Return(nil).Once()
						}
					}
				}
			}

			req := &org_iam.AttachOrgRoleByOrgUserReservedRequest{
				UserId:     userID,
				UserIdpKey: testUserIdpKey,
			}
			_, err := svc.AttachOrgRoleByOrgUserReserved(context.Background(), req)
			if c.wantErr && err == nil {
				t.Fatalf("want error but got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
