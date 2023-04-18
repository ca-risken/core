package report

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/report"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

/*
 * Report
 */

func TestGetReportFinding(t *testing.T) {
	cases := []struct {
		name         string
		input        *report.GetReportFindingRequest
		want         *report.GetReportFindingResponse
		mockResponce *[]model.ReportFinding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &report.GetReportFindingRequest{ProjectId: 1},
			want:         &report.GetReportFindingResponse{ReportFinding: []*report.ReportFinding{{ReportFindingId: 1001}, {ReportFindingId: 1002}}},
			mockResponce: &[]model.ReportFinding{{ReportFindingID: 1001}, {ReportFindingID: 1002}},
		},
		{
			name:      "NG Record not found",
			input:     &report.GetReportFindingRequest{ProjectId: 1},
			want:      &report.GetReportFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetReportFinding", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetReportFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestGetReportFindingAll(t *testing.T) {
	cases := []struct {
		name         string
		input        *report.GetReportFindingAllRequest
		want         *report.GetReportFindingAllResponse
		mockResponce *[]model.ReportFinding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &report.GetReportFindingAllRequest{},
			want:         &report.GetReportFindingAllResponse{ReportFinding: []*report.ReportFinding{{ReportFindingId: 1001}, {ReportFindingId: 1002}}},
			mockResponce: &[]model.ReportFinding{{ReportFindingID: 1001}, {ReportFindingID: 1002}},
		},
		{
			name:      "NG Record not found",
			input:     &report.GetReportFindingAllRequest{ProjectId: 1},
			want:      &report.GetReportFindingAllResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetReportFindingAll", test.RepeatMockAnything(5)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetReportFindingAll(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPurgeReportFinding(t *testing.T) {
	cases := []struct {
		name      string
		wantErr   bool
		mockError error
	}{
		{
			name:      "OK",
			wantErr:   false,
			mockError: nil,
		},
		{
			name:      "NG DB Error",
			wantErr:   true,
			mockError: errors.New("DB Error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewReportRepository(t)
			svc := ReportService{repository: mockDB}
			mockDB.On("PurgeReportFinding", mock.Anything).Return(c.mockError).Once()

			_, err := svc.PurgeReportFinding(ctx, &empty.Empty{})
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}
