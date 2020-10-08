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

/*
 * Alert
 */

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
 * AlertHistory
 */

func TestListAlertHistory(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *alert.ListAlertHistoryRequest
		want         *alert.ListAlertHistoryResponse
		mockResponce *[]model.AlertHistory
		mockError    error
	}{
		{
			name:         "OK",
			input:        &alert.ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}},
			want:         &alert.ListAlertHistoryResponse{AlertHistory: []*alert.AlertHistory{{AlertHistoryId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}, {AlertHistoryId: 1002, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}}},
			mockResponce: &[]model.AlertHistory{{AlertHistoryID: 1001, CreatedAt: now, UpdatedAt: now}, {AlertHistoryID: 1002, CreatedAt: now, UpdatedAt: now}},
		},
		{
			name:      "NG Record not found",
			input:     &alert.ListAlertHistoryRequest{ProjectId: 1, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}},
			want:      &alert.ListAlertHistoryResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertHistory").Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.ListAlertHistory(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListAlertHistoryRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *alert.ListAlertHistoryRequest
		want  *alert.ListAlertHistoryRequest
	}{
		{
			name:  "OK full-set",
			input: &alert.ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &alert.ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListAlertHistoryRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetAlertHistory(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertHistory").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertHistory").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertHistoryRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertHistoryRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlertHistory").Return(c.mockErr).Once()
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
		name  string
		input *model.AlertHistory
		want  *alert.AlertHistory
	}{
		{
			name:  "OK convert unix time",
			input: &model.AlertHistory{AlertHistoryID: 1001, CreatedAt: now, UpdatedAt: now},
			want:  &alert.AlertHistory{AlertHistoryId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &alert.AlertHistory{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertAlertHistory(c.input)
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListRelAlertFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRelAlertFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertRelAlertFinding").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteRelAlertFindingRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001, FindingId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteRelAlertFindingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001, FindingId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteRelAlertFinding").Return(c.mockErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondition").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondition").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondition").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertConditionRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertConditionRequest{ProjectId: 1, AlertConditionId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertConditionRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertConditionRequest{ProjectId: 1, AlertConditionId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlertCondition").Return(c.mockErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertRule").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertRule").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertRule").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertRuleRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertRuleRequest{ProjectId: 1, AlertRuleId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertRuleRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertRuleRequest{ProjectId: 1, AlertRuleId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlertRule").Return(c.mockErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondRule").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondRule").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondRule").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertCondRuleRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertCondRuleRequest{ProjectId: 1, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertCondRuleRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertCondRuleRequest{ProjectId: 1, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlertCondRule").Return(c.mockErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListNotification").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetNotification").Return(c.mockResponce, c.mockError).Once()
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

func TestPutNotification(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			name:       "OK Upsert",
			input:      &alert.PutNotificationRequest{Notification: &alert.NotificationForUpsert{ProjectId: 1001, Name: "name", Type: "type", NotifySetting: `{"hoge":"fuga"}`}},
			want:       &alert.PutNotificationResponse{Notification: &alert.Notification{ProjectId: 1001, Name: "name", Type: "type", NotifySetting: `{"hoge":"fuga"}`, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpResp: &model.Notification{ProjectID: 1001, Name: "name", Type: "type", NotifySetting: `{"hoge":"fuga"}`, CreatedAt: now, UpdatedAt: now},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertNotification").Return(c.mockUpResp, c.mockUpErr).Once()
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

func TestDeleteNotification(t *testing.T) {
	var ctx context.Context
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteNotificationRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteNotificationRequest{ProjectId: 1, NotificationId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteNotificationRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteNotificationRequest{ProjectId: 1, NotificationId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteNotification").Return(c.mockErr).Once()
			_, err := svc.DeleteNotification(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestConvertNotification(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.Notification
		want  *alert.Notification
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
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertNotification(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAlertCondNotification").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAlertCondNotification").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}

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
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertAlertCondNotification").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockAlertRepository{}
	svc := alertService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *alert.DeleteAlertCondNotificationRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &alert.DeleteAlertCondNotificationRequest{ProjectId: 1, AlertConditionId: 1001, NotificationId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &alert.DeleteAlertCondNotificationRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &alert.DeleteAlertCondNotificationRequest{ProjectId: 1, AlertConditionId: 1001, NotificationId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteAlertCondNotification").Return(c.mockErr).Once()
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

func (m *mockAlertRepository) ListFinding(uint32) (*[]model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Finding), args.Error(1)
}

func (m *mockAlertRepository) ListFindingTag(uint32, uint64) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}
