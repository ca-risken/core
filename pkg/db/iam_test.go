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
