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

func TestGetRecommendByDataSourceType(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		dataSource    string
		recommendType string
	}
	cases := []struct {
		name       string
		input      args
		want       *model.Recommend
		wantErr    bool
		mockResult *sqlmock.Rows
		mockErr    error
	}{
		{
			name:  "OK",
			input: args{dataSource: "ds", recommendType: "type"},
			want: &model.Recommend{
				RecommendID:    1,
				DataSource:     "ds",
				Type:           "type",
				Risk:           "risk",
				Recommendation: "recommendation",
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			wantErr: false,
			mockResult: sqlmock.NewRows([]string{
				"recommend_id", "data_source", "type", "risk", "recommendation", "created_at", "updated_at"}).
				AddRow(1, "ds", "type", "risk", "recommendation", now, now),
		},
		{
			name:    "NG DB error",
			input:   args{dataSource: "ds", recommendType: "type"},
			want:    nil,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockResult != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRecommendByDataSourceType)).WillReturnRows(c.mockResult)
			} else if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRecommendByDataSourceType)).WillReturnError(c.mockErr)
			}
			got, err := f.GetRecommendByDataSourceType(ctx, c.input.dataSource, c.input.recommendType)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBulkUpsertRecommend(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.Recommend
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.Recommend{
				{RecommendID: 1, DataSource: "ds", Type: "type1", Risk: "risk", Recommendation: "recommend"},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO recommend
  (recommend_id, data_source, type, risk, recommendation)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.Recommend{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertRecommend(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

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

func TestBulkUpsertRecommendFinding(t *testing.T) {
	f, mock, err := newMockFindingDB()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		input   []*model.RecommendFinding
		mockSQL string
	}{
		{
			name: "OK",
			input: []*model.RecommendFinding{
				{FindingID: 1, RecommendID: 1, ProjectID: 1},
			},
			mockSQL: regexp.QuoteMeta(`
INSERT INTO recommend_finding
  (finding_id, recommend_id, project_id)
VALUES`),
		},
		{
			name:    "No data",
			input:   []*model.RecommendFinding{},
			mockSQL: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockSQL != "" {
				mock.ExpectExec(c.mockSQL).WillReturnResult(sqlmock.NewResult(int64(len(c.input)), int64(len(c.input))))
			}
			err := f.BulkUpsertRecommendFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
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
