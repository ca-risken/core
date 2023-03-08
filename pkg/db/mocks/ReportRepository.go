// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/ca-risken/core/pkg/model"
)

// ReportRepository is an autogenerated mock type for the ReportRepository type
type ReportRepository struct {
	mock.Mock
}

// CollectReportFinding provides a mock function with given fields: ctx
func (_m *ReportRepository) CollectReportFinding(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetReportFinding provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *ReportRepository) GetReportFinding(_a0 context.Context, _a1 uint32, _a2 []string, _a3 string, _a4 string, _a5 float32) (*[]model.ReportFinding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.ReportFinding
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []string, string, string, float32) *[]model.ReportFinding); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ReportFinding)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint32, []string, string, string, float32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReportFindingAll provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *ReportRepository) GetReportFindingAll(_a0 context.Context, _a1 []string, _a2 string, _a3 string, _a4 float32) (*[]model.ReportFinding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	var r0 *[]model.ReportFinding
	if rf, ok := ret.Get(0).(func(context.Context, []string, string, string, float32) *[]model.ReportFinding); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ReportFinding)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, string, string, float32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReportRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewReportRepository creates a new instance of ReportRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReportRepository(t mockConstructorTestingTNewReportRepository) *ReportRepository {
	mock := &ReportRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
