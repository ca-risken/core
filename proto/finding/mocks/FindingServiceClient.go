// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	finding "github.com/ca-risken/core/proto/finding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// FindingServiceClient is an autogenerated mock type for the FindingServiceClient type
type FindingServiceClient struct {
	mock.Mock
}

// BatchListFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) BatchListFinding(ctx context.Context, in *finding.BatchListFindingRequest, opts ...grpc.CallOption) (*finding.BatchListFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.BatchListFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.BatchListFindingRequest, ...grpc.CallOption) *finding.BatchListFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.BatchListFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.BatchListFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClearScore provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ClearScore(ctx context.Context, in *finding.ClearScoreRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ClearScoreRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ClearScoreRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) DeleteFinding(ctx context.Context, in *finding.DeleteFindingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.DeleteFindingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.DeleteFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteFindingSetting provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) DeleteFindingSetting(ctx context.Context, in *finding.DeleteFindingSettingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.DeleteFindingSettingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.DeleteFindingSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePendFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) DeletePendFinding(ctx context.Context, in *finding.DeletePendFindingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.DeletePendFindingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.DeletePendFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) DeleteResource(ctx context.Context, in *finding.DeleteResourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.DeleteResourceRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.DeleteResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAISummary provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetAISummary(ctx context.Context, in *finding.GetAISummaryRequest, opts ...grpc.CallOption) (*finding.GetAISummaryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetAISummaryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetAISummaryRequest, ...grpc.CallOption) *finding.GetAISummaryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetAISummaryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetAISummaryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetFinding(ctx context.Context, in *finding.GetFindingRequest, opts ...grpc.CallOption) (*finding.GetFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetFindingRequest, ...grpc.CallOption) *finding.GetFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFindingSetting provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetFindingSetting(ctx context.Context, in *finding.GetFindingSettingRequest, opts ...grpc.CallOption) (*finding.GetFindingSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetFindingSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetFindingSettingRequest, ...grpc.CallOption) *finding.GetFindingSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetFindingSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetFindingSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPendFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetPendFinding(ctx context.Context, in *finding.GetPendFindingRequest, opts ...grpc.CallOption) (*finding.GetPendFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetPendFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetPendFindingRequest, ...grpc.CallOption) *finding.GetPendFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetPendFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetPendFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecommend provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetRecommend(ctx context.Context, in *finding.GetRecommendRequest, opts ...grpc.CallOption) (*finding.GetRecommendResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetRecommendResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetRecommendRequest, ...grpc.CallOption) *finding.GetRecommendResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetRecommendResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetRecommendRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) GetResource(ctx context.Context, in *finding.GetResourceRequest, opts ...grpc.CallOption) (*finding.GetResourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.GetResourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.GetResourceRequest, ...grpc.CallOption) *finding.GetResourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.GetResourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.GetResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListFinding(ctx context.Context, in *finding.ListFindingRequest, opts ...grpc.CallOption) (*finding.ListFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingRequest, ...grpc.CallOption) *finding.ListFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingSetting provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListFindingSetting(ctx context.Context, in *finding.ListFindingSettingRequest, opts ...grpc.CallOption) (*finding.ListFindingSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListFindingSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingSettingRequest, ...grpc.CallOption) *finding.ListFindingSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListFindingSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTag provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListFindingTag(ctx context.Context, in *finding.ListFindingTagRequest, opts ...grpc.CallOption) (*finding.ListFindingTagResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListFindingTagResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagRequest, ...grpc.CallOption) *finding.ListFindingTagResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListFindingTagResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFindingTagName provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListFindingTagName(ctx context.Context, in *finding.ListFindingTagNameRequest, opts ...grpc.CallOption) (*finding.ListFindingTagNameResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListFindingTagNameResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListFindingTagNameRequest, ...grpc.CallOption) *finding.ListFindingTagNameResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListFindingTagNameResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListFindingTagNameRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListResource(ctx context.Context, in *finding.ListResourceRequest, opts ...grpc.CallOption) (*finding.ListResourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListResourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceRequest, ...grpc.CallOption) *finding.ListResourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListResourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTag provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListResourceTag(ctx context.Context, in *finding.ListResourceTagRequest, opts ...grpc.CallOption) (*finding.ListResourceTagResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListResourceTagResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagRequest, ...grpc.CallOption) *finding.ListResourceTagResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListResourceTagResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListResourceTagName provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) ListResourceTagName(ctx context.Context, in *finding.ListResourceTagNameRequest, opts ...grpc.CallOption) (*finding.ListResourceTagNameResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.ListResourceTagNameResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.ListResourceTagNameRequest, ...grpc.CallOption) *finding.ListResourceTagNameResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.ListResourceTagNameResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.ListResourceTagNameRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutFinding(ctx context.Context, in *finding.PutFindingRequest, opts ...grpc.CallOption) (*finding.PutFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.PutFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutFindingRequest, ...grpc.CallOption) *finding.PutFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.PutFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutFindingBatch provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutFindingBatch(ctx context.Context, in *finding.PutFindingBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutFindingBatchRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutFindingBatchRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutFindingSetting provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutFindingSetting(ctx context.Context, in *finding.PutFindingSettingRequest, opts ...grpc.CallOption) (*finding.PutFindingSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.PutFindingSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutFindingSettingRequest, ...grpc.CallOption) *finding.PutFindingSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.PutFindingSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutFindingSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPendFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutPendFinding(ctx context.Context, in *finding.PutPendFindingRequest, opts ...grpc.CallOption) (*finding.PutPendFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.PutPendFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutPendFindingRequest, ...grpc.CallOption) *finding.PutPendFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.PutPendFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutPendFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRecommend provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutRecommend(ctx context.Context, in *finding.PutRecommendRequest, opts ...grpc.CallOption) (*finding.PutRecommendResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.PutRecommendResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutRecommendRequest, ...grpc.CallOption) *finding.PutRecommendResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.PutRecommendResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutRecommendRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutResource(ctx context.Context, in *finding.PutResourceRequest, opts ...grpc.CallOption) (*finding.PutResourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.PutResourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutResourceRequest, ...grpc.CallOption) *finding.PutResourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.PutResourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutResourceBatch provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) PutResourceBatch(ctx context.Context, in *finding.PutResourceBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.PutResourceBatchRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.PutResourceBatchRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) TagFinding(ctx context.Context, in *finding.TagFindingRequest, opts ...grpc.CallOption) (*finding.TagFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.TagFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.TagFindingRequest, ...grpc.CallOption) *finding.TagFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.TagFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.TagFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TagResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) TagResource(ctx context.Context, in *finding.TagResourceRequest, opts ...grpc.CallOption) (*finding.TagResourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *finding.TagResourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *finding.TagResourceRequest, ...grpc.CallOption) *finding.TagResourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*finding.TagResourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.TagResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UntagFinding provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) UntagFinding(ctx context.Context, in *finding.UntagFindingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.UntagFindingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.UntagFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UntagResource provides a mock function with given fields: ctx, in, opts
func (_m *FindingServiceClient) UntagResource(ctx context.Context, in *finding.UntagResourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *finding.UntagResourceRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *finding.UntagResourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFindingServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewFindingServiceClient creates a new instance of FindingServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFindingServiceClient(t mockConstructorTestingTNewFindingServiceClient) *FindingServiceClient {
	mock := &FindingServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
