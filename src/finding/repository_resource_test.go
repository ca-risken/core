package main

import (
	"reflect"
	"testing"

	"github.com/ca-risken/core/src/finding/model"
)

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
