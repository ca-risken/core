package db

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/pkg/model"
	"gorm.io/gorm"
)

func TestListAlert(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID   uint32
		status      []string
		severity    []string
		description string
		fromAt      int64
		toAt        int64
	}
	cases := []struct {
		name        string
		args        args
		want        *[]model.Alert
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, status: []string{"ACTIVE"}, severity: []string{"HIGH"}, description: "test", fromAt: now.Unix(), toAt: now.Unix()},
			want: &[]model.Alert{
				{AlertID: 1, ProjectID: 1, Status: "ACTIVE", Severity: "HIGH", Description: "test", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from alert where project_id = ? and updated_at between ? and ? and severity in (?) and status in (?) and description = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"alert_id", "project_id", "status", "severity", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "ACTIVE", "HIGH", "test", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, status: []string{"ACTIVE"}, severity: []string{"HIGH"}, description: "test", fromAt: now.Unix(), toAt: now.Unix()},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from alert where project_id = ? and updated_at between ? and ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListAlert(ctx, c.args.projectID, c.args.status, c.args.severity, c.args.description, c.args.fromAt, c.args.toAt)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAlert(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID uint32
		alertID   uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Alert
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, alertID: 1},
			want: &model.Alert{AlertID: 1, ProjectID: 1, Status: "ACTIVE", Severity: "HIGH", Description: "test", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert` WHERE project_id = ? AND alert_id = ? ORDER BY `alert`.`alert_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"alert_id", "project_id", "status", "severity", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "ACTIVE", "HIGH", "test", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, alertID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert` WHERE project_id = ? AND alert_id = ? ORDER BY `alert`.`alert_id` LIMIT 1")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetAlert(ctx, c.args.projectID, c.args.alertID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsertAlert(t *testing.T) {
	now := time.Now()
	type args struct {
		alert *model.Alert
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Alert
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{alert: &model.Alert{AlertID: 1, ProjectID: 1, Status: "ACTIVE", Severity: "HIGH", Description: "test"}},
			want: &model.Alert{AlertID: 1, ProjectID: 1, Status: "ACTIVE", Severity: "HIGH", Description: "test", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert` WHERE project_id = ? AND alert_id = ?")).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `alert`")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{alert: &model.Alert{AlertID: 1, ProjectID: 1, Status: "ACTIVE", Severity: "HIGH", Description: "test"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert` WHERE project_id = ? AND alert_id = ?")).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `alert`")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertAlert(ctx, c.args.alert)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateAlertFirstViewedAt(t *testing.T) {
	type args struct {
		projectID uint32
		alertID   uint32
		viewedAt  int64
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{projectID: 1, alertID: 1, viewedAt: time.Now().Unix()},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `alert` SET `first_viewed_at`=? WHERE project_id = ? AND alert_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, alertID: 1, viewedAt: time.Now().Unix()},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `alert` SET `first_viewed_at`=? WHERE project_id = ? AND alert_id = ?")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.UpdateAlertFirstViewedAt(ctx, c.args.projectID, c.args.alertID, c.args.viewedAt)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAlert(t *testing.T) {
	type args struct {
		projectID uint32
		alertID   uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{projectID: 1, alertID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `alert` WHERE project_id = ? AND alert_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, alertID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `alert` WHERE project_id = ? AND alert_id = ?")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeleteAlert(ctx, c.args.projectID, c.args.alertID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListAlertCondition(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID uint32
		severity  []string
		enabled   bool
		fromAt    int64
		toAt      int64
	}
	cases := []struct {
		name        string
		args        args
		want        *[]model.AlertCondition
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, severity: []string{"HIGH"}, enabled: true, fromAt: now.Unix(), toAt: now.Unix()},
			want: &[]model.AlertCondition{
				{AlertConditionID: 1, ProjectID: 1, Severity: "HIGH", Enabled: true, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from alert_condition where project_id = ? and updated_at between ? and ?")).WillReturnRows(sqlmock.NewRows([]string{
					"alert_condition_id", "project_id", "severity", "enabled", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "HIGH", true, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, severity: []string{"HIGH"}, enabled: true, fromAt: now.Unix(), toAt: now.Unix()},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from alert_condition where project_id = ? and updated_at between ? and ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListAlertCondition(ctx, c.args.projectID, c.args.severity, c.args.enabled, c.args.fromAt, c.args.toAt)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAlertCondition(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID        uint32
		alertConditionID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.AlertCondition
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, alertConditionID: 1},
			want: &model.AlertCondition{AlertConditionID: 1, ProjectID: 1, Severity: "HIGH", Enabled: true, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ? ORDER BY `alert_condition`.`alert_condition_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"alert_condition_id", "project_id", "severity", "enabled", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "HIGH", true, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, alertConditionID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ? ORDER BY `alert_condition`.`alert_condition_id` LIMIT 1")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetAlertCondition(ctx, c.args.projectID, c.args.alertConditionID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsertAlertCondition(t *testing.T) {
	now := time.Now()
	type args struct {
		alertCondition *model.AlertCondition
	}
	cases := []struct {
		name        string
		args        args
		want        *model.AlertCondition
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{alertCondition: &model.AlertCondition{AlertConditionID: 1, ProjectID: 1, Severity: "HIGH", Enabled: true}},
			want: &model.AlertCondition{AlertConditionID: 1, ProjectID: 1, Severity: "HIGH", Enabled: true, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `alert_condition`")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ? ORDER BY `alert_condition`.`alert_condition_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"alert_condition_id", "project_id", "severity", "enabled", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "HIGH", true, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{alertCondition: &model.AlertCondition{AlertConditionID: 1, ProjectID: 1, Severity: "HIGH", Enabled: true}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `alert_condition`")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertAlertCondition(ctx, c.args.alertCondition)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAlertCondition(t *testing.T) {
	type args struct {
		projectID        uint32
		alertConditionID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{projectID: 1, alertConditionID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, alertConditionID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeleteAlertCondition(ctx, c.args.projectID, c.args.alertConditionID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListNotification(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID uint32
		name      string
		fromAt    int64
		toAt      int64
	}
	cases := []struct {
		name        string
		args        args
		want        *[]model.Notification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, name: "notification1", fromAt: now.Unix(), toAt: now.Unix()},
			want: &[]model.Notification{
				{NotificationID: 1, ProjectID: 1, Name: "notification1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from notification where project_id = ? and updated_at between ? and ?")).WillReturnRows(sqlmock.NewRows([]string{
					"notification_id", "project_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "notification1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, name: "notification1", fromAt: now.Unix(), toAt: now.Unix()},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from notification where project_id = ? and updated_at between ? and ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListNotification(ctx, c.args.projectID, c.args.name, c.args.fromAt, c.args.toAt)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetNotification(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID      uint32
		notificationID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Notification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1, notificationID: 1},
			want: &model.Notification{NotificationID: 1, ProjectID: 1, Name: "notification1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notification` WHERE project_id = ? AND notification_id = ? ORDER BY `notification`.`notification_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"notification_id", "project_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "notification1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, notificationID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notification` WHERE project_id = ? AND notification_id = ? ORDER BY `notification`.`notification_id` LIMIT 1")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetNotification(ctx, c.args.projectID, c.args.notificationID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsertNotification(t *testing.T) {
	now := time.Now()
	type args struct {
		notification *model.Notification
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Notification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{notification: &model.Notification{NotificationID: 1, ProjectID: 1, Name: "notification1"}},
			want: &model.Notification{NotificationID: 1, ProjectID: 1, Name: "notification1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `notification`")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notification` WHERE project_id = ? AND notification_id = ? ORDER BY `notification`.`notification_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"notification_id", "project_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "notification1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{notification: &model.Notification{NotificationID: 1, ProjectID: 1, Name: "notification1"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `notification`")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertNotification(ctx, c.args.notification)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	type args struct {
		projectID      uint32
		notificationID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{projectID: 1, notificationID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `notification` WHERE project_id = ? AND notification_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, notificationID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `notification` WHERE project_id = ? AND notification_id = ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeleteNotification(ctx, c.args.projectID, c.args.notificationID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeactivateAlert(t *testing.T) {
	now := time.Now()
	type args struct {
		alert *model.Alert
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{alert: &model.Alert{AlertID: 1, ProjectID: 1, Status: "DEACTIVE", UpdatedAt: now}},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `alert` SET")).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{alert: &model.Alert{AlertID: 1, ProjectID: 1, Status: "DEACTIVE", UpdatedAt: now}},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `alert` SET")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeactivateAlert(ctx, c.args.alert)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
