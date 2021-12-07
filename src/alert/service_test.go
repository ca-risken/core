package main

import (
	"context"

	"github.com/ca-risken/core/src/alert/model"
	"github.com/stretchr/testify/mock"
)

/*
 * Mock Repository
 */
type mockAlertRepository struct {
	mock.Mock
}

// Alert

func (m *mockAlertRepository) ListAlert(context.Context, uint32, []string, []string, string, int64, int64) (*[]model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Alert), args.Error(1)
}
func (m *mockAlertRepository) GetAlert(context.Context, uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) GetAlertByAlertConditionID(context.Context, uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlert(context.Context, *model.Alert) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlert(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertHistory(context.Context, uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) GetAlertHistory(context.Context, uint32, uint32) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertHistory(context.Context, *model.AlertHistory) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertHistory(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListRelAlertFinding(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) GetRelAlertFinding(context.Context, uint32, uint32, uint32) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) UpsertRelAlertFinding(context.Context, *model.RelAlertFinding) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) DeleteRelAlertFinding(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondition(context.Context, uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondition(context.Context, uint32, uint32) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondition(context.Context, *model.AlertCondition) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondition(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertRule(context.Context, uint32, float32, float32, int64, int64) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) GetAlertRule(context.Context, uint32, uint32) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertRule(context.Context, *model.AlertRule) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertRule(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondRule(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondRule(context.Context, uint32, uint32, uint32) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondRule(context.Context, *model.AlertCondRule) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondRule(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListNotification(context.Context, uint32, string, int64, int64) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *mockAlertRepository) GetNotification(context.Context, uint32, uint32) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *mockAlertRepository) UpsertNotification(context.Context, *model.Notification) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *mockAlertRepository) DeleteNotification(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondNotification(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondNotification(context.Context, uint32, uint32, uint32) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondNotification(context.Context, *model.AlertCondNotification) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondNotification(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockAlertRepository) ListAlertRuleByAlertConditionID(context.Context, uint32, uint32) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) ListNotificationByAlertConditionID(context.Context, uint32, uint32) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *mockAlertRepository) DeactivateAlert(context.Context, *model.Alert) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) GetAlertByAlertConditionIDStatus(context.Context, uint32, uint32, []string) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}

func (m *mockAlertRepository) ListEnabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}

func (m *mockAlertRepository) ListDisabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}
