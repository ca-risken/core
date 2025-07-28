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

func TestPurgeReportFinding(t *testing.T) {
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			wantErr: false,
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectExec("delete from report_finding").WillReturnError(c.mockErr)
			} else {
				mock.ExpectExec("delete from report_finding").WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}
			err := client.PurgeReportFinding(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}

func TestListReport(t *testing.T) {
	now := time.Now()
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		projectID uint32
	}
	cases := []struct {
		name       string
		input      args
		want       *[]model.Report
		wantErr    bool
		mockResult *sqlmock.Rows
		mockErr    error
	}{
		{
			name:  "OK",
			input: args{projectID: 1},
			want: &[]model.Report{
				{ReportID: 1, ProjectID: 1, Name: "report1", Type: "type1", Status: "active", Content: "content1", CreatedAt: now, UpdatedAt: now},
				{ReportID: 2, ProjectID: 1, Name: "report2", Type: "type2", Status: "active", Content: "content2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockResult: sqlmock.NewRows([]string{
				"report_id", "project_id", "name", "type", "status", "content", "created_at", "updated_at"}).
				AddRow(uint32(1), uint32(1), "report1", "type1", "active", "content1", now, now).
				AddRow(uint32(2), uint32(1), "report2", "type2", "active", "content2", now, now),
		},
		{
			name:    "NG DB error",
			input:   args{projectID: 1},
			want:    nil,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockResult != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListReport)).WillReturnRows(c.mockResult)
			} else if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListReport)).WillReturnError(c.mockErr)
			}
			got, err := client.ListReport(ctx, c.input.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetReport(t *testing.T) {
	now := time.Now()
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		projectID uint32
		reportID  uint32
	}
	cases := []struct {
		name       string
		input      args
		want       *model.Report
		wantErr    bool
		mockResult *sqlmock.Rows
		mockErr    error
	}{
		{
			name:  "OK",
			input: args{projectID: 1, reportID: 1},
			want: &model.Report{
				ReportID: 1, ProjectID: 1, Name: "report1", Type: "type1", Status: "active", Content: "content1", CreatedAt: now, UpdatedAt: now,
			},
			wantErr: false,
			mockResult: sqlmock.NewRows([]string{
				"report_id", "project_id", "name", "type", "status", "content", "created_at", "updated_at"}).
				AddRow(uint32(1), uint32(1), "report1", "type1", "active", "content1", now, now),
		},
		{
			name:    "NG DB error",
			input:   args{projectID: 1, reportID: 1},
			want:    nil,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockResult != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetReport)).WillReturnRows(c.mockResult)
			} else if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetReport)).WillReturnError(c.mockErr)
			}
			got, err := client.GetReport(ctx, c.input.projectID, c.input.reportID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutReport(t *testing.T) {
	now := time.Now()
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		report *model.Report
	}
	cases := []struct {
		name       string
		input      args
		want       *model.Report
		wantErr    bool
		mockExecErr error
		mockQueryResult *sqlmock.Rows
		mockQueryErr    error
	}{
		{
			name: "OK Insert",
			input: args{
				report: &model.Report{
					ProjectID: 1,
					Name:      "new_report",
					Type:      "type1",
					Status:    "active",
					Content:   "content1",
				},
			},
			want: &model.Report{
				ReportID:  1,
				ProjectID: 1,
				Name:      "new_report",
				Type:      "type1",
				Status:    "active",
				Content:   "content1",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			mockQueryResult: sqlmock.NewRows([]string{
				"report_id", "project_id", "name", "type", "status", "content", "created_at", "updated_at"}).
				AddRow(uint32(1), uint32(1), "new_report", "type1", "active", "content1", now, now),
		},
		{
			name: "OK Update",
			input: args{
				report: &model.Report{
					ReportID:  1,
					ProjectID: 1,
					Name:      "updated_report",
					Type:      "type2",
					Status:    "inactive",
					Content:   "content2",
				},
			},
			want: &model.Report{
				ReportID:  1,
				ProjectID: 1,
				Name:      "updated_report",
				Type:      "type2",
				Status:    "inactive",
				Content:   "content2",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			mockQueryResult: sqlmock.NewRows([]string{
				"report_id", "project_id", "name", "type", "status", "content", "created_at", "updated_at"}).
				AddRow(uint32(1), uint32(1), "updated_report", "type2", "inactive", "content2", now, now),
		},
		{
			name: "NG Exec error",
			input: args{
				report: &model.Report{
					ProjectID: 1,
					Name:      "report1",
					Type:      "type1",
					Status:    "active",
					Content:   "content1",
				},
			},
			want:        nil,
			wantErr:     true,
			mockExecErr: errors.New("DB error on exec"),
		},
		{
			name: "NG getReportByName error",
			input: args{
				report: &model.Report{
					ProjectID: 1,
					Name:      "report1",
					Type:      "type1",
					Status:    "active",
					Content:   "content1",
				},
			},
			want:         nil,
			wantErr:      true,
			mockQueryErr: errors.New("DB error on query"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockExecErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(insertPutReport)).WillReturnError(c.mockExecErr)
			} else {
				mock.ExpectExec(regexp.QuoteMeta(insertPutReport)).WillReturnResult(sqlmock.NewResult(1, 1))
				if c.mockQueryErr != nil {
					mock.ExpectQuery(regexp.QuoteMeta(selectGetReportByName)).WillReturnError(c.mockQueryErr)
				} else {
					mock.ExpectQuery(regexp.QuoteMeta(selectGetReportByName)).WillReturnRows(c.mockQueryResult)
				}
			}
			got, err := client.PutReport(ctx, c.input.report)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("Expected error but got nil")
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
