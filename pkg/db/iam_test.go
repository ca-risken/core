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

func TestCreateUser(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *model.User
		want        *model.User
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false,
			},
			want:    &model.User{UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false, UpdatedAt: now, CreatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(insertUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub", "name", "user_idp_key", true, false, now, now))
			},
		},
		{
			name: "NG failed to insert user",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(insertUser)).WillReturnError(errors.New("DB error"))
				//				mock.ExpectQuery(regexp.QuoteMeta(insertPutUser)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name: "NG failed to get user",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(insertUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnError(errors.New("DB error"))
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
			got, err := db.CreateUser(ctx, c.input)
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

func TestPutUser(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *model.User
		want        *model.User
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false,
			},
			want:    &model.User{UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false, UpdatedAt: now, CreatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub", "name", "user_idp_key", true, false, now, now))
			},
		},
		{
			name: "NG failed to insert user",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnError(errors.New("DB error"))
				//				mock.ExpectQuery(regexp.QuoteMeta(insertPutUser)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name: "NG failed to get user",
			input: &model.User{
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, IsAdmin: false,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnError(errors.New("DB error"))
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
			got, err := db.PutUser(ctx, c.input)
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

func TestListUserReserved(t *testing.T) {
	now := time.Now()
	type args struct {
		projectID  uint32
		userIdpKey string
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.UserReserved
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{projectID: 1},
			want: &[]model.UserReserved{
				{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1", CreatedAt: now, UpdatedAt: now},
				{ReservedID: 2, RoleID: 1, UserIdpKey: "uik2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listUserReserved)).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "role_id", "user_idp_key", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "uik1", now, now).
					AddRow(uint32(2), uint32(1), "uik2", now, now))
			},
		},
		{
			name: "OK input UserIdpKey is exist",
			args: args{projectID: 1, userIdpKey: "uik1"},
			want: &[]model.UserReserved{
				{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listUserReserved + " and ur.user_idp_key like ? escape '*'")).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "role_id", "user_idp_key", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "uik1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listUserReserved)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListUserReserved(ctx, c.args.projectID, c.args.userIdpKey)
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

func TestListUserReservedWithProjectID(t *testing.T) {
	type args struct {
		userIdpKey string
	}

	cases := []struct {
		name        string
		args        args
		want        *[]UserReservedWithProjectID
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{userIdpKey: "uik"},
			want: &[]UserReservedWithProjectID{
				{ProjectID: 1, RoleID: 1},
				{ProjectID: 1, RoleID: 2},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listUserReservedWithProjectID)).WillReturnRows(sqlmock.NewRows([]string{
					"project_id", "role_id"}).
					AddRow(uint32(1), uint32(1)).
					AddRow(uint32(1), uint32(2)))
			},
		},
		{
			name:    "NG DB error",
			args:    args{userIdpKey: "uik"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listUserReservedWithProjectID)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListUserReservedWithProjectID(ctx, c.args.userIdpKey)
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

func TestPutUserReserved(t *testing.T) {
	now := time.Now()
	type args struct {
		data *model.UserReserved
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK update",
			args: args{
				data: &model.UserReserved{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1"},
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_reserved` WHERE reserved_id = ? ORDER BY `user_reserved`.`reserved_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "role_id", "user_idp_key", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "uik1", now, now))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_reserved`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "OK insert",
			args: args{
				data: &model.UserReserved{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_reserved`")).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "role_id", "user_idp_key", "created_at", "updated_at"}))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `user_reserved`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "NG DB error",
			args: args{
				data: &model.UserReserved{ReservedID: 1, RoleID: 1, UserIdpKey: "uik1"},
			},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_reserved` WHERE reserved_id = ? ORDER BY `user_reserved`.`reserved_id` LIMIT 1")).WillReturnError(errors.New("DB error"))
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
			_, err = db.PutUserReserved(ctx, c.args.data)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteUserReserved(t *testing.T) {
	type args struct {
		projectID  uint32
		reservedID uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{projectID: 1, reservedID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteUserReserved)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{projectID: 1, reservedID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteUserReserved)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteUserReserved(ctx, c.args.projectID, c.args.reservedID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateUserAdmin(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		userID      uint32
		isAdmin     bool
		want        *model.User
		wantError   bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK Update to admin",
			userID:  1,
			isAdmin: true,
			want:    &model.User{UserID: 1, Sub: "test@example.com", Name: "Test User", Activated: true, IsAdmin: true, CreatedAt: now, UpdatedAt: now},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user` SET `is_admin`=\\?,`updated_at`=\\? WHERE user_id = \\?").
					WithArgs(true, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				rows := sqlmock.NewRows([]string{"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "Test User", nil, true, true, now, now)
				mock.ExpectQuery("select \\* from user where activated = 'true'").
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:    "OK Update to non-admin",
			userID:  1,
			isAdmin: false,
			want:    &model.User{UserID: 1, Sub: "test@example.com", Name: "Test User", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user` SET `is_admin`=\\?,`updated_at`=\\? WHERE user_id = \\?").
					WithArgs(false, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				rows := sqlmock.NewRows([]string{"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "Test User", nil, true, false, now, now)
				mock.ExpectQuery("select \\* from user where activated = 'true'").
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:      "NG Update error",
			userID:    1,
			isAdmin:   true,
			wantError: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user` SET `is_admin`=\\?,`updated_at`=\\? WHERE user_id = \\?").
					WithArgs(true, sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to create mock client: %v", err)
			}
			c.mockClosure(mock)
			got, err := client.UpdateUserAdmin(ctx, c.userID, c.isAdmin)
			if c.wantError {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet expectations: %v", err)
			}
		})
	}
}

func TestListUser(t *testing.T) {
	now := time.Now()
	type args struct {
		activated      bool
		projectID      uint32
		organizationID uint32
		name           string
		userID         uint32
		admin          bool
		userIdpKey     string
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.User
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK - basic filter",
			args: args{activated: true, projectID: 1, organizationID: 0, name: "", userID: 0, admin: false, userIdpKey: ""},
			want: &[]model.User{
				{UserID: 1, Sub: "sub1", Name: "user1", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
				{UserID: 2, Sub: "sub2", Name: "user2", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedQuery := `
select
  u.*
from
  user u
where
  activated = ? and exists (select * from user_role ur inner join role r using(role_id, project_id) where ur.user_id = u.user_id and ur.project_id = ?)`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "", true, false, now, now).
					AddRow(uint32(2), "sub2", "user2", "", true, false, now, now))
			},
		},
		{
			name: "OK - with organization filter",
			args: args{activated: true, projectID: 0, organizationID: 1, name: "", userID: 0, admin: false, userIdpKey: ""},
			want: &[]model.User{
				{UserID: 1, Sub: "sub1", Name: "user1", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedQuery := `
select
  u.*
from
  user u
where
  activated = ? and exists (select * from user_organization_role uor inner join organization_role r using(role_id) where uor.user_id = u.user_id and r.organization_id = ?)`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "", true, false, now, now))
			},
		},
		{
			name: "OK - with name filter",
			args: args{activated: true, projectID: 0, organizationID: 0, name: "test", userID: 0, admin: false, userIdpKey: ""},
			want: &[]model.User{
				{UserID: 1, Sub: "sub1", Name: "testuser", Activated: true, IsAdmin: false, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedQuery := `
select
  u.*
from
  user u
where
  activated = ? and (u.name like ? escape '*' or u.user_idp_key like ? escape '*' )`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "is_admin", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "testuser", "", true, false, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{activated: true, projectID: 1, organizationID: 0, name: "", userID: 0, admin: false, userIdpKey: ""},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedQuery := `
select
  u.*
from
  user u
where
  activated = ? and exists (select * from user_role ur inner join role r using(role_id, project_id) where ur.user_id = u.user_id and ur.project_id = ?)`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListUser(ctx, c.args.activated, c.args.projectID, c.args.organizationID, c.args.name, c.args.userID, c.args.admin, c.args.userIdpKey)
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
