package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"gorm.io/gorm"
)

func TestGetPendFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *finding.PutPendFindingRequest
		want     *finding.PutPendFindingResponse
		wantErr  bool
		mockResp *model.PendFinding
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			want:     &finding.PutPendFindingResponse{PendFinding: &finding.PendFinding{FindingId: 1, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResp: &model.PendFinding{FindingID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{}},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &finding.PutPendFindingRequest{ProjectId: 1, PendFinding: &finding.PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("UpsertPendFinding").Return(c.mockResp, c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
