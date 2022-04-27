package iam

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func TestListPolicy(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.ListPolicyRequest
		want         *iam.ListPolicyResponse
		wantErr      bool
		mockResponce *[]model.Policy
		mockError    error
	}{
		{
			name:  "OK",
			input: &iam.ListPolicyRequest{ProjectId: 1, Name: "nm", RoleId: 1},
			want:  &iam.ListPolicyResponse{PolicyId: []uint32{1, 2, 3}},
			mockResponce: &[]model.Policy{
				{PolicyID: 1, Name: "nm"},
				{PolicyID: 2, Name: "nm"},
				{PolicyID: 3, Name: "nm"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &iam.ListPolicyRequest{ProjectId: 1, Name: "nm", RoleId: 1},
			want:      &iam.ListPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &iam.ListPolicyRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &iam.ListPolicyRequest{ProjectId: 1, Name: "nm"},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("ListPolicy").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetPolicy(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name         string
		input        *iam.GetPolicyRequest
		want         *iam.GetPolicyResponse
		wantErr      bool
		mockResponce *model.Policy
		mockError    error
	}{
		{
			name:         "OK",
			input:        &iam.GetPolicyRequest{PolicyId: 111, ProjectId: 123},
			want:         &iam.GetPolicyResponse{Policy: &iam.Policy{PolicyId: 111, Name: "nm", ProjectId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Policy{PolicyID: 111, Name: "nm", ProjectID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &iam.GetPolicyRequest{PolicyId: 111, ProjectId: 123},
			want:      &iam.GetPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &iam.GetPolicyRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &iam.GetPolicyRequest{PolicyId: 111, ProjectId: 123},
			wantErr:   true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mock.On("GetPolicy").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutPolicy(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name        string
		input       *iam.PutPolicyRequest
		want        *iam.PutPolicyResponse
		wantErr     bool
		mockGetResp *model.Policy
		mockGetErr  error
		mockUpdResp *model.Policy
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &iam.PutPolicyRequest{ProjectId: 123, Policy: &iam.PolicyForUpsert{Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"}},
			want:        &iam.PutPolicyResponse{Policy: &iam.Policy{PolicyId: 1, Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.Policy{PolicyID: 1, Name: "nm", ProjectID: 123, ActionPtn: ".*", ResourcePtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &iam.PutPolicyRequest{ProjectId: 123, Policy: &iam.PolicyForUpsert{Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"}},
			want:        &iam.PutPolicyResponse{Policy: &iam.Policy{PolicyId: 1, Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Policy{PolicyID: 1, Name: "nm", ProjectID: 123, ActionPtn: ".*", ResourcePtn: ".*", CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.Policy{PolicyID: 1, Name: "nm", ProjectID: 123, ActionPtn: ".*", ResourcePtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &iam.PutPolicyRequest{ProjectId: 999, Policy: &iam.PolicyForUpsert{Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"}},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetPolkicyByName)",
			input:      &iam.PutPolicyRequest{ProjectId: 123, Policy: &iam.PolicyForUpsert{Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"}},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutPolicy)",
			input:      &iam.PutPolicyRequest{ProjectId: 123, Policy: &iam.PolicyForUpsert{Name: "nm", ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpdErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetPolicyByName").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutPolicy").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DeletePolicyRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &iam.DeletePolicyRequest{ProjectId: 1, PolicyId: 1},
			wantErr: false,
		},
		{
			name:    "NG Invalid parameters",
			input:   &iam.DeletePolicyRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &iam.DeletePolicyRequest{ProjectId: 1, PolicyId: 1},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DeletePolicy").Return(c.mockErr).Once()
			_, err := svc.DeletePolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachPolicy(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name     string
		input    *iam.AttachPolicyRequest
		want     *iam.AttachPolicyResponse
		wantErr  bool
		mockResp *model.RolePolicy
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &iam.AttachPolicyRequest{ProjectId: 1, RoleId: 2, PolicyId: 3},
			want:     &iam.AttachPolicyResponse{RolePolicy: &iam.RolePolicy{RoleId: 2, PolicyId: 3, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.RolePolicy{ProjectID: 1, RoleID: 2, PolicyID: 3, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.AttachPolicyRequest{ProjectId: 1, RoleId: 2},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &iam.AttachPolicyRequest{ProjectId: 1, RoleId: 2, PolicyId: 3},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("AttachPolicy").Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.AttachPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachPolicy(t *testing.T) {
	var ctx context.Context
	mock := mocks.MockIAMRepository{}
	svc := IAMService{repository: &mock}
	cases := []struct {
		name    string
		input   *iam.DetachPolicyRequest
		wantErr bool
		mockErr error
	}{
		{
			name:  "OK",
			input: &iam.DetachPolicyRequest{ProjectId: 1, RoleId: 2, PolicyId: 3},
		},
		{
			name:    "NG Invalid parameter",
			input:   &iam.DetachPolicyRequest{RoleId: 2, PolicyId: 3},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &iam.DetachPolicyRequest{ProjectId: 1, RoleId: 2, PolicyId: 3},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DetachPolicy").Return(c.mockErr).Once()
			_, err := svc.DetachPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
