package report

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/report"
	"gorm.io/gorm"
)

func TestGetReport(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *report.GetReportRequest
		want         *report.GetReportResponse
		mockResponse *model.Report
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &report.GetReportRequest{ProjectId: 1, ReportId: 1001},
			want: &report.GetReportResponse{
				Report: &report.Report{
					ReportId:  1001,
					ProjectId: 1,
					Name:      "test report",
					Type:      "Markdown",
					Status:    "OK",
					Content:   "test content",
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
			mockResponse: &model.Report{
				ReportID:  1001,
				ProjectID: 1,
				Name:      "test report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "test content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name:      "NG Record not found",
			input:     &report.GetReportRequest{ProjectId: 1, ReportId: 9999},
			want:      nil,
			mockError: gorm.ErrRecordNotFound,
			wantErr:   true,
		},
		{
			name:    "NG Validation error",
			input:   &report.GetReportRequest{ProjectId: 0, ReportId: 1001},
			want:    nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}

			if c.input.ProjectId != 0 && (c.mockResponse != nil || c.mockError != nil) {
				mockDB.On("GetReport", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.GetReport(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestListReport(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *report.ListReportRequest
		want         *report.ListReportResponse
		mockResponse *[]model.Report
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK multiple reports",
			input: &report.ListReportRequest{ProjectId: 1},
			want: &report.ListReportResponse{
				Report: []*report.Report{
					{
						ReportId:  1001,
						ProjectId: 1,
						Name:      "report1",
						Type:      "Markdown",
						Status:    "OK",
						Content:   "content1",
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
					{
						ReportId:  1002,
						ProjectId: 1,
						Name:      "report2",
						Type:      "HTML",
						Status:    "IN_PROGRESS",
						Content:   "content2",
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
			},
			mockResponse: &[]model.Report{
				{
					ReportID:  1001,
					ProjectID: 1,
					Name:      "report1",
					Type:      "Markdown",
					Status:    "OK",
					Content:   "content1",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ReportID:  1002,
					ProjectID: 1,
					Name:      "report2",
					Type:      "HTML",
					Status:    "IN_PROGRESS",
					Content:   "content2",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		{
			name:         "OK empty list",
			input:        &report.ListReportRequest{ProjectId: 2},
			want:         &report.ListReportResponse{},
			mockResponse: &[]model.Report{},
		},
		{
			name:      "NG DB error",
			input:     &report.ListReportRequest{ProjectId: 1},
			want:      nil,
			mockError: errors.New("DB error"),
			wantErr:   true,
		},
		{
			name:    "NG Validation error",
			input:   &report.ListReportRequest{ProjectId: 0},
			want:    nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}

			if c.input.ProjectId != 0 && (c.mockResponse != nil || c.mockError != nil) {
				mockDB.On("ListReport", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListReport(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutReport(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name             string
		input            *report.PutReportRequest
		want             *report.PutReportResponse
		mockGetResponse  *model.Report
		mockGetError     error
		mockPutResponse  *model.Report
		mockPutError     error
		wantErr          bool
	}{
		{
			name: "OK create new report",
			input: &report.PutReportRequest{
				ProjectId: 1,
				ReportId:  0,
				Name:      "new report",
				Type:      "Markdown",
				Status:    "IN_PROGRESS",
				Content:   "new content",
			},
			want: &report.PutReportResponse{
				Report: &report.Report{
					ReportId:  2001,
					ProjectId: 1,
					Name:      "new report",
					Type:      "Markdown",
					Status:    "IN_PROGRESS",
					Content:   "new content",
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
			mockPutResponse: &model.Report{
				ReportID:  2001,
				ProjectID: 1,
				Name:      "new report",
				Type:      "Markdown",
				Status:    "IN_PROGRESS",
				Content:   "new content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "OK update existing report",
			input: &report.PutReportRequest{
				ProjectId: 1,
				ReportId:  1001,
				Name:      "updated report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "updated content",
			},
			want: &report.PutReportResponse{
				Report: &report.Report{
					ReportId:  1001,
					ProjectId: 1,
					Name:      "updated report",
					Type:      "Markdown",
					Status:    "OK",
					Content:   "updated content",
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
			mockGetResponse: &model.Report{
				ReportID:  1001,
				ProjectID: 1,
				Name:      "old report",
				Type:      "Markdown",
				Status:    "IN_PROGRESS",
				Content:   "old content",
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockPutResponse: &model.Report{
				ReportID:  1001,
				ProjectID: 1,
				Name:      "updated report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "updated content",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "NG update non-existent report",
			input: &report.PutReportRequest{
				ProjectId: 1,
				ReportId:  9999,
				Name:      "report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "content",
			},
			want:         nil,
			mockGetError: gorm.ErrRecordNotFound,
			wantErr:      true,
		},
		{
			name: "NG DB error on put",
			input: &report.PutReportRequest{
				ProjectId: 1,
				ReportId:  0,
				Name:      "report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "content",
			},
			want:         nil,
			mockPutError: errors.New("DB error"),
			wantErr:      true,
		},
		{
			name: "NG Validation error",
			input: &report.PutReportRequest{
				ProjectId: 0,
				ReportId:  1001,
				Name:      "report",
				Type:      "Markdown",
				Status:    "OK",
				Content:   "content",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}

			if c.input.ProjectId != 0 && c.input.ReportId != 0 {
				mockDB.On("GetReport", test.RepeatMockAnything(3)...).Return(c.mockGetResponse, c.mockGetError).Once()
			}
			if c.input.ProjectId != 0 && c.mockGetError == nil && (c.mockPutResponse != nil || c.mockPutError != nil) {
				mockDB.On("PutReport", test.RepeatMockAnything(2)...).Return(c.mockPutResponse, c.mockPutError).Once()
			}
			result, err := svc.PutReport(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}