// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/ca-risken/core/pkg/model"
)

// AlertRepository is an autogenerated mock type for the AlertRepository type
type AlertRepository struct {
	mock.Mock
}

// DeactivateAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) DeactivateAlert(_a0 context.Context, _a1 *model.Alert) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Alert) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlert provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) DeleteAlert(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlertCondNotification provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) DeleteAlertCondNotification(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlertCondRule provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) DeleteAlertCondRule(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlertCondition provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) DeleteAlertCondition(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlertHistory provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) DeleteAlertHistory(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlertRule provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) DeleteAlertRule(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteNotification provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) DeleteNotification(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRelAlertFinding provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) DeleteRelAlertFinding(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint64) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint64) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAlert provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) GetAlert(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.Alert, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.Alert, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.Alert); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertByAlertConditionIDStatus provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) GetAlertByAlertConditionIDStatus(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 []string) (*model.Alert, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, []string) (*model.Alert, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, []string) *model.Alert); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, []string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondNotification provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) GetAlertCondNotification(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32) (*model.AlertCondNotification, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.AlertCondNotification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) (*model.AlertCondNotification, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) *model.AlertCondNotification); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondNotification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondRule provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) GetAlertCondRule(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32) (*model.AlertCondRule, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.AlertCondRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) (*model.AlertCondRule, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32) *model.AlertCondRule); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertCondition provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) GetAlertCondition(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.AlertCondition, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.AlertCondition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.AlertCondition, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.AlertCondition); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertHistory provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) GetAlertHistory(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.AlertHistory, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.AlertHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.AlertHistory, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.AlertHistory); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlertRule provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) GetAlertRule(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.AlertRule, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.AlertRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.AlertRule, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.AlertRule); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotification provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) GetNotification(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.Notification, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.Notification, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.Notification); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelAlertFinding provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *AlertRepository) GetRelAlertFinding(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint64) (*model.RelAlertFinding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *model.RelAlertFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint64) (*model.RelAlertFinding, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint64) *model.RelAlertFinding); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RelAlertFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlert provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *AlertRepository) ListAlert(_a0 context.Context, _a1 uint32, _a2 []string, _a3 []string, _a4 string, _a5 int64, _a6 int64) (*[]model.Alert, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 *[]model.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []string, []string, string, int64, int64) (*[]model.Alert, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []string, []string, string, int64, int64) *[]model.Alert); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, []string, []string, string, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondNotification provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *AlertRepository) ListAlertCondNotification(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32, _a4 int64, _a5 int64) (*[]model.AlertCondNotification, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.AlertCondNotification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32, int64, int64) *[]model.AlertCondNotification); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertCondNotification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondRule provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *AlertRepository) ListAlertCondRule(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint32, _a4 int64, _a5 int64) (*[]model.AlertCondRule, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.AlertCondRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint32, int64, int64) *[]model.AlertCondRule); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertCondRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint32, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertCondition provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *AlertRepository) ListAlertCondition(_a0 context.Context, _a1 uint32, _a2 []string, _a3 bool, _a4 int64, _a5 int64) (*[]model.AlertCondition, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.AlertCondition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []string, bool, int64, int64) *[]model.AlertCondition); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertCondition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, []string, bool, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertHistory provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *AlertRepository) ListAlertHistory(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 string, _a4 uint32) (*[]model.AlertHistory, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	var r0 *[]model.AlertHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, string, uint32) (*[]model.AlertHistory, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, string, uint32) *[]model.AlertHistory); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, string, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertRule provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *AlertRepository) ListAlertRule(_a0 context.Context, _a1 uint32, _a2 float32, _a3 float32, _a4 int64, _a5 int64) (*[]model.AlertRule, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.AlertRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, float32, float32, int64, int64) (*[]model.AlertRule, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, float32, float32, int64, int64) *[]model.AlertRule); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, float32, float32, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlertRuleByAlertConditionID provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) ListAlertRuleByAlertConditionID(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.AlertRule, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.AlertRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.AlertRule, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.AlertRule); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDisabledAlertCondition provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) ListDisabledAlertCondition(_a0 context.Context, _a1 uint32, _a2 []uint32) (*[]model.AlertCondition, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.AlertCondition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []uint32) *[]model.AlertCondition); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertCondition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, []uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListEnabledAlertCondition provides a mock function with given fields: _a0, _a1, _a2
func (_m *AlertRepository) ListEnabledAlertCondition(_a0 context.Context, _a1 uint32, _a2 []uint32) (*[]model.AlertCondition, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.AlertCondition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, []uint32) *[]model.AlertCondition); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.AlertCondition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, []uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNotification provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *AlertRepository) ListNotification(_a0 context.Context, _a1 uint32, _a2 string, _a3 int64, _a4 int64) (*[]model.Notification, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	var r0 *[]model.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, int64, int64) (*[]model.Notification, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string, int64, int64) *[]model.Notification); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRelAlertFinding provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *AlertRepository) ListRelAlertFinding(_a0 context.Context, _a1 uint32, _a2 uint32, _a3 uint64, _a4 int64, _a5 int64) (*[]model.RelAlertFinding, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 *[]model.RelAlertFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint64, int64, int64) (*[]model.RelAlertFinding, error)); ok {
		return rf(_a0, _a1, _a2, _a3, _a4, _a5)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32, uint64, int64, int64) *[]model.RelAlertFinding); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.RelAlertFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32, uint64, int64, int64) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlert provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlert(_a0 context.Context, _a1 *model.Alert) (*model.Alert, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.Alert
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Alert) (*model.Alert, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Alert) *model.Alert); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Alert)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Alert) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlertCondNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlertCondNotification(_a0 context.Context, _a1 *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.AlertCondNotification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondNotification) (*model.AlertCondNotification, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondNotification) *model.AlertCondNotification); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondNotification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AlertCondNotification) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlertCondRule provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlertCondRule(_a0 context.Context, _a1 *model.AlertCondRule) (*model.AlertCondRule, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.AlertCondRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondRule) (*model.AlertCondRule, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondRule) *model.AlertCondRule); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AlertCondRule) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlertCondition provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlertCondition(_a0 context.Context, _a1 *model.AlertCondition) (*model.AlertCondition, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.AlertCondition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondition) (*model.AlertCondition, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertCondition) *model.AlertCondition); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertCondition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AlertCondition) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlertHistory provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlertHistory(_a0 context.Context, _a1 *model.AlertHistory) (*model.AlertHistory, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.AlertHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertHistory) (*model.AlertHistory, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertHistory) *model.AlertHistory); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AlertHistory) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertAlertRule provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertAlertRule(_a0 context.Context, _a1 *model.AlertRule) (*model.AlertRule, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.AlertRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertRule) (*model.AlertRule, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertRule) *model.AlertRule); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AlertRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AlertRule) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertNotification provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertNotification(_a0 context.Context, _a1 *model.Notification) (*model.Notification, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.Notification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Notification) (*model.Notification, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Notification) *model.Notification); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Notification)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Notification) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertRelAlertFinding provides a mock function with given fields: _a0, _a1
func (_m *AlertRepository) UpsertRelAlertFinding(_a0 context.Context, _a1 *model.RelAlertFinding) (*model.RelAlertFinding, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.RelAlertFinding
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.RelAlertFinding) (*model.RelAlertFinding, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.RelAlertFinding) *model.RelAlertFinding); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RelAlertFinding)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.RelAlertFinding) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAlertRepository creates a new instance of AlertRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlertRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlertRepository {
	mock := &AlertRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
