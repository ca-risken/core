package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	findingmock "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jarcoal/httpmock"
	"gorm.io/gorm"
)

/*
 * Alert
 */

func TestAnalyzeAlert(t *testing.T) {
	now := time.Now()
	mockDB := mocks.MockAlertRepository{}
	mockProject := projectmock.ProjectServiceClient{}
	svc := AlertService{repository: &mockDB, projectClient: &mockProject}
	cases := []struct {
		name                              string
		input                             *alert.AnalyzeAlertRequest
		want                              *empty.Empty
		wantErr                           bool
		mockListProject                   *project.ListProjectResponse
		mockListProjectErr                error
		mockListAlertCondition            *[]model.AlertCondition
		mockListAlertConditionErr         error
		mockListAlertRuleErr              error
		mockListDisabledAlertCondition    *[]model.AlertCondition
		mockListDisabledAlertConditionErr error
	}{
		{
			name:    "OK",
			input:   &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:    &empty.Empty{},
			wantErr: false,
			mockListProject: &project.ListProjectResponse{Project: []*project.Project{
				{ProjectId: 1001, Name: "project1"},
			}},
			mockListAlertCondition:            &[]model.AlertCondition{},
			mockListDisabledAlertCondition:    &[]model.AlertCondition{},
			mockListDisabledAlertConditionErr: nil,
		},
		{
			name:               "NG ListProjectErr",
			input:              &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:               nil,
			wantErr:            true,
			mockListProject:    nil,
			mockListProjectErr: errors.New("Something error occurred LListProject"),
		},
		{
			name:    "NG ListAlertConditionErr",
			input:   &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:    nil,
			wantErr: true,
			mockListProject: &project.ListProjectResponse{Project: []*project.Project{
				{ProjectId: 1001, Name: "project1"},
			}},
			mockListAlertCondition:    nil,
			mockListAlertConditionErr: errors.New("Something error occured listAlertCondition"),
			mockListAlertRuleErr:      nil,
		},
		{
			name:    "NG AlertAnalyzeError",
			input:   &alert.AnalyzeAlertRequest{ProjectId: 1001},
			want:    nil,
			wantErr: true,
			mockListProject: &project.ListProjectResponse{Project: []*project.Project{
				{ProjectId: 1001, Name: "project1"},
			}},
			mockListAlertCondition:    &[]model.AlertCondition{{AlertConditionID: 1001, CreatedAt: now, UpdatedAt: now}},
			mockListAlertConditionErr: nil,
			mockListAlertRuleErr:      errors.New("Something error occured ListAlertRule"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockDB = mocks.MockAlertRepository{}
			mockDB.On("ListEnabledAlertCondition").Return(c.mockListAlertCondition, c.mockListAlertConditionErr).Once()
			mockDB.On("ListAlertRuleByAlertConditionID").Return(&[]model.AlertRule{}, c.mockListAlertRuleErr).Once()
			mockDB.On("ListDisabledAlertCondition").Return(c.mockListAlertCondition, c.mockListAlertConditionErr).Once()
			mockProject = projectmock.ProjectServiceClient{}
			mockProject.On("ListProject", ctx, &project.ListProjectRequest{ProjectId: c.input.ProjectId}).
				Return(c.mockListProject, c.mockListProjectErr).Once()
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
		alert         *model.Alert
		project       *project.Project
		wantErr       bool
	}{
		{
			name:          "OK",
			notifySetting: `{"webhook_url":"http://hogehoge.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       false,
		},
		{
			name:          "NG Json.Marshal Error",
			notifySetting: `{"webhook_url":http://hogehoge.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       true,
		},
		{
			name:          "Warn webhook_url not set",
			notifySetting: `{}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       false,
		},
		{
			name:          "HTTP Error",
			notifySetting: `{"webhook_url":"http://fugafuga.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sendSlackNotification("unused", c.notifySetting, c.alert, c.project, &[]model.AlertRule{})
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
	mockDB := mocks.MockAlertRepository{}
	svc := AlertService{repository: &mockDB}

	httpmock.RegisterResponder("POST", "http://hogehoge.com", httpmock.NewStringResponder(200, "mocked"))
	httpmock.RegisterResponder("POST", "http://fugafuga.com", httpmock.NewErrorResponder(errors.New("Something Wrong")))
	cases := []struct {
		name                               string
		alertCondition                     *model.AlertCondition
		alert                              *model.Alert
		wantErr                            bool
		mockListAlertCondNotification      *[]model.AlertCondNotification
		mockListAlertCondNotificationErr   error
		mockGetNotification                *model.Notification
		mockGetNotificationErr             error
		mockGetProject                     *project.Project
		mockGetProjectErr                  error
		mockUpsertAlertCondNotification    *model.AlertCondNotification
		mockUpsertAlertCondNotificationErr error
	}{
		{
			name:                             "OK 0 AlertCondNotification",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alert:                            &model.Alert{},
			wantErr:                          false,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{},
			mockListAlertCondNotificationErr: nil,
		},
		{
			name:                               "OK Notification Success",
			alertCondition:                     &model.AlertCondition{AlertConditionID: 1},
			alert:                              &model.Alert{},
			wantErr:                            false,
			mockListAlertCondNotification:      &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr:   nil,
			mockGetNotification:                &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://hogehoge.com"}`},
			mockGetNotificationErr:             nil,
			mockGetProject:                     &project.Project{},
			mockGetProjectErr:                  nil,
			mockUpsertAlertCondNotification:    &model.AlertCondNotification{},
			mockUpsertAlertCondNotificationErr: nil,
		},
		{
			name:                             "OK Don't send Notification caused NotifedAt",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alert:                            &model.Alert{},
			wantErr:                          false,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1, CacheSecond: 30, NotifiedAt: now}},
			mockListAlertCondNotificationErr: nil,
			mockGetNotification:              &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://fugafuga.com"}`},
			mockGetNotificationErr:           nil,
		},
		{
			name:                             "Error ListAlertCondNotification Failed",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alert:                            &model.Alert{},
			wantErr:                          true,
			mockListAlertCondNotification:    nil,
			mockListAlertCondNotificationErr: errors.New("Somethinng error occured"),
		},
		{
			name:                             "Error GetNotification Failed",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alert:                            &model.Alert{},
			wantErr:                          true,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr: nil,
			mockGetNotification:              nil,
			mockGetNotificationErr:           errors.New("Somethinng error occured"),
		},
		{
			name:                             "Error GetNotification Failed",
			alertCondition:                   &model.AlertCondition{AlertConditionID: 1},
			alert:                            &model.Alert{},
			wantErr:                          true,
			mockListAlertCondNotification:    &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr: nil,
			mockGetNotification:              nil,
			mockGetNotificationErr:           errors.New("Somethinng error occured"),
			mockGetProject:                   nil,
			mockGetProjectErr:                errors.New("Somethinng error occured"),
		},
		{
			name:                               "Error UpsertAlertCondNotification Failed",
			alertCondition:                     &model.AlertCondition{AlertConditionID: 1},
			alert:                              &model.Alert{},
			wantErr:                            true,
			mockListAlertCondNotification:      &[]model.AlertCondNotification{{AlertConditionID: 1, NotificationID: 1}},
			mockListAlertCondNotificationErr:   nil,
			mockGetNotification:                &model.Notification{Type: "slack", NotifySetting: `{"webhook_url":"http://hogehoge.com"}`},
			mockGetNotificationErr:             nil,
			mockGetProject:                     &project.Project{},
			mockGetProjectErr:                  nil,
			mockUpsertAlertCondNotification:    nil,
			mockUpsertAlertCondNotificationErr: errors.New("Somethinng error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mocks.MockAlertRepository{}
			mockDB.On("ListAlertCondNotification").Return(c.mockListAlertCondNotification, c.mockListAlertCondNotificationErr).Once()
			mockDB.On("GetNotification").Return(c.mockGetNotification, c.mockGetNotificationErr).Once()
			mockDB.On("UpsertAlertCondNotification").Return(c.mockUpsertAlertCondNotification, c.mockUpsertAlertCondNotificationErr).Once()
			mockDB.On("GetProject").Return(c.mockGetProject, c.mockGetProjectErr).Once()
			got := svc.NotificationAlert(context.Background(), c.alertCondition, c.alert, &[]model.AlertRule{}, &project.Project{})
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestAnalyzeAlertByRule(t *testing.T) {
	now := time.Now()
	mockFinding := findingmock.FindingServiceClient{}
	svc := AlertService{findingClient: &mockFinding}
	cases := []struct {
		name                   string
		inputAlertRule         *model.AlertRule
		wantBool               bool
		wantIntArr             *[]uint64
		wantErr                bool
		mockListFindingRequest *finding.BatchListFindingRequest
		mockListFinding        *finding.BatchListFindingResponse
		mockListFindingErr     error
	}{
		{
			name:                   "OK Not Match 0 Findings",
			inputAlertRule:         &model.AlertRule{Score: 1.0, CreatedAt: now, UpdatedAt: now, FindingCnt: 1},
			wantBool:               false,
			wantIntArr:             &[]uint64{},
			wantErr:                false,
			mockListFindingRequest: &finding.BatchListFindingRequest{FromScore: 1.0, Status: finding.FindingStatus_FINDING_ACTIVE},
			mockListFinding:        &finding.BatchListFindingResponse{FindingId: []uint64{}, Total: 0, Count: 0},
			mockListFindingErr:     nil,
		},
		{
			name:                   "OK FindingCnt <= Match Findings",
			inputAlertRule:         &model.AlertRule{Score: 0.1, CreatedAt: now, UpdatedAt: now, FindingCnt: 2},
			wantBool:               true,
			wantIntArr:             &[]uint64{1, 2},
			wantErr:                false,
			mockListFindingRequest: &finding.BatchListFindingRequest{FromScore: 0.1, Status: finding.FindingStatus_FINDING_ACTIVE},
			mockListFinding:        &finding.BatchListFindingResponse{FindingId: []uint64{1, 2}, Total: 2, Count: 2},
			mockListFindingErr:     nil,
		},
		{
			name:                   "OK FindingCnt > Match Findings",
			inputAlertRule:         &model.AlertRule{Score: 0.1, CreatedAt: now, UpdatedAt: now, FindingCnt: 2},
			wantBool:               false,
			wantIntArr:             &[]uint64{1},
			wantErr:                false,
			mockListFindingRequest: &finding.BatchListFindingRequest{FromScore: 0.1, Status: finding.FindingStatus_FINDING_ACTIVE},
			mockListFinding:        &finding.BatchListFindingResponse{FindingId: []uint64{1}, Total: 1, Count: 1},
			mockListFindingErr:     nil,
		},
		{
			name:           "NG DB Error",
			inputAlertRule: &model.AlertRule{Score: 0.1, ResourceName: "hoge", Tag: "fuga", CreatedAt: now, UpdatedAt: now, FindingCnt: 1},
			wantBool:       false,
			wantIntArr:     &[]uint64{},
			wantErr:        true,
			mockListFindingRequest: &finding.BatchListFindingRequest{FromScore: 0.1,
				ResourceName: []string{"hoge"}, Tag: []string{"fuga"}, Status: finding.FindingStatus_FINDING_ACTIVE},
			mockListFinding:    nil,
			mockListFindingErr: errors.New("something error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()

			mockFinding.On("BatchListFinding", ctx, c.mockListFindingRequest).
				Return(c.mockListFinding, c.mockListFindingErr).Once()
			gotBool, gotArr, err := svc.analyzeAlertByRule(ctx, c.inputAlertRule)
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
	mockDB := mocks.MockAlertRepository{}
	svc := AlertService{repository: &mockDB}
	cases := []struct {
		name                                    string
		alertCondition                          *model.AlertCondition
		wantErr                                 bool
		mockGetAlertByAlertConditionIDStatus    *model.Alert
		mockGetAlertByAlertConditionIDStatusErr error
		mockDeactivateAlertErr                  error
		mockUpsertAlertHistory                  *model.AlertHistory
		mockUpsertAlertHistoryErr               error
		mockListRelAlertFinding                 *[]model.RelAlertFinding
		mockListRelAlertFindingErr              error
		mockListDeleteAlertFindingErr           error
	}{
		{
			name:                                    "OK 0 Alert",
			alertCondition:                          &model.AlertCondition{ProjectID: 1, AlertConditionID: 1},
			wantErr:                                 false,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: gorm.ErrRecordNotFound,
		},
		{
			name:                                    "OK Deactivate Alert Success",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 false,
			mockGetAlertByAlertConditionIDStatus:    &model.Alert{AlertID: 1},
			mockGetAlertByAlertConditionIDStatusErr: nil,
			mockDeactivateAlertErr:                  nil,
			mockUpsertAlertHistory:                  &model.AlertHistory{},
			mockUpsertAlertHistoryErr:               nil,
			mockListRelAlertFinding:                 &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:              nil,
			mockListDeleteAlertFindingErr:           nil,
		},
		{
			name:                                    "Error GetAlertByAlertConditionIDStatus",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: errors.New("Something error occured"),
		},
		{
			name:                                    "Error DeactivateAlert",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    &model.Alert{AlertID: 1, Status: "ACTIVE"},
			mockGetAlertByAlertConditionIDStatusErr: nil,
			mockDeactivateAlertErr:                  gorm.ErrInvalidDB,
		},
		{
			name:                                    "Error UpsertAlertHistory",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    &model.Alert{AlertID: 1, Status: "ACTIVE"},
			mockGetAlertByAlertConditionIDStatusErr: nil,
			mockDeactivateAlertErr:                  nil,
			mockUpsertAlertHistory:                  nil,
			mockUpsertAlertHistoryErr:               errors.New("Something error occured"),
		},
		{
			name:                                    "Error ListRelAlertFinding",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    &model.Alert{AlertID: 1, Status: "ACTIVE"},
			mockGetAlertByAlertConditionIDStatusErr: nil,
			mockDeactivateAlertErr:                  nil,
			mockUpsertAlertHistory:                  &model.AlertHistory{},
			mockUpsertAlertHistoryErr:               nil,
			mockListRelAlertFinding:                 nil,
			mockListRelAlertFindingErr:              errors.New("Something error occured"),
		},
		{
			name:                                    "Error DeleteAlertFinding",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    &model.Alert{AlertID: 1, Status: "ACTIVE"},
			mockGetAlertByAlertConditionIDStatusErr: nil,
			mockDeactivateAlertErr:                  nil,
			mockUpsertAlertHistory:                  &model.AlertHistory{},
			mockUpsertAlertHistoryErr:               nil,
			mockListRelAlertFinding:                 &[]model.RelAlertFinding{{AlertID: 1, FindingID: 1, ProjectID: 1}},
			mockListRelAlertFindingErr:              nil,
			mockListDeleteAlertFindingErr:           errors.New("Something error occured"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mocks.MockAlertRepository{}
			mockDB.On("GetAlertByAlertConditionIDStatus").Return(c.mockGetAlertByAlertConditionIDStatus, c.mockGetAlertByAlertConditionIDStatusErr).Once()
			mockDB.On("DeactivateAlert").Return(c.mockDeactivateAlertErr).Once()
			mockDB.On("UpsertAlertHistory").Return(c.mockUpsertAlertHistory, c.mockUpsertAlertHistoryErr).Once()
			mockDB.On("ListRelAlertFinding").Return(c.mockListRelAlertFinding, c.mockListRelAlertFindingErr).Once()
			mockDB.On("DeleteRelAlertFinding").Return(c.mockListDeleteAlertFindingErr).Once()
			got := svc.DeleteAlertByAnalyze(context.Background(), c.alertCondition)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestRegistAlertByAnalyze(t *testing.T) {
	//	now := time.Now()
	mockDB := mocks.MockAlertRepository{}
	svc := AlertService{repository: &mockDB}
	cases := []struct {
		name                                    string
		alertCondition                          *model.AlertCondition
		findingIDs                              []uint64
		want                                    *model.Alert
		wantErr                                 bool
		mockGetAlertByAlertConditionIDStatus    *model.Alert
		mockGetAlertByAlertConditionIDStatusErr error
		mockUpsertAlert                         *model.Alert
		mockUpsertAlertErr                      error
		mockUpsertAlertHistory                  *model.AlertHistory
		mockUpsertAlertHistoryErr               error
		mockListRelAlertFinding                 *[]model.RelAlertFinding
		mockListRelAlertFindingErr              error
		mockUpsertRelAlertFinding               *model.RelAlertFinding
		mockUpsertRelAlertFindingErr            error
	}{
		{
			name:                                    "OK RegistAlert Success",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			findingIDs:                              []uint64{1},
			want:                                    &model.Alert{AlertID: 1},
			wantErr:                                 false,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: gorm.ErrRecordNotFound,
			mockUpsertAlert:                         &model.Alert{AlertID: 1},
			mockUpsertAlertErr:                      nil,
			mockUpsertAlertHistory:                  &model.AlertHistory{},
			mockUpsertAlertHistoryErr:               nil,
			mockListRelAlertFinding:                 &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:              nil,
			mockUpsertRelAlertFinding:               &model.RelAlertFinding{},
			mockUpsertRelAlertFindingErr:            nil,
		},
		{
			name:                                    "Error GetAlertByAlertConditionIDStatus",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			findingIDs:                              []uint64{1},
			want:                                    nil,
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: errors.New("Something error occured"),
		},
		{
			name:                                    "Error UpsertAlert",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			findingIDs:                              []uint64{1},
			want:                                    nil,
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: gorm.ErrRecordNotFound,
			mockUpsertAlert:                         nil,
			mockUpsertAlertErr:                      errors.New("Something error occured"),
		},
		{
			name:                                    "Error UpsertRelAlertFinding",
			alertCondition:                          &model.AlertCondition{AlertConditionID: 1},
			findingIDs:                              []uint64{1},
			want:                                    nil,
			wantErr:                                 true,
			mockGetAlertByAlertConditionIDStatus:    nil,
			mockGetAlertByAlertConditionIDStatusErr: gorm.ErrRecordNotFound,
			mockUpsertAlert:                         &model.Alert{AlertID: 1},
			mockUpsertAlertErr:                      nil,
			mockUpsertAlertHistory:                  &model.AlertHistory{},
			mockUpsertAlertHistoryErr:               nil,
			mockUpsertRelAlertFinding:               nil,
			mockUpsertRelAlertFindingErr:            errors.New("Something error occured"),
			mockListRelAlertFinding:                 &[]model.RelAlertFinding{},
			mockListRelAlertFindingErr:              nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB = mocks.MockAlertRepository{}
			mockDB.On("GetAlertByAlertConditionIDStatus").Return(c.mockGetAlertByAlertConditionIDStatus, c.mockGetAlertByAlertConditionIDStatusErr).Once()
			mockDB.On("UpsertAlert").Return(c.mockUpsertAlert, c.mockUpsertAlertErr).Once()
			mockDB.On("UpsertAlertHistory").Return(c.mockUpsertAlertHistory, c.mockUpsertAlertHistoryErr).Once()
			mockDB.On("ListRelAlertFinding").Return(c.mockListRelAlertFinding, c.mockListRelAlertFindingErr).Once()
			mockDB.On("UpsertRelAlertFinding").Return(c.mockUpsertRelAlertFinding, c.mockUpsertRelAlertFindingErr).Once()
			got, err := svc.RegistAlertByAnalyze(context.Background(), c.alertCondition, c.findingIDs)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", got)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGenerateRuleList(t *testing.T) {
	cases := []struct {
		name  string
		input *[]model.AlertRule
		want  string
	}{
		{
			name: "1 line",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
			},
			want: "- aaa",
		},
		{
			name: "Multi lines",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
				{AlertRuleID: 2, Name: "bbb"},
				{AlertRuleID: 3, Name: "ccc"},
			},
			want: "- aaa\n- bbb\n- ccc",
		},
		{
			name:  "Nil input",
			input: nil,
			want:  "",
		},
		{
			name: "Too many rules",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
				{AlertRuleID: 2, Name: "bbb"},
				{AlertRuleID: 3, Name: "ccc"},
				{AlertRuleID: 4, Name: "ddd"},
				{AlertRuleID: 5, Name: "eee"},
				{AlertRuleID: 6, Name: "fff"},
				{AlertRuleID: 7, Name: "ggg"},
				{AlertRuleID: 8, Name: "hhh"},
				{AlertRuleID: 9, Name: "iii"},
				{AlertRuleID: 10, Name: "jjj"},
			},
			want: "- aaa\n- bbb\n- ccc\n- ddd\n- eee\n- ...",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := generateRuleList(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}