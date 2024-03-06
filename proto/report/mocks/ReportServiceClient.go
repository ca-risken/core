// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	report "github.com/ca-risken/core/proto/report"
)

// ReportServiceClient is an autogenerated mock type for the ReportServiceClient type
type ReportServiceClient struct {
	mock.Mock
}

// CollectReportFinding provides a mock function with given fields: ctx, in, opts
func (_m *ReportServiceClient) CollectReportFinding(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CollectReportFinding")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReportFinding provides a mock function with given fields: ctx, in, opts
func (_m *ReportServiceClient) GetReportFinding(ctx context.Context, in *report.GetReportFindingRequest, opts ...grpc.CallOption) (*report.GetReportFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetReportFinding")
	}

	var r0 *report.GetReportFindingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingRequest, ...grpc.CallOption) (*report.GetReportFindingResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingRequest, ...grpc.CallOption) *report.GetReportFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*report.GetReportFindingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *report.GetReportFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReportFindingAll provides a mock function with given fields: ctx, in, opts
func (_m *ReportServiceClient) GetReportFindingAll(ctx context.Context, in *report.GetReportFindingAllRequest, opts ...grpc.CallOption) (*report.GetReportFindingAllResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetReportFindingAll")
	}

	var r0 *report.GetReportFindingAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingAllRequest, ...grpc.CallOption) (*report.GetReportFindingAllResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *report.GetReportFindingAllRequest, ...grpc.CallOption) *report.GetReportFindingAllResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*report.GetReportFindingAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *report.GetReportFindingAllRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PurgeReportFinding provides a mock function with given fields: ctx, in, opts
func (_m *ReportServiceClient) PurgeReportFinding(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PurgeReportFinding")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReportServiceClient creates a new instance of ReportServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReportServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReportServiceClient {
	mock := &ReportServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
