package main

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/src/finding/model"
)

func TestListFindingTagByFindingID(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		projectID uint32
		findingID uint64
	}
	cases := []struct {
		name       string
		input      args
		want       *[]model.FindingTag
		wantErr    bool
		mockResult *sqlmock.Rows
		mockErr    error
	}{
		{
			name:  "OK",
			input: args{projectID: 1, findingID: 1},
			want: &[]model.FindingTag{
				{FindingTagID: 1, ProjectID: 1, FindingID: 1, Tag: "tag1", CreatedAt: now, UpdatedAt: now},
				{FindingTagID: 2, ProjectID: 1, FindingID: 1, Tag: "tag2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockResult: sqlmock.NewRows([]string{
				"finding_tag_id", "finding_id", "project_id", "tag", "created_at", "updated_at"}).
				AddRow(uint64(1), uint64(1), uint32(1), "tag1", now, now).
				AddRow(uint64(2), uint64(1), uint32(1), "tag2", now, now),
		},
		{
			name:    "NG DB error",
			input:   args{projectID: 1, findingID: 1},
			want:    nil,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockResult != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListFindingTagByFindingID)).WillReturnRows(c.mockResult)
			} else if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListFindingTagByFindingID)).WillReturnError(c.mockErr)
			}
			got, err := f.ListFindingTagByFindingID(ctx, c.input.projectID, c.input.findingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBulkUpsertFinding(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.Finding
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.Finding{
				{FindingID: 1, Description: "desc", DataSource: "ds", DataSourceID: "1", ResourceName: "r", ProjectID: 1, OriginalScore: 1, Score: 1, Data: "data"},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO finding
  (finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.Finding{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestGenerateBulkUpsertFindingSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.Finding
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.Finding{
				{FindingID: 1, Description: "desc", DataSource: "ds", DataSourceID: "1", ResourceName: "r", ProjectID: 1, OriginalScore: 1, Score: 1, Data: "data"},
			},
			wantSQL: `
INSERT INTO finding
  (finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  description=VALUES(description),
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  original_score=VALUES(original_score),
  score=VALUES(score),
  data=VALUES(data),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), "desc", "ds", "1", "r", uint32(1), float32(1), float32(1), "data",
			},
		},
		{
			name: "Multi",
			input: []*model.Finding{
				{FindingID: 1, Description: "desc", DataSource: "ds", DataSourceID: "1", ResourceName: "r", ProjectID: 1, OriginalScore: 1, Score: 1, Data: "data"},
				{FindingID: 2, Description: "desc", DataSource: "ds", DataSourceID: "2", ResourceName: "r", ProjectID: 1, OriginalScore: 1, Score: 1, Data: "data"},
				{FindingID: 3, Description: "desc", DataSource: "ds", DataSourceID: "3", ResourceName: "r", ProjectID: 1, OriginalScore: 1, Score: 1, Data: "data"},
			},
			wantSQL: `
INSERT INTO finding
  (finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?),
  (?, ?, ?, ?, ?, ?, ?, ?, ?),
  (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  description=VALUES(description),
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  original_score=VALUES(original_score),
  score=VALUES(score),
  data=VALUES(data),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), "desc", "ds", "1", "r", uint32(1), float32(1), float32(1), "data",
				uint64(2), "desc", "ds", "2", "r", uint32(1), float32(1), float32(1), "data",
				uint64(3), "desc", "ds", "3", "r", uint32(1), float32(1), float32(1), "data",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertFindingSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}

func TestBulkUpsertFindingTag(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.FindingTag
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.FindingTag{
				{FindingTagID: 1, FindingID: 1, ProjectID: 1, Tag: "t1"},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.FindingTag{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertFindingTag(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestGenerateBulkUpsertFindingTagSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.FindingTag
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.FindingTag{
				{FindingTagID: 1, FindingID: 1, ProjectID: 1, Tag: "t1"},
			},
			wantSQL: `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  finding_id=VALUES(finding_id),
  project_id=VALUES(project_id),
  tag=VALUES(tag),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), uint64(1), uint32(1), "t1",
			},
		},
		{
			name: "Multi",
			input: []*model.FindingTag{
				{FindingTagID: 1, FindingID: 1, ProjectID: 1, Tag: "t1"},
				{FindingTagID: 2, FindingID: 1, ProjectID: 1, Tag: "t2"},
				{FindingID: 1, ProjectID: 1, Tag: "t3"},
			},
			wantSQL: `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES
  (?, ?, ?, ?),
  (?, ?, ?, ?),
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  finding_id=VALUES(finding_id),
  project_id=VALUES(project_id),
  tag=VALUES(tag),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), uint64(1), uint32(1), "t1",
				uint64(2), uint64(1), uint32(1), "t2",
				uint64(0), uint64(1), uint32(1), "t3",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertFindingTagSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}
