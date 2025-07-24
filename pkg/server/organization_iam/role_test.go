package organization_iam

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
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const (
	length65string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListOrganizationRole(t *testing.T) {
	cases := []struct {
		name         string
		input        *organization_iam.ListOrganizationRoleRequest
		want         *organization_iam.ListOrganizationRoleResponse
		wantErr      bool
		mockResponce []*model.OrganizationRole
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization_iam.ListOrganizationRoleRequest{OrganizationId: 1, Name: "nm", UserId: 1},
			want:  &organization_iam.ListOrganizationRoleResponse{RoleId: []uint32{1, 2, 3}},
			mockResponce: []*model.OrganizationRole{
				{RoleID: 1, Name: "nm"},
				{RoleID: 2, Name: "nm"},
				{RoleID: 3, Name: "nm"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &organization_iam.ListOrganizationRoleRequest{OrganizationId: 1, Name: "nm", UserId: 1},
			want:      &organization_iam.ListOrganizationRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &organization_iam.ListOrganizationRoleRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid SQL error",
			input:     &organization_iam.ListOrganizationRoleRequest{OrganizationId: 1, Name: "nm"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if len(c.mockResponce) > 0 || c.mockError != nil {
				mock.On("ListOrganizationRole", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrganizationRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.GetOrganizationRoleRequest
		want         *organization_iam.GetOrganizationRoleResponse
		wantErr      bool
		mockResponce *model.OrganizationRole
		mockError    error
	}{
		{
			name:         "OK",
			input:        &organization_iam.GetOrganizationRoleRequest{RoleId: 111, OrganizationId: 123},
			want:         &organization_iam.GetOrganizationRoleResponse{Role: &organization_iam.OrganizationRole{RoleId: 111, Name: "nm", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.OrganizationRole{RoleID: 111, Name: "nm", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &organization_iam.GetOrganizationRoleRequest{RoleId: 111, OrganizationId: 123},
			want:      &organization_iam.GetOrganizationRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.GetOrganizationRoleRequest{},
			wantErr: true,
		},
		{
			name:      "invalid DB error",
			input:     &organization_iam.GetOrganizationRoleRequest{RoleId: 111, OrganizationId: 123},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetOrganizationRole", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrganizationRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *organization_iam.PutOrganizationRoleRequest
		want        *organization_iam.PutOrganizationRoleResponse
		wantErr     bool
		mockGetResp *model.OrganizationRole
		mockGetErr  error
		mockUpdResp *model.OrganizationRole
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &organization_iam.PutOrganizationRoleRequest{Name: "nm", OrganizationId: 123},
			want:        &organization_iam.PutOrganizationRoleResponse{Role: &organization_iam.OrganizationRole{RoleId: 1, Name: "nm", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.OrganizationRole{RoleID: 1, Name: "nm", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &organization_iam.PutOrganizationRoleRequest{Name: "after", OrganizationId: 123},
			want:        &organization_iam.PutOrganizationRoleResponse{Role: &organization_iam.OrganizationRole{RoleId: 1, Name: "after", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.OrganizationRole{RoleID: 1, Name: "before", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.OrganizationRole{RoleID: 1, Name: "after", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &organization_iam.PutOrganizationRoleRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetOrganizationRoleByName)",
			input:      &organization_iam.PutOrganizationRoleRequest{Name: "nm", OrganizationId: 123},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutOrganizationRole)",
			input:      &organization_iam.PutOrganizationRoleRequest{Name: "nm", OrganizationId: 123},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetOrganizationRoleByName", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutOrganizationRole", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrganizationRole(t *testing.T) {
	cases := []struct {
		name     string
		input    *organization_iam.DeleteOrganizationRoleRequest
		wantErr  bool
		mockCall bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &organization_iam.DeleteOrganizationRoleRequest{OrganizationId: 1, RoleId: 1},
			wantErr:  false,
			mockCall: true,
		},
		{
			name:     "NG Invalid parameters",
			input:    &organization_iam.DeleteOrganizationRoleRequest{OrganizationId: 1},
			wantErr:  true,
			mockCall: false,
		},
		{
			name:     "Invalid DB error",
			input:    &organization_iam.DeleteOrganizationRoleRequest{OrganizationId: 1, RoleId: 1},
			wantErr:  true,
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.mockCall {
				mock.On("DeleteOrganizationRole", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachOrganizationRole(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.AttachOrganizationRoleRequest
		want         *organization_iam.AttachOrganizationRoleResponse
		mockResponse *model.OrganizationRole
		mockErr      error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &organization_iam.AttachOrganizationRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			want: &organization_iam.AttachOrganizationRoleResponse{
				Role: &organization_iam.OrganizationRole{
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
			input: &organization_iam.AttachOrganizationRoleRequest{
				OrganizationId: 1,
				RoleId:         0,
				UserId:         1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &organization_iam.AttachOrganizationRoleRequest{
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
			mockDB := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mockDB}
			if c.mockErr != nil || c.mockResponse != nil {
				mockDB.On("AttachOrganizationRole", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockErr).Once()
			}
			got, err := svc.AttachOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrganizationRole(t *testing.T) {
	cases := []struct {
		name     string
		input    *organization_iam.DetachOrganizationRoleRequest
		mockErr  error
		wantErr  bool
		mockCall bool
	}{
		{
			name: "OK",
			input: &organization_iam.DetachOrganizationRoleRequest{
				OrganizationId: 1,
				RoleId:         1,
				UserId:         1,
			},
			mockCall: true,
		},
		{
			name: "NG Invalid param",
			input: &organization_iam.DetachOrganizationRoleRequest{
				OrganizationId: 1,
				RoleId:         0,
				UserId:         1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &organization_iam.DetachOrganizationRoleRequest{
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
			mockDB := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mockDB}
			if c.mockCall {
				mockDB.On("DetachOrganizationRole", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DetachOrganizationRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
		})
	}
}

func TestAttachOrganizationRoleByOrganizationUserReserved(t *testing.T) {
	testUserIdpKey := "uik"
	userID := uint32(100)

	cases := []struct {
		name               string
		mockListResp       *[]db.UserReservedWithOrganizationID
		mockListErr        error
		mockAttachErrIndex int // -1なら全て成功
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
			wantErr:            false,
		},
		{
			name:         "ListOrganizationUserReservedWithOrganizationID error",
			mockListResp: nil,
			mockListErr:  errors.New("list error"),
			wantErr:      true,
		},
		{
			name: "AttachOrganizationRole error",
			mockListResp: &[]db.UserReservedWithOrganizationID{
				{OrganizationID: 1, ReservedID: 1, RoleID: 10},
				{OrganizationID: 2, ReservedID: 2, RoleID: 20},
			},
			mockListErr:        nil,
			mockAttachErrIndex: 1,
			wantErr:            true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repoMock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: repoMock}

			repoMock.On("ListOrganizationUserReservedWithOrganizationID", mock.Anything, testUserIdpKey).Return(c.mockListResp, c.mockListErr).Once()

			if c.mockListErr == nil && c.mockListResp != nil {
				for i, u := range *c.mockListResp {
					if c.mockAttachErrIndex == i {
						repoMock.On("AttachOrganizationRole", mock.Anything, u.OrganizationID, u.RoleID, userID).Return(nil, gorm.ErrInvalidDB).Once()
					} else {
						repoMock.On("AttachOrganizationRole", mock.Anything, u.OrganizationID, u.RoleID, userID).Return(&model.OrganizationRole{}, nil).Once()
					}
				}
			}

			req := &organization_iam.AttachOrganizationRoleByOrganizationUserReservedRequest{
				UserId:     userID,
				UserIdpKey: testUserIdpKey,
			}
			_, err := svc.AttachOrganizationRoleByOrganizationUserReserved(context.Background(), req)
			if c.wantErr && err == nil {
				t.Fatalf("want error but got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
