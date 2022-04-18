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

func TestListResourceTagByResourceID(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		projectID  uint32
		resourceID uint64
	}
	cases := []struct {
		name       string
		input      args
		want       *[]model.ResourceTag
		wantErr    bool
		mockResult *sqlmock.Rows
		mockErr    error
	}{
		{
			name:  "OK",
			input: args{projectID: 1, resourceID: 1},
			want: &[]model.ResourceTag{
				{ResourceTagID: 1, ProjectID: 1, ResourceID: 1, Tag: "tag1", CreatedAt: now, UpdatedAt: now},
				{ResourceTagID: 2, ProjectID: 1, ResourceID: 1, Tag: "tag2", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockResult: sqlmock.NewRows([]string{
				"resource_tag_id", "resource_id", "project_id", "tag", "created_at", "updated_at"}).
				AddRow(uint64(1), uint64(1), uint32(1), "tag1", now, now).
				AddRow(uint64(2), uint64(1), uint32(1), "tag2", now, now),
		},
		{
			name:    "NG DB error",
			input:   args{projectID: 1, resourceID: 1},
			want:    nil,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockResult != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListResourceTagByResourceID)).WillReturnRows(c.mockResult)
			} else if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectListResourceTagByResourceID)).WillReturnError(c.mockErr)
			}
			got, err := f.ListResourceTagByResourceID(ctx, c.input.projectID, c.input.resourceID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBulkUpsertResource(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.Resource
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.Resource{
				{ResourceID: 1, ResourceName: "name1", ProjectID: 1},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.Resource{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertResource(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestGenerateBulkUpsertResourceSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.Resource
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.Resource{
				{ResourceID: 1, ResourceName: "name1", ProjectID: 1},
			},
			wantSQL: `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), "name1", uint32(1),
			},
		},
		{
			name: "Multi",
			input: []*model.Resource{
				{ResourceID: 1, ResourceName: "name1", ProjectID: 1},
				{ResourceName: "name2", ProjectID: 1},
			},
			wantSQL: `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES
  (?, ?, ?),
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), "name1", uint32(1),
				uint64(0), "name2", uint32(1),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertResourceSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}

func TestBulkUpsertResourceTag(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.ResourceTag
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.ResourceTag{
				{ResourceTagID: 1, ResourceID: 1, ProjectID: 1, Tag: "tag1"},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.ResourceTag{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertResourceTag(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestGenerateBulkUpsertResourceTagSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.ResourceTag
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.ResourceTag{
				{ResourceTagID: 1, ResourceID: 1, ProjectID: 1, Tag: "tag1"},
			},
			wantSQL: `
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)`,
			wantParam: []interface{}{
				uint64(1), uint64(1), uint32(1), "tag1",
			},
		},
		{
			name: "Multi",
			input: []*model.ResourceTag{
				{ResourceTagID: 1, ResourceID: 1, ProjectID: 1, Tag: "tag1"},
				{ResourceTagID: 2, ResourceID: 1, ProjectID: 1, Tag: "tag2"},
				{ResourceID: 1, ProjectID: 1, Tag: "tag3"},
			},
			wantSQL: `
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES
  (?, ?, ?, ?),
  (?, ?, ?, ?),
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)`,
			wantParam: []interface{}{
				uint64(1), uint64(1), uint32(1), "tag1",
				uint64(2), uint64(1), uint32(1), "tag2",
				uint64(0), uint64(1), uint32(1), "tag3",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertResourceTagSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}
