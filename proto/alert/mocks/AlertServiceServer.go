// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	alert "github.com/ca-risken/core/proto/alert"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"
)

// AlertServiceServer is an autogenerated mock type for the AlertServiceServer type
type AlertServiceServer struct {
	mock.Mock
}

// AnalyzeAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) AnalyzeAlert(_a0 context.Context, _a1 *alert.AnalyzeAlertRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AnalyzeAlert")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.AnalyzeAlertRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.AnalyzeAlertRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.AnalyzeAlertRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyzeAlertAll provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) AnalyzeAlertAll(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AnalyzeAlertAll")
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

// DeleteAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlert(_a0 context.Context, _a1 *alert.DeleteAlertRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlert")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlertCondNotification(_a0 context.Context, _a1 *alert.DeleteAlertCondNotificationRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlertCondNotification")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondNotificationRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondNotificationRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertCondNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlertCondRule(_a0 context.Context, _a1 *alert.DeleteAlertCondRuleRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlertCondRule")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondRuleRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertCondRuleRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertCondRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertCondition provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlertCondition(_a0 context.Context, _a1 *alert.DeleteAlertConditionRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlertCondition")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertConditionRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertConditionRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertConditionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertHistory provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlertHistory(_a0 context.Context, _a1 *alert.DeleteAlertHistoryRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlertHistory")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertHistoryRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertHistoryRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAlertRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteAlertRule(_a0 context.Context, _a1 *alert.DeleteAlertRuleRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlertRule")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRuleRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteAlertRuleRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteAlertRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteNotification(_a0 context.Context, _a1 *alert.DeleteNotificationRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNotification")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteNotificationRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteNotificationRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRelAlertFinding provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) DeleteRelAlertFinding(_a0 context.Context, _a1 *alert.DeleteRelAlertFindingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRelAlertFinding")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteRelAlertFindingRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.DeleteRelAlertFindingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.DeleteRelAlertFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlert(_a0 context.Context, _a1 *alert.GetAlertRequest) (*alert.GetAlertResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlert")
	}

	var r0 *alert.GetAlertResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRequest) (*alert.GetAlertResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRequest) *alert.GetAlertResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlertCondNotification(_a0 context.Context, _a1 *alert.GetAlertCondNotificationRequest) (*alert.GetAlertCondNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlertCondNotification")
	}

	var r0 *alert.GetAlertCondNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondNotificationRequest) (*alert.GetAlertCondNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondNotificationRequest) *alert.GetAlertCondNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertCondNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertCondNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlertCondRule(_a0 context.Context, _a1 *alert.GetAlertCondRuleRequest) (*alert.GetAlertCondRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlertCondRule")
	}

	var r0 *alert.GetAlertCondRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondRuleRequest) (*alert.GetAlertCondRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertCondRuleRequest) *alert.GetAlertCondRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertCondRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertCondRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondition provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlertCondition(_a0 context.Context, _a1 *alert.GetAlertConditionRequest) (*alert.GetAlertConditionResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlertCondition")
	}

	var r0 *alert.GetAlertConditionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertConditionRequest) (*alert.GetAlertConditionResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertConditionRequest) *alert.GetAlertConditionResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertConditionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertConditionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertHistory provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlertHistory(_a0 context.Context, _a1 *alert.GetAlertHistoryRequest) (*alert.GetAlertHistoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlertHistory")
	}

	var r0 *alert.GetAlertHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertHistoryRequest) (*alert.GetAlertHistoryResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertHistoryRequest) *alert.GetAlertHistoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetAlertRule(_a0 context.Context, _a1 *alert.GetAlertRuleRequest) (*alert.GetAlertRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAlertRule")
	}

	var r0 *alert.GetAlertRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRuleRequest) (*alert.GetAlertRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetAlertRuleRequest) *alert.GetAlertRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetAlertRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetAlertRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetNotification(_a0 context.Context, _a1 *alert.GetNotificationRequest) (*alert.GetNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetNotification")
	}

	var r0 *alert.GetNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetNotificationRequest) (*alert.GetNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetNotificationRequest) *alert.GetNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelAlertFinding provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) GetRelAlertFinding(_a0 context.Context, _a1 *alert.GetRelAlertFindingRequest) (*alert.GetRelAlertFindingResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetRelAlertFinding")
	}

	var r0 *alert.GetRelAlertFindingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetRelAlertFindingRequest) (*alert.GetRelAlertFindingResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.GetRelAlertFindingRequest) *alert.GetRelAlertFindingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.GetRelAlertFindingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.GetRelAlertFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlert(_a0 context.Context, _a1 *alert.ListAlertRequest) (*alert.ListAlertResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlert")
	}

	var r0 *alert.ListAlertResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRequest) (*alert.ListAlertResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRequest) *alert.ListAlertResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlertCondNotification(_a0 context.Context, _a1 *alert.ListAlertCondNotificationRequest) (*alert.ListAlertCondNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlertCondNotification")
	}

	var r0 *alert.ListAlertCondNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondNotificationRequest) (*alert.ListAlertCondNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondNotificationRequest) *alert.ListAlertCondNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertCondNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertCondNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlertCondRule(_a0 context.Context, _a1 *alert.ListAlertCondRuleRequest) (*alert.ListAlertCondRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlertCondRule")
	}

	var r0 *alert.ListAlertCondRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondRuleRequest) (*alert.ListAlertCondRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertCondRuleRequest) *alert.ListAlertCondRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertCondRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertCondRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondition provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlertCondition(_a0 context.Context, _a1 *alert.ListAlertConditionRequest) (*alert.ListAlertConditionResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlertCondition")
	}

	var r0 *alert.ListAlertConditionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertConditionRequest) (*alert.ListAlertConditionResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertConditionRequest) *alert.ListAlertConditionResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertConditionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertConditionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertHistory provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlertHistory(_a0 context.Context, _a1 *alert.ListAlertHistoryRequest) (*alert.ListAlertHistoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlertHistory")
	}

	var r0 *alert.ListAlertHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertHistoryRequest) (*alert.ListAlertHistoryResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertHistoryRequest) *alert.ListAlertHistoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListAlertRule(_a0 context.Context, _a1 *alert.ListAlertRuleRequest) (*alert.ListAlertRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListAlertRule")
	}

	var r0 *alert.ListAlertRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRuleRequest) (*alert.ListAlertRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListAlertRuleRequest) *alert.ListAlertRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListAlertRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListAlertRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListNotification(_a0 context.Context, _a1 *alert.ListNotificationRequest) (*alert.ListNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListNotification")
	}

	var r0 *alert.ListNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListNotificationRequest) (*alert.ListNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListNotificationRequest) *alert.ListNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNotificationForInternal provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListNotificationForInternal(_a0 context.Context, _a1 *alert.ListNotificationForInternalRequest) (*alert.ListNotificationForInternalResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListNotificationForInternal")
	}

	var r0 *alert.ListNotificationForInternalResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListNotificationForInternalRequest) (*alert.ListNotificationForInternalResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListNotificationForInternalRequest) *alert.ListNotificationForInternalResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListNotificationForInternalResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListNotificationForInternalRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRelAlertFinding provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) ListRelAlertFinding(_a0 context.Context, _a1 *alert.ListRelAlertFindingRequest) (*alert.ListRelAlertFindingResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListRelAlertFinding")
	}

	var r0 *alert.ListRelAlertFindingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListRelAlertFindingRequest) (*alert.ListRelAlertFindingResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.ListRelAlertFindingRequest) *alert.ListRelAlertFindingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.ListRelAlertFindingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.ListRelAlertFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlert(_a0 context.Context, _a1 *alert.PutAlertRequest) (*alert.PutAlertResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlert")
	}

	var r0 *alert.PutAlertResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRequest) (*alert.PutAlertResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRequest) *alert.PutAlertResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertCondNotification(_a0 context.Context, _a1 *alert.PutAlertCondNotificationRequest) (*alert.PutAlertCondNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertCondNotification")
	}

	var r0 *alert.PutAlertCondNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondNotificationRequest) (*alert.PutAlertCondNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondNotificationRequest) *alert.PutAlertCondNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertCondNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertCondNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertCondRule(_a0 context.Context, _a1 *alert.PutAlertCondRuleRequest) (*alert.PutAlertCondRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertCondRule")
	}

	var r0 *alert.PutAlertCondRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondRuleRequest) (*alert.PutAlertCondRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertCondRuleRequest) *alert.PutAlertCondRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertCondRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertCondRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertCondition provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertCondition(_a0 context.Context, _a1 *alert.PutAlertConditionRequest) (*alert.PutAlertConditionResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertCondition")
	}

	var r0 *alert.PutAlertConditionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertConditionRequest) (*alert.PutAlertConditionResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertConditionRequest) *alert.PutAlertConditionResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertConditionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertConditionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertFirstViewedAt provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertFirstViewedAt(_a0 context.Context, _a1 *alert.PutAlertFirstViewedAtRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertFirstViewedAt")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertFirstViewedAtRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertFirstViewedAtRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertFirstViewedAtRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertHistory provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertHistory(_a0 context.Context, _a1 *alert.PutAlertHistoryRequest) (*alert.PutAlertHistoryResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertHistory")
	}

	var r0 *alert.PutAlertHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertHistoryRequest) (*alert.PutAlertHistoryResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertHistoryRequest) *alert.PutAlertHistoryResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertHistoryRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutAlertRule provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutAlertRule(_a0 context.Context, _a1 *alert.PutAlertRuleRequest) (*alert.PutAlertRuleResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutAlertRule")
	}

	var r0 *alert.PutAlertRuleResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRuleRequest) (*alert.PutAlertRuleResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutAlertRuleRequest) *alert.PutAlertRuleResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutAlertRuleResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutAlertRuleRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutNotification(_a0 context.Context, _a1 *alert.PutNotificationRequest) (*alert.PutNotificationResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutNotification")
	}

	var r0 *alert.PutNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutNotificationRequest) (*alert.PutNotificationResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutNotificationRequest) *alert.PutNotificationResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRelAlertFinding provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) PutRelAlertFinding(_a0 context.Context, _a1 *alert.PutRelAlertFindingRequest) (*alert.PutRelAlertFindingResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PutRelAlertFinding")
	}

	var r0 *alert.PutRelAlertFindingResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutRelAlertFindingRequest) (*alert.PutRelAlertFindingResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.PutRelAlertFindingRequest) *alert.PutRelAlertFindingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*alert.PutRelAlertFindingResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.PutRelAlertFindingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RequestProjectRoleNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) RequestProjectRoleNotification(_a0 context.Context, _a1 *alert.RequestProjectRoleNotificationRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RequestProjectRoleNotification")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.RequestProjectRoleNotificationRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.RequestProjectRoleNotificationRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.RequestProjectRoleNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TestNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertServiceServer) TestNotification(_a0 context.Context, _a1 *alert.TestNotificationRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TestNotification")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *alert.TestNotificationRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *alert.TestNotificationRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *alert.TestNotificationRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAlertServiceServer creates a new instance of AlertServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlertServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlertServiceServer {
	mock := &AlertServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
