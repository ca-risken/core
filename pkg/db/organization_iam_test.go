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

func TestListOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
		userID         uint32
		accessTokenID  uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK - no filters",
			args: args{organizationID: 1, name: "", userID: 0},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
				{RoleID: 2, OrganizationID: 1, Name: "role2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now).
					AddRow(uint32(2), uint32(1), "role2", now, now))
			},
		},
		{
			name: "OK - with userID filter",
			args: args{organizationID: 1, name: "", userID: 123},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ? and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name: "OK - with accessToken filter",
			args: args{organizationID: 1, name: "", userID: 0, accessTokenID: 10},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ? and exists (select * from organization_access_token_role oatr where oatr.role_id = r.role_id and oatr.access_token_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name: "OK - with userID and accessToken filters",
			args: args{organizationID: 1, name: "", userID: 123, accessTokenID: 10},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ? and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? ) and exists (select * from organization_access_token_role oatr where oatr.role_id = r.role_id and oatr.access_token_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name: "OK - with name filter",
			args: args{organizationID: 1, name: "admin", userID: 0, accessTokenID: 0},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "admin", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ? and r.name = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "admin", now, now))
			},
		},
		{
			name: "OK - with both filters",
			args: args{organizationID: 1, name: "admin", userID: 123, accessTokenID: 0},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "admin", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ? and r.name = ? and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "admin", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, name: "", userID: 0, accessTokenID: 0},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationRole + " and r.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganizationRole(ctx, c.args.organizationID, c.args.name, c.args.userID, c.args.accessTokenID)
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

func TestGetOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		roleID         uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1},
			want:    &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id = ? and r.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, roleID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id = ? and r.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationRole(ctx, c.args.organizationID, c.args.roleID)
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

func TestGetOrganizationRoleByName(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, name: "role1"},
			want:    &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRoleByName)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, name: "role1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRoleByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationRoleByName(ctx, c.args.organizationID, c.args.name)
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

func TestPutOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		role *model.OrganizationRole
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{role: &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1"}},
			want:    &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRoleByName)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name:    "NG failed to insert role",
			args:    args{role: &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationRole)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get role",
			args:    args{role: &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRoleByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.PutOrganizationRole(ctx, c.args.role)
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

func TestDeleteOrganizationRole(t *testing.T) {
	type args struct {
		organizationID uint32
		roleID         uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, roleID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationRole)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrganizationRole(ctx, c.args.organizationID, c.args.roleID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListOrganizationPolicy(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
		roleID         uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK - no filters",
			args: args{organizationID: 1, name: "", roleID: 0},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
				{PolicyID: 2, OrganizationID: 1, Name: "policy2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now).
					AddRow(uint32(2), uint32(1), "policy2", now, now))
			},
		},
		{
			name: "OK - with roleID filter",
			args: args{organizationID: 1, name: "", roleID: 456},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationPolicy + " and exists(select * from organization_role_policy orp where orp.policy_id = op.policy_id and orp.role_id = ?)")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name: "OK - with name filter",
			args: args{organizationID: 1, name: "read-only", roleID: 0},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "read-only", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationPolicy + " and op.name = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "read-only", now, now))
			},
		},
		{
			name: "OK - with both filters",
			args: args{organizationID: 1, name: "read-only", roleID: 456},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "read-only", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationPolicy + " and op.name = ? and exists(select * from organization_role_policy orp where orp.policy_id = op.policy_id and orp.role_id = ?)")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "read-only", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, name: "", roleID: 0},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(ListOrganizationPolicy)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganizationPolicy(ctx, c.args.organizationID, c.args.name, c.args.roleID)
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

func TestGetOrganizationPolicy(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		policyID       uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, policyID: 1},
			want:    &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationPolicy(ctx, c.args.organizationID, c.args.policyID)
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

func TestGetOrganizationPolicyByName(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, name: "policy1"},
			want:    &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByName)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, name: "policy1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationPolicyByName(ctx, c.args.organizationID, c.args.name)
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

func TestGetOrganizationPolicyByUserID(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		userID         uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *[]model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, userID: 1},
			want:    &[]model.OrganizationPolicy{{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now}},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByUserID)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, userID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByUserID)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationPolicyByUserID(ctx, c.args.organizationID, c.args.userID)
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

