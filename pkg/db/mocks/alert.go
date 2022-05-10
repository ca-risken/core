package mocks

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/stretchr/testify/mock"
)

/*
 * Mock Repository
 */
type MockAlertRepository struct {
	mock.Mock
}

// Alert

func (m *MockAlertRepository) ListAlert(context.Context, uint32, []string, []string, string, int64, int64) (*[]model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Alert), args.Error(1)
}
func (m *MockAlertRepository) GetAlert(context.Context, uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *MockAlertRepository) GetAlertByAlertConditionID(context.Context, uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlert(context.Context, *model.Alert) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlert(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListAlertHistory(context.Context, uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertHistory), args.Error(1)
}
func (m *MockAlertRepository) GetAlertHistory(context.Context, uint32, uint32) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlertHistory(context.Context, *model.AlertHistory) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlertHistory(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListRelAlertFinding(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelAlertFinding), args.Error(1)
}
func (m *MockAlertRepository) GetRelAlertFinding(context.Context, uint32, uint32, uint32) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *MockAlertRepository) UpsertRelAlertFinding(context.Context, *model.RelAlertFinding) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *MockAlertRepository) DeleteRelAlertFinding(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListAlertCondition(context.Context, uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}
func (m *MockAlertRepository) GetAlertCondition(context.Context, uint32, uint32) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlertCondition(context.Context, *model.AlertCondition) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlertCondition(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListAlertRule(context.Context, uint32, float32, float32, int64, int64) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *MockAlertRepository) GetAlertRule(context.Context, uint32, uint32) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlertRule(context.Context, *model.AlertRule) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlertRule(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListAlertCondRule(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondRule), args.Error(1)
}
func (m *MockAlertRepository) GetAlertCondRule(context.Context, uint32, uint32, uint32) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlertCondRule(context.Context, *model.AlertCondRule) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlertCondRule(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListNotification(context.Context, uint32, string, int64, int64) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *MockAlertRepository) GetNotification(context.Context, uint32, uint32) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *MockAlertRepository) UpsertNotification(context.Context, *model.Notification) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *MockAlertRepository) DeleteNotification(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) ListAlertCondNotification(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondNotification), args.Error(1)
}
func (m *MockAlertRepository) GetAlertCondNotification(context.Context, uint32, uint32, uint32) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *MockAlertRepository) UpsertAlertCondNotification(context.Context, *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *MockAlertRepository) DeleteAlertCondNotification(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAlertRepository) ListAlertRuleByAlertConditionID(context.Context, uint32, uint32) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *MockAlertRepository) ListNotificationByAlertConditionID(context.Context, uint32, uint32) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *MockAlertRepository) DeactivateAlert(context.Context, *model.Alert) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAlertRepository) GetAlertByAlertConditionIDStatus(context.Context, uint32, uint32, []string) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}

func (m *MockAlertRepository) ListEnabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}

func (m *MockAlertRepository) ListDisabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}
