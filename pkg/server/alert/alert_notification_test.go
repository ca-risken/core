package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
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
	now := time.Now()
	projectRelation := model.AlertCondNotification{ProjectID: 1, AlertConditionID: 2, NotificationID: 1}
	projectSuppressed := projectRelation
	projectSuppressed.CacheSecond = 1800
	projectSuppressed.NotifiedAt = now
	orgTarget := &db.OrgAlertNotificationTarget{OrganizationID: 10, OrganizationName: "org-name", ProjectID: 1, AlertConditionID: 2, NotificationID: 1, Type: "slack", NotifySetting: "org-10"}
	orgSuppressed := *orgTarget
	orgSuppressed.CacheSecond = 1800
	orgSuppressed.NotifiedAt = now
	cases := []struct {
		name                 string
		projectRelations     []model.AlertCondNotification
		projectListErr       error
		projectNotifications map[uint32]*model.Notification
		projectGetErr        map[uint32]error
		orgTargets           []*db.OrgAlertNotificationTarget
		orgListErr           error
		orgExists            map[uint32]bool
		orgLatest            map[uint32]*db.OrgAlertNotificationTarget
		sendErr              map[string]error
		projectUpdateErr     error
		orgUpdateErr         map[uint32]error
		existsNewFindings    bool
		wantSent             []string
		wantOrganizationName []string
		wantProjectUpdates   int
		wantOrgRechecks      []uint32
		wantOrgUpdates       []uint32
		wantErr              bool
	}{
		{
			name:                 "OK - project and organization targets with overlapping IDs",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:            map[uint32]bool{10: true},
			wantSent:             []string{"project", "org-10"},
			wantOrganizationName: []string{"", "org-name"},
			wantProjectUpdates:   1,
			wantOrgRechecks:      []uint32{10},
			wantOrgUpdates:       []uint32{10},
		},
		{
			name:            "OK - organization only",
			orgTargets:      []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:       map[uint32]bool{10: true},
			wantSent:        []string{"org-10"},
			wantOrgRechecks: []uint32{10},
			wantOrgUpdates:  []uint32{10},
		},
		{
			name:              "OK - new finding organization notification does not update project window",
			orgTargets:        []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:         map[uint32]bool{10: true},
			existsNewFindings: true,
			wantSent:          []string{"org-10"},
			wantOrgRechecks:   []uint32{10},
		},
		{
			name:                 "OK - project only",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			wantSent:             []string{"project"},
			wantProjectUpdates:   1,
		},
		{name: "OK - no targets"},
		{
			name:                 "OK - project suppression does not suppress organization",
			projectRelations:     []model.AlertCondNotification{projectSuppressed},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:            map[uint32]bool{10: true},
			wantSent:             []string{"org-10"},
			wantOrgRechecks:      []uint32{10},
			wantOrgUpdates:       []uint32{10},
		},
		{
			name:                 "OK - organization suppression does not suppress project",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{&orgSuppressed},
			orgExists:            map[uint32]bool{10: true},
			wantSent:             []string{"project"},
			wantProjectUpdates:   1,
			wantOrgRechecks:      []uint32{10},
		},
		{
			name:                 "NG - project send failure does not stop organization",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:            map[uint32]bool{10: true},
			sendErr:              map[string]error{"project": errors.New("send error")},
			wantSent:             []string{"project", "org-10"},
			wantOrgRechecks:      []uint32{10},
			wantOrgUpdates:       []uint32{10},
			wantErr:              true,
		},
		{
			name:                 "NG - organization send failure does not stop project",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:            map[uint32]bool{10: true},
			sendErr:              map[string]error{"org-10": errors.New("send error")},
			wantSent:             []string{"project", "org-10"},
			wantProjectUpdates:   1,
			wantOrgRechecks:      []uint32{10},
			wantErr:              true,
		},
		{
			name: "OK - all organizations are sent",
			orgTargets: []*db.OrgAlertNotificationTarget{
				orgTarget,
				{OrganizationID: 20, ProjectID: 1, AlertConditionID: 2, NotificationID: 1, Type: "slack", NotifySetting: "org-20"},
			},
			orgExists:       map[uint32]bool{10: true, 20: true},
			wantSent:        []string{"org-10", "org-20"},
			wantOrgRechecks: []uint32{10, 20},
			wantOrgUpdates:  []uint32{10, 20},
		},
		{
			name:            "OK - send-time removed organization target is skipped",
			orgTargets:      []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:       map[uint32]bool{10: false},
			wantOrgRechecks: []uint32{10},
		},
		{
			name:            "OK - send-time organization suppression uses latest state",
			orgTargets:      []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:       map[uint32]bool{10: true},
			orgLatest:       map[uint32]*db.OrgAlertNotificationTarget{10: &orgSuppressed},
			wantOrgRechecks: []uint32{10},
		},
		{
			name:                 "OK - unsupported targets do not update suppression state",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "email", NotifySetting: "project-email"}},
			orgTargets:           []*db.OrgAlertNotificationTarget{{OrganizationID: 10, ProjectID: 1, AlertConditionID: 2, NotificationID: 1, Type: "email", NotifySetting: "org-email"}},
			orgExists:            map[uint32]bool{10: true},
			wantOrgRechecks:      []uint32{10},
		},
		{
			name:            "NG - project resolution error still sends organization",
			projectListErr:  errors.New("DB error"),
			orgTargets:      []*db.OrgAlertNotificationTarget{orgTarget},
			orgExists:       map[uint32]bool{10: true},
			wantSent:        []string{"org-10"},
			wantOrgRechecks: []uint32{10},
			wantOrgUpdates:  []uint32{10},
			wantErr:         true,
		},
		{
			name:                 "NG - organization resolution error still sends project",
			projectRelations:     []model.AlertCondNotification{projectRelation},
			projectNotifications: map[uint32]*model.Notification{1: {NotificationID: 1, Type: "slack", NotifySetting: "project"}},
			orgListErr:           errors.New("DB error"),
			wantSent:             []string{"project"},
			wantProjectUpdates:   1,
			wantErr:              true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			projectRelations := c.projectRelations
			mockDB.On("ListAlertCondNotificationForNotification", mock.Anything, uint32(1), uint32(2)).Return(&projectRelations, c.projectListErr).Once()
			for _, relation := range c.projectRelations {
				mockDB.On("GetNotification", mock.Anything, uint32(1), relation.NotificationID).Return(c.projectNotifications[relation.NotificationID], c.projectGetErr[relation.NotificationID]).Once()
			}
			mockDB.On("ListOrgAlertNotificationTarget", mock.Anything, uint32(1), uint32(2)).Return(c.orgTargets, c.orgListErr).Once()
			for _, organizationID := range c.wantOrgRechecks {
				var latest *db.OrgAlertNotificationTarget
				for _, candidate := range c.orgTargets {
					if candidate.OrganizationID == organizationID {
						latest = candidate
						break
					}
				}
				if candidate := c.orgLatest[organizationID]; candidate != nil {
					latest = candidate
				}
				if c.orgExists[organizationID] {
					mockDB.On("GetOrgAlertProjectNotificationTarget", mock.Anything, organizationID, uint32(1), uint32(2), uint32(1)).Return(latest, nil).Once()
				} else {
					mockDB.On("GetOrgAlertProjectNotificationTarget", mock.Anything, organizationID, uint32(1), uint32(2), uint32(1)).Return(nil, gorm.ErrRecordNotFound).Once()
				}
			}
			for i := 0; i < c.wantProjectUpdates; i++ {
				mockDB.On("UpsertAlertCondNotification", mock.Anything, mock.Anything).Return(&model.AlertCondNotification{}, c.projectUpdateErr).Once()
			}
			for _, organizationID := range c.wantOrgUpdates {
				mockDB.On("UpdateOrgAlertProjectNotificationNotifiedAt", mock.Anything, organizationID, uint32(1), uint32(1), mock.Anything).Return(c.orgUpdateErr[organizationID]).Once()
			}
			var sent []string
			var organizationNames []string
			svc.sendAlertNotification = func(_ context.Context, setting, organizationName string, _ *model.Alert, _ *project.Project, _ *[]model.AlertRule, _ *findingDetail, _ notificationReason) error {
				sent = append(sent, setting)
				organizationNames = append(organizationNames, organizationName)
				return c.sendErr[setting]
			}
			findingIDs := []uint64{}
			err := svc.NotificationAlert(context.Background(), &model.AlertCondition{ProjectID: 1, AlertConditionID: 2}, &model.Alert{}, &[]model.AlertRule{}, &project.Project{ProjectId: 1}, &findingIDs, c.existsNewFindings)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(sent, c.wantSent) {
				t.Fatalf("Unexpected sends: want=%v, got=%v", c.wantSent, sent)
			}
			if c.wantOrganizationName != nil && !reflect.DeepEqual(organizationNames, c.wantOrganizationName) {
				t.Fatalf("Unexpected organization names: want=%v, got=%v", c.wantOrganizationName, organizationNames)
			}
		})
	}
}

