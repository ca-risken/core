package db

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/pkg/model"
	orgalert "github.com/ca-risken/core/proto/org_alert"
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
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectFirstOrCreate)).
					WithArgs(uint32(1), uint32(0)).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id"}))
				mock.ExpectExec(regexp.QuoteMeta(insertFirstOrCreate)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertOrgAlertCondNotificationByOrgNotification)).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG - DB error on select",
			args:    args{data: &model.OrganizationNotification{NotificationID: 0, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectFirstOrCreate)).
					WithArgs(uint32(1), uint32(0)).
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
		{
			name:    "NG - relation cross product rolls back notification",
			args:    args{data: &model.OrganizationNotification{NotificationID: 0, Name: "notif1", OrganizationID: 1, Type: "slack", NotifySetting: "{}"}},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectFirstOrCreate)).
					WithArgs(uint32(1), uint32(0)).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id"}))
				mock.ExpectExec(regexp.QuoteMeta(insertFirstOrCreate)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(insertOrgAlertCondNotificationByOrgNotification)).WillReturnError(errors.New("DB error"))
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
			name:    "OK - relation deleted before notification",
			args:    args{organizationID: 1, notificationID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgAlertCondNotificationByOrgNotification)).WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgNotification)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, notificationID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgAlertCondNotificationByOrgNotification)).WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgNotification)).WillReturnError(errors.New("DB error"))
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

func TestListOrgNotificationByProjectID(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID uint32
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
			args: args{projectID: 1},
			want: []*model.OrganizationNotification{
				{NotificationID: 1, Name: "notif1", OrganizationID: 10, Type: "slack", NotifySetting: "{}", CreatedAt: now, UpdatedAt: now},
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgNotificationByProjectID)).
					WithArgs(uint32(1)).
					WillReturnRows(sqlmock.NewRows([]string{
						"notification_id", "name", "organization_id", "type", "notify_setting", "created_at", "updated_at"}).
						AddRow(uint32(1), "notif1", uint32(10), "slack", "{}", now, now))
			},
		},
		{
			name: "OK - empty",
			args: args{projectID: 999},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgNotificationByProjectID)).
					WithArgs(uint32(999)).
					WillReturnRows(sqlmock.NewRows([]string{
						"notification_id", "name", "organization_id", "type", "notify_setting", "created_at", "updated_at"}))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgNotificationByProjectID)).
					WithArgs(uint32(1)).
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
			got, err := db.ListOrgNotificationByProjectID(ctx, c.args.projectID)
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

func TestListOrgAlertCondNotification(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	epoch := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)
	type args struct {
		organizationID   uint32
		projectID        uint32
		alertConditionID uint32
		notificationID   uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*orgalert.OrgAlertCondNotification
		wantErr     bool
		mockClosure func(sqlmock.Sqlmock)
	}{
		{
			name: "OK - epoch is API zero",
			args: args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4},
			want: []*orgalert.OrgAlertCondNotification{{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, Enabled: true, CacheSecond: 1800, NotifiedAt: 0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockClosure: func(mock sqlmock.Sqlmock) {
				query := listOrgAlertCondNotification + " and project_id = ? and alert_condition_id = ? and notification_id = ? order by project_id, alert_condition_id, notification_id limit ? offset ?"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4), uint32(101), uint32(0)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(1800), epoch, now, now))
			},
		},
		{
			name: "OK - organization scope",
			args: args{organizationID: 1},
			want: []*orgalert.OrgAlertCondNotification{{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, Enabled: true, CacheSecond: 1800, NotifiedAt: 0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockClosure: func(mock sqlmock.Sqlmock) {
				query := listOrgAlertCondNotification + " order by project_id, alert_condition_id, notification_id limit ? offset ?"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint32(1), uint32(101), uint32(0)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(1800), epoch, now, now))
			},
		},
		{
			name: "OK - optional subset",
			args: args{organizationID: 1, alertConditionID: 3},
			want: []*orgalert.OrgAlertCondNotification{},
			mockClosure: func(mock sqlmock.Sqlmock) {
				query := listOrgAlertCondNotification + " and alert_condition_id = ? order by project_id, alert_condition_id, notification_id limit ? offset ?"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint32(1), uint32(3), uint32(101), uint32(0)).
					WillReturnRows(orgAlertCondNotificationRows())
			},
		},
		{
			name:    "NG - DB error",
			args:    args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				query := listOrgAlertCondNotification + " and project_id = ? and alert_condition_id = ? and notification_id = ? order by project_id, alert_condition_id, notification_id limit ? offset ?"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4), uint32(101), uint32(0)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Unexpected mock DB error: %v", err)
			}
			c.mockClosure(mock)
			got, err := database.ListOrgAlertCondNotification(context.Background(), c.args.organizationID, c.args.projectID, c.args.alertConditionID, c.args.notificationID, 101, 0)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetOrgAlertCondNotification(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	type args struct {
		organizationID   uint32
		projectID        uint32
		alertConditionID uint32
		notificationID   uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *orgalert.OrgAlertCondNotification
		wantErr     bool
		mockClosure func(sqlmock.Sqlmock)
	}{
		{
			name: "OK - all four keys",
			args: args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4},
			want: &orgalert.OrgAlertCondNotification{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, Enabled: true, CacheSecond: 60, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotification)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(60), now, now, now))
			},
		},
		{
			name:    "NG - DB error",
			args:    args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotification)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Unexpected mock DB error: %v", err)
			}
			c.mockClosure(mock)
			got, err := database.GetOrgAlertCondNotification(context.Background(), c.args.organizationID, c.args.projectID, c.args.alertConditionID, c.args.notificationID)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateOrgAlertCondNotificationCache(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	type args struct {
		organizationID   uint32
		projectID        uint32
		alertConditionID uint32
		notificationID   uint32
		cacheSecond      uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *orgalert.OrgAlertCondNotification
		wantErr     bool
		mockClosure func(sqlmock.Sqlmock)
	}{
		{
			name: "OK - membership and relation are locked with all four keys",
			args: args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4, cacheSecond: 300},
			want: &orgalert.OrgAlertCondNotification{OrganizationId: 1, ProjectId: 2, AlertConditionId: 3, NotificationId: 4, Enabled: true, CacheSecond: 300, NotifiedAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotificationForUpdate)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(1800), now, now, now))
				mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertCondNotificationCache)).
					WithArgs(uint32(300), uint32(1), uint32(2), uint32(3), uint32(4)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotification)).
					WithArgs(uint32(1), uint32(2), uint32(3), uint32(4)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(300), now, now, now))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG - update error",
			args:    args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4, cacheSecond: 300},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotificationForUpdate)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(1800), now, now, now))
				mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertCondNotificationCache)).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
		{
			name:    "NG - membership or relation is missing",
			args:    args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 999, cacheSecond: 300},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotificationForUpdate)).WillReturnRows(orgAlertCondNotificationRows())
				mock.ExpectRollback()
			},
		},
		{
			name:    "NG - relation is deleted before response reload",
			args:    args{organizationID: 1, projectID: 2, alertConditionID: 3, notificationID: 4, cacheSecond: 300},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotificationForUpdate)).
					WillReturnRows(orgAlertCondNotificationRows().AddRow(uint32(1), uint32(2), uint32(3), uint32(4), true, uint32(1800), now, now, now))
				mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertCondNotificationCache)).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertCondNotification)).WillReturnRows(orgAlertCondNotificationRows())
				mock.ExpectRollback()
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Unexpected mock DB error: %v", err)
			}
			c.mockClosure(mock)
			got, err := database.UpdateOrgAlertCondNotificationCache(context.Background(), c.args.organizationID, c.args.projectID, c.args.alertConditionID, c.args.notificationID, c.args.cacheSecond)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateOrgAlertProjectNotificationEnabled(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	cases := []struct {
		name    string
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{name: "OK - all project conditions are disabled", rows: orgAlertCondNotificationRows().AddRow(1, 2, 3, 4, true, 1800, now, now, now)},
		{name: "NG - project notification relation is missing", rows: orgAlertCondNotificationRows(), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertProjectNotificationForUpdate)).WithArgs(uint32(1), uint32(2), uint32(4)).WillReturnRows(c.rows)
			if c.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertProjectNotificationEnabled)).WithArgs(false, uint32(1), uint32(2), uint32(4)).WillReturnResult(sqlmock.NewResult(0, 1))
				query := listOrgAlertCondNotification + " and project_id = ? and notification_id = ? order by alert_condition_id"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uint32(1), uint32(2), uint32(4)).WillReturnRows(orgAlertCondNotificationRows().AddRow(1, 2, 3, 4, false, 1800, now, now, now))
				mock.ExpectCommit()
			}
			got, err := database.UpdateOrgAlertProjectNotificationEnabled(context.Background(), 1, 2, 4, false)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !c.wantErr && (len(got) != 1 || got[0].Enabled) {
				t.Fatalf("Unexpected relations: %+v", got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateOrgAlertProjectNotificationCache(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	cases := []struct {
		name        string
		cacheSecond uint32
	}{
		{name: "OK - project cache is applied to all conditions", cacheSecond: 300},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(selectOrgAlertProjectNotificationForUpdate)).WithArgs(uint32(1), uint32(2), uint32(4)).WillReturnRows(orgAlertCondNotificationRows().AddRow(1, 2, 3, 4, true, 1800, now, now, now))
			mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertProjectNotificationCache)).WithArgs(c.cacheSecond, uint32(1), uint32(2), uint32(4)).WillReturnResult(sqlmock.NewResult(0, 1))
			query := listOrgAlertCondNotification + " and project_id = ? and notification_id = ? order by alert_condition_id"
			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(uint32(1), uint32(2), uint32(4)).WillReturnRows(orgAlertCondNotificationRows().AddRow(1, 2, 3, 4, true, c.cacheSecond, now, now, now))
			mock.ExpectCommit()
			got, err := database.UpdateOrgAlertProjectNotificationCache(context.Background(), 1, 2, 4, c.cacheSecond)
			if err != nil || len(got) != 1 || got[0].CacheSecond != c.cacheSecond {
				t.Fatalf("Unexpected result: got=%+v, err=%v", got, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func orgAlertCondNotificationRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"organization_id", "project_id", "alert_condition_id", "notification_id", "enabled", "cache_second", "notified_at", "created_at", "updated_at"})
}

func TestListOrgAlertNotificationTarget(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	cases := []struct {
		name        string
		want        []*OrgAlertNotificationTarget
		wantErr     bool
		mockClosure func(sqlmock.Sqlmock)
	}{
		{
			name: "OK - all organizations",
			want: []*OrgAlertNotificationTarget{
				{OrganizationID: 10, OrganizationName: "org-one", ProjectID: 1, AlertConditionID: 2, NotificationID: 100, CacheSecond: 30, NotifiedAt: now, Type: "slack", NotifySetting: "one"},
				{OrganizationID: 20, OrganizationName: "org-two", ProjectID: 1, AlertConditionID: 2, NotificationID: 200, CacheSecond: 60, NotifiedAt: now, Type: "slack", NotifySetting: "two"},
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listOrgAlertNotificationTarget)).WithArgs(uint32(1), uint32(2)).WillReturnRows(
					sqlmock.NewRows([]string{"organization_id", "organization_name", "project_id", "alert_condition_id", "notification_id", "cache_second", "notified_at", "type", "notify_setting"}).
						AddRow(uint32(10), "org-one", uint32(1), uint32(2), uint32(100), uint32(30), now, "slack", "one").
						AddRow(uint32(20), "org-two", uint32(1), uint32(2), uint32(200), uint32(60), now, "slack", "two"))
			},
		},
		{
			name:    "NG - DB error",
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listOrgAlertNotificationTarget)).WithArgs(uint32(1), uint32(2)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			c.mockClosure(mock)
			got, err := database.ListOrgAlertNotificationTarget(context.Background(), 1, 2)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected targets: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetOrgAlertNotificationTarget(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	target := &OrgAlertNotificationTarget{OrganizationID: 10, OrganizationName: "org-name", ProjectID: 1, AlertConditionID: 2, NotificationID: 100, CacheSecond: 1800, NotifiedAt: now, Type: "slack", NotifySetting: "latest"}
	cases := []struct {
		name     string
		row      *sqlmock.Rows
		queryErr error
		want     *OrgAlertNotificationTarget
		wantErr  bool
	}{
		{name: "OK - latest setting", row: sqlmock.NewRows([]string{"organization_id", "organization_name", "project_id", "alert_condition_id", "notification_id", "cache_second", "notified_at", "type", "notify_setting"}).AddRow(10, "org-name", 1, 2, 100, 1800, now, "slack", "latest"), want: target},
		{name: "NG - removed", row: sqlmock.NewRows([]string{"organization_id"}), wantErr: true},
		{name: "NG - DB error", queryErr: errors.New("DB error"), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			expectation := mock.ExpectQuery(regexp.QuoteMeta(getOrgAlertNotificationTarget)).WithArgs(uint32(10), uint32(1), uint32(2), uint32(100))
			if c.queryErr != nil {
				expectation.WillReturnError(c.queryErr)
			} else {
				expectation.WillReturnRows(c.row)
			}
			got, err := database.GetOrgAlertNotificationTarget(context.Background(), 10, 1, 2, 100)
			if (err != nil) != c.wantErr || !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: got=%v, err=%v", got, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOrgAlertNotificationTargetRequiresEnabledRelation(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{
		{name: "list targets", query: listOrgAlertNotificationTarget},
		{name: "recheck target before sending", query: getOrgAlertNotificationTarget},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(strings.ToLower(c.query), "oacn.enabled = true") {
				t.Fatalf("disabled project relation must not be notified: %s", c.query)
			}
		})
	}
}

