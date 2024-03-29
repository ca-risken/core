package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func TestGetPendFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.GetPendFindingRequest
		want         *finding.GetPendFindingResponse
		mockResponce *model.PendFinding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.GetPendFindingRequest{ProjectId: 1, FindingId: 1},
			want:         &finding.GetPendFindingResponse{PendFinding: &finding.PendFinding{FindingId: 1, ProjectId: 1, ExpiredAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.PendFinding{FindingID: 1, ProjectID: 1, ExpiredAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &finding.GetPendFindingRequest{ProjectId: 1, FindingId: 1},
			want:      &finding.GetPendFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetPendFinding", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetPendFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutPendFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.PutPendFindingRequest
		want         *finding.PutPendFindingResponse
		wantErr      bool
		mockGetResp  *model.Finding
		mockGetErr   error
		mockPendResp *model.PendFinding
		mockPendErr  error
	}{
		{
			name:         "OK",
			input:        &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1, PendUserId: 1}},
			want:         &finding.PutPendFindingResponse{PendFinding: &finding.PendFinding{FindingId: 1, ProjectId: 1, PendUserId: 1, ExpiredAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp:  &model.Finding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			mockPendResp: &model.PendFinding{FindingID: 1, ProjectID: 1, PendUserID: 1, ExpiredAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{}},
			wantErr: true,
		},
		{
			name:       "Record not found",
			input:      &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1, PendUserId: 1}},
			wantErr:    true,
			mockGetErr: gorm.ErrRecordNotFound,
		},
		{
			name:       "Invalid DB error(Get)",
			input:      &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1, PendUserId: 1}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "Invalid DB error(Pend)",
			input:       &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1, PendUserId: 1}},
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			mockPendErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetFinding", test.RepeatMockAnything(4)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockPendResp != nil || c.mockPendErr != nil {
				mockDB.On("UpsertPendFinding", test.RepeatMockAnything(7)...).Return(c.mockPendResp, c.mockPendErr).Once()
			}
			got, err := svc.PutPendFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeletePendFinding(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                  string
		input                 *finding.DeletePendFindingRequest
		wantErr               bool
		mockErr               error
		callDeletePendFinding bool
	}{
		{
			name:                  "OK",
			input:                 &finding.DeletePendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr:               false,
			callDeletePendFinding: true,
		},
		{
			name:                  "NG validation error",
			input:                 &finding.DeletePendFindingRequest{ProjectId: 1},
			wantErr:               true,
			callDeletePendFinding: false,
		},
		{
			name:                  "Invalid DB error",
			input:                 &finding.DeletePendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr:               true,
			mockErr:               gorm.ErrInvalidDB,
			callDeletePendFinding: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.callDeletePendFinding {
				mockDB.On("DeletePendFinding", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeletePendFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
