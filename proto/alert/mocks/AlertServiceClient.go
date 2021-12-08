// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	alert "github.com/ca-risken/core/proto/alert"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// AlertServiceClient is an autogenerated mock type for the AlertServiceClient type
type AlertServiceClient struct {
	mock.Mock
}

// AnalyzeAlert provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) AnalyzeAlert(ctx context.Context, in *alert.AnalyzeAlertRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.AnalyzeAlertRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.AnalyzeAlertRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyzeAlertAll provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) AnalyzeAlertAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlert provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlert(ctx context.Context, in *alert.DeleteAlertRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlertCondNotification(ctx context.Context, in *alert.DeleteAlertCondNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondNotificationRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertCondNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlertCondRule(ctx context.Context, in *alert.DeleteAlertCondRuleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondRuleRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertCondRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondition provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlertCondition(ctx context.Context, in *alert.DeleteAlertConditionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertConditionRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertConditionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertHistory provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlertHistory(ctx context.Context, in *alert.DeleteAlertHistoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertHistoryRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteAlertRule(ctx context.Context, in *alert.DeleteAlertRuleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRuleRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteNotification(ctx context.Context, in *alert.DeleteNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteNotificationRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRelAlertFinding provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) DeleteRelAlertFinding(ctx context.Context, in *alert.DeleteRelAlertFindingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteRelAlertFindingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteRelAlertFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlert provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlert(ctx context.Context, in *alert.GetAlertRequest, opts ...grpc.CallOption) (*alert.GetAlertResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRequest, ...grpc.CallOption) *alert.GetAlertResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlertCondNotification(ctx context.Context, in *alert.GetAlertCondNotificationRequest, opts ...grpc.CallOption) (*alert.GetAlertCondNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertCondNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondNotificationRequest, ...grpc.CallOption) *alert.GetAlertCondNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertCondNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertCondNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlertCondRule(ctx context.Context, in *alert.GetAlertCondRuleRequest, opts ...grpc.CallOption) (*alert.GetAlertCondRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertCondRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondRuleRequest, ...grpc.CallOption) *alert.GetAlertCondRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertCondRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertCondRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondition provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlertCondition(ctx context.Context, in *alert.GetAlertConditionRequest, opts ...grpc.CallOption) (*alert.GetAlertConditionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertConditionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertConditionRequest, ...grpc.CallOption) *alert.GetAlertConditionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertConditionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertConditionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertHistory provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlertHistory(ctx context.Context, in *alert.GetAlertHistoryRequest, opts ...grpc.CallOption) (*alert.GetAlertHistoryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertHistoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertHistoryRequest, ...grpc.CallOption) *alert.GetAlertHistoryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetAlertRule(ctx context.Context, in *alert.GetAlertRuleRequest, opts ...grpc.CallOption) (*alert.GetAlertRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetAlertRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRuleRequest, ...grpc.CallOption) *alert.GetAlertRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetNotification(ctx context.Context, in *alert.GetNotificationRequest, opts ...grpc.CallOption) (*alert.GetNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetNotificationRequest, ...grpc.CallOption) *alert.GetNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelAlertFinding provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) GetRelAlertFinding(ctx context.Context, in *alert.GetRelAlertFindingRequest, opts ...grpc.CallOption) (*alert.GetRelAlertFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.GetRelAlertFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetRelAlertFindingRequest, ...grpc.CallOption) *alert.GetRelAlertFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetRelAlertFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetRelAlertFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlert provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlert(ctx context.Context, in *alert.ListAlertRequest, opts ...grpc.CallOption) (*alert.ListAlertResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRequest, ...grpc.CallOption) *alert.ListAlertResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlertCondNotification(ctx context.Context, in *alert.ListAlertCondNotificationRequest, opts ...grpc.CallOption) (*alert.ListAlertCondNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertCondNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondNotificationRequest, ...grpc.CallOption) *alert.ListAlertCondNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertCondNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertCondNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlertCondRule(ctx context.Context, in *alert.ListAlertCondRuleRequest, opts ...grpc.CallOption) (*alert.ListAlertCondRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertCondRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondRuleRequest, ...grpc.CallOption) *alert.ListAlertCondRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertCondRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertCondRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondition provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlertCondition(ctx context.Context, in *alert.ListAlertConditionRequest, opts ...grpc.CallOption) (*alert.ListAlertConditionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertConditionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertConditionRequest, ...grpc.CallOption) *alert.ListAlertConditionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertConditionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertConditionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertHistory provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlertHistory(ctx context.Context, in *alert.ListAlertHistoryRequest, opts ...grpc.CallOption) (*alert.ListAlertHistoryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertHistoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertHistoryRequest, ...grpc.CallOption) *alert.ListAlertHistoryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListAlertRule(ctx context.Context, in *alert.ListAlertRuleRequest, opts ...grpc.CallOption) (*alert.ListAlertRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListAlertRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRuleRequest, ...grpc.CallOption) *alert.ListAlertRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListNotification(ctx context.Context, in *alert.ListNotificationRequest, opts ...grpc.CallOption) (*alert.ListNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListNotificationRequest, ...grpc.CallOption) *alert.ListNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRelAlertFinding provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) ListRelAlertFinding(ctx context.Context, in *alert.ListRelAlertFindingRequest, opts ...grpc.CallOption) (*alert.ListRelAlertFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.ListRelAlertFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListRelAlertFindingRequest, ...grpc.CallOption) *alert.ListRelAlertFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListRelAlertFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListRelAlertFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlert provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlert(ctx context.Context, in *alert.PutAlertRequest, opts ...grpc.CallOption) (*alert.PutAlertResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRequest, ...grpc.CallOption) *alert.PutAlertResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlertCondNotification(ctx context.Context, in *alert.PutAlertCondNotificationRequest, opts ...grpc.CallOption) (*alert.PutAlertCondNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertCondNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondNotificationRequest, ...grpc.CallOption) *alert.PutAlertCondNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertCondNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertCondNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlertCondRule(ctx context.Context, in *alert.PutAlertCondRuleRequest, opts ...grpc.CallOption) (*alert.PutAlertCondRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertCondRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondRuleRequest, ...grpc.CallOption) *alert.PutAlertCondRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertCondRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertCondRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondition provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlertCondition(ctx context.Context, in *alert.PutAlertConditionRequest, opts ...grpc.CallOption) (*alert.PutAlertConditionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertConditionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertConditionRequest, ...grpc.CallOption) *alert.PutAlertConditionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertConditionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertConditionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertHistory provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlertHistory(ctx context.Context, in *alert.PutAlertHistoryRequest, opts ...grpc.CallOption) (*alert.PutAlertHistoryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertHistoryResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertHistoryRequest, ...grpc.CallOption) *alert.PutAlertHistoryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertHistoryResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertRule provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutAlertRule(ctx context.Context, in *alert.PutAlertRuleRequest, opts ...grpc.CallOption) (*alert.PutAlertRuleResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutAlertRuleResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRuleRequest, ...grpc.CallOption) *alert.PutAlertRuleResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertRuleResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertRuleRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutNotification(ctx context.Context, in *alert.PutNotificationRequest, opts ...grpc.CallOption) (*alert.PutNotificationResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutNotificationResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutNotificationRequest, ...grpc.CallOption) *alert.PutNotificationResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutNotificationResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRelAlertFinding provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) PutRelAlertFinding(ctx context.Context, in *alert.PutRelAlertFindingRequest, opts ...grpc.CallOption) (*alert.PutRelAlertFindingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *alert.PutRelAlertFindingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutRelAlertFindingRequest, ...grpc.CallOption) *alert.PutRelAlertFindingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutRelAlertFindingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutRelAlertFindingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TestNotification provides a mock function with given fields: ctx, in, opts
func (_m *AlertServiceClient) TestNotification(ctx context.Context, in *alert.TestNotificationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *alert.TestNotificationRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *alert.TestNotificationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}