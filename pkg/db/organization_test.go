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

func TestListOrganization(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		want        []*model.Organization
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*model.Organization{
				{OrganizationID: 1, Name: "org1", Description: "desc1", CreatedAt: now, UpdatedAt: now},
				{OrganizationID: 2, Name: "org2", Description: "desc2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization o where 1 = 1 order by o.organization_id")).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "name", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), "org1", "desc1", now, now).
					AddRow(uint32(2), "org2", "desc2", now, now))
			},
		},
		{
			name:    "NG DB error",
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization o where 1 = 1 order by o.organization_id")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganization(ctx, 0, "")
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



func TestGetOrganizationByName(t *testing.T) {
	now := time.Now()
	type args struct {
		name string
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Organization
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{name: "org1"},
			want: &model.Organization{OrganizationID: 1, Name: "org1", Description: "desc1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "name", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), "org1", "desc1", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{name: "org1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetOrganizationByName(ctx, c.args.name)
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

func TestCreateOrganization(t *testing.T) {
	now := time.Now()
	type args struct {
		name        string
		description string
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Organization
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{name: "org1", description: "desc1"},
			want: &model.Organization{OrganizationID: 1, Name: "org1", Description: "desc1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectExec(regexp.QuoteMeta(insertCreateOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "name", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), "org1", "desc1", now, now))
			},
		},
		{
			name:    "NG failed to insert organization",
			args:    args{name: "org1", description: "desc1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectExec(regexp.QuoteMeta(insertCreateOrganization)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get organization",
			args:    args{name: "org1", description: "desc1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectExec(regexp.QuoteMeta(insertCreateOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.CreateOrganization(ctx, c.args.name, c.args.description)
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

func TestUpdateOrganization(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		name           string
		description    string
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Organization
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1, name: "org1", description: "desc1"},
			want: &model.Organization{OrganizationID: 1, Name: "org1", Description: "desc1", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUpdateOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "name", "description", "created_at", "updated_at"}).
					AddRow(uint32(1), "org1", "desc1", now, now))
			},
		},
		{
			name:    "NG failed to update organization",
			args:    args{organizationID: 1, name: "org1", description: "desc1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUpdateOrganization)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get organization",
			args:    args{organizationID: 1, name: "org1", description: "desc1"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateUpdateOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationByName)).WillReturnError(errors.New("DB error"))
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
			got, err := db.UpdateOrganization(ctx, c.args.organizationID, c.args.name, c.args.description)
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

func TestDeleteOrganization(t *testing.T) {
	type args struct {
		organizationID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganization)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrganization(ctx, c.args.organizationID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListProjectsInOrganization(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.Project
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1},
			want: []*model.Project{
				{ProjectID: 1, Name: "project1", CreatedAt: now, UpdatedAt: now},
				{ProjectID: 2, Name: "project2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listProjectsInOrganization)).WillReturnRows(sqlmock.NewRows([]string{
					"project_id", "name", "created_at", "updated_at"}).
					AddRow(uint32(1), "project1", now, now).
					AddRow(uint32(2), "project2", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(listProjectsInOrganization)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListProjectsInOrganization(ctx, c.args.organizationID)
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

func TestRemoveProjectsInOrganization(t *testing.T) {
	type args struct {
		organizationID uint32
		projectID      uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, projectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(removeProjectsInOrganization)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(removeProjectsInOrganization)).WillReturnError(errors.New("DB error"))
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
			err = db.RemoveProjectsInOrganization(ctx, c.args.organizationID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListOrganizationInvitation(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
	}
	cases := []struct {
		name        string
		args        args
		want        []*model.OrganizationInvitation
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1},
			want: []*model.OrganizationInvitation{
				{OrganizationID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
				{OrganizationID: 1, ProjectID: 2, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_invitation oi where 1=1 and oi.organization_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "project_id", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), now, now).
					AddRow(uint32(1), uint32(2), now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from organization_invitation oi where 1=1 and oi.organization_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOrganizationInvitation(ctx, c.args.organizationID, 0)
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

func TestPutOrganizationInvitation(t *testing.T) {
	now := time.Now()
	type args struct {
		organizationID uint32
		projectID      uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.OrganizationInvitation
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{organizationID: 1, projectID: 1},
			want: &model.OrganizationInvitation{OrganizationID: 1, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationInvitation)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationInvitation)).WillReturnRows(sqlmock.NewRows([]string{
					"organization_id", "project_id", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), now, now))
			},
		},
		{
			name:    "NG failed to insert invitation",
			args:    args{organizationID: 1, projectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationInvitation)).WillReturnError(errors.New("DB error"))
			},
		},
		{
			name:    "NG failed to get invitation",
			args:    args{organizationID: 1, projectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(putOrganizationInvitation)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetOrganizationInvitation)).WillReturnError(errors.New("DB error"))
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
			got, err := db.PutOrganizationInvitation(ctx, c.args.organizationID, c.args.projectID, "pending")
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

func TestDeleteOrganizationInvitation(t *testing.T) {
	type args struct {
		organizationID uint32
		projectID      uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{organizationID: 1, projectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationInvitation)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{organizationID: 1, projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteOrganizationInvitation)).WillReturnError(errors.New("DB error"))
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
			err = db.DeleteOrganizationInvitation(ctx, c.args.organizationID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
