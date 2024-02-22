// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	report "github.com/ca-risken/core/proto/report"
)

// ReportServiceServer is an autogenerated mock type for the ReportServiceServer type
type ReportServiceServer struct {
	mock.Mock
}

// CollectReportFinding provides a mock function with given fields: _a0, _a1
func (_m *ReportServiceServer) CollectReportFinding(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CollectReportFinding")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReportFinding provides a mock function with given fields: _a0, _a1
func (_m *ReportServiceServer) GetReportFinding(_a0 context.Context, _a1 *report.GetReportFindingRequest) (*report.GetReportFindingResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetReportFinding")
	}

	var r0 *report.GetReportFindingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingRequest) (*report.GetReportFindingResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingRequest) *report.GetReportFindingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*report.GetReportFindingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *report.GetReportFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReportFindingAll provides a mock function with given fields: _a0, _a1
func (_m *ReportServiceServer) GetReportFindingAll(_a0 context.Context, _a1 *report.GetReportFindingAllRequest) (*report.GetReportFindingAllResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetReportFindingAll")
	}

	var r0 *report.GetReportFindingAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingAllRequest) (*report.GetReportFindingAllResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingAllRequest) *report.GetReportFindingAllResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*report.GetReportFindingAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *report.GetReportFindingAllRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PurgeReportFinding provides a mock function with given fields: _a0, _a1
func (_m *ReportServiceServer) PurgeReportFinding(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PurgeReportFinding")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReportServiceServer creates a new instance of ReportServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReportServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReportServiceServer {
	mock := &ReportServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
