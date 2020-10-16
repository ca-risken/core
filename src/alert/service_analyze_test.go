package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jarcoal/httpmock"
)

/*
 * Alert
 */

func TestAnalyzeAlert(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name                              string
		input                             *alert.AnalyzeAlertRequest
		want                              *empty.Empty
		wantErr                           bool
		mockListAlertCondition            *[]model.AlertCondition
		mockListAlertConditionErr         error
		mockListFinding                   *[]model.Finding
		mockListFindingErr                error
		mockListAlertRuleErr              error
		mockListDisabledAlertCondition    *[]model.AlertCondition
		mockListDisabledAlertConditionErr error
	}{
		{
			name:                              "OK",
			input:                             &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:                              &empty.Empty{},
			wantErr:                           false,
			mockListAlertCondition:            &[]model.AlertCondition{},
			mockListFinding:                   &[]model.Finding{},
			mockListDisabledAlertCondition:    &[]model.AlertCondition{},
			mockListDisabledAlertConditionErr: nil,
		},
		{
			name:                      "NG ListAlertConditionErr",
			input:                     &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:                      nil,
			wantErr:                   true,
			mockListAlertCondition:    nil,
			mockListAlertConditionErr: errors.New("Something error occured listAlertCondition"),
			mockListFinding:           &[]model.Finding{},
			mockListFindingErr:        nil,
			mockListAlertRuleErr:      nil,
		},
		{
			name:                      "NG ListFindingErr",
			input:                     &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:                      nil,
			wantErr:                   true,
			mockListAlertCondition:    &[]model.AlertCondition{},
			mockListAlertConditionErr: nil,
			mockListFinding:           nil,
			mockListFindingErr:        errors.New("Something error occured listFinding"),
			mockListAlertRuleErr:      nil,
		},
		{
			name:                      "NG AlertAnalyzeError",
			input:                     &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:                      nil,
			wantErr:                   true,
			mockListAlertCondition:    &[]model.AlertCondition{{AlertConditionID: 1001, CreatedAt: now, UpdatedAt: now}},
			mockListAlertConditionErr: nil,
			mockListFinding:           &[]model.Finding{},
			mockListFindingErr:        nil,
			mockListAlertRuleErr:      errors.New("Something error occured ListAlertRule"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("ListAlertCondition").Return(c.mockListAlertCondition, c.mockListAlertConditionErr).Once()
			mockDB.On("ListFinding").Return(c.mockListFinding, c.mockListFindingErr).Once()
			mockDB.On("ListAlertRuleByAlertConditionID").Return(&[]model.AlertRule{}, c.mockListAlertRuleErr).Once()
			mockDB.On("ListDisabledAlertCondition").Return(c.mockListAlertCondition, c.mockListAlertConditionErr).Once()
			got, err := svc.AnalyzeAlert(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestSendSlackNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://hogehoge.com", httpmock.NewStringResponder(200, "mocked"))
	httpmock.RegisterResponder("POST", "http://fugafuga.com", httpmock.NewErrorResponder(errors.New("Something Wrong")))
	cases := []struct {
		name          string
		notifySetting string
		alertID       uint32
		wantErr       bool
	}{
		{
			name:          "OK",
			notifySetting: `{"webhook_url":"http://hogehoge.com"}`,
			alertID:       1,
			wantErr:       false,
		},
		{
			name:          "NG Json.Marshal Error",
			notifySetting: `{"webhook_url":http://hogehoge.com"}`,
			alertID:       1,
			wantErr:       true,
		},
		{
			name:          "Warn webhook_url not set",
			notifySetting: `{}`,
			alertID:       1,
			wantErr:       false,
		},
		{
			name:          "HTTP Error",
			notifySetting: `{"webhook_url":"http://fugafuga.com"}`,
			alertID:       1,
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sendSlackNotification(c.notifySetting, c.alertID)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestNotificationAlert(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

	httpmock.RegisterResponder("POST", "http://hogehoge.com", httpmock.NewStringResponder(200, "mocked"))
	httpmock.RegisterResponder("POST", "http://fugafuga.com", httpmock.NewErrorResponder(errors.New("Something Wrong")))
	cases := []struct {
		name                               string
		alertCondition                     *model.AlertCondition
		alertID                            uint32
		wantErr                            bool
		mockListAlertCondNotification      *[]model.AlertCondNotification
		mockListAlertCondNotificationErr   error
		mockGetNotification                *model.Notification
		mockGetNotificationErr             error
		mockUpsertAlertCondNotification    *model.AlertCondNotification
		mockUpsertAlertCondNotificationErr error
	}{
		{
			name:                             "OK 0 AlertCondNotification",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alertID:                          1,
			wantErr:                          false,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{},
			mockListAlertCondNotificationErr: nil,
		},
		{
			name:                               "OK Notification Success",
			alertCondition:                     &model.AlertCondition{AlertConditionID: 1},
			alertID:                            1,
			wantErr:                            false,
			mockListAlertCondNotification:      &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr:   nil,
			mockGetNotification:                &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://hogehoge.com"}`},
			mockGetNotificationErr:             nil,
			mockUpsertAlertCondNotification:    &model.AlertCondNotification{},
			mockUpsertAlertCondNotificationErr: nil,
		},
		{
			name:                             "OK Don't send Notification caused NotifedAt",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alertID:                          1,
			wantErr:                          false,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1, CacheSecond: 30, NotifiedAt: now}},
			mockListAlertCondNotificationErr: nil,
			mockGetNotification:              &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://fugafuga.com"}`},
			mockGetNotificationErr:           nil,
		},
		{
			name:                             "Error ListAlertCondNotification Failed",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alertID:                          1,
			wantErr:                          true,
			mockListAlertCondNotification:    nil,
			mockListAlertCondNotificationErr: errors.New("Somethinng error occured"),
		},
		{
			name:                             "Error GetNotification Failed",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alertID:                          1,
			wantErr:                          true,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr: nil,
			mockGetNotification:              nil,
			mockGetNotificationErr:           errors.New("Somethinng error occured"),
		},
		{
			name:                               "Error UpsertAlertCondNotification Failed",
			alertCondition:                     &model.AlertCondition{AlertConditionID: 1},
			alertID:                            1,
			wantErr:                            true,
			mockListAlertCondNotification:      &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr:   nil,
			mockGetNotification:                &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://hogehoge.com"}`},
			mockGetNotificationErr:             nil,
			mockUpsertAlertCondNotification:    nil,
			mockUpsertAlertCondNotificationErr: errors.New("Somethinng error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("ListAlertCondNotification").Return(c.mockListAlertCondNotification, c.mockListAlertCondNotificationErr).Once()
			mockDB.On("GetNotification").Return(c.mockGetNotification, c.mockGetNotificationErr).Once()
			mockDB.On("UpsertAlertCondNotification").Return(c.mockUpsertAlertCondNotification, c.mockUpsertAlertCondNotificationErr).Once()
			got := svc.NotificationAlert(c.alertCondition, c.alertID)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestCheckMatchAlertRuleFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name                  string
		inputAlertRule        *model.AlertRule
		inputFinding          *model.Finding
		want                  bool
		wantErr               bool
		mockListFindingTag    *[]model.FindingTag
		mockListFindingTagErr error
	}{
		{
			name:                  "OK Not Match Score",
			inputAlertRule:        &model.AlertRule{Score: 1.0, CreatedAt: now, UpdatedAt: now},
			inputFinding:          &model.Finding{Score: 0.5, CreatedAt: now, UpdatedAt: now},
			want:                  false,
			wantErr:               false,
			mockListFindingTag:    nil,
			mockListFindingTagErr: nil,
		},
		{
			name:                  "OK Not Match ResourceName",
			inputAlertRule:        &model.AlertRule{Score: 0.1, ResourceName: "piyo", CreatedAt: now, UpdatedAt: now},
			inputFinding:          &model.Finding{Score: 0.5, ResourceName: "hogefuga", CreatedAt: now, UpdatedAt: now},
			want:                  false,
			wantErr:               false,
			mockListFindingTag:    nil,
			mockListFindingTagErr: nil,
		},
		{
			name:                  "OK Not Match FindingTag",
			inputAlertRule:        &model.AlertRule{Score: 0.1, ResourceName: "hoge", Tag: "hoge", CreatedAt: now, UpdatedAt: now},
			inputFinding:          &model.Finding{FindingID: 1, Score: 0.5, ResourceName: "hogefuga", CreatedAt: now, UpdatedAt: now},
			want:                  false,
			wantErr:               false,
			mockListFindingTag:    &[]model.FindingTag{{FindingID: 1, Tag: "fuga"}, {FindingID: 1, Tag: "piyo"}},
			mockListFindingTagErr: nil,
		},
		{
			name:                  "OK Match All",
			inputAlertRule:        &model.AlertRule{Score: 0.1, ResourceName: "hoge", Tag: "hoge", CreatedAt: now, UpdatedAt: now},
			inputFinding:          &model.Finding{FindingID: 1, Score: 0.5, ResourceName: "hogefuga", CreatedAt: now, UpdatedAt: now},
			want:                  true,
			wantErr:               false,
			mockListFindingTag:    &[]model.FindingTag{{FindingID: 1, Tag: "fuga"}, {FindingID: 1, Tag: "hoge"}},
			mockListFindingTagErr: nil,
		},
		{
			name:                  "NG ListFinding Error",
			inputAlertRule:        &model.AlertRule{Score: 0.1, ResourceName: "hoge", Tag: "hoge", CreatedAt: now, UpdatedAt: now},
			inputFinding:          &model.Finding{FindingID: 1, Score: 0.5, ResourceName: "hogefuga", CreatedAt: now, UpdatedAt: now},
			want:                  false,
			wantErr:               true,
			mockListFindingTag:    nil,
			mockListFindingTagErr: errors.New("Somthing error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("ListFindingTag").Return(c.mockListFindingTag, c.mockListFindingTagErr).Once()
			got, err := svc.checkMatchAlertRuleFinding(ctx, c.inputAlertRule, c.inputFinding)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestAnalyzeAlertByRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name                  string
		inputAlertRule        *model.AlertRule
		inputFindings         *[]model.Finding
		wantBool              bool
		wantIntArr            *[]uint64
		wantErr               bool
		mockListFindingTag    *[]model.FindingTag
		mockListFindingTagErr error
	}{
		{
			name:                  "OK Not Match 0 Findings",
			inputAlertRule:        &model.AlertRule{Score: 1.0, CreatedAt: now, UpdatedAt: now, FindingCnt: 1},
			inputFindings:         &[]model.Finding{},
			wantBool:              false,
			wantIntArr:            &[]uint64{},
			wantErr:               false,
			mockListFindingTag:    nil,
			mockListFindingTagErr: nil,
		},
		{
			name:                  "OK FindingCnt <= Match Findings",
			inputAlertRule:        &model.AlertRule{Score: 0.1, CreatedAt: now, UpdatedAt: now, FindingCnt: 2},
			inputFindings:         &[]model.Finding{{FindingID: 1, Score: 0.5}, {FindingID: 2, Score: 1.0}},
			wantBool:              true,
			wantIntArr:            &[]uint64{1, 2},
			wantErr:               false,
			mockListFindingTag:    nil,
			mockListFindingTagErr: nil,
		},
		{
			name:                  "OK FindingCnt > Match Findings",
			inputAlertRule:        &model.AlertRule{Score: 0.1, CreatedAt: now, UpdatedAt: now, FindingCnt: 2},
			inputFindings:         &[]model.Finding{{FindingID: 1, Score: 0.5}},
			wantBool:              false,
			wantIntArr:            &[]uint64{1},
			wantErr:               false,
			mockListFindingTag:    nil,
			mockListFindingTagErr: nil,
		},
		{
			name:                  "NG DB Error",
			inputAlertRule:        &model.AlertRule{Score: 0.1, ResourceName: "hoge", Tag: "fuga", CreatedAt: now, UpdatedAt: now, FindingCnt: 1},
			inputFindings:         &[]model.Finding{{FindingID: 1, Score: 0.5, ResourceName: "hoge"}},
			wantBool:              false,
			wantIntArr:            &[]uint64{},
			wantErr:               true,
			mockListFindingTag:    nil,
			mockListFindingTagErr: errors.New("something error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("ListFindingTag").Return(c.mockListFindingTag, c.mockListFindingTagErr).Once()
			gotBool, gotArr, err := svc.analyzeAlertByRule(ctx, c.inputAlertRule, c.inputFindings)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(gotBool, c.wantBool) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.wantBool, gotBool)
			}
			if !reflect.DeepEqual(*gotArr, *c.wantIntArr) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", *c.wantIntArr, *gotArr)
			}
		})
	}
}

func TestDeleteAlertByAnalyze(t *testing.T) {
	//	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name                                           string
		alertCondition                                 *model.AlertCondition
		wantErr                                        bool
		mockGetAlertByAlertConditionIDWithActivated    *model.Alert
		mockGetAlertByAlertConditionIDWithActivatedErr error
		mockDeactivateAlertErr                         error
		mockUpsertAlertHistory                         *model.AlertHistory
		mockUpsertAlertHistoryErr                      error
		mockListRelAlertFinding                        *[]model.RelAlertFinding
		mockListRelAlertFindingErr                     error
		mockListDeleteAlertFindingErr                  error
	}{
		{
			name:           "OK 0 Alert",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        false,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
		},
		{
			name:           "OK Deactivate Alert Success",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        false,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{AlertID: 1},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockDeactivateAlertErr:                         nil,
			mockUpsertAlertHistory:                         &model.AlertHistory{},
			mockUpsertAlertHistoryErr:                      nil,
			mockListRelAlertFinding:                        &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:                     nil,
			mockListDeleteAlertFindingErr:                  nil,
		},
		{
			name:           "Error GetAlertByAlertConditionIDWithActivated",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    nil,
			mockGetAlertByAlertConditionIDWithActivatedErr: errors.New("Something error occured"),
		},
		{
			name:           "Error DeactivateAlert",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{AlertID: 1, Activated: true},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockDeactivateAlertErr:                         errors.New("Something error occured"),
		},
		{
			name:           "Error UpsertAlertHistory",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{AlertID: 1, Activated: true},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockDeactivateAlertErr:                         nil,
			mockUpsertAlertHistory:                         nil,
			mockUpsertAlertHistoryErr:                      errors.New("Something error occured"),
		},
		{
			name:           "Error ListRelAlertFinding",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{AlertID: 1, Activated: true},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockDeactivateAlertErr:                         nil,
			mockUpsertAlertHistory:                         &model.AlertHistory{},
			mockUpsertAlertHistoryErr:                      nil,
			mockListRelAlertFinding:                        nil,
			mockListRelAlertFindingErr:                     errors.New("Something error occured"),
		},
		{
			name:           "Error DeleteAlertFinding",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{AlertID: 1, Activated: true},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockDeactivateAlertErr:                         nil,
			mockUpsertAlertHistory:                         &model.AlertHistory{},
			mockUpsertAlertHistoryErr:                      nil,
			mockListRelAlertFinding:                        &[]model.RelAlertFinding{{AlertID: 1, FindingID: 1, ProjectID: 1}},
			mockListRelAlertFindingErr:                     nil,
			mockListDeleteAlertFindingErr:                  errors.New("Something error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("GetAlertByAlertConditionIDWithActivated").Return(c.mockGetAlertByAlertConditionIDWithActivated, c.mockGetAlertByAlertConditionIDWithActivatedErr).Once()
			mockDB.On("DeactivateAlert").Return(c.mockDeactivateAlertErr).Once()
			mockDB.On("UpsertAlertHistory").Return(c.mockUpsertAlertHistory, c.mockUpsertAlertHistoryErr).Once()
			mockDB.On("ListRelAlertFinding").Return(c.mockListRelAlertFinding, c.mockListRelAlertFindingErr).Once()
			mockDB.On("DeleteRelAlertFinding").Return(c.mockListDeleteAlertFindingErr).Once()
			got := svc.DeleteAlertByAnalyze(c.alertCondition)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestRegistAlertByAnalyze(t *testing.T) {
	//	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name                                           string
		alertCondition                                 *model.AlertCondition
		findingIDs                                     []uint64
		want                                           uint32
		wantErr                                        bool
		mockGetAlertByAlertConditionIDWithActivated    *model.Alert
		mockGetAlertByAlertConditionIDWithActivatedErr error
		mockUpsertAlert                                *model.Alert
		mockUpsertAlertErr                             error
		mockUpsertAlertHistory                         *model.AlertHistory
		mockUpsertAlertHistoryErr                      error
		mockListRelAlertFinding                        *[]model.RelAlertFinding
		mockListRelAlertFindingErr                     error
		mockUpsertRelAlertFinding                      *model.RelAlertFinding
		mockUpsertRelAlertFindingErr                   error
	}{
		{
			name:           "OK Register Alert Success",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			findingIDs:     []uint64{1},
			want:           1,
			wantErr:        false,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockUpsertAlert:              &model.Alert{AlertID: 1},
			mockUpsertAlertErr:           nil,
			mockUpsertAlertHistory:       &model.AlertHistory{},
			mockUpsertAlertHistoryErr:    nil,
			mockListRelAlertFinding:      &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:   nil,
			mockUpsertRelAlertFinding:    &model.RelAlertFinding{},
			mockUpsertRelAlertFindingErr: nil,
		},
		{
			name:           "Error GetAlertByAlertConditionIDWithActivated",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			findingIDs:     []uint64{1},
			want:           0,
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    nil,
			mockGetAlertByAlertConditionIDWithActivatedErr: errors.New("Something error occured"),
		},
		{
			name:           "Error UpsertAlert",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			findingIDs:     []uint64{1},
			want:           0,
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockUpsertAlert:    nil,
			mockUpsertAlertErr: errors.New("Something error occured"),
		},
		{
			name:           "Error UpsertAlert",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			findingIDs:     []uint64{1},
			want:           0,
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockUpsertAlert:           &model.Alert{AlertID: 1},
			mockUpsertAlertErr:        nil,
			mockUpsertAlertHistory:    nil,
			mockUpsertAlertHistoryErr: errors.New("Something error occured"),
		},
		{
			name:           "Error UpsertRelAlertFinding",
			alertCondition: &model.AlertCondition{AlertConditionID: 1},
			findingIDs:     []uint64{1},
			want:           0,
			wantErr:        true,
			mockGetAlertByAlertConditionIDWithActivated:    &model.Alert{},
			mockGetAlertByAlertConditionIDWithActivatedErr: nil,
			mockUpsertAlert:              &model.Alert{AlertID: 1},
			mockUpsertAlertErr:           nil,
			mockUpsertAlertHistory:       &model.AlertHistory{},
			mockUpsertAlertHistoryErr:    nil,
			mockUpsertRelAlertFinding:    nil,
			mockUpsertRelAlertFindingErr: errors.New("Something error occured"),
			mockListRelAlertFinding:      &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:   nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mockAlertRepository{}
			mockDB.On("GetAlertByAlertConditionIDWithActivated").Return(c.mockGetAlertByAlertConditionIDWithActivated, c.mockGetAlertByAlertConditionIDWithActivatedErr).Once()
			mockDB.On("UpsertAlert").Return(c.mockUpsertAlert, c.mockUpsertAlertErr).Once()
			mockDB.On("UpsertAlertHistory").Return(c.mockUpsertAlertHistory, c.mockUpsertAlertHistoryErr).Once()
			mockDB.On("ListRelAlertFinding").Return(c.mockListRelAlertFinding, c.mockListRelAlertFindingErr).Once()
			mockDB.On("UpsertRelAlertFinding").Return(c.mockUpsertRelAlertFinding, c.mockUpsertRelAlertFindingErr).Once()
			got, err := svc.RegistAlertByAnalyze(c.alertCondition, c.findingIDs)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", got)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
