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

func TestListOrganizationPolicy(t *testing.T) {
	cases := []struct {
		name         string
		input        *organization_iam.ListOrganizationPolicyRequest
		want         *organization_iam.ListOrganizationPolicyResponse
		wantErr      bool
		mockResponce []*model.OrganizationPolicy
		mockError    error
	}{
		{
			name:  "OK",
			input: &organization_iam.ListOrganizationPolicyRequest{OrganizationId: 1, RoleId: 1},
			want:  &organization_iam.ListOrganizationPolicyResponse{PolicyId: []uint32{1, 2, 3}},
			mockResponce: []*model.OrganizationPolicy{
				{PolicyID: 1, Name: "nm1", OrganizationID: 1, ActionPtn: ".*"},
				{PolicyID: 2, Name: "nm2", OrganizationID: 1, ActionPtn: ".*"},
				{PolicyID: 3, Name: "nm3", OrganizationID: 1, ActionPtn: ".*"},
			},
		},
		{
			name:      "OK empty reponse",
			input:     &organization_iam.ListOrganizationPolicyRequest{OrganizationId: 1, Name: "nm", RoleId: 1},
			want:      &organization_iam.ListOrganizationPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid param",
			input:   &organization_iam.ListOrganizationPolicyRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization_iam.ListOrganizationPolicyRequest{OrganizationId: 1, Name: "nm"},
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
				mock.On("ListOrganizationPolicy", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOrganizationPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.GetOrganizationPolicyRequest
		want         *organization_iam.GetOrganizationPolicyResponse
		wantErr      bool
		mockResponce *model.OrganizationPolicy
		mockError    error
	}{
		{
			name:         "OK",
			input:        &organization_iam.GetOrganizationPolicyRequest{PolicyId: 111, OrganizationId: 123},
			want:         &organization_iam.GetOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 111, Name: "nm", OrganizationId: 123, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.OrganizationPolicy{PolicyID: 111, Name: "nm", OrganizationID: 123, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record Not Found",
			input:     &organization_iam.GetOrganizationPolicyRequest{PolicyId: 111, OrganizationId: 123},
			want:      &organization_iam.GetOrganizationPolicyResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &organization_iam.GetOrganizationPolicyRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &organization_iam.GetOrganizationPolicyRequest{PolicyId: 111, OrganizationId: 123},
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
				mock.On("GetOrganizationPolicy", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOrganizationPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *organization_iam.PutOrganizationPolicyRequest
		want        *organization_iam.PutOrganizationPolicyResponse
		wantErr     bool
		mockGetResp *model.OrganizationPolicy
		mockGetErr  error
		mockUpdResp *model.OrganizationPolicy
		mockUpdErr  error
	}{
		{
			name:        "OK Insert",
			input:       &organization_iam.PutOrganizationPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
			want:        &organization_iam.PutOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 1, Name: "nm", OrganizationId: 123, ActionPtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpdResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &organization_iam.PutOrganizationPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
			want:        &organization_iam.PutOrganizationPolicyResponse{Policy: &organization_iam.OrganizationPolicy{PolicyId: 1, Name: "nm", OrganizationId: 123, ActionPtn: ".*", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".+", CreatedAt: now, UpdatedAt: now},
			mockUpdResp: &model.OrganizationPolicy{PolicyID: 1, Name: "nm", OrganizationID: 123, ActionPtn: ".*", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid param",
			input:   &organization_iam.PutOrganizationPolicyRequest{Name: "nm", OrganizationId: 0, ActionPtn: ".*"},
			wantErr: true,
		},
		{
			name:       "NG DB error(GetPolicyByName)",
			input:      &organization_iam.PutOrganizationPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
			mockGetErr: gorm.ErrInvalidTransaction,
			wantErr:    true,
		},
		{
			name:       "NG DB error(PutPolicy)",
			input:      &organization_iam.PutOrganizationPolicyRequest{Name: "nm", OrganizationId: 123, ActionPtn: ".*"},
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
				mock.On("GetOrganizationPolicyByName", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mock.On("PutOrganizationPolicy", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOrganizationPolicy(t *testing.T) {
	cases := []struct {
		name     string
		input    *organization_iam.DeleteOrganizationPolicyRequest
		wantErr  bool
		callMock bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &organization_iam.DeleteOrganizationPolicyRequest{OrganizationId: 1, PolicyId: 1},
			wantErr:  false,
			callMock: true,
		},
		{
			name:     "NG Invalid parameters",
			input:    &organization_iam.DeleteOrganizationPolicyRequest{OrganizationId: 1},
			wantErr:  true,
			callMock: false,
		},
		{
			name:     "Invalid DB error",
			input:    &organization_iam.DeleteOrganizationPolicyRequest{OrganizationId: 1, PolicyId: 1},
			wantErr:  true,
			callMock: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mock := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mock}

			if c.callMock {
				mock.On("DeleteOrganizationPolicy", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestAttachOrganizationPolicy(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *organization_iam.AttachOrganizationPolicyRequest
		want         *organization_iam.AttachOrganizationPolicyResponse
		mockResponse *model.OrganizationPolicy
		mockErr      error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &organization_iam.AttachOrganizationPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			want: &organization_iam.AttachOrganizationPolicyResponse{
				Policy: &organization_iam.OrganizationPolicy{
					OrganizationId: 1,
					PolicyId:       1,
					Name:           "test-policy",
					ActionPtn:      "test:*",
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
				},
			},
			mockResponse: &model.OrganizationPolicy{
				OrganizationID: 1,
				PolicyID:       1,
				Name:           "test-policy",
				ActionPtn:      "test:*",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		{
			name: "NG Invalid param",
			input: &organization_iam.AttachOrganizationPolicyRequest{
				OrganizationId: 0,
				RoleId:         1,
				PolicyId:       1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &organization_iam.AttachOrganizationPolicyRequest{
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
			mockDB := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mockDB}
			if c.mockErr != nil || c.mockResponse != nil {
				mockDB.On("AttachOrganizationPolicy", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockErr).Once()
			}
			got, err := svc.AttachOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachOrganizationPolicy(t *testing.T) {
	cases := []struct {
		name     string
		input    *organization_iam.DetachOrganizationPolicyRequest
		mockErr  error
		wantErr  bool
		mockCall bool
	}{
		{
			name: "OK",
			input: &organization_iam.DetachOrganizationPolicyRequest{
				OrganizationId: 1,
				RoleId:         1,
				PolicyId:       1,
			},
			mockCall: true,
		},
		{
			name: "NG Invalid param",
			input: &organization_iam.DetachOrganizationPolicyRequest{
				OrganizationId: 0,
				RoleId:         1,
				PolicyId:       1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error",
			input: &organization_iam.DetachOrganizationPolicyRequest{
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
			mockDB := mocks.NewOrganizationIAMRepository(t)
			svc := OrganizationIAMService{repository: mockDB}
			if c.mockCall {
				mockDB.On("DetachOrganizationPolicy", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DetachOrganizationPolicy(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v, wantErr: %+v", err, c.wantErr)
			}
		})
	}
}
