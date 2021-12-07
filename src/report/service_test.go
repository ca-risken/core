package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/core/proto/report"
	"github.com/ca-risken/core/src/report/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

/*
 * Report
 */

func TestGetReportFinding(t *testing.T) {
	var ctx context.Context
	mockDB := mockReportRepository{}
	svc := reportService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetReportFinding").Return(c.mockResponce, c.mockError).Once()
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
	var ctx context.Context
	mockDB := mockReportRepository{}
	svc := reportService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetReportFindingAll").Return(c.mockResponce, c.mockError).Once()
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

/*
 * Mock Repository
 */
type mockReportRepository struct {
	mock.Mock
}

// Report

func (m *mockReportRepository) GetReportFinding(context.Context, uint32, []string, string, string, float32) (*[]model.ReportFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ReportFinding), args.Error(1)
}
func (m *mockReportRepository) GetReportFindingAll(context.Context, []string, string, string, float32) (*[]model.ReportFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ReportFinding), args.Error(1)
}
func (m *mockReportRepository) CollectReportFinding(ctx context.Context) error {
	args := m.Called()
	return args.Error(1)
}
