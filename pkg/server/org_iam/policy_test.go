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

func TestListOrgPolicy(t *testing.T) {
	cases := []struct {
		name         string
		input        *org_iam.ListOrgPolicyRequest
		want         *org_iam.ListOrgPolicyResponse
		wantErr      bool
		mockResponce []*model.OrganizationPolicy
		mockError    error
	}{
		{
			name:  "OK",
			input: &org_iam.ListOrgPolicyRequest{OrganizationId: 1, RoleId: 1},
			want:  &org_iam.ListOrgPolicyResponse{PolicyId: []uint32{1, 2, 3}},
			mockResponce: []*model.OrganizationPolicy{
				{PolicyID: 1, Name: "nm1", OrganizationID: 1, ActionPtn: ".*", ProjectPtn: ".*"},
				{PolicyID: 2, Name: "nm2", OrganizationID: 1, ActionPtn: ".*", ProjectPtn: ".*"},
				{PolicyID: 3, Name: "nm3", OrganizationID: 1, ActionPtn: ".*", ProjectPtn: ".*"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &org_iam.ListOrgPolicyRequest{OrganizationId: 1, Name: "nm", RoleId: 1},
			want:      &org_iam.ListOrgPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &org_iam.ListOrgPolicyRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &org_iam.ListOrgPolicyRequest{OrganizationId: 1, Name: "nm"},
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
				mock.On("ListOrgPolicy", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrgPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.GetOrgPolicyRequest
		want         *org_iam.GetOrgPolicyResponse
		wantErr      bool
		mockResponce *model.OrganizationPolicy
		mockError    error
	}{
		{
			name:         "OK",
			input:        &org_iam.GetOrgPolicyRequest{PolicyId: 111, OrganizationId: 123},
			want:         &org_iam.GetOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 111, Name: "nm", ActionPtn: ".*", ProjectPtn: ".*", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.OrganizationPolicy{PolicyID: 111, Name: "nm", ActionPtn: ".*", ProjectPtn: ".*", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &org_iam.GetOrgPolicyRequest{PolicyId: 111, OrganizationId: 123},
			want:      &org_iam.GetOrgPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &org_iam.GetOrgPolicyRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &org_iam.GetOrgPolicyRequest{PolicyId: 111, OrganizationId: 123},
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
				mock.On("GetOrgPolicy", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrgPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *org_iam.PutOrgPolicyRequest
		want        *org_iam.PutOrgPolicyResponse
		wantErr     bool
		mockGetResp *model.OrganizationPolicy
		mockGetErr  error
		mockUpdResp *model.OrganizationPolicy
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*", ProjectPtn: ".*"},
			want:        &org_iam.PutOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 1, Name: "nm", OrganizationId: 123, ActionPtn: ".*", ProjectPtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".*", ProjectPtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*", ProjectPtn: ".*"},
			want:        &org_iam.PutOrgPolicyResponse{Policy: &org_iam.OrgPolicy{PolicyId: 1, Name: "nm", OrganizationId: 123, ActionPtn: ".*", ProjectPtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".+", ProjectPtn: ".*", CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".*", ProjectPtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 0, ActionPtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Invalid project_ptn regex",
			input:   &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*", ProjectPtn: "[invalid"},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetPolicyByName)",
			input:      &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutPolicy)",
			input:      &org_iam.PutOrgPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
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
				mock.On("GetOrgPolicyByName", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutOrgPolicy", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrgPolicy(t *testing.T) {
	cases := []struct {
		name     string
		input    *org_iam.DeleteOrgPolicyRequest
		wantErr  bool
		callMock bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &org_iam.DeleteOrgPolicyRequest{OrganizationId: 1, PolicyId: 1},
			wantErr:  false,
			callMock: true,
		},
		{
			name:     "NG Invalid parameters",
			input:    &org_iam.DeleteOrgPolicyRequest{OrganizationId: 1},
			wantErr:  true,
			callMock: false,
		},
		{
			name:     "Invalid DB error",
			input:    &org_iam.DeleteOrgPolicyRequest{OrganizationId: 1, PolicyId: 1},
			wantErr:  true,
			callMock: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mock}

			if c.callMock {
				mock.On("DeleteOrgPolicy", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachOrgPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *org_iam.AttachOrgPolicyRequest
		want         *org_iam.AttachOrgPolicyResponse
		mockResponse *model.OrganizationPolicy
		mockErr      error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &org_iam.AttachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			want: &org_iam.AttachOrgPolicyResponse{
				Policy: &org_iam.OrgPolicy{
					PolicyId:       1,
					Name:           "test-policy",
					ActionPtn:      "test:*",
					ProjectPtn:     ".*",
					OrganizationId: 1,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationPolicy{
				PolicyID:       1,
				Name:           "test-policy",
				ActionPtn:      "test:*",
				ProjectPtn:     ".*",
				OrganizationID: 1,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "NG Invalid param",
			input: &org_iam.AttachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         0,
				PolicyId:       1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &org_iam.AttachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mockDB}
			if c.mockErr != nil || c.mockResponse != nil {
				mockDB.On("AttachOrgPolicy", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockErr).Once()
			}
			got, err := svc.AttachOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrgPolicy(t *testing.T) {
	cases := []struct {
		name     string
		input    *org_iam.DetachOrgPolicyRequest
		mockErr  error
		wantErr  bool
		mockCall bool
	}{
		{
			name: "OK",
			input: &org_iam.DetachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			mockCall: true,
		},
		{
			name: "NG Invalid param",
			input: &org_iam.DetachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         0,
				PolicyId:       1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &org_iam.DetachOrgPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
			wantErr:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOrgIAMRepository(t)
			svc := OrgIAMService{repository: mockDB}
			if c.mockCall {
				mockDB.On("DetachOrgPolicy", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DetachOrgPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
		})
	}
}