func TestNotificationAlertDefersOrgProjectWindowUpdate(t *testing.T) {
	mockDB := mocks.NewAlertRepository(t)
	svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
	conditions := []*model.AlertCondition{
		{ProjectID: 1, AlertConditionID: 2},
		{ProjectID: 1, AlertConditionID: 3},
	}
	alertData := &model.Alert{}
	rules := &[]model.AlertRule{}
	projectData := &project.Project{ProjectId: 1}
	findingIDs := []uint64{}
	orgTarget := &db.OrgAlertNotificationTarget{
		OrganizationID:   10,
		OrganizationName: "org-name",
		ProjectID:        1,
		AlertConditionID: 2,
		NotificationID:   1,
		Type:             "slack",
		NotifySetting:    "org-10",
	}
	collector := newOrgProjectNotificationUpdateCollector()

	for _, condition := range conditions {
		target := *orgTarget
		target.AlertConditionID = condition.AlertConditionID
		mockDB.On("ListAlertCondNotificationForNotification", mock.Anything, condition.ProjectID, condition.AlertConditionID).Return(&[]model.AlertCondNotification{}, nil).Once()
		mockDB.On("ListOrgAlertNotificationTarget", mock.Anything, condition.ProjectID, condition.AlertConditionID).Return([]*db.OrgAlertNotificationTarget{&target}, nil).Once()
		mockDB.On("GetOrgAlertProjectNotificationTarget", mock.Anything, target.OrganizationID, target.ProjectID, target.AlertConditionID, target.NotificationID).Return(&target, nil).Once()
	}
	mockDB.On("UpdateOrgAlertProjectNotificationNotifiedAt", mock.Anything, uint32(10), uint32(1), uint32(1), mock.Anything).Return(nil).Once()

	var sent []string
	svc.sendAlertNotification = func(_ context.Context, setting, _ string, _ *model.Alert, _ *project.Project, _ *[]model.AlertRule, _ *findingDetail, _ notificationReason) error {
		sent = append(sent, setting)
		return nil
	}

	for _, condition := range conditions {
		if err := svc.notificationAlert(context.Background(), condition, alertData, rules, projectData, &findingIDs, false, collector); err != nil {
			t.Fatalf("Unexpected notification error: %v", err)
		}
	}
	if !reflect.DeepEqual(sent, []string{"org-10", "org-10"}) {
		t.Fatalf("Unexpected sends: got=%v", sent)
	}
	if err := collector.flush(context.Background(), mockDB); err != nil {
		t.Fatalf("Unexpected flush error: %v", err)
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
			name:  "OK getFinding API error skipped",
			input: inputParam{ProjectID: 1, FindingIDs: &[]uint64{1, 2}},
			want: &findingDetail{
				FindingCount: 2,
				Exampls: []*findingExample{
					{FindingID: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9, Tags: []string{"tag2"}},
				},
			},
			wantErr: false,
		},
		{
			name:  "OK listFindingTag API error skipped",
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
			want: &findingDetail{
				FindingCount: 3,
				Exampls: []*findingExample{
					{FindingID: 1, Description: "desc", ResourceName: "rn", DataSource: "ds", Score: 1.0},
				},
			},
			wantErr: false,
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
			if c.name == "OK getFinding API error skipped" {
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 1
				})).Return(nil, errors.New("api error")).Once()
				mockFinding.On("GetFinding", mock.Anything, mock.MatchedBy(func(req *finding.GetFindingRequest) bool {
					return req.ProjectId == 1 && req.FindingId == 2
				})).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 2, Description: "desc-2", ResourceName: "rn-2", DataSource: "ds", Score: 0.9},
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
