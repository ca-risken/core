package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/jinzhu/gorm"
)

func TestListRole(t *testing.T) {
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.ListRoleRequest
		want         *iam.ListRoleResponse
		wantErr      bool
		mockResponce *[]model.Role
		mockError    error
	}{
		{
			name:  "OK",
			input: &iam.ListRoleRequest{ProjectId: 1, Name: "nm", UserId: 1},
			want:  &iam.ListRoleResponse{RoleId: []uint32{1, 2, 3}},
			mockResponce: &[]model.Role{
				{RoleID: 1, Name: "nm"},
				{RoleID: 2, Name: "nm"},
				{RoleID: 3, Name: "nm"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &iam.ListRoleRequest{ProjectId: 1, Name: "nm", UserId: 1},
			want:      &iam.ListRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &iam.ListRoleRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:      "NG SQL error",
			input:     &iam.ListRoleRequest{ProjectId: 1, Name: "nm"},
			wantErr:   true,
			mockError: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListRole").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetRole(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.GetRoleRequest
		want         *iam.GetRoleResponse
		wantErr      bool
		mockResponce *model.Role
		mockError    error
	}{
		{
			name:         "OK",
			input:        &iam.GetRoleRequest{RoleId: 111, ProjectId: 123},
			want:         &iam.GetRoleResponse{Role: &iam.Role{RoleId: 111, Name: "nm", ProjectId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Role{RoleID: 111, Name: "nm", ProjectID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &iam.GetRoleRequest{RoleId: 111, ProjectId: 123},
			want:      &iam.GetRoleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &iam.GetRoleRequest{},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &iam.GetRoleRequest{RoleId: 111, ProjectId: 123},
			wantErr:   true,
			mockError: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetRole").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutRole(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name        string
		input       *iam.PutRoleRequest
		want        *iam.PutRoleResponse
		wantErr     bool
		mockGetResp *model.Role
		mockGetErr  error
		mockUpdResp *model.Role
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &iam.PutRoleRequest{Role: &iam.RoleForUpsert{Name: "nm", ProjectId: 123}},
			want:        &iam.PutRoleResponse{Role: &iam.Role{RoleId: 1, Name: "nm", ProjectId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.Role{RoleID: 1, Name: "nm", ProjectID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &iam.PutRoleRequest{Role: &iam.RoleForUpsert{Name: "after", ProjectId: 123}},
			want:        &iam.PutRoleResponse{Role: &iam.Role{RoleId: 1, Name: "after", ProjectId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Role{RoleID: 1, Name: "before", ProjectID: 123, CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.Role{RoleID: 1, Name: "after", ProjectID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &iam.PutRoleRequest{Role: &iam.RoleForUpsert{Name: "nm"}},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetRoleByName)",
			input:      &iam.PutRoleRequest{Role: &iam.RoleForUpsert{Name: "nm", ProjectId: 123}},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutRole)",
			input:      &iam.PutRoleRequest{Role: &iam.RoleForUpsert{Name: "nm", ProjectId: 123}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetRoleByName").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutRole").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteRole(t *testing.T) {
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DeleteRoleRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &iam.DeleteRoleRequest{ProjectId: 1, RoleId: 1},
			wantErr: false,
		},
		{
			name:    "NG Invalid parameters",
			input:   &iam.DeleteRoleRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &iam.DeleteRoleRequest{ProjectId: 1, RoleId: 1},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DeleteRole").Return(c.mockErr).Once()
			_, err := svc.DeleteRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachRole(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name     string
		input    *iam.AttachRoleRequest
		want     *iam.AttachRoleResponse
		wantErr  bool
		mockResp *model.UserRole
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &iam.AttachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
			want:     &iam.AttachRoleResponse{UserRole: &iam.UserRole{ProjectId: 123, UserId: 1, RoleId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.UserRole{ProjectID: 123, UserID: 1, RoleID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.AttachRoleRequest{UserId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &iam.AttachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
			mockErr: gorm.ErrInvalidSQL,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("AttachRole").Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.AttachRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachRole(t *testing.T) {
	var ctx context.Context
	mock := mockIAMRepository{}
	svc := iamService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DetachRoleRequest
		wantErr bool
		mockErr error
	}{
		{
			name:  "OK",
			input: &iam.DetachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.DetachRoleRequest{UserId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &iam.DetachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
			mockErr: gorm.ErrInvalidSQL,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DetachRole").Return(c.mockErr).Once()
			_, err := svc.DetachRole(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
