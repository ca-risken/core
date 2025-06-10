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

func TestListOrganizationRole(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
		userID         uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationRole
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1, name: "", userID: 0},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
				{RoleID: 2, OrganizationID: 1, Name: "role2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role or where 1=1 and r.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now).
					AddRow(uint32(2), uint32(1), "role2", now, now))
			},
		},
		{
			name: "OK with userID",
			args: args{organizationID: 1, name: "", userID: 123},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role or where 1=1 and r.organization_id = ? and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name: "OK with name",
			args: args{organizationID: 1, name: "admin", userID: 0},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "admin", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role or where 1=1 and r.organization_id = ? and r.name = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "admin", now, now))
			},
		},
		{
			name: "OK with name and userID",
			args: args{organizationID: 1, name: "admin", userID: 123},
			want: []*model.OrganizationRole{
				{RoleID: 1, OrganizationID: 1, Name: "admin", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role or where 1=1 and r.organization_id = ? and r.name = ? and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? )")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "admin", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, name: "", userID: 0},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role or where 1=1 and r.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganizationRole(ctx, c.args.organizationID, c.args.name, c.args.userID)
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
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =? and r.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
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
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =? and r.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
			name: "OK",
			args: args{organizationID: 1, name: "", roleID: 0},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
				{PolicyID: 2, OrganizationID: 1, Name: "policy2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from policy p where p.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now).
					AddRow(uint32(2), uint32(1), "policy2", now, now))
			},
		},
		{
			name: "OK with roleID",
			args: args{organizationID: 1, name: "", roleID: 456},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from policy p where p.organization_id = ? and exists(select * from organization_role_policy orp where orp.policy_id = p.policy_id and orp.role_id = ?)")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
			},
		},
		{
			name: "OK with name",
			args: args{organizationID: 1, name: "read-only", roleID: 0},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "read-only", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from policy p where p.organization_id = ? and p.name = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "read-only", now, now))
			},
		},
		{
			name: "OK with name and roleID",
			args: args{organizationID: 1, name: "read-only", roleID: 456},
			want: []*model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1, Name: "read-only", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from policy p where p.organization_id = ? and p.name = ? and exists(select * from organization_role_policy orp where orp.policy_id = p.policy_id and orp.role_id = ?)")).WillReturnRows(sqlmock.NewRows([]string{
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
				mock.ExpectQuery(regexp.QuoteMeta("select * from policy p where p.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
		roleID uint32
		userID uint32
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
			args:    args{roleID: 1, userID: 1},
			want:    &model.OrganizationRole{RoleID: 1, OrganizationID: 1, Name: "role1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
			},
		},
		{
			name:    "NG failed to attach role",
			args:    args{roleID: 1, userID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationRole)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get role",
			args:    args{roleID: 1, userID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta("select * from user where activated = 'true' and user_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"user_id", "sub", "name", "user_idp_key", "activated", "created_at", "updated_at"}).
					AddRow(uint32(1), "sub1", "user1", "key1", true, now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationRole)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.AttachOrganizationRole(ctx, c.args.roleID, c.args.userID)
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
		roleID uint32
		userID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{roleID: 1, userID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			args:    args{roleID: 1, userID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			err = db.DetachOrganizationRole(ctx, c.args.roleID, c.args.userID)
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
		roleID   uint32
		policyID uint32
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
			args:    args{roleID: 1, policyID: 1},
			want:    &model.OrganizationPolicy{PolicyID: 1, OrganizationID: 1, Name: "policy1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			args:    args{roleID: 1, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			args:    args{roleID: 1, policyID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
					"role_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "role1", now, now))
				mock.ExpectQuery(regexp.QuoteMeta(getOrganizationPolicy)).WillReturnRows(sqlmock.NewRows([]string{
					"policy_id", "organization_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), "policy1", now, now))
				mock.ExpectExec(regexp.QuoteMeta(insertAttachOrganizationPolicy)).WillReturnResult(sqlmock.NewResult(1, 1))
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
			got, err := db.AttachOrganizationPolicy(ctx, c.args.roleID, c.args.policyID)
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
		roleID   uint32
		policyID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{roleID: 1, policyID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			args:    args{roleID: 1, policyID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_role r where role_id =?")).WillReturnRows(sqlmock.NewRows([]string{
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
			err = db.DetachOrganizationPolicy(ctx, c.args.roleID, c.args.policyID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
