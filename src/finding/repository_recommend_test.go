package main

import (
	"reflect"
	"testing"

	"github.com/ca-risken/core/src/finding/model"
)

func TestGenerateBulkUpsertRecommendSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.Recommend
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.Recommend{
				{RecommendID: 1, DataSource: "ds", Type: "type1", Risk: "risk", Recommendation: "recommend"},
			},
			wantSQL: `
INSERT INTO recommend
  (recommend_id, data_source, type, risk, recommendation)
VALUES
  (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  data_source=VALUES(data_source),
  type=VALUES(type),
  risk=VALUES(risk),
  recommendation=VALUES(recommendation),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint32(1), "ds", "type1", "risk", "recommend",
			},
		},
		{
			name: "Multi",
			input: []*model.Recommend{
				{RecommendID: 1, DataSource: "ds", Type: "type1", Risk: "risk", Recommendation: "recommend"},
				{RecommendID: 2, DataSource: "ds", Type: "type2", Risk: "risk", Recommendation: "recommend"},
				{DataSource: "ds", Type: "type3", Risk: "risk", Recommendation: "recommend"},
			},
			wantSQL: `
INSERT INTO recommend
  (recommend_id, data_source, type, risk, recommendation)
VALUES
  (?, ?, ?, ?, ?),
  (?, ?, ?, ?, ?),
  (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  data_source=VALUES(data_source),
  type=VALUES(type),
  risk=VALUES(risk),
  recommendation=VALUES(recommendation),
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint32(1), "ds", "type1", "risk", "recommend",
				uint32(2), "ds", "type2", "risk", "recommend",
				uint32(0), "ds", "type3", "risk", "recommend",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertRecommendSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}

func TestGenerateBulkUpsertRecommendFindingSQL(t *testing.T) {
	cases := []struct {
		name      string
		input     []*model.RecommendFinding
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single",
			input: []*model.RecommendFinding{
				{FindingID: 1, RecommendID: 1, ProjectID: 1},
			},
			wantSQL: `
INSERT INTO recommend_finding
  (finding_id, recommend_id, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), uint32(1), uint32(1),
			},
		},
		{
			name: "Multi",
			input: []*model.RecommendFinding{
				{FindingID: 1, RecommendID: 1, ProjectID: 1},
				{FindingID: 1, RecommendID: 2, ProjectID: 1},
			},
			wantSQL: `
INSERT INTO recommend_finding
  (finding_id, recommend_id, project_id)
VALUES
  (?, ?, ?),
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  updated_at=NOW()`,
			wantParam: []interface{}{
				uint64(1), uint32(1), uint32(1),
				uint64(1), uint32(2), uint32(1),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generateBulkUpsertRecommendFindingSQL(c.input)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%+v, got=%+v", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}
