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

func TestListOrgNotification(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationNotification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1},
			want: []*model.OrganizationNotification{
				{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_notification where organization_id = ? order by notification_id")).WillReturnRows(sqlmock.NewRows([]string{
					"notification_id", "name", "organization_id", "type", "notify_setting", "created_at", "updated_at"}).
					AddRow(uint32(1), "notif1", uint32(1), "slack", "{}", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_notification where organization_id = ? order by notification_id")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrgNotification(ctx, c.args.organizationID)
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

func TestGetOrgNotification(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		notificationID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationNotification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, notificationID: 1},
			want:    &model.OrganizationNotification{NotificationID: 1, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgNotification)).WillReturnRows(sqlmock.NewRows([]string{
					"notification_id", "name", "organization_id", "type", "notify_setting", "created_at", "updated_at"}).
					AddRow(uint32(1), "notif1", uint32(1), "slack", "{}", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, notificationID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgNotification)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrgNotification(ctx, c.args.organizationID, c.args.notificationID)
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

func TestUpsertOrgNotification(t *testing.T) {
	type args struct {
		data *model.OrganizationNotification
	}

	// FirstOrCreate: SELECT → (INSERT if not found)
	selectFirstOrCreate := `SELECT * FROM ` + "`organization_notification`" + ` WHERE organization_id = ? AND notification_id = ? ORDER BY ` + "`organization_notification`" + `.` + "`notification_id`" + ` LIMIT 1`
	insertFirstOrCreate := "INSERT INTO `organization_notification`"

	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationNotification
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK - create new",
			args:    args{data: &model.OrganizationNotification{NotificationID: 0, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}"}},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectFirstOrCreate)).
					WithArgs(uint32(1), uint32(0)).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id"}))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(insertFirstOrCreate)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG - DB error on select",
			args:    args{data: &model.OrganizationNotification{NotificationID: 0, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectFirstOrCreate)).
					WithArgs(uint32(1), uint32(0)).
					WillReturnError(errors.New("DB error"))
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
			got, err := db.UpsertOrgNotification(ctx, c.args.data)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !c.wantErr && got == nil {
				t.Fatal("Expected non-nil result, got nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteOrgNotification(t *testing.T) {
	type args struct {
		organizationID uint32
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
			args:    args{organizationID: 1, notificationID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgNotification)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, notificationID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgNotification)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrgNotification(ctx, c.args.organizationID, c.args.notificationID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

