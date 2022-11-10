package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	findingmock "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/ca-risken/core/proto/project"
	"github.com/jarcoal/httpmock"
)

func TestNotificationAlert(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	now := time.Now()
	mockDB := mocks.MockAlertRepository{}
	svc := AlertService{repository: &mockDB}
	testFindingIDs := []uint64{}

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
			got := svc.NotificationAlert(context.Background(), c.alertCondition, c.alert, &[]model.AlertRule{}, &project.Project{}, &testFindingIDs)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestGetFindingDetailsForNotification(t *testing.T) {
	type inputParam struct {
		ProjectID  uint32
		FindingIDs *[]uint64
	}
	type mockGetFinding struct {
		Resp *finding.GetFindingResponse
		Err  error
	}
	type mockListFindingTag struct {
		Resp *finding.ListFindingTagResponse
		Err  error
	}
	cases := []struct {
		name           string
		input          inputParam
		getFinding     mockGetFinding
		listFindingTag mockListFindingTag

		want    *findingDetail
		wantErr bool
	}{
		{
			name:  "OK single data",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1}},
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			listFindingTag: mockListFindingTag{
				Resp: &finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{
						{FindingTagId: 1, Tag: "tag1"},
					},
				},
				Err: nil,
			},
			want: &findingDetail{
				FindingCount: 1,
				Exampls: []*findingExample{
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1"}},
				},
			},
			wantErr: false,
		},
		{
			name:  "OK multi datas",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 1, 1}},
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			listFindingTag: mockListFindingTag{
				Resp: &finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{
						{FindingTagId: 1, Tag: "tag1"},
						{FindingTagId: 2, Tag: "tag2"},
					},
				},
				Err: nil,
			},
			want: &findingDetail{
				FindingCount: 3,
				Exampls: []*findingExample{
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
				},
			},
			wantErr: false,
		},
		{
			name:  "OK over max findings(max=3)",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 1, 1, 1}},
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			listFindingTag: mockListFindingTag{
				Resp: &finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{
						{FindingTagId: 1, Tag: "tag1"},
						{FindingTagId: 2, Tag: "tag2"},
					},
				},
				Err: nil,
			},
			want: &findingDetail{
				FindingCount: 4,
				Exampls: []*findingExample{
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1", "tag2"}},
				},
			},
			wantErr: false,
		},
		{
			name:  "NG getFinding API error",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 1, 1}},
			getFinding: mockGetFinding{
				Resp: nil,
				Err:  errors.New("api error"),
			},
			listFindingTag: mockListFindingTag{},
			want:           nil,
			wantErr:        true,
		},
		{
			name:  "NG listFindingTag API error",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 1, 1}},
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			listFindingTag: mockListFindingTag{
				Resp: nil,
				Err:  errors.New("api error"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockFinding := findingmock.FindingServiceClient{}
			svc := AlertService{findingClient: &mockFinding}
			mockFinding.On("GetFinding", mock.Anything, mock.Anything).Return(c.getFinding.Resp, c.getFinding.Err)
			mockFinding.On("ListFindingTag", mock.Anything, mock.Anything).Return(c.listFindingTag.Resp, c.listFindingTag.Err)
			got, err := svc.getFindingDetailsForNotification(context.TODO(), c.input.ProjectID, c.input.FindingIDs)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: got=%+v, want=%+v", got, c.want)
			}
		})
	}
}
