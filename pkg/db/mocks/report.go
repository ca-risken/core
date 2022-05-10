package mocks

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/stretchr/testify/mock"
)

/*
 * Mock Repository
 */
type MockReportRepository struct {
	mock.Mock
}

// Report

func (m *MockReportRepository) GetReportFinding(context.Context, uint32, []string, string, string, float32) (*[]model.ReportFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ReportFinding), args.Error(1)
}
func (m *MockReportRepository) GetReportFindingAll(context.Context, []string, string, string, float32) (*[]model.ReportFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.ReportFinding), args.Error(1)
}
func (m *MockReportRepository) CollectReportFinding(ctx context.Context) error {
	args := m.Called()
	return args.Error(1)
}
