package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/pkg/model"
	"gorm.io/gorm"
)

func TestGetPendFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mocks.MockFindingRepository{}
	svc := FindingService{repository: &mockDB}
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
			want:         &finding.GetPendFindingResponse{PendFinding: &finding.PendFinding{FindingId: 1, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.PendFinding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetPendFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mocks.MockFindingRepository{}
	svc := FindingService{repository: &mockDB}
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
			input:        &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			want:         &finding.PutPendFindingResponse{PendFinding: &finding.PendFinding{FindingId: 1, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp:  &model.Finding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			mockPendResp: &model.PendFinding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{}},
			wantErr: true,
		},
		{
			name:       "Record not found",
			input:      &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			wantErr:    true,
			mockGetErr: gorm.ErrRecordNotFound,
		},
		{
			name:       "Invalid DB error(Get)",
			input:      &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "Invalid DB error(Pend)",
			input:       &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			mockPendErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetFinding").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockPendResp != nil || c.mockPendErr != nil {
				mockDB.On("UpsertPendFinding").Return(c.mockPendResp, c.mockPendErr).Once()
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
	mockDB := mocks.MockFindingRepository{}
	svc := FindingService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *finding.DeletePendFindingRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.DeletePendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr: false,
		},
		{
			name:    "NG validation error",
			input:   &finding.DeletePendFindingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &finding.DeletePendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeletePendFinding").Return(c.mockErr).Once()
			_, err := svc.DeletePendFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
