package main

import (
	"reflect"
	"testing"

	"github.com/ca-risken/core/src/finding/model"
)

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