func TestUpdateOrgAlertCondNotificationNotifiedAt(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	cases := []struct {
		name      string
		updateErr error
		wantErr   bool
	}{
		{name: "OK - four-key success update"},
		{name: "NG - DB error", updateErr: errors.New("DB error"), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			database, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			expectation := mock.ExpectExec(regexp.QuoteMeta(updateOrgAlertCondNotificationNotifiedAt)).WithArgs(now, uint32(10), uint32(1), uint32(2), uint32(100))
			if c.updateErr != nil {
				expectation.WillReturnError(c.updateErr)
			} else {
				expectation.WillReturnResult(sqlmock.NewResult(0, 1))
			}
			err = database.UpdateOrgAlertCondNotificationNotifiedAt(context.Background(), 10, 1, 2, 100, now)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLifecycleRelationInsertContract(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{
		{name: "organization project", query: insertOrgAlertCondNotificationByOrganizationProject},
		{name: "alert condition", query: insertOrgAlertCondNotificationByAlertCondition},
		{name: "organization notification", query: insertOrgAlertCondNotificationByOrgNotification},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query := strings.ToLower(c.query)
			if !strings.Contains(query, "insert ignore") {
				t.Fatalf("relation insert must be idempotent: %s", c.query)
			}
			if strings.Contains(query, "notified_at") || strings.Contains(query, "update") {
				t.Fatalf("relation insert must preserve retries and use defaults after genuine re-add: %s", c.query)
			}
			if c.name == "alert condition" && (!strings.Contains(query, "enabled") || !strings.Contains(query, "cache_second")) {
				t.Fatalf("new alert conditions must inherit the project notification setting: %s", c.query)
			}
		})
	}
}
