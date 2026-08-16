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
)

func TestListAlertCondNotificationForNotification(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	rows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"alert_condition_id", "notification_id", "project_id", "cache_second", "notified_at", "created_at", "updated_at"})
	}
	cases := []struct {
		name    string
		want    *[]model.AlertCondNotification
		dbErr   error
		wantErr bool
	}{
		{
			name: "OK - exact project and condition without time bound",
			want: &[]model.AlertCondNotification{{AlertConditionID: 2, NotificationID: 3, ProjectID: 1, CacheSecond: 1800, NotifiedAt: now, CreatedAt: now, UpdatedAt: now.Add(time.Second)}},
		},
		{name: "OK - empty", want: &[]model.AlertCondNotification{}},
		{name: "NG - DB error", dbErr: errors.New("DB error"), wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			expectation := mock.ExpectQuery(regexp.QuoteMeta(listAlertCondNotificationForNotification)).WithArgs(uint32(1), uint32(2))
			if tc.dbErr != nil {
				expectation.WillReturnError(tc.dbErr)
			} else if len(*tc.want) == 0 {
				expectation.WillReturnRows(rows())
			} else {
				relation := (*tc.want)[0]
				expectation.WillReturnRows(rows().AddRow(relation.AlertConditionID, relation.NotificationID, relation.ProjectID, relation.CacheSecond, relation.NotifiedAt, relation.CreatedAt, relation.UpdatedAt))
			}
			got, err := database.ListAlertCondNotificationForNotification(context.Background(), 1, 2)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Unexpected data: want=%+v got=%+v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpsertAlertConditionLifecycle(t *testing.T) {
	selectAlertCondition := "SELECT * FROM `alert_condition` WHERE project_id = ? AND alert_condition_id = ? ORDER BY `alert_condition`.`alert_condition_id` LIMIT 1"
	cases := []struct {
		name        string
		relationErr error
		wantErr     bool
	}{
		{name: "OK - condition and cross product commit together"},
		{name: "NG - relation failure rolls back condition", relationErr: errors.New("DB error"), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(selectAlertCondition)).WithArgs(uint32(1), uint32(0)).WillReturnRows(sqlmock.NewRows([]string{"alert_condition_id"}))
			mock.ExpectExec("INSERT INTO `alert_condition`").WillReturnResult(sqlmock.NewResult(1, 1))
			relation := mock.ExpectExec(regexp.QuoteMeta(insertOrgAlertCondNotificationByAlertCondition)).WithArgs(uint32(1), uint32(1))
			if c.relationErr != nil {
				relation.WillReturnError(c.relationErr)
				mock.ExpectRollback()
			} else {
				relation.WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			}
			_, err = database.upsertAlertCondition(context.Background(), &model.AlertCondition{ProjectID: 1})
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteAlertConditionLifecycle(t *testing.T) {
	cases := []struct {
		name         string
		conditionErr error
		wantErr      bool
	}{
		{name: "OK - relation deleted before condition"},
		{name: "NG - condition failure rolls back relation", conditionErr: errors.New("DB error"), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(deleteOrgAlertCondNotificationByAlertCondition)).WithArgs(uint32(1), uint32(2)).WillReturnResult(sqlmock.NewResult(0, 2))
			condition := mock.ExpectExec("DELETE FROM `alert_condition`").WithArgs(uint32(1), uint32(2))
			if c.conditionErr != nil {
				condition.WillReturnError(c.conditionErr)
				mock.ExpectRollback()
			} else {
				condition.WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			err = database.deleteAlertCondition(context.Background(), 1, 2)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
