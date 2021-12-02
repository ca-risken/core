package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/vikyd/zero"
)

func (f *findingDB) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*[]model.Finding, error) {
	query := "select finding.* from finding inner join finding f_alias using(finding_id) "
	cond, params := generateListFindingCondition(req.ProjectId,
		req.FromScore, req.ToScore, req.FromAt, req.ToAt,
		req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	query += cond
	query += fmt.Sprintf(" order by f_alias.%s %s", req.Sort, req.Direction)
	query += fmt.Sprintf(" limit %d, %d", req.Offset, req.Limit)
	var data []model.Finding
	if err := f.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) BatchListFinding(ctx context.Context, req *finding.BatchListFindingRequest) (*[]model.Finding, error) {
	query := "select finding.* from finding "
	cond, params := generateListFindingCondition(req.ProjectId,
		req.FromScore, req.ToScore, req.FromAt, req.ToAt,
		req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	query += cond
	var data []model.Finding
	if err := f.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) ListFindingCount(
	ctx context.Context,
	projectID uint32,
	fromScore, toScore float32,
	fromAt, toAt int64,
	findingID uint64,
	dataSources, resourceNames, tags []string,
	status finding.FindingStatus) (int64, error) {
	query := "select count(*) from finding "
	cond, params := generateListFindingCondition(projectID,
		fromScore, toScore, fromAt, toAt,
		findingID, dataSources, resourceNames, tags, status)
	query += cond
	var count int64
	if err := f.Slave.WithContext(ctx).Raw(query, params...).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func generateListFindingCondition(
	projectID uint32,
	fromScore, toScore float32,
	fromAt, toAt int64,
	findingID uint64,
	dataSources, resourceNames, tags []string,
	status finding.FindingStatus) (string, []interface{}) {
	join := ""
	query := `
where
  finding.project_id = ?
  and finding.score between ? and ?
  and finding.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, projectID, fromScore, toScore, time.Unix(fromAt, 0), time.Unix(toAt, 0))
	if !zero.IsZeroVal(findingID) {
		query += " and finding.finding_id = ?"
		params = append(params, findingID)
	}
	if len(dataSources) > 0 {
		query += " and finding.data_source regexp ?"
		params = append(params, strings.Join(dataSources, "|"))
	}
	if len(resourceNames) > 0 {
		query += " and finding.resource_name regexp ?"
		params = append(params, strings.Join(resourceNames, "|"))
	}
	// EXISTS and NOT EXISTS subquery cause performance slow so used join clause instead
	if len(tags) > 0 {
		join += " inner join finding_tag ft using(finding_id)"
		query += " and ft.tag in (?)"
		params = append(params, tags)
	}
	if status == finding.FindingStatus_FINDING_ACTIVE {
		join += " left join pend_finding pf using(finding_id)"
		query += " and pf.finding_id is null"
	}
	if status == finding.FindingStatus_FINDING_PENDING {
		join += " inner join pend_finding using(finding_id)"
	}
	return join + query, params
}

const selectGetFinding = `select * from finding where project_id = ? and finding_id = ?`

func (f *findingDB) GetFinding(ctx context.Context, projectID uint32, findingID uint64) (*model.Finding, error) {
	var data model.Finding
	if err := f.Slave.WithContext(ctx).Raw(selectGetFinding, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertFinding = `
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
	updated_at=NOW()
`

func (f *findingDB) UpsertFinding(ctx context.Context, data *model.Finding) (*model.Finding, error) {
	if err := f.Master.WithContext(ctx).Exec(insertUpsertFinding,
		data.FindingID, data.Description, data.DataSource, data.DataSourceID, data.ResourceName,
		data.ProjectID, data.OriginalScore, data.Score, data.Data).Error; err != nil {
		return nil, err
	}
	return f.GetFindingByDataSource(ctx, data.ProjectID, data.DataSource, data.DataSourceID)
}

const selectGetFindingByDataSource = `select * from finding where project_id = ? and data_source = ? and data_source_id = ?`

func (f *findingDB) GetFindingByDataSource(ctx context.Context, projectID uint32, dataSource, dataSourceID string) (*model.Finding, error) {
	var result model.Finding
	if err := f.Master.WithContext(ctx).Raw(selectGetFindingByDataSource,
		projectID, dataSource, dataSourceID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

const deleteDeleteFinding = `delete from finding where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteFinding(ctx context.Context, projectID uint32, findingID uint64) error {
	if err := f.Master.WithContext(ctx).Exec(deleteDeleteFinding, projectID, findingID).Error; err != nil {
		return err
	}
	return f.DeleteTagByFindingID(ctx, projectID, findingID)
}

const deleteDeleteTagByFindingID = `delete from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) error {
	return f.Master.WithContext(ctx).Exec(deleteDeleteTagByFindingID, projectID, findingID).Error
}

const selectListFindingTag = `select * from finding_tag where project_id = ? and finding_id = ? order by %s %s limit %d, %d`

func (f *findingDB) ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := f.Slave.WithContext(ctx).Raw(
		fmt.Sprintf(selectListFindingTag, param.Sort, param.Direction, param.Offset, param.Limit),
		param.ProjectId, param.FindingId).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListFindingTagCount = `select count(*) from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error) {
	var count int64
	if err := f.Slave.WithContext(ctx).Raw(selectListFindingTagCount, param.ProjectId, param.FindingId).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const selectListFindingTagName = `
select
  distinct tag
from
  finding_tag
where
  project_id = ?
  and updated_at between ? and ?
order by %s %s limit %d, %d
`

func (f *findingDB) ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]tagName, error) {
	var data []tagName
	if err := f.Slave.WithContext(ctx).Raw(
		fmt.Sprintf(selectListFindingTagName, param.Sort, param.Direction, param.Offset, param.Limit),
		param.ProjectId, time.Unix(param.FromAt, 0), time.Unix(param.ToAt, 0)).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListFindingTagNameCount = `
select count(*) from (
  select tag
  from finding_tag
  where project_id = ? and updated_at between ? and ?
	group by project_id, tag
) tag
`

func (f *findingDB) ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error) {
	var count int64
	if err := f.Slave.WithContext(ctx).Raw(selectListFindingTagNameCount,
		param.ProjectId, time.Unix(param.FromAt, 0), time.Unix(param.ToAt, 0)).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const insertTagFinding = `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)
`

func (f *findingDB) TagFinding(ctx context.Context, tag *model.FindingTag) (*model.FindingTag, error) {
	if err := f.Master.WithContext(ctx).Exec(insertTagFinding,
		tag.FindingTagID, tag.FindingID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return f.GetFindingTagByKey(ctx, tag.ProjectID, tag.FindingID, tag.Tag)
}

const selectGetFindingTagByKey = `select * from finding_tag where project_id = ? and finding_id = ? and tag = ?`

func (f *findingDB) GetFindingTagByKey(ctx context.Context, projectID uint32, findingID uint64, tag string) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.Master.WithContext(ctx).Raw(selectGetFindingTagByKey, projectID, findingID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteUntagFinding = `delete from finding_tag where project_id = ? and finding_tag_id = ?`

func (f *findingDB) UntagFinding(ctx context.Context, projectID uint32, findingTagID uint64) error {
	return f.Master.WithContext(ctx).Exec(deleteUntagFinding, projectID, findingTagID).Error
}

func (f *findingDB) ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	var params []interface{}
	sql := `update finding f left outer join finding_tag ft using(finding_id) set f.score=0.0 where f.score > 0.0 and f.data_source = ?`

	params = append(params, req.DataSource)
	if !zero.IsZeroVal(req.ProjectId) {
		sql += " and f.project_id = ?"
		params = append(params, req.ProjectId)
	}
	if len(req.Tag) > 0 {
		sql += " and ft.tag in (?)"
		params = append(params, req.Tag)
	}
	if !zero.IsZeroVal(req.FindingId) {
		sql += " and f.finding_id = ?"
		params = append(params, req.FindingId)
	}
	return f.Master.WithContext(ctx).Exec(sql, params...).Error
}
