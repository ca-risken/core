package main

import (
	"context"
	"strings"

	"github.com/ca-risken/core/src/finding/model"
)

const selectGetRecommend = `
select r.* 
from recommend r 
  inner join recommend_finding rf using(recommend_id) 
where rf.project_id=? and rf.finding_id=?
`

func (f *findingDB) GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error) {
	var data model.Recommend
	if err := f.Slave.WithContext(ctx).Raw(selectGetRecommend, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	var ret model.Recommend
	if err := f.Master.WithContext(ctx).Where("data_source=? AND type=?", data.DataSource, data.Type).Assign(data).FirstOrCreate(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (f *findingDB) UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	var ret model.RecommendFinding
	if err := f.Master.WithContext(ctx).Where("finding_id=?", data.FindingID).Assign(data).FirstOrCreate(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

const selectGetRecommendByDataSourceType = `
select r.* 
from recommend r 
where r.data_source=? and r.type=?
`

func (f *findingDB) GetRecommendByDataSourceType(ctx context.Context, dataSource, recommendType string) (*model.Recommend, error) {
	var data model.Recommend
	if err := f.Master.WithContext(ctx).Raw(selectGetRecommendByDataSourceType, dataSource, recommendType).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertRecommendSQL(data)
	return f.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertRecommendSQL(data []*model.Recommend) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO recommend
  (recommend_id, data_source, type, risk, recommendation)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?, ?, ?),`
		params = append(params, d.RecommendID, d.DataSource, d.Type, d.Risk, d.Recommendation)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  data_source=VALUES(data_source),
  type=VALUES(type),
  risk=VALUES(risk),
  recommendation=VALUES(recommendation),
  updated_at=NOW()`
	return sql, params
}

func (f *findingDB) BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertRecommendFindingSQL(data)
	return f.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertRecommendFindingSQL(data []*model.RecommendFinding) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO recommend_finding
  (finding_id, recommend_id, project_id)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?),`
		params = append(params, d.FindingID, d.RecommendID, d.ProjectID)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  updated_at=NOW()`
	return sql, params
}
