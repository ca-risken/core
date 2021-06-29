package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

/*
 * Mock Repository
 */
type mockAlertRepository struct {
	mock.Mock
}

// Alert

func (m *mockAlertRepository) ListAlert(uint32, []string, []string, string, int64, int64) (*[]model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Alert), args.Error(1)
}
func (m *mockAlertRepository) GetAlert(uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) GetAlertByAlertConditionID(uint32, uint32) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlert(*model.Alert) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlert(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertHistory(uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) GetAlertHistory(uint32, uint32) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertHistory(*model.AlertHistory) (*model.AlertHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertHistory), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertHistory(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListRelAlertFinding(uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) GetRelAlertFinding(uint32, uint32, uint32) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) UpsertRelAlertFinding(*model.RelAlertFinding) (*model.RelAlertFinding, error) {
	args := m.Called()
	return args.Get(0).(*model.RelAlertFinding), args.Error(1)
}
func (m *mockAlertRepository) DeleteRelAlertFinding(uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondition(uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondition(uint32, uint32) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondition(*model.AlertCondition) (*model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondition), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondition(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertRule(uint32, float32, float32, int64, int64) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) GetAlertRule(uint32, uint32) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertRule(*model.AlertRule) (*model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertRule(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondRule(uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondRule(uint32, uint32, uint32) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondRule(*model.AlertCondRule) (*model.AlertCondRule, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondRule), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondRule(uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListNotification(uint32, string, int64, int64) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *mockAlertRepository) GetNotification(uint32, uint32) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *mockAlertRepository) UpsertNotification(*model.Notification) (*model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*model.Notification), args.Error(1)
}
func (m *mockAlertRepository) DeleteNotification(uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) ListAlertCondNotification(uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) GetAlertCondNotification(uint32, uint32, uint32) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) UpsertAlertCondNotification(*model.AlertCondNotification) (*model.AlertCondNotification, error) {
	args := m.Called()
	return args.Get(0).(*model.AlertCondNotification), args.Error(1)
}
func (m *mockAlertRepository) DeleteAlertCondNotification(uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockAlertRepository) ListAlertRuleByAlertConditionID(uint32, uint32) (*[]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertRule), args.Error(1)
}
func (m *mockAlertRepository) ListNotificationByAlertConditionID(uint32, uint32) (*[]model.Notification, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Notification), args.Error(1)
}
func (m *mockAlertRepository) DeactivateAlert(*model.Alert) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockAlertRepository) GetAlertByAlertConditionIDStatus(uint32, uint32, []string) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}

func (m *mockAlertRepository) ListFinding(uint32) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}

func (m *mockAlertRepository) ListFindingTag(uint32, uint64) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}

func (m *mockAlertRepository) ListEnabledAlertCondition(uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}

func (m *mockAlertRepository) ListDisabledAlertCondition(uint32, []uint32) (*[]model.AlertCondition, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AlertCondition), args.Error(1)
}

func (m *mockAlertRepository) GetProject(uint32) (*model.Project, error) {
	args := m.Called()
	return args.Get(0).(*model.Project), args.Error(1)
}

type mockFindingClient struct {
	mock.Mock
}

func (m *mockFindingClient) ListFinding(context.Context, *finding.ListFindingRequest, ...grpc.CallOption) (*finding.ListFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingResponse), args.Error(1)
}
func (m *mockFindingClient) BatchListFinding(context.Context, *finding.BatchListFindingRequest, ...grpc.CallOption) (*finding.BatchListFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.BatchListFindingResponse), args.Error(1)
}
func (m *mockFindingClient) GetFinding(context.Context, *finding.GetFindingRequest, ...grpc.CallOption) (*finding.GetFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetFindingResponse), args.Error(1)
}
func (m *mockFindingClient) PutFinding(context.Context, *finding.PutFindingRequest, ...grpc.CallOption) (*finding.PutFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutFindingResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteFinding(context.Context, *finding.DeleteFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListFindingTag(context.Context, *finding.ListFindingTagRequest, ...grpc.CallOption) (*finding.ListFindingTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingTagResponse), args.Error(1)
}
func (m *mockFindingClient) ListFindingTagName(context.Context, *finding.ListFindingTagNameRequest, ...grpc.CallOption) (*finding.ListFindingTagNameResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingTagNameResponse), args.Error(1)
}
func (m *mockFindingClient) TagFinding(context.Context, *finding.TagFindingRequest, ...grpc.CallOption) (*finding.TagFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagFindingResponse), args.Error(1)
}
func (m *mockFindingClient) UntagFinding(context.Context, *finding.UntagFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResource(context.Context, *finding.ListResourceRequest, ...grpc.CallOption) (*finding.ListResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceResponse), args.Error(1)
}
func (m *mockFindingClient) GetResource(context.Context, *finding.GetResourceRequest, ...grpc.CallOption) (*finding.GetResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetResourceResponse), args.Error(1)
}
func (m *mockFindingClient) PutResource(context.Context, *finding.PutResourceRequest, ...grpc.CallOption) (*finding.PutResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutResourceResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteResource(context.Context, *finding.DeleteResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResourceTag(context.Context, *finding.ListResourceTagRequest, ...grpc.CallOption) (*finding.ListResourceTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceTagResponse), args.Error(1)
}
func (m *mockFindingClient) ListResourceTagName(context.Context, *finding.ListResourceTagNameRequest, ...grpc.CallOption) (*finding.ListResourceTagNameResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceTagNameResponse), args.Error(1)
}
func (m *mockFindingClient) TagResource(context.Context, *finding.TagResourceRequest, ...grpc.CallOption) (*finding.TagResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagResourceResponse), args.Error(1)
}
func (m *mockFindingClient) UntagResource(context.Context, *finding.UntagResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) GetPendFinding(context.Context, *finding.GetPendFindingRequest, ...grpc.CallOption) (*finding.GetPendFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetPendFindingResponse), args.Error(1)
}
func (m *mockFindingClient) PutPendFinding(context.Context, *finding.PutPendFindingRequest, ...grpc.CallOption) (*finding.PutPendFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutPendFindingResponse), args.Error(1)
}
func (m *mockFindingClient) DeletePendFinding(context.Context, *finding.DeletePendFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListFindingSetting(context.Context, *finding.ListFindingSettingRequest, ...grpc.CallOption) (*finding.ListFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) GetFindingSetting(context.Context, *finding.GetFindingSettingRequest, ...grpc.CallOption) (*finding.GetFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) PutFindingSetting(context.Context, *finding.PutFindingSettingRequest, ...grpc.CallOption) (*finding.PutFindingSettingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutFindingSettingResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteFindingSetting(context.Context, *finding.DeleteFindingSettingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}