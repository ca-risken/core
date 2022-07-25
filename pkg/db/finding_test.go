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

func TestListFindingTagByFindingID(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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

func TestGetRecommendByDataSourceType(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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

func TestListResourceTagByResourceID(t *testing.T) {
	now := time.Now()
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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
	f, mock, err := newMockClient()
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

func TestGeneratePrefixMatchSQLStatement(t *testing.T) {
	type args struct {
		column string
		params []string
	}
	cases := []struct {
		name      string
		input     args
		wantSQL   string
		wantParam []interface{}
	}{
		{
			name: "Single param",
			input: args{
				column: "name",
				params: []string{"aaa"},
			},
			wantSQL:   "name like ? escape '*'",
			wantParam: []interface{}{"aaa%"},
		},
		{
			name: "Multi params",
			input: args{
				column: "name",
				params: []string{"aaa", "bbb"},
			},
			wantSQL:   "name like ? escape '*' or name like ? escape '*'",
			wantParam: []interface{}{"aaa%", "bbb%"},
		},
		{
			name: "Blank param",
			input: args{
				column: "name",
				params: []string{"aaa", "bbb", ""},
			},
			wantSQL:   "name like ? escape '*' or name like ? escape '*'",
			wantParam: []interface{}{"aaa%", "bbb%"},
		},
		{
			name: "All blank param",
			input: args{
				column: "name",
				params: []string{"", "", ""},
			},
			wantSQL:   "",
			wantParam: nil,
		},
		{
			name: "Escaping param",
			input: args{
				column: "name",
				params: []string{"%", "_"},
			},
			wantSQL:   "name like ? escape '*' or name like ? escape '*'",
			wantParam: []interface{}{"*%%", "*_%"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sql, param := generatePrefixMatchSQLStatement(c.input.column, c.input.params)
			if !reflect.DeepEqual(sql, c.wantSQL) {
				t.Fatalf("Unexpected SQL response: want=%s, got=%s", c.wantSQL, sql)
			}
			if !reflect.DeepEqual(param, c.wantParam) {
				t.Fatalf("Unexpected param response: want=%+v, got=%+v", c.wantParam, param)
			}
		})
	}
}
