package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/test"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	findingmock "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/ca-risken/core/proto/iam"
	iammock "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
	"github.com/jarcoal/httpmock"
)

func TestNotificationAlert(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	now := time.Now()
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
			name:                               "Error UpsertAlertCondNotification Failed",
			alertCondition:                     &model.AlertCondition{AlertConditionID: 1},
			alert:                              &model.Alert{},
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
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockListAlertCondNotification != nil || c.mockListAlertCondNotificationErr != nil {
				mockDB.On("ListAlertCondNotification", test.RepeatMockAnything(6)...).Return(c.mockListAlertCondNotification, c.mockListAlertCondNotificationErr).Once()
			}
			if c.mockGetNotification != nil || c.mockGetNotificationErr != nil {
				mockDB.On("GetNotification", test.RepeatMockAnything(3)...).Return(c.mockGetNotification, c.mockGetNotificationErr).Once()
			}
			if c.mockUpsertAlertCondNotification != nil || c.mockUpsertAlertCondNotificationErr != nil {
				mockDB.On("UpsertAlertCondNotification", test.RepeatMockAnything(2)...).Return(c.mockUpsertAlertCondNotification, c.mockUpsertAlertCondNotificationErr).Once()
			}
			got := svc.NotificationAlert(context.Background(), c.alertCondition, c.alert, &[]model.AlertRule{}, &project.Project{}, &testFindingIDs, false)
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
	type mockGetAlertAISummary struct {
		Resp *finding.GetAlertAISummaryResponse
		Err  error
	}
	cases := []struct {
		name              string
		input             inputParam
		aiSummaryEnabled  bool
		summaryLanguage   string
		getFinding        mockGetFinding
		getAlertAISummary *mockGetAlertAISummary
		listFindingTag    mockListFindingTag

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
			name:             "OK with alert AI summary",
			input:            inputParam{ProjectID: 1, FindingIDs: &[]uint64{1}},
			aiSummaryEnabled: true,
			summaryLanguage:  "ja",
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			getAlertAISummary: &mockGetAlertAISummary{
				Resp: &finding.GetAlertAISummaryResponse{AiSummary: "summary"},
				Err:  nil,
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
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0, Tags: []string{"tag1"}, AISummary: "summary"},
				},
			},
			wantErr: false,
		},
		{
			name:             "OK alert AI summary failure continues",
			input:            inputParam{ProjectID: 1, FindingIDs: &[]uint64{1}},
			aiSummaryEnabled: true,
			summaryLanguage:  "ja",
			getFinding: mockGetFinding{
				Resp: &finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
				Err: nil,
			},
			getAlertAISummary: &mockGetAlertAISummary{
				Resp: nil,
				Err:  errors.New("ai error"),
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
				},
			},
			wantErr: false,
		},
		{
			name:  "OK over max findings(max=1)",
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
				},
			},
			wantErr: false,
		},
		{
			name:  "OK sort by score and created_at desc",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 2, 3}},
			want: &findingDetail{
				FindingCount: 3,
				Exampls: []*findingExample{
					{FindingID: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9, CreatedAt: 200, Tags: []string{"tag2"}},
				},
			},
			wantErr: false,
		},
		{
			name:  "OK getFinding API error skipped",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 2}},
			want: &findingDetail{
				FindingCount: 2,
				Exampls: []*findingExample{
					{FindingID: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9, CreatedAt: 200, Tags: []string{"tag2"}},
				},
			},
			wantErr: false,
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
			svc := AlertService{
				findingClient:    &mockFinding,
				logger:           logging.NewLogger(),
				aiSummaryEnabled: c.aiSummaryEnabled,
				summaryLanguage:  c.summaryLanguage,
			}
			if c.name == "OK sort by score and created_at desc" {
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 1
				})).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1, Description: "desc-1", ResourceName: "rn-1", DataSource: "ds", Score: 0.8, CreatedAt: 300},
				}, nil).Once()
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 2
				})).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9, CreatedAt: 200},
				}, nil).Once()
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 3
				})).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 3, Description: "desc-3", ResourceName: "rn-3", DataSource: "ds", Score: 0.9, CreatedAt: 100},
				}, nil).Once()
				mockFinding.On("ListFindingTag", mock.Anything, mock.MatchedBy(func(req *finding.ListFindingTagRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 2
				})).Return(&finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{{FindingTagId: 2, Tag: "tag2"}},
				}, nil).Once()
			} else if c.name == "OK getFinding API error skipped" {
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 1
				})).Return(nil, errors.New("api error")).Once()
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 2
				})).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9, CreatedAt: 200},
				}, nil).Once()
				mockFinding.On("ListFindingTag", mock.Anything, mock.MatchedBy(func(req *finding.ListFindingTagRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 2
				})).Return(&finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{{FindingTagId: 2, Tag: "tag2"}},
				}, nil).Once()
			} else {
				mockFinding.On("GetFinding", mock.Anything, mock.Anything).Return(c.getFinding.Resp, c.getFinding.Err)
				if c.getAlertAISummary != nil {
					mockFinding.On("GetAlertAISummary", mock.Anything, mock.Anything).Return(c.getAlertAISummary.Resp, c.getAlertAISummary.Err)
				}
				mockFinding.On("ListFindingTag", mock.Anything, mock.Anything).Return(c.listFindingTag.Resp, c.listFindingTag.Err)
			}
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

func TestPutNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name        string
		input       *alert.PutNotificationRequest
		want        *alert.PutNotificationResponse
		wantErr     bool
		mockGetResp *model.Notification
		mockGetErr  error
		mockUpResp  *model.Notification
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &alert.PutNotificationRequest{Notification: &alert.NotificationForUpsert{ProjectId: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`}},
			want:       &alert.PutNotificationResponse{Notification: &alert.Notification{ProjectId: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.Notification{ProjectID: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &alert.PutNotificationRequest{Notification: &alert.NotificationForUpsert{NotificationId: 1001, ProjectId: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`}},
			want:        &alert.PutNotificationResponse{Notification: &alert.Notification{NotificationId: 1001, ProjectId: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url":"https://e**********","channel_id":"","data":{},"locale":""}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Notification{NotificationID: 1001, ProjectID: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
			mockUpResp:  &model.Notification{NotificationID: 1001, ProjectID: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "NG Update (Notification Not Found)",
			input:       &alert.PutNotificationRequest{Notification: &alert.NotificationForUpsert{NotificationId: 1001, ProjectId: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`}},
			want:        &alert.PutNotificationResponse{},
			wantErr:     true,
			mockGetResp: &model.Notification{},
			mockGetErr:  gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetNotification", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertNotification", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}

			got, err := svc.PutNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestRequestProjectRoleNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// External HTTP mock for slack notification
	httpmock.RegisterResponder("POST", "https://example.com",
		httpmock.NewStringResponder(200, "ok"))

	var ctx context.Context
	type mockListNotification struct {
		Resp *[]model.Notification
		Err  error
	}
	type mockListProject struct {
		Resp *project.ListProjectResponse
		Err  error
	}
	type mockGetUser struct {
		Resp *iam.GetUserResponse
		Err  error
	}
	now := time.Now()
	cases := []struct {
		name             string
		input            *alert.RequestProjectRoleNotificationRequest
		wantErr          bool
		listNotification mockListNotification
		listProject      mockListProject
		getUser          mockGetUser
	}{
		{
			name:    "OK Request project role",
			input:   &alert.RequestProjectRoleNotificationRequest{ProjectId: 1001, UserId: 1001},
			wantErr: false,
			listNotification: mockListNotification{
				Resp: &[]model.Notification{
					{
						ProjectID:     1001,
						Name:          "name",
						Type:          "slack",
						NotifySetting: `{"webhook_url": "https://example.com"}`,
						CreatedAt:     now,
						UpdatedAt:     now,
					},
				},
				Err: nil,
			},
			getUser: mockGetUser{
				Resp: &iam.GetUserResponse{
					User: &iam.User{UserId: 1001, Name: "userName"},
				},
				Err: nil,
			},
			listProject: mockListProject{
				Resp: &project.ListProjectResponse{
					Project: []*project.Project{
						{ProjectId: 1001, Name: "projectName"},
					},
				},
				Err: nil,
			},
		},
		{
			name:    "NG unimplemented notification type",
			input:   &alert.RequestProjectRoleNotificationRequest{ProjectId: 1001, UserId: 1001},
			wantErr: true,
			listNotification: mockListNotification{
				Resp: &[]model.Notification{
					{ProjectID: 1001, Name: "name", Type: "unimplemented", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
				},
				Err: nil,
			},
			listProject: mockListProject{
				Resp: &project.ListProjectResponse{
					Project: []*project.Project{
						{ProjectId: 1001, Name: "projectName"},
					},
				},
				Err: nil,
			},
			getUser: mockGetUser{
				Resp: &iam.GetUserResponse{
					User: &iam.User{UserId: 1001, Name: "userName"},
				},
				Err: nil,
			},
		},
		{
			name:    "NG ListNotification (Notification Not Found)",
			input:   &alert.RequestProjectRoleNotificationRequest{ProjectId: 1001, UserId: 1001},
			wantErr: true,
			listNotification: mockListNotification{
				Resp: &[]model.Notification{},
				Err:  gorm.ErrRecordNotFound,
			},
		},
		{
			name:    "NG ListProject (API Error)",
			input:   &alert.RequestProjectRoleNotificationRequest{ProjectId: 1001, UserId: 1001},
			wantErr: true,
			listNotification: mockListNotification{
				Resp: &[]model.Notification{
					{ProjectID: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
				},
				Err: nil,
			},
			listProject: mockListProject{
				Resp: &project.ListProjectResponse{
					Project: []*project.Project{
						{ProjectId: 1001, Name: "projectName"},
					},
				},
				Err: errors.New("api error"),
			},
		},
		{
			name:    "NG GetUser (API Error)",
			input:   &alert.RequestProjectRoleNotificationRequest{ProjectId: 1001, UserId: 1001},
			wantErr: true,
			listNotification: mockListNotification{
				Resp: &[]model.Notification{
					{ProjectID: 1001, Name: "name", Type: "slack", NotifySetting: `{"webhook_url": "https://example.com"}`, CreatedAt: now, UpdatedAt: now},
				},
				Err: errors.New("api error"),
			},
			listProject: mockListProject{
				Resp: &project.ListProjectResponse{
					Project: []*project.Project{
						{ProjectId: 1001, Name: "projectName"},
					},
				},
				Err: nil,
			},
			getUser: mockGetUser{
				Resp: &iam.GetUserResponse{
					User: &iam.User{UserId: 1001, Name: "userName"},
				},
				Err: gorm.ErrRecordNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			mockDB.On("ListNotification", test.RepeatMockAnything(5)...).Return(c.listNotification.Resp, c.listNotification.Err).Once()
			mockProject := projectmock.ProjectServiceClient{}
			mockProject.On("ListProject", mock.Anything, mock.Anything).Return(c.listProject.Resp, c.listProject.Err)
			mockIAM := iammock.IAMServiceClient{}
			mockIAM.On("GetUser", mock.Anything, mock.Anything).Return(c.getUser.Resp, c.getUser.Err)

			svc := AlertService{projectClient: &mockProject, iamClient: &mockIAM, repository: mockDB, logger: logging.NewLogger()}
			_, err := svc.RequestProjectRoleNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
