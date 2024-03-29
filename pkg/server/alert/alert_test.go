package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/alert"
	"gorm.io/gorm"
)

/*
 * Alert
 */

func TestListAlert(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListAlertRequest
		want         *alert.ListAlertResponse
		mockResponce *[]model.Alert
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertRequest{ProjectId: 1, Status: []alert.Status{alert.Status_ACTIVE}, Severity: []string{"high"}, Description: ""},
			want:         &alert.ListAlertResponse{Alert: []*alert.Alert{{AlertId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.Alert{{AlertID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertRequest{ProjectId: 1, Status: []alert.Status{alert.Status_ACTIVE}, Severity: []string{"high"}, Description: ""},
			want:      &alert.ListAlertResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlert", test.RepeatMockAnything(7)...).Return(c.mockResponce, c.mockError).Once()
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
			input: &alert.ListAlertRequest{ProjectId: 1, Severity: []string{"high"}, Description: "desc", Status: []alert.Status{alert.Status_ACTIVE}, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertRequest{ProjectId: 1, Severity: []string{"high"}, Description: "desc", Status: []alert.Status{alert.Status_ACTIVE}, FromAt: now.Unix(), ToAt: now.Unix()},
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
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlert", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
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

	cases := []struct {
		name                        string
		input                       *alert.PutAlertRequest
		want                        *alert.PutAlertResponse
		wantErr                     bool
		mockGetResp                 *model.Alert
		mockGetErr                  error
		mockUpResp                  *model.Alert
		mockUpErr                   error
		mockListRelAlertFindingResp *[]model.RelAlertFinding
		mockListRelAlertFindingErr  error
		mockUpHistoryResp           *model.AlertHistory
		mockUpHistoryErr            error
	}{
		{
			name:                        "OK Insert",
			input:                       &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE}},
			want:                        &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp:                  &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			mockListRelAlertFindingResp: &[]model.RelAlertFinding{{FindingID: 1001, ProjectID: 1001}},
			mockUpHistoryResp:           &model.AlertHistory{AlertHistoryID: 1001, HistoryType: "created"},
		},
		{
			name:                        "OK Update",
			input:                       &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE}},
			want:                        &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp:                 &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			mockUpResp:                  &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			mockListRelAlertFindingResp: &[]model.RelAlertFinding{{FindingID: 1001, ProjectID: 1001}},
			mockListRelAlertFindingErr:  nil,
			mockUpHistoryResp:           &model.AlertHistory{AlertHistoryID: 1001, HistoryType: "updated"},
			mockUpHistoryErr:            nil,
		},
		{
			name:       "NG No record found with alertID",
			input:      &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE}},
			want:       nil,
			wantErr:    true,
			mockGetErr: gorm.ErrRecordNotFound,
		},
		{
			name:                        "NG failed listRelAlertFinding",
			input:                       &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE}},
			want:                        &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			wantErr:                     true,
			mockUpResp:                  &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			mockListRelAlertFindingResp: nil,
			mockListRelAlertFindingErr:  errors.New("something error"),
		},
		{
			name:                        "NG failed putAlertHistory",
			input:                       &alert.PutAlertRequest{Alert: &alert.AlertForUpsert{ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE}},
			want:                        &alert.PutAlertResponse{Alert: &alert.Alert{AlertId: 1001, ProjectId: 1001, AlertConditionId: 1001, Description: "desc", Severity: "high", Status: alert.Status_ACTIVE, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			wantErr:                     true,
			mockUpResp:                  &model.Alert{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now},
			mockListRelAlertFindingResp: &[]model.RelAlertFinding{{FindingID: 1001, ProjectID: 1001}},
			mockListRelAlertFindingErr:  nil,
			mockUpHistoryResp:           nil,
			mockUpHistoryErr:            errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetAlert", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlert", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			if c.mockListRelAlertFindingResp != nil || c.mockListRelAlertFindingErr != nil {
				mockDB.On("ListRelAlertFinding", test.RepeatMockAnything(6)...).Return(c.mockListRelAlertFindingResp, c.mockListRelAlertFindingErr).Once()
			}
			if c.mockUpHistoryResp != nil || c.mockUpHistoryErr != nil {
				mockDB.On("UpsertAlertHistory", test.RepeatMockAnything(2)...).Return(c.mockUpHistoryResp, c.mockUpHistoryErr).Once()
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

func TestPutAlertFirstViewedAt(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name         string
		input        *alert.PutAlertFirstViewedAtRequest
		wantErr      bool
		mockListResp *[]model.Alert
		mockListErr  error
		callUpdate   bool
		mockUpErr    error
	}{
		{
			name:         "OK",
			input:        &alert.PutAlertFirstViewedAtRequest{ProjectId: 1001, AlertId: 1001},
			mockListResp: &[]model.Alert{{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now}},
			callUpdate:   true,
		},
		{
			name:         "OK Already set",
			input:        &alert.PutAlertFirstViewedAtRequest{ProjectId: 1001, AlertId: 1001},
			mockListResp: &[]model.Alert{{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now, FirstViewedAt: &now}},
		},
		{
			name:    "NG Validation Error",
			input:   &alert.PutAlertFirstViewedAtRequest{AlertId: 1001},
			wantErr: true,
		},
		{
			name:        "NG GetAlert Error",
			input:       &alert.PutAlertFirstViewedAtRequest{ProjectId: 1001, AlertId: 1001},
			mockListErr: errors.New("something error"),
			wantErr:     true,
		},
		{
			name:         "NG Update Error",
			input:        &alert.PutAlertFirstViewedAtRequest{ProjectId: 1001, AlertId: 1001},
			mockListResp: &[]model.Alert{{AlertID: 1001, ProjectID: 1001, AlertConditionID: 1001, Description: "desc", Severity: "high", Status: "ACTIVE", CreatedAt: now, UpdatedAt: now}},
			callUpdate:   true,
			mockUpErr:    errors.New("something error"),
			wantErr:      true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockListResp != nil || c.mockListErr != nil {
				mockDB.On("ListAlert", test.RepeatMockAnything(7)...).Return(c.mockListResp, c.mockListErr).Once()
			}
			if c.callUpdate {
				mockDB.On("UpdateAlertFirstViewedAt", test.RepeatMockAnything(4)...).Return(c.mockUpErr).Once()
			}
			_, err := svc.PutAlertFirstViewedAt(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestDeleteAlert(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name            string
		input           *alert.DeleteAlertRequest
		wantErr         bool
		deleteAlertCall bool
		mockErr         error
	}{
		{
			name:            "OK",
			input:           &alert.DeleteAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr:         false,
			deleteAlertCall: true,
			mockErr:         nil,
		},
		{
			name:            "NG validation error",
			input:           &alert.DeleteAlertRequest{ProjectId: 1},
			deleteAlertCall: false,
			wantErr:         true,
		},
		{
			name:            "Invalid DB error",
			input:           &alert.DeleteAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr:         true,
			deleteAlertCall: true,
			mockErr:         gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteAlertCall {
				mockDB.On("DeleteAlert", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
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
 * AlertHistory
 */

func TestListAlertHistory(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name                    string
		input                   *alert.ListAlertHistoryRequest
		want                    *alert.ListAlertHistoryResponse
		wantErr                 bool
		mockListResponce        *[]model.AlertHistory
		mockListError           error
		mockListCreatedResponce *[]model.AlertHistory
		mockListCreatedError    error
	}{
		{
			name:  "OK first list result contains created",
			input: &alert.ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001},
			want: &alert.ListAlertHistoryResponse{AlertHistory: []*alert.AlertHistory{
				{AlertHistoryId: 1001, HistoryType: "created", FindingHistory: `{"count":3}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{AlertHistoryId: 1002, HistoryType: "updated", FindingHistory: `{"count":0}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockListResponce: &[]model.AlertHistory{
				{AlertHistoryID: 1001, HistoryType: "created", FindingHistory: `{"finding_id":[1,2,3]}`, CreatedAt: now, UpdatedAt: now},
				{AlertHistoryID: 1002, HistoryType: "updated", FindingHistory: `{}`, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:  "OK first list result doesn't contain created",
			input: &alert.ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001},
			want: &alert.ListAlertHistoryResponse{AlertHistory: []*alert.AlertHistory{
				{AlertHistoryId: 1001, HistoryType: "updated", FindingHistory: `{"count":1}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{AlertHistoryId: 1002, HistoryType: "updated", FindingHistory: `{"count":0}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{AlertHistoryId: 1003, HistoryType: "created", FindingHistory: `{"count":2}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockListResponce: &[]model.AlertHistory{
				{AlertHistoryID: 1001, HistoryType: "updated", FindingHistory: `{"finding_id":[1]}`, CreatedAt: now, UpdatedAt: now},
				{AlertHistoryID: 1002, HistoryType: "updated", FindingHistory: `{}`, CreatedAt: now, UpdatedAt: now},
			},
			mockListCreatedResponce: &[]model.AlertHistory{
				{AlertHistoryID: 1003, HistoryType: "created", FindingHistory: `{"finding_id":[1,2]}`, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:          "NG error listAlertHistory",
			input:         &alert.ListAlertHistoryRequest{ProjectId: 1, AlertId: 1001},
			want:          nil,
			wantErr:       true,
			mockListError: errors.New("somethingError"),
		},
		{
			name:    "NG error listAlertHistory only Created",
			input:   &alert.ListAlertHistoryRequest{ProjectId: 1, AlertId: 1001},
			want:    nil,
			wantErr: true,
			mockListResponce: &[]model.AlertHistory{
				{AlertHistoryID: 1001, HistoryType: "updated", FindingHistory: `{"finding_id":[1]}`, CreatedAt: now, UpdatedAt: now},
				{AlertHistoryID: 1002, HistoryType: "updated", FindingHistory: `{}`, CreatedAt: now, UpdatedAt: now},
			},
			mockListCreatedError: errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockListResponce != nil || c.mockListError != nil {
				mockDB.On("ListAlertHistory", test.RepeatMockAnything(8)...).Return(c.mockListResponce, c.mockListError).Once()
			}
			if c.mockListCreatedResponce != nil || c.mockListCreatedError != nil {
				mockDB.On("ListAlertHistory", test.RepeatMockAnything(8)...).Return(c.mockListCreatedResponce, c.mockListCreatedError).Once()
			}
			result, err := svc.ListAlertHistory(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestGetAlertHistory(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetAlertHistoryRequest
		want         *alert.GetAlertHistoryResponse
		mockResponce *model.AlertHistory
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			want:         &alert.GetAlertHistoryResponse{AlertHistory: &alert.AlertHistory{AlertHistoryId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AlertHistory{AlertHistoryID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 9999},
			want:      &alert.GetAlertHistoryResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertHistory", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlertHistory(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlertHistory(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutAlertHistoryRequest
		want        *alert.PutAlertHistoryResponse
		wantErr     bool
		mockGetResp *model.AlertHistory
		mockGetErr  error
		mockUpResp  *model.AlertHistory
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &alert.PutAlertHistoryRequest{AlertHistory: &alert.AlertHistoryForUpsert{ProjectId: 1001, AlertId: 1001, Description: "desc", Severity: "high", HistoryType: "created"}},
			want:       &alert.PutAlertHistoryResponse{AlertHistory: &alert.AlertHistory{AlertHistoryId: 1001, ProjectId: 1001, AlertId: 1001, Description: "desc", Severity: "high", HistoryType: "created", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.AlertHistory{AlertHistoryID: 1001, ProjectID: 1001, AlertID: 1001, Description: "desc", Severity: "high", HistoryType: "created", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:       "OK Update",
			input:      &alert.PutAlertHistoryRequest{AlertHistory: &alert.AlertHistoryForUpsert{AlertHistoryId: 1001, ProjectId: 1001, AlertId: 1001, Description: "desc", Severity: "high", HistoryType: "created"}},
			want:       &alert.PutAlertHistoryResponse{AlertHistory: &alert.AlertHistory{AlertHistoryId: 1001, ProjectId: 1001, AlertId: 1001, Description: "desc", Severity: "high", HistoryType: "created", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.AlertHistory{AlertHistoryID: 1001, ProjectID: 1001, AlertID: 1001, Description: "desc", Severity: "high", HistoryType: "created", CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertHistory", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlertHistory(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlertHistory(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                   string
		input                  *alert.DeleteAlertHistoryRequest
		wantErr                bool
		deleteAlertHistoryCall bool
		mockErr                error
	}{
		{
			name:                   "OK",
			input:                  &alert.DeleteAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr:                false,
			deleteAlertHistoryCall: true,
			mockErr:                nil,
		},
		{
			name:                   "NG validation error",
			input:                  &alert.DeleteAlertHistoryRequest{ProjectId: 1},
			deleteAlertHistoryCall: false,
			wantErr:                true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &alert.DeleteAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr:                true,
			deleteAlertHistoryCall: true,
			mockErr:                gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		mockDB := mocks.NewAlertRepository(t)
		svc := AlertService{repository: mockDB}
		t.Run(c.name, func(t *testing.T) {
			if c.deleteAlertHistoryCall {
				mockDB.On("DeleteAlertHistory", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteAlertHistory(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlertHistory(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		input    *model.AlertHistory
		getCount bool
		want     *alert.AlertHistory
		wantErr  bool
	}{
		{
			name:  "OK not convert finding history",
			input: &model.AlertHistory{AlertHistoryID: 1001, FindingHistory: `{"finding_id":[1,2,3]}`, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertHistory{AlertHistoryId: 1001, FindingHistory: `{"finding_id":[1,2,3]}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:     "OK convert finding history",
			input:    &model.AlertHistory{AlertHistoryID: 1001, FindingHistory: `{"finding_id":[1,2,3]}`, CreatedAt: now, UpdatedAt: now},
			getCount: true,
			want:     &alert.AlertHistory{AlertHistoryId: 1001, FindingHistory: `{"count":3}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertHistory{},
		},
		{
			name:     "NG convert finding history error",
			input:    &model.AlertHistory{AlertHistoryID: 1001, FindingHistory: `invalid`},
			getCount: true,
			want:     nil,
			wantErr:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := convertAlertHistory(c.input, c.getCount)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: err=%+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestConvertIDsToCount(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "OK convert",
			input: `{"finding_id":[1,2,3]}`,
			want:  `{"count":3}`,
		},
		{
			name:  "OK empty",
			input: `{}`,
			want:  `{"count":0}`,
		},
		{
			name:    "NG",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := convertIDsToCount(c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: err=%+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * RelAlertFinding
 */

func TestListRelAlertFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListRelAlertFindingRequest
		want         *alert.ListRelAlertFindingResponse
		mockResponce *[]model.RelAlertFinding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001},
			want:         &alert.ListRelAlertFindingResponse{RelAlertFinding: []*alert.RelAlertFinding{{AlertId: 1001, FindingId: 1001, ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertId: 1002, FindingId: 1001, ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.RelAlertFinding{{AlertID: 1001, FindingID: 1001, ProjectID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertID: 1002, FindingID: 1001, ProjectID: 1001, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001},
			want:      &alert.ListRelAlertFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		mockDB := mocks.NewAlertRepository(t)
		svc := AlertService{repository: mockDB}
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListRelAlertFinding", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListRelAlertFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListRelAlertFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListRelAlertFindingRequest
		want  *alert.ListRelAlertFindingRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListRelAlertFindingRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetRelAlertFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetRelAlertFindingRequest
		want         *alert.GetRelAlertFindingResponse
		mockResponce *model.RelAlertFinding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001},
			want:         &alert.GetRelAlertFindingResponse{RelAlertFinding: &alert.RelAlertFinding{ProjectId: 1001, AlertId: 1001, FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.RelAlertFinding{ProjectID: 1001, AlertID: 1001, FindingID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetRelAlertFindingRequest{ProjectId: 1001, AlertId: 9999, FindingId: 9999},
			want:      &alert.GetRelAlertFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRelAlertFinding", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetRelAlertFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutRelAlertFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutRelAlertFindingRequest
		want        *alert.PutRelAlertFindingResponse
		wantErr     bool
		mockGetResp *model.RelAlertFinding
		mockGetErr  error
		mockUpResp  *model.RelAlertFinding
		mockUpErr   error
	}{
		{
			name:       "OK Upsert",
			input:      &alert.PutRelAlertFindingRequest{RelAlertFinding: &alert.RelAlertFindingForUpsert{ProjectId: 1001, AlertId: 1001, FindingId: 1001}},
			want:       &alert.PutRelAlertFindingResponse{RelAlertFinding: &alert.RelAlertFinding{ProjectId: 1001, AlertId: 1001, FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.RelAlertFinding{ProjectID: 1001, AlertID: 1001, FindingID: 1001, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertRelAlertFinding", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutRelAlertFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteRelAlertFinding(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                      string
		input                     *alert.DeleteRelAlertFindingRequest
		wantErr                   bool
		deleteRelAlertFindingCall bool
		mockErr                   error
	}{
		{
			name:                      "OK",
			input:                     &alert.DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001, FindingId: 1001},
			wantErr:                   false,
			deleteRelAlertFindingCall: true,
			mockErr:                   nil,
		},
		{
			name:                      "NG validation error",
			input:                     &alert.DeleteRelAlertFindingRequest{ProjectId: 1},
			deleteRelAlertFindingCall: false,
			wantErr:                   true,
		},
		{
			name:                      "Invalid DB error",
			input:                     &alert.DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001, FindingId: 1001},
			wantErr:                   true,
			deleteRelAlertFindingCall: true,
			mockErr:                   gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteRelAlertFindingCall {
				mockDB.On("DeleteRelAlertFinding", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteRelAlertFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertRelAlertFinding(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.RelAlertFinding
		want  *alert.RelAlertFinding
	}{
		{
			name:  "OK convert unix time",
			input: &model.RelAlertFinding{ProjectID: 1001, AlertID: 1001, FindingID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.RelAlertFinding{ProjectId: 1001, AlertId: 1001, FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.RelAlertFinding{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertRelAlertFinding(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * AlertCondition
 */

func TestListAlertCondition(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListAlertConditionRequest
		want         *alert.ListAlertConditionResponse
		mockResponce *[]model.AlertCondition
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high"}, Enabled: true},
			want:         &alert.ListAlertConditionResponse{AlertCondition: []*alert.AlertCondition{{AlertConditionId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertConditionId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.AlertCondition{{AlertConditionID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertConditionID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high"}, Enabled: true},
			want:      &alert.ListAlertConditionResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondition", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlertCondition(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertConditionRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertConditionRequest
		want  *alert.ListAlertConditionRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertConditionRequest{ProjectId: 1001, Enabled: true, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertConditionRequest{ProjectId: 1001, Enabled: true, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertConditionRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlertCondition(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetAlertConditionRequest
		want         *alert.GetAlertConditionResponse
		mockResponce *model.AlertCondition
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertConditionRequest{ProjectId: 1001, AlertConditionId: 1001},
			want:         &alert.GetAlertConditionResponse{AlertCondition: &alert.AlertCondition{AlertConditionId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AlertCondition{AlertConditionID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertConditionRequest{ProjectId: 1, AlertConditionId: 9999},
			want:      &alert.GetAlertConditionResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondition", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlertCondition(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlertCondition(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutAlertConditionRequest
		want        *alert.PutAlertConditionResponse
		wantErr     bool
		mockGetResp *model.AlertCondition
		mockGetErr  error
		mockUpResp  *model.AlertCondition
		mockUpErr   error
	}{
		{
			name:       "OK Upsert",
			input:      &alert.PutAlertConditionRequest{AlertCondition: &alert.AlertConditionForUpsert{ProjectId: 1001, Description: "desc", Severity: "high", Enabled: true, AndOr: "and"}},
			want:       &alert.PutAlertConditionResponse{AlertCondition: &alert.AlertCondition{ProjectId: 1001, Description: "desc", Severity: "high", Enabled: true, AndOr: "and", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.AlertCondition{ProjectID: 1001, Description: "desc", Severity: "high", Enabled: true, AndOr: "and", CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondition", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlertCondition(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlertCondition(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                            string
		input                           *alert.DeleteAlertConditionRequest
		wantErr                         bool
		deleteAlertConditionCall        bool
		listAlertCondRuleCall           bool
		deleteAlertCondRuleCall         bool
		listAlertCondNotificationCall   bool
		deleteAlertCondNotificationCall bool
		mockErr                         error
	}{
		{
			name:                            "OK",
			input:                           &alert.DeleteAlertConditionRequest{ProjectId: 1, AlertConditionId: 1001},
			wantErr:                         false,
			deleteAlertConditionCall:        true,
			listAlertCondRuleCall:           true,
			deleteAlertCondRuleCall:         true,
			listAlertCondNotificationCall:   true,
			deleteAlertCondNotificationCall: true,
			mockErr:                         nil,
		},
		{
			name:                     "NG validation error",
			input:                    &alert.DeleteAlertConditionRequest{ProjectId: 1},
			deleteAlertConditionCall: false,
			wantErr:                  true,
		},
		{
			name:                            "Invalid DB error",
			input:                           &alert.DeleteAlertConditionRequest{ProjectId: 1, AlertConditionId: 1001},
			wantErr:                         true,
			deleteAlertConditionCall:        false,
			listAlertCondRuleCall:           true,
			deleteAlertCondRuleCall:         true,
			listAlertCondNotificationCall:   false,
			deleteAlertCondNotificationCall: false,
			mockErr:                         gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.listAlertCondRuleCall {
				mockDB.On("ListAlertCondRule", test.RepeatMockAnything(7)...).Return(&[]model.AlertCondRule{{AlertConditionID: 1}}, nil)
			}
			if c.deleteAlertCondRuleCall {
				mockDB.On("DeleteAlertCondRule", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			if c.listAlertCondNotificationCall {
				mockDB.On("ListAlertCondNotification", test.RepeatMockAnything(6)...).Return(&[]model.AlertCondNotification{{AlertConditionID: 1}}, nil)
			}
			if c.deleteAlertCondNotificationCall {
				mockDB.On("DeleteAlertCondNotification", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			if c.deleteAlertConditionCall {
				mockDB.On("DeleteAlertCondition", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}

			_, err := svc.DeleteAlertCondition(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlertCondition(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.AlertCondition
		want  *alert.AlertCondition
	}{
		{
			name:  "OK convert unix time",
			input: &model.AlertCondition{AlertConditionID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertCondition{AlertConditionId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertCondition{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlertCondition(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * AlertRule
 */

func TestListAlertRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListAlertRuleRequest
		want         *alert.ListAlertRuleResponse
		mockResponce *[]model.AlertRule
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0},
			want:         &alert.ListAlertRuleResponse{AlertRule: []*alert.AlertRule{{AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertRuleId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.AlertRule{{AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertRuleID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0},
			want:      &alert.ListAlertRuleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertRule", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlertRule(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertRuleRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertRuleRequest
		want  *alert.ListAlertRuleRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 0.5, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 0.5, FromAt: now.Unix(), ToAt: now.Unix()},
		},
		{
			name:  "OK Convert ToValue",
			input: &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 0.0, FromAt: now.Unix(), ToAt: 0},
			want:  &alert.ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertRuleRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlertRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetAlertRuleRequest
		want         *alert.GetAlertRuleResponse
		mockResponce *model.AlertRule
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertRuleRequest{ProjectId: 1001, AlertRuleId: 1001},
			want:         &alert.GetAlertRuleResponse{AlertRule: &alert.AlertRule{AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AlertRule{AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertRuleRequest{ProjectId: 1, AlertRuleId: 9999},
			want:      &alert.GetAlertRuleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertRule", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlertRule(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlertRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutAlertRuleRequest
		want        *alert.PutAlertRuleResponse
		wantErr     bool
		mockGetResp *model.AlertRule
		mockGetErr  error
		mockUpResp  *model.AlertRule
		mockUpErr   error
	}{
		{
			name:       "OK Upsert",
			input:      &alert.PutAlertRuleRequest{AlertRule: &alert.AlertRuleForUpsert{ProjectId: 1001, Name: "name", Score: 0.1, ResourceName: "rn", Tag: "tag", FindingCnt: 1}},
			want:       &alert.PutAlertRuleResponse{AlertRule: &alert.AlertRule{ProjectId: 1001, Name: "name", Score: 0.1, ResourceName: "rn", Tag: "tag", FindingCnt: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.AlertRule{ProjectID: 1001, Name: "name", Score: 0.1, ResourceName: "rn", Tag: "tag", FindingCnt: 1, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertRule", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlertRule(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlertRule(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                string
		input               *alert.DeleteAlertRuleRequest
		wantErr             bool
		deleteAlertRuleCall bool
		mockErr             error
	}{
		{
			name:                "OK",
			input:               &alert.DeleteAlertRuleRequest{ProjectId: 1, AlertRuleId: 1001},
			wantErr:             false,
			deleteAlertRuleCall: true,
			mockErr:             nil,
		},
		{
			name:                "NG validation error",
			input:               &alert.DeleteAlertRuleRequest{ProjectId: 1},
			deleteAlertRuleCall: false,
			wantErr:             true,
		},
		{
			name:                "Invalid DB error",
			input:               &alert.DeleteAlertRuleRequest{ProjectId: 1, AlertRuleId: 1001},
			wantErr:             true,
			deleteAlertRuleCall: true,
			mockErr:             gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteAlertRuleCall {
				mockDB.On("ListAlertCondRule", test.RepeatMockAnything(6)...).Return(&[]model.AlertCondRule{{AlertConditionID: 1}}, nil)
				mockDB.On("DeleteAlertCondRule", test.RepeatMockAnything(4)...).Return(nil)
				mockDB.On("DeleteAlertRule", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteAlertRule(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlertRule(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.AlertRule
		want  *alert.AlertRule
	}{
		{
			name:  "OK convert unix time",
			input: &model.AlertRule{AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertRule{AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertRule{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlertRule(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * AlertCondRule
 */

func TestListAlertCondRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListAlertCondRuleRequest
		want         *alert.ListAlertCondRuleResponse
		mockResponce *[]model.AlertCondRule
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			want:         &alert.ListAlertCondRuleResponse{AlertCondRule: []*alert.AlertCondRule{{AlertConditionId: 1001, AlertRuleId: 1001, ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertConditionId: 1002, AlertRuleId: 1001, ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.AlertCondRule{{AlertConditionID: 1001, AlertRuleID: 1001, ProjectID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertConditionID: 1002, AlertRuleID: 1001, ProjectID: 1001, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			want:      &alert.ListAlertCondRuleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondRule", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlertCondRule(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertCondRuleRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertCondRuleRequest
		want  *alert.ListAlertCondRuleRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertCondRuleRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlertCondRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetAlertCondRuleRequest
		want         *alert.GetAlertCondRuleResponse
		mockResponce *model.AlertCondRule
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			want:         &alert.GetAlertCondRuleResponse{AlertCondRule: &alert.AlertCondRule{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AlertCondRule{ProjectID: 1001, AlertConditionID: 1001, AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 9999, AlertRuleId: 9999},
			want:      &alert.GetAlertCondRuleResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondRule", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlertCondRule(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlertCondRule(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutAlertCondRuleRequest
		want        *alert.PutAlertCondRuleResponse
		wantErr     bool
		mockGetResp *model.AlertCondRule
		mockGetErr  error
		mockUpResp  *model.AlertCondRule
		mockUpErr   error
	}{
		{
			name:       "OK Upsert",
			input:      &alert.PutAlertCondRuleRequest{AlertCondRule: &alert.AlertCondRuleForUpsert{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001}},
			want:       &alert.PutAlertCondRuleResponse{AlertCondRule: &alert.AlertCondRule{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.AlertCondRule{ProjectID: 1001, AlertConditionID: 1001, AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondRule", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlertCondRule(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlertCondRule(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                    string
		input                   *alert.DeleteAlertCondRuleRequest
		wantErr                 bool
		deleteAlertCondRuleCall bool
		mockErr                 error
	}{
		{
			name:                    "OK",
			input:                   &alert.DeleteAlertCondRuleRequest{ProjectId: 1, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr:                 false,
			deleteAlertCondRuleCall: true,
			mockErr:                 nil,
		},
		{
			name:                    "NG validation error",
			input:                   &alert.DeleteAlertCondRuleRequest{ProjectId: 1},
			deleteAlertCondRuleCall: false,
			wantErr:                 true,
		},
		{
			name:                    "Invalid DB error",
			input:                   &alert.DeleteAlertCondRuleRequest{ProjectId: 1, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr:                 true,
			deleteAlertCondRuleCall: true,
			mockErr:                 gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteAlertCondRuleCall {
				mockDB.On("DeleteAlertCondRule", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteAlertCondRule(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlertCondRule(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.AlertCondRule
		want  *alert.AlertCondRule
	}{
		{
			name:  "OK convert unix time",
			input: &model.AlertCondRule{ProjectID: 1001, AlertConditionID: 1001, AlertRuleID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertCondRule{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertCondRule{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlertCondRule(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * Notification
 */

func TestListNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListNotificationRequest
		want         *alert.ListNotificationResponse
		mockResponce *[]model.Notification
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListNotificationRequest{ProjectId: 1001, Type: "type"},
			want:         &alert.ListNotificationResponse{Notification: []*alert.Notification{{NotificationId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {NotificationId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.Notification{{NotificationID: 1001, CreatedAt: now, UpdatedAt: now}, {NotificationID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListNotificationRequest{ProjectId: 1001, Type: "type"},
			want:      &alert.ListNotificationResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListNotification", test.RepeatMockAnything(5)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListNotification(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListNotificationRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListNotificationRequest
		want  *alert.ListNotificationRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListNotificationRequest{ProjectId: 1001, Type: "type", FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListNotificationRequest{ProjectId: 1001, Type: "type", FromAt: now.Unix(), ToAt: now.Unix()},
		},
		{
			name:  "OK Convert ToValue",
			input: &alert.ListNotificationRequest{ProjectId: 1001, Type: "type", FromAt: now.Unix(), ToAt: 0},
			want:  &alert.ListNotificationRequest{ProjectId: 1001, Type: "type", FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListNotificationRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetNotificationRequest
		want         *alert.GetNotificationResponse
		mockResponce *model.Notification
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetNotificationRequest{ProjectId: 1001, NotificationId: 1001},
			want:         &alert.GetNotificationResponse{Notification: &alert.Notification{NotificationId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Notification{NotificationID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetNotificationRequest{ProjectId: 1, NotificationId: 9999},
			want:      &alert.GetNotificationResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetNotification", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetNotification(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestReplaceSlackNotifySetting(t *testing.T) {
	a := AlertService{logger: logging.NewLogger()}
	cases := []struct {
		name        string
		inputExist  string
		inputUpdate string
		want        slackNotifySetting
		wantErr     bool
	}{
		{
			name:        "OK no replacing",
			inputExist:  "{\"webhook_url\":\"hoge1\", \"data\":{\"channel\":\"ch\"}}",
			inputUpdate: "{\"data\":{\"channel\":\"ch\"}}",
			want: slackNotifySetting{
				WebhookURL: "hoge1",
				Data:       slackNotifyOption{Channel: "ch"},
			},
			wantErr: false,
		},
		{
			name:        "OK replace webhook_url",
			inputExist:  "{\"webhook_url\":\"hoge1\",\"data\":{\"hoge\":\"fuga\"}}",
			inputUpdate: "{\"webhook_url\":\"hoge2\"}",
			want:        slackNotifySetting{WebhookURL: "hoge2"},
			wantErr:     false,
		},
		{
			name:        "OK blank",
			inputExist:  "{}",
			inputUpdate: "{}",
			want:        slackNotifySetting{},
			wantErr:     false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := a.replaceSlackNotifySetting(context.Background(), c.inputExist, c.inputUpdate)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                   string
		input                  *alert.DeleteNotificationRequest
		wantErr                bool
		deleteNotificationCall bool
		mockErr                error
	}{
		{
			name:                   "OK",
			input:                  &alert.DeleteNotificationRequest{ProjectId: 1, NotificationId: 1001},
			wantErr:                false,
			deleteNotificationCall: true,
			mockErr:                nil,
		},
		{
			name:                   "NG validation error",
			input:                  &alert.DeleteNotificationRequest{ProjectId: 1},
			deleteNotificationCall: false,
			wantErr:                true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &alert.DeleteNotificationRequest{ProjectId: 1, NotificationId: 1001},
			wantErr:                true,
			deleteNotificationCall: true,
			mockErr:                gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteNotificationCall {
				mockDB.On("ListAlertCondNotification", test.RepeatMockAnything(6)...).Return(&[]model.AlertCondNotification{{ProjectID: 1, AlertConditionID: 1, NotificationID: 1}}, nil)
				mockDB.On("DeleteAlertCondNotification", test.RepeatMockAnything(4)...).Return(nil)
				mockDB.On("DeleteNotification", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertNotification(t *testing.T) {
	a := AlertService{logger: logging.NewLogger()}
	now := time.Now()
	cases := []struct {
		name    string
		input   *model.Notification
		want    *alert.Notification
		wantErr bool
	}{
		{
			name:  "OK convert unix time",
			input: &model.Notification{NotificationID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.Notification{NotificationId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.Notification{},
		},
		{
			name:    "NG json marshal error",
			input:   &model.Notification{NotificationID: 1001, CreatedAt: now, UpdatedAt: now, Type: "slack", NotifySetting: "{\"data\": {\"aaa\"}, \"webhook_url\": \"http://hogehoge.com\"}"},
			want:    &alert.Notification{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := a.convertNotification(context.Background(), c.input, true)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

/*
 * AlertCondNotification
 */

func TestListAlertCondNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.ListAlertCondNotificationRequest
		want         *alert.ListAlertCondNotificationResponse
		mockResponce *[]model.AlertCondNotification
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			want:         &alert.ListAlertCondNotificationResponse{AlertCondNotification: []*alert.AlertCondNotification{{AlertConditionId: 1001, NotificationId: 1001, ProjectId: 1001, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertConditionId: 1002, NotificationId: 1001, ProjectId: 1001, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.AlertCondNotification{{AlertConditionID: 1001, NotificationID: 1001, ProjectID: 1001, NotifiedAt: now, CreatedAt: now, UpdatedAt: now}, {AlertConditionID: 1002, NotificationID: 1001, ProjectID: 1001, NotifiedAt: now, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			want:      &alert.ListAlertCondNotificationResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB, logger: logging.NewLogger()}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondNotification", test.RepeatMockAnything(6)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlertCondNotification(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertCondNotificationRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertCondNotificationRequest
		want  *alert.ListAlertCondNotificationRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertCondNotificationRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlertCondNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *alert.GetAlertCondNotificationRequest
		want         *alert.GetAlertCondNotificationResponse
		mockResponce *model.AlertCondNotification
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.GetAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			want:         &alert.GetAlertCondNotificationResponse{AlertCondNotification: &alert.AlertCondNotification{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.AlertCondNotification{ProjectID: 1001, AlertConditionID: 1001, NotificationID: 1001, NotifiedAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &alert.GetAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 9999, NotificationId: 9999},
			want:      &alert.GetAlertCondNotificationResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondNotification", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetAlertCondNotification(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutAlertCondNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name        string
		input       *alert.PutAlertCondNotificationRequest
		want        *alert.PutAlertCondNotificationResponse
		wantErr     bool
		mockGetResp *model.AlertCondNotification
		mockGetErr  error
		mockUpResp  *model.AlertCondNotification
		mockUpErr   error
	}{
		{
			name:       "OK Upsert",
			input:      &alert.PutAlertCondNotificationRequest{AlertCondNotification: &alert.AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, CacheSecond: 1, NotifiedAt: now.Unix()}},
			want:       &alert.PutAlertCondNotificationResponse{AlertCondNotification: &alert.AlertCondNotification{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, CacheSecond: 1, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.AlertCondNotification{ProjectID: 1001, AlertConditionID: 1001, NotificationID: 1001, CacheSecond: 1, NotifiedAt: now, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondNotification", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutAlertCondNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAlertCondNotification(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name                            string
		input                           *alert.DeleteAlertCondNotificationRequest
		wantErr                         bool
		deleteAlertCondNotificationCall bool
		mockErr                         error
	}{
		{
			name:                            "OK",
			input:                           &alert.DeleteAlertCondNotificationRequest{ProjectId: 1, AlertConditionId: 1001, NotificationId: 1001},
			wantErr:                         false,
			deleteAlertCondNotificationCall: true,
			mockErr:                         nil,
		},
		{
			name:                            "NG validation error",
			input:                           &alert.DeleteAlertCondNotificationRequest{ProjectId: 1},
			deleteAlertCondNotificationCall: false,
			wantErr:                         true,
		},
		{
			name:                            "Invalid DB error",
			input:                           &alert.DeleteAlertCondNotificationRequest{ProjectId: 1, AlertConditionId: 1001, NotificationId: 1001},
			wantErr:                         true,
			deleteAlertCondNotificationCall: true,
			mockErr:                         gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAlertRepository(t)
			svc := AlertService{repository: mockDB}
			if c.deleteAlertCondNotificationCall {
				mockDB.On("DeleteAlertCondNotification", test.RepeatMockAnything(4)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteAlertCondNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertAlertCondNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.AlertCondNotification
		want  *alert.AlertCondNotification
	}{
		{
			name:  "OK convert unix time",
			input: &model.AlertCondNotification{ProjectID: 1001, AlertConditionID: 1001, NotificationID: 1001, NotifiedAt: now, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertCondNotification{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertCondNotification{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlertCondNotification(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
