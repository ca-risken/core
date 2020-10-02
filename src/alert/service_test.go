package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

func TestListAlert(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *alert.ListAlertRequest
		want         *alert.ListAlertResponse
		mockResponce *[]model.Alert
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertRequest{ProjectId: 1, Activated: false, Severity: []string{"high"}, Description: ""},
			want:         &alert.ListAlertResponse{Alert: []*alert.Alert{{AlertId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.Alert{{AlertID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertRequest{ProjectId: 1, Activated: false, Severity: []string{"high"}, Description: ""},
			want:      &alert.ListAlertResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlert").Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlert(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertRequest
		want  *alert.ListAlertRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertRequest{ProjectId: 1, Severity: []string{"high"}, Description: "desc", Activated: true, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertRequest{ProjectId: 1, Severity: []string{"high"}, Description: "desc", Activated: true, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlert(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *alert.GetAlertRequest
		want         *alert.GetAlertResponse
		mockResponce *model.Alert
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertRequest{ProjectId: 1, AlertId: 1001},
			want:         &alert.GetAlertResponse{Alert: &alert.Alert{AlertId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Alert{AlertID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertRequest{ProjectId: 1, AlertId: 9999},
			want:      &alert.GetAlertResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlert").Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlert(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlert(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

	cases := []struct {
		name        string
		input       *alert.PutAlertRequest
		want        *alert.PutAlertResponse
		wantErr     bool
		mockGetResp *model.Alert
		mockGetErr  error
		mockUpResp  *model.Alert
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Activated: true}},
			want:       &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Activated: true}},
			want:        &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Activated: true, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Activated: true, CreatedAt: now, UpdatedAt: now},
			mockUpResp:  &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Activated: true, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetAlertByAlertConditionID").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlert").Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlert(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlert(t *testing.T) {
	var ctx context.Context
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlert").Return(c.mockErr).Once()
			_, err := svc.DeleteAlert(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlert(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.Alert
		want  *alert.Alert
	}{
		{
			name:  "OK convert unix time",
			input: &model.Alert{AlertID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.Alert{AlertId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.Alert{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlert(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * Mock Repository
 */
type mockAlertRepository struct {
	mock.Mock
}

// Alert

func (m *mockAlertRepository) ListAlert(uint32, bool, []string, string, int64, int64) (*[]model.Alert, error) {
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
func (m *mockAlertRepository) GetAlertByAlertConditionIDWithActivated(uint32, uint32, bool) (*model.Alert, error) {
	args := m.Called()
	return args.Get(0).(*model.Alert), args.Error(1)
}
