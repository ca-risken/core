package main

import (
	"context"

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