func TestPutOrganizationPolicy(t *testing.T) {
	now := time.Now()
	type args struct {
		policy *model.OrganizationPolicy
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{policy: &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1"}},
			want:    &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByName)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name:    "NG failed to insert policy",
			args:    args{policy: &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationPolicy)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get policy",
			args:    args{policy: &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicyByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.PutOrganizationPolicy(ctx, c.args.policy)
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

func TestDeleteOrganizationPolicy(t *testing.T) {
	type args struct {
		organizationID uint32
		policyID       uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, policyID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, policyID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationPolicy)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrganizationPolicy(ctx, c.args.organizationID, c.args.policyID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAttachOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		roleID         uint32
		userID         uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1, userID: 1},
			want:    &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name:    "NG failed to attach role",
			args:    args{organizationID: 1, roleID: 1, userID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationRole)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG role not found",
			args:    args{organizationID: 1, roleID: 999, userID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG user not found",
			args:    args{organizationID: 1, roleID: 1, userID: 999},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnError(gorm.ErrRecordNotFound)
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
			got, err := db.AttachOrganizationRole(ctx, c.args.organizationID, c.args.roleID, c.args.userID)
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

func TestDetachOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		roleID         uint32
		userID         uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1, userID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, roleID: 1, userID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrganizationRole)).WillReturnError(errors.New("DB error"))
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
			err = db.DetachOrganizationRole(ctx, c.args.organizationID, c.args.roleID, c.args.userID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAttachOrganizationPolicy(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		roleID         uint32
		policyID       uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationPolicy
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1, policyID: 1},
			want:    &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name:    "NG failed to attach policy",
			args:    args{organizationID: 1, roleID: 1, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationPolicy)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get policy",
			args:    args{organizationID: 1, roleID: 1, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG role not found",
			args:    args{organizationID: 1, roleID: 999, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG policy not found",
			args:    args{organizationID: 1, roleID: 1, policyID: 999},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnError(gorm.ErrRecordNotFound)
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
			got, err := db.AttachOrganizationPolicy(ctx, c.args.organizationID, c.args.policyID, c.args.roleID)
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

func TestDetachOrganizationPolicy(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		roleID         uint32
		policyID       uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, roleID: 1, policyID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, roleID: 1, policyID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole)).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrganizationPolicy)).WillReturnError(errors.New("DB error"))
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
			err = db.DetachOrganizationPolicy(ctx, c.args.organizationID, c.args.policyID, c.args.roleID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListOrganizationUserReserved(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		userIdpKey     string
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationUserReserved
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK - no filter",
			args: args{organizationID: 1, userIdpKey: ""},
			want: []*model.OrganizationUserReserved{
				{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
				{ReservedID: 2, UserIdpKey: "key2", RoleID: 11, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listOrganizationUserReserved)).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "user_idp_key", "role_id", "created_at", "updated_at"}).
					AddRow(uint32(1), "key1", uint32(10), now, now).
					AddRow(uint32(2), "key2", uint32(11), now, now))
			},
		},
		{
			name: "OK - with userIDPKey partial match",
			args: args{organizationID: 1, userIdpKey: "key"},
			want: []*model.OrganizationUserReserved{
				{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
				{ReservedID: 2, UserIdpKey: "key2", RoleID: 11, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listOrganizationUserReserved + " and ur.user_idp_key like ? escape '*' ")).WillReturnRows(sqlmock.NewRows([]string{
					"reserved_id", "user_idp_key", "role_id", "created_at", "updated_at"}).
					AddRow(uint32(1), "key1", uint32(10), now, now).
					AddRow(uint32(2), "key2", uint32(11), now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, userIdpKey: "key1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listOrganizationUserReserved + " and ur.user_idp_key like ? escape '*' ")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganizationUserReserved(ctx, c.args.organizationID, c.args.userIdpKey)
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

func TestPutOrganizationUserReserved(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *model.OrganizationUserReserved
		want        *model.OrganizationUserReserved
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			input:   &model.OrganizationUserReserved{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
			want:    &model.OrganizationUserReserved{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationUserReserved)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationUserReserved)).WillReturnRows(sqlmock.NewRows([]string{"reserved_id", "user_idp_key", "role_id", "created_at", "updated_at"}).AddRow(uint32(1), "key1", uint32(10), now, now))
			},
		},
		{
			name:    "NG failed to put organization user reserved",
			input:   &model.OrganizationUserReserved{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationUserReserved)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get organization user reserved",
			input:   &model.OrganizationUserReserved{ReservedID: 1, UserIdpKey: "key1", RoleID: 10, CreatedAt: now, UpdatedAt: now},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationUserReserved)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationUserReserved)).WillReturnError(errors.New("DB error"))
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
			got, err := db.PutOrganizationUserReserved(ctx, c.input)
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

func TestDeleteOrganizationUserReserved(t *testing.T) {
	cases := []struct {
		name           string
		organizationID uint32
		reservedID     uint32
		wantErr        bool
		mockClosure    func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "OK",
			organizationID: 1,
			reservedID:     1,
			wantErr:        false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationUserReserved)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "NG DB error",
			organizationID: 1,
			reservedID:     1,
			wantErr:        true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationUserReserved)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrganizationUserReserved(ctx, c.organizationID, c.reservedID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListOrgAccessToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	expected := sqlmock.NewRows([]string{
		"access_token_id", "token_hash", "name", "description", "organization_id", "expired_at", "last_updated_user_id", "created_at", "updated_at",
	}).AddRow(uint32(1), "hash", "token", "desc", uint32(1), now, uint32(2), now, now)
	mock.ExpectQuery(regexp.QuoteMeta("select * from organization_access_token a where a.organization_id = ?")).
		WithArgs(uint32(1)).
		WillReturnRows(expected)

	got, err := client.ListOrgAccessToken(ctx, 1, "", 0)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if len(*got) != 1 {
		t.Fatalf("Unexpected result length: %d", len(*got))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPutOrgAccessToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	token := &model.OrgAccessToken{
		TokenHash:         "hash",
		Name:              "token",
		Description:       "desc",
		OrgID:             1,
		ExpiredAt:         now,
		LastUpdatedUserID: 2,
	}

	mock.ExpectExec(regexp.QuoteMeta(insertPutOrgAccessToken)).
		WithArgs(token.AccessTokenID, token.TokenHash, token.Name, token.Description, token.OrgID, sqlmock.AnyArg(), token.LastUpdatedUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{
		"access_token_id", "token_hash", "name", "description", "organization_id", "expired_at", "last_updated_user_id", "created_at", "updated_at",
	}).AddRow(uint32(1), "hash", "token", "desc", uint32(1), now, uint32(2), now, now)
	mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByUniqueKey)).
		WithArgs(uint32(1), "token").
		WillReturnRows(rows)

	got, err := client.PutOrgAccessToken(ctx, token)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if got.AccessTokenID != 1 {
		t.Fatalf("Unexpected access token id: %d", got.AccessTokenID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteOrgAccessToken(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		orgID   uint32
		tokenID uint32
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			orgID:   1,
			tokenID: 2,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgAccessToken)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:    "NG DB error",
			orgID:   1,
			tokenID: 2,
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrgAccessToken)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(errors.New("delete error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			err = client.DeleteOrgAccessToken(ctx, c.orgID, c.tokenID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetActiveOrgAccessTokenHash(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cases := []struct {
		name    string
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"access_token_id", "token_hash", "name", "description", "organization_id", "expired_at", "last_updated_user_id", "created_at", "updated_at",
				}).AddRow(uint32(2), "hash", "token", "desc", uint32(1), now, uint32(9), now, now)
				mock.ExpectQuery(regexp.QuoteMeta(selectGetActiveOrgAccessTokenHash)).
					WithArgs(uint32(1), uint32(2), "hashed").
					WillReturnRows(rows)
			},
		},
		{
			name:    "NG record not found",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetActiveOrgAccessTokenHash)).
					WithArgs(uint32(1), uint32(2), "hashed").
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetActiveOrgAccessTokenHash)).
					WithArgs(uint32(1), uint32(2), "hashed").
					WillReturnError(errors.New("select error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			got, err := client.GetActiveOrgAccessTokenHash(ctx, 1, 2, "hashed")
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if err == nil && got.AccessTokenID != 2 {
				t.Fatalf("unexpected access token id: %d", got.AccessTokenID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestExistsOrgAccessTokenMaintainer(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		want    bool
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name: "OK exists",
			want: true,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id"}).AddRow(uint32(10))
				mock.ExpectQuery(regexp.QuoteMeta(selectExistsOrgAccessTokenMaintainer)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnRows(rows)
			},
		},
		{
			name: "OK not found",
			want: false,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectExistsOrgAccessTokenMaintainer)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectExistsOrgAccessTokenMaintainer)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(errors.New("select error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			got, err := client.ExistsOrgAccessTokenMaintainer(ctx, 1, 2)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !c.wantErr && got != c.want {
				t.Fatalf("unexpected result: want=%t, got=%t", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetOrgTokenPolicy(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cases := []struct {
		name    string
		wantLen int
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			wantLen: 1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"policy_id", "name", "organization_id", "action_ptn", "created_at", "updated_at"}).
					AddRow(uint32(1), "policy", uint32(2), "organization/.*", now, now)
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgTokenPolicy)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnRows(rows)
			},
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgTokenPolicy)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(errors.New("select error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			got, err := client.GetOrgTokenPolicy(ctx, 1, 2)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !c.wantErr && len(*got) != c.wantLen {
				t.Fatalf("unexpected result length: want=%d, got=%d", c.wantLen, len(*got))
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAttachOrgAccessTokenRole(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cases := []struct {
		name    string
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				expectOrganizationRoleExists(mock, now, 1, 3)
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				rows := sqlmock.NewRows([]string{"access_token_id", "role_id", "created_at", "updated_at"}).
					AddRow(uint32(2), uint32(3), now, now)
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnRows(rows)
			},
		},
		{
			name:    "NG access token not found",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByID)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG access token lookup error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByID)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(errors.New("token lookup error"))
			},
		},
		{
			name:    "NG role not found",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole+" and r.organization_id = ?")).
					WithArgs(uint32(3), uint32(1)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG role lookup error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole+" and r.organization_id = ?")).
					WithArgs(uint32(3), uint32(1)).
					WillReturnError(errors.New("role lookup error"))
			},
		},
		{
			name:    "NG insert error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				expectOrganizationRoleExists(mock, now, 1, 3)
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnError(errors.New("insert error"))
			},
		},
		{
			name:    "NG fetch relation error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				expectOrganizationRoleExists(mock, now, 1, 3)
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnError(errors.New("select relation error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			got, err := client.AttachOrgAccessTokenRole(ctx, 1, 3, 2)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if err == nil && (got.AccessTokenID != 2 || got.RoleID != 3) {
				t.Fatalf("unexpected relation: %+v", got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDetachOrgAccessTokenRole(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cases := []struct {
		name    string
		wantErr bool
		mock    func(sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				expectOrganizationRoleExists(mock, now, 1, 3)
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG access token not found",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByID)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG access token lookup error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByID)).
					WithArgs(uint32(1), uint32(2)).
					WillReturnError(errors.New("lookup error"))
			},
		},
		{
			name:    "NG role not found",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole+" and r.organization_id = ?")).
					WithArgs(uint32(3), uint32(1)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "NG role lookup error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole+" and r.organization_id = ?")).
					WithArgs(uint32(3), uint32(1)).
					WillReturnError(errors.New("role lookup error"))
			},
		},
		{
			name:    "NG delete error",
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				expectOrgAccessTokenExists(mock, now, 1, 2)
				expectOrganizationRoleExists(mock, now, 1, 3)
				mock.ExpectExec(regexp.QuoteMeta(deleteDetachOrgAccessTokenRole)).
					WithArgs(uint32(2), uint32(3)).
					WillReturnError(errors.New("delete error"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("failed to prepare mock db: %v", err)
			}
			c.mock(mock)
			err = client.DetachOrgAccessTokenRole(ctx, 1, 3, 2)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func expectOrgAccessTokenExists(mock sqlmock.Sqlmock, now time.Time, orgID, accessTokenID uint32) {
	rows := sqlmock.NewRows([]string{
		"access_token_id", "token_hash", "name", "description", "organization_id", "expired_at", "last_updated_user_id", "created_at", "updated_at",
	}).AddRow(accessTokenID, "hash", "token", "desc", orgID, now, uint32(9), now, now)
	mock.ExpectQuery(regexp.QuoteMeta(selectGetOrgAccessTokenByID)).
		WithArgs(orgID, accessTokenID).
		WillReturnRows(rows)
}

func expectOrganizationRoleExists(mock sqlmock.Sqlmock, now time.Time, orgID, roleID uint32) {
	rows := sqlmock.NewRows([]string{"role_id", "organization_id", "name", "created_at", "updated_at"}).
		AddRow(roleID, orgID, "role", now, now)
	mock.ExpectQuery(regexp.QuoteMeta(getOrganizationRole+" and r.organization_id = ?")).
		WithArgs(roleID, orgID).
		WillReturnRows(rows)
}
