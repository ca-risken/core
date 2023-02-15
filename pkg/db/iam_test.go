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
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true,
			},
			want:    &model.User{UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, UpdatedAt: now, CreatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(insertUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub", "name", "user_idp_key", true, now, now))
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
				UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true,
			},
			want:    &model.User{UserID: 1, Sub: "sub", Name: "name", UserIdpKey: "user_idp_key", Activated: true, UpdatedAt: now, CreatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetUserBySub)).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub", "name", "user_idp_key", true, now, now))
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
				mock.ExpectExec(regexp.QuoteMeta(updateUser)).WillReturnError(errors.New("DB error"))
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
