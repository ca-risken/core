package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/cenkalti/backoff/v4"
	"github.com/vikyd/zero"
)

const (
	escapeString = "*"
)

type FindingRepository interface {
	// Finding
	ListFinding(context.Context, *finding.ListFindingRequest) (*[]model.Finding, error)
	BatchListFinding(context.Context, *finding.BatchListFindingRequest) (*[]model.Finding, error)
	ListFindingCount(
		ctx context.Context,
		projectID, alertID uint32,
		fromScore, toScore float32,
		fromAt, toAt int64,
		findingID uint64,
		dataSources, resourceNames, tags []string,
		status finding.FindingStatus,
	) (int64, error)
	GetFinding(context.Context, uint32, uint64, bool) (*model.Finding, error)
	GetFindingByDataSource(context.Context, uint32, string, string) (*model.Finding, error)
	UpsertFinding(context.Context, *model.Finding) (*model.Finding, error)
	DeleteFinding(context.Context, uint32, uint64) error
	ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error)
	ListFindingTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) (*[]model.FindingTag, error)
	ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error)
	ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]TagName, error)
	ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error)
	GetFindingTagByKey(context.Context, uint32, uint64, string) (*model.FindingTag, error)
	TagFinding(context.Context, *model.FindingTag) (*model.FindingTag, error)
	UntagFinding(context.Context, uint32, uint64) error
	ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error
	BulkUpsertFinding(ctx context.Context, data []*model.Finding) error
	BulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error

	// Resource
	ListResource(context.Context, *finding.ListResourceRequest) (*[]model.Resource, error)
	ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error)
	GetResource(context.Context, uint32, uint64) (*model.Resource, error)
	GetResourceByName(context.Context, uint32, string) (*model.Resource, error)
	UpsertResource(context.Context, *model.Resource) (*model.Resource, error)
	DeleteResource(context.Context, uint32, uint64) error
	ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error)
	ListResourceTagByResourceID(ctx context.Context, projectID uint32, resourceID uint64) (*[]model.ResourceTag, error)
	ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error)
	ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]TagName, error)
	ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error)
	GetResourceTagByKey(context.Context, uint32, uint64, string) (*model.ResourceTag, error)
	TagResource(context.Context, *model.ResourceTag) (*model.ResourceTag, error)
	UntagResource(context.Context, uint32, uint64) error
	BulkUpsertResource(ctx context.Context, data []*model.Resource) error
	BulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error

	// PendFinding
	GetPendFinding(ctx context.Context, projectID uint32, findingID uint64) (*model.PendFinding, error)
	UpsertPendFinding(ctx context.Context, pend *finding.PendFindingForUpsert) (*model.PendFinding, error)
	DeletePendFinding(ctx context.Context, projectID uint32, findingID uint64) error

	// FindingSetting
	ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error)
	GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error)
	GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error)
	UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error)
	DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error

	// Recommend
	GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error)
	UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error)
	UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error)
	GetRecommendByDataSourceType(ctx context.Context, dataSource, recommendType string) (*model.Recommend, error)
	BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error
	BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error
}

var _ FindingRepository = (*Client)(nil)

func (c *Client) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*[]model.Finding, error) {
	query := "select finding.* from finding inner join finding f_alias using(finding_id) "
	cond, params := generateListFindingCondition(
		req.ProjectId, req.AlertId,
		req.FromScore, req.ToScore, req.FromAt, req.ToAt,
		req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	query += cond
	query += fmt.Sprintf(" order by f_alias.%s %s", req.Sort, req.Direction)
	query += fmt.Sprintf(" limit %d, %d", req.Offset, req.Limit)
	var data []model.Finding
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) BatchListFinding(ctx context.Context, req *finding.BatchListFindingRequest) (*[]model.Finding, error) {
	query := "select finding.* from finding "
	cond, params := generateListFindingCondition(
		req.ProjectId, req.AlertId,
		req.FromScore, req.ToScore, req.FromAt, req.ToAt,
		req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	query += cond
	var data []model.Finding
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListFindingCount(
	ctx context.Context,
	projectID, alertID uint32,
	fromScore, toScore float32,
	fromAt, toAt int64,
	findingID uint64,
	dataSources, resourceNames, tags []string,
	status finding.FindingStatus) (int64, error) {
	query := "select count(*) from finding "
	cond, params := generateListFindingCondition(
		projectID, alertID,
		fromScore, toScore, fromAt, toAt,
		findingID, dataSources, resourceNames, tags, status)
	query += cond
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func generateListFindingCondition(
	projectID, alertID uint32,
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
	if findingID != 0 {
		query += " and finding.finding_id = ?"
		params = append(params, findingID)
	}
	if alertID != 0 {
		query += " and exists(select * from rel_alert_finding raf where raf.finding_id=finding.finding_id and raf.alert_id = ?)"
		params = append(params, alertID)
	}
	if len(dataSources) > 0 {
		sql, sqlParams := generatePrefixMatchSQLStatement("finding.data_source", dataSources)
		if sql != "" {
			query += fmt.Sprintf(" and (%s)", sql)
			params = append(params, sqlParams...)
		}
	}
	if len(resourceNames) > 0 {
		sql, sqlParams := generatePrefixMatchSQLStatement("finding.resource_name", resourceNames)
		if sql != "" {
			query += fmt.Sprintf(" and (%s)", sql)
			params = append(params, sqlParams...)
		}
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

func escapeLikeParam(s string) string {
	s = strings.ReplaceAll(s, "*", escapeString+"*")
	s = strings.ReplaceAll(s, "_", escapeString+"_")
	s = strings.ReplaceAll(s, "%", escapeString+"%")
	return s
}

func generatePrefixMatchSQLStatement(column string, params []string) (sql string, sqlParams []interface{}) {
	for _, p := range params {
		if p == "" {
			continue
		}
		if sql != "" {
			sql += " or "
		}
		sql += fmt.Sprintf("%s like ? escape '%s'", column, escapeString)
		sqlParams = append(sqlParams, escapeLikeParam(p)+"%") // prefix match
	}
	return sql, sqlParams
}

const selectGetFinding = `select * from finding where project_id = ? and finding_id = ?`

func (c *Client) GetFinding(ctx context.Context, projectID uint32, findingID uint64, immediately bool) (*model.Finding, error) {
	var data model.Finding
	var err error
	if immediately {
		err = c.Master.WithContext(ctx).Raw(selectGetFinding, projectID, findingID).First(&data).Error
	} else {
		err = c.Slave.WithContext(ctx).Raw(selectGetFinding, projectID, findingID).First(&data).Error
	}
	if err != nil {
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

func (c *Client) UpsertFinding(ctx context.Context, data *model.Finding) (*model.Finding, error) {
	operation := func() (*model.Finding, error) {
		return c.upsertFinding(ctx, data)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "UpsertFinding"))
}

func (c *Client) upsertFinding(ctx context.Context, data *model.Finding) (*model.Finding, error) {
	if err := c.Master.WithContext(ctx).Exec(insertUpsertFinding,
		data.FindingID, data.Description, data.DataSource, data.DataSourceID, data.ResourceName,
		data.ProjectID, data.OriginalScore, data.Score, data.Data).Error; err != nil {
		return nil, err
	}
	return c.GetFindingByDataSource(ctx, data.ProjectID, data.DataSource, data.DataSourceID)
}

const selectGetFindingByDataSource = `select * from finding where project_id = ? and data_source = ? and data_source_id = ?`

func (c *Client) GetFindingByDataSource(ctx context.Context, projectID uint32, dataSource, dataSourceID string) (*model.Finding, error) {
	var result model.Finding
	if err := c.Master.WithContext(ctx).Raw(selectGetFindingByDataSource,
		projectID, dataSource, dataSourceID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

const deleteDeleteFinding = `delete from finding where project_id = ? and finding_id = ?`

func (c *Client) DeleteFinding(ctx context.Context, projectID uint32, findingID uint64) error {
	if err := c.Master.WithContext(ctx).Exec(deleteDeleteFinding, projectID, findingID).Error; err != nil {
		return err
	}
	return c.DeleteTagByFindingID(ctx, projectID, findingID)
}

const deleteDeleteTagByFindingID = `delete from finding_tag where project_id = ? and finding_id = ?`

func (c *Client) DeleteTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) error {
	return c.Master.WithContext(ctx).Exec(deleteDeleteTagByFindingID, projectID, findingID).Error
}

const selectListFindingTag = `select * from finding_tag where project_id = ? and finding_id = ? order by %s %s limit %d, %d`

func (c *Client) ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := c.Slave.WithContext(ctx).Raw(
		fmt.Sprintf(selectListFindingTag, param.Sort, param.Direction, param.Offset, param.Limit),
		param.ProjectId, param.FindingId).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListFindingTagByFindingID = `select * from finding_tag where project_id = ? and finding_id = ? `

func (c *Client) ListFindingTagByFindingID(ctx context.Context, projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := c.Master.WithContext(ctx).Raw(selectListFindingTagByFindingID, projectID, findingID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListFindingTagCount = `select count(*) from finding_tag where project_id = ? and finding_id = ?`

func (c *Client) ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error) {
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(selectListFindingTagCount, param.ProjectId, param.FindingId).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const selectListFindingTagName = `
select
  tag
from
  finding_tag
where
  project_id = ?
  and updated_at between ? and ?
group by project_id, tag
order by %s %s limit %d, %d
`

func (c *Client) ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]TagName, error) {
	var data []TagName
	if err := c.Slave.WithContext(ctx).Raw(
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

func (c *Client) ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error) {
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(selectListFindingTagNameCount,
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

func (c *Client) TagFinding(ctx context.Context, tag *model.FindingTag) (*model.FindingTag, error) {
	operation := func() (*model.FindingTag, error) {
		return c.tagFinding(ctx, tag)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "TagFinding"))
}

func (c *Client) tagFinding(ctx context.Context, tag *model.FindingTag) (*model.FindingTag, error) {
	if err := c.Master.WithContext(ctx).Exec(insertTagFinding,
		tag.FindingTagID, tag.FindingID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return c.GetFindingTagByKey(ctx, tag.ProjectID, tag.FindingID, tag.Tag)
}

const selectGetFindingTagByKey = `select * from finding_tag where project_id = ? and finding_id = ? and tag = ?`

func (c *Client) GetFindingTagByKey(ctx context.Context, projectID uint32, findingID uint64, tag string) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := c.Master.WithContext(ctx).Raw(selectGetFindingTagByKey, projectID, findingID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteUntagFinding = `delete from finding_tag where project_id = ? and finding_tag_id = ?`

func (c *Client) UntagFinding(ctx context.Context, projectID uint32, findingTagID uint64) error {
	return c.Master.WithContext(ctx).Exec(deleteUntagFinding, projectID, findingTagID).Error
}

func (c *Client) ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	operation := func() error {
		return c.clearScoreFinding(ctx, req)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "ClearScoreFinding"))
}

func (c *Client) clearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error {
	var params []interface{}
	sql := `update finding f left outer join finding_tag ft using(finding_id) set f.score=0.0 where f.score > 0.0 and f.data_source = ?`

	params = append(params, req.DataSource)
	if req.ProjectId != 0 {
		sql += " and f.project_id = ?"
		params = append(params, req.ProjectId)
	}
	if len(req.Tag) > 0 {
		sql += " and ft.tag in (?)"
		params = append(params, req.Tag)
	}
	if req.FindingId != 0 {
		sql += " and f.finding_id = ?"
		params = append(params, req.FindingId)
	}
	if req.BeforeAt != 0 {
		sql += " and f.updated_at < ?"
		params = append(params, time.Unix(req.BeforeAt, 0))
	}
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
}

func (c *Client) BulkUpsertFinding(ctx context.Context, data []*model.Finding) error {
	operation := func() error {
		return c.bulkUpsertFinding(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertFinding"))
}

func (c *Client) bulkUpsertFinding(ctx context.Context, data []*model.Finding) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertFindingSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertFindingSQL(data []*model.Finding) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO finding
  (finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?, ?, ?, ?, ?, ?, ?),`
		params = append(params, d.FindingID, d.Description, d.DataSource, d.DataSourceID,
			d.ResourceName, d.ProjectID, d.OriginalScore, d.Score, d.Data)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  description=VALUES(description),
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  original_score=VALUES(original_score),
  score=VALUES(score),
  data=VALUES(data),
  updated_at=NOW()`
	return sql, params
}

func (c *Client) BulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error {
	operation := func() error {
		return c.bulkUpsertFindingTag(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertFindingTag"))
}

func (c *Client) bulkUpsertFindingTag(ctx context.Context, data []*model.FindingTag) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertFindingTagSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertFindingTagSQL(data []*model.FindingTag) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?, ?),`
		params = append(params, d.FindingTagID, d.FindingID, d.ProjectID, d.Tag)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  finding_id=VALUES(finding_id),
  project_id=VALUES(project_id),
  tag=VALUES(tag),
  updated_at=NOW()`
	return sql, params
}

func (c *Client) ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error) {
	var param []interface{}
	query := "select * from finding_setting where project_id=?"
	param = append(param, req.ProjectId)
	if req.Status != finding.FindingSettingStatus_SETTING_UNKNOWN {
		query += " and status=?"
		param = append(param, getStatusString(req.Status))
	}
	var data []model.FindingSetting
	if err := c.Slave.WithContext(ctx).Raw(query, param...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func getStatusString(status finding.FindingSettingStatus) string {
	switch status {
	case finding.FindingSettingStatus_SETTING_ACTIVE:
		return "ACTIVE"
	case finding.FindingSettingStatus_SETTING_DEACTIVE:
		return "DEACTIVE"
	default:
		return "UNKNOWN"
	}
}

const selectGetFindingSetting = `select * from finding_setting where project_id=? and finding_setting_id=?`

func (c *Client) GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := c.Slave.WithContext(ctx).Raw(selectGetFindingSetting, projectID, findingSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingSettingByResource = `select * from finding_setting where project_id=? and resource_name=?`

func (c *Client) GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error) {
	var data model.FindingSetting
	if err := c.Master.WithContext(ctx).Raw(selectGetFindingSettingByResource, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error) {
	var retData model.FindingSetting
	if err := c.Master.WithContext(ctx).Where("project_id=? AND resource_name=?", data.ProjectID, data.ResourceName).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

const deleteDeleteFindingSetting = `delete from finding_setting where project_id = ? and finding_setting_id = ?`

func (c *Client) DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteDeleteFindingSetting, projectID, findingSettingID).Error
}

func (c *Client) ListResource(ctx context.Context, req *finding.ListResourceRequest) (*[]model.Resource, error) {
	query := `
select 
  r.*
from
  resource r
where
  r.project_id = ?
  and r.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if !zero.IsZeroVal(req.ResourceId) {
		query += " and r.resource_id = ?"
		params = append(params, req.ResourceId)
	}
	if len(req.ResourceName) > 0 {
		sql, sqlParams := generatePrefixMatchSQLStatement("r.resource_name", req.ResourceName)
		if sql != "" {
			query += fmt.Sprintf(" and (%s)", sql)
			params = append(params, sqlParams...)
		}
	}
	if len(req.Tag) > 0 {
		for _, tag := range req.Tag {
			query += " and exists (select * from resource_tag rt where rt.resource_id=r.resource_id and rt.tag = ?)"
			params = append(params, tag)
		}
	}
	if req.FromSumScore > 0 {
		query += " and exists (select resource_name from finding where resource_name=r.resource_name group by resource_name having sum(COALESCE(score, 0)) between ? and ?)"
		params = append(params, req.FromSumScore, req.ToSumScore)
	}
	query += fmt.Sprintf(" order by %s %s", req.Sort, req.Direction)
	query += fmt.Sprintf(" limit %d, %d", req.Offset, req.Limit)
	var data []model.Resource
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error) {
	query := `
select count(*) from (
  select r.*
  from
    resource r
  where
    r.project_id = ?
    and r.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if !zero.IsZeroVal(req.ResourceId) {
		query += " and r.resource_id = ?"
		params = append(params, req.ResourceId)
	}
	if len(req.ResourceName) > 0 {
		sql, sqlParams := generatePrefixMatchSQLStatement("r.resource_name", req.ResourceName)
		if sql != "" {
			query += fmt.Sprintf(" and (%s)", sql)
			params = append(params, sqlParams...)
		}
	}
	if len(req.Tag) > 0 {
		for _, tag := range req.Tag {
			query += " and exists (select * from resource_tag rt where rt.resource_id=r.resource_id and rt.tag = ?)"
			params = append(params, tag)
		}
	}
	if req.FromSumScore > 0 {
		query += " and exists (select resource_name from finding where resource_name=r.resource_name group by resource_name having sum(COALESCE(score, 0)) between ? and ?)"
		params = append(params, req.FromSumScore, req.ToSumScore)
	}
	query += ") as resource"
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const selectGetResource = `select * from resource where project_id = ? and resource_id = ?`

func (c *Client) GetResource(ctx context.Context, projectID uint32, resourceID uint64) (*model.Resource, error) {
	var data model.Resource
	if err := c.Slave.WithContext(ctx).Raw(selectGetResource, projectID, resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDeleteResource = `delete from resource where project_id = ? and resource_id = ?`

func (c *Client) DeleteResource(ctx context.Context, projectID uint32, resourceID uint64) error {
	if err := c.Master.WithContext(ctx).Exec(deleteDeleteResource, projectID, resourceID).Error; err != nil {
		return err
	}
	return c.DeleteTagByResourceID(ctx, projectID, resourceID)
}

const deleteDeleteTagByResourceID = `delete from resource_tag where project_id = ? and resource_id = ?`

func (c *Client) DeleteTagByResourceID(ctx context.Context, projectID uint32, resourceID uint64) error {
	return c.Master.WithContext(ctx).Exec(deleteDeleteTagByResourceID, projectID, resourceID).Error
}

const selectListResourceTag = `select * from resource_tag where project_id = ? and resource_id = ? order by %s %s limit %d, %d`

func (c *Client) ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error) {
	var data []model.ResourceTag
	if err := c.Slave.WithContext(ctx).Raw(
		fmt.Sprintf(selectListResourceTag, param.Sort, param.Direction, param.Offset, param.Limit),
		param.ProjectId, param.ResourceId).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListResourceTagByResourceID = `select * from resource_tag where project_id = ? and resource_id = ?`

func (c *Client) ListResourceTagByResourceID(ctx context.Context, projectID uint32, resourceID uint64) (*[]model.ResourceTag, error) {
	var data []model.ResourceTag
	if err := c.Master.WithContext(ctx).Raw(selectListResourceTagByResourceID, projectID, resourceID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListResourceTagCount = `select count(*) from resource_tag where project_id = ? and resource_id = ?`

func (c *Client) ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error) {
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(selectListResourceTagCount, param.ProjectId, param.ResourceId).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const selectListResourceTagName = `
select
  tag
from
  resource_tag
where
  project_id = ?
  and updated_at between ? and ?
group by project_id, tag
order by %s %s
limit %d, %d
`

func (c *Client) ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]TagName, error) {
	var data []TagName
	if err := c.Slave.WithContext(ctx).Raw(
		fmt.Sprintf(selectListResourceTagName, param.Sort, param.Direction, param.Offset, param.Limit),
		param.ProjectId, time.Unix(param.FromAt, 0), time.Unix(param.ToAt, 0)).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListResourceTagNameCount = `
select count(*) from (
  select tag
  from resource_tag
  where project_id = ? and updated_at between ? and ?
  group by project_id, tag
) tag
`

func (c *Client) ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error) {
	var count int64
	if err := c.Slave.WithContext(ctx).Raw(selectListResourceTagNameCount,
		param.ProjectId, time.Unix(param.FromAt, 0), time.Unix(param.ToAt, 0)).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

const insertUpsertResource = `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
	project_id=VALUES(project_id),
	updated_at=NOW()
`

func (c *Client) UpsertResource(ctx context.Context, data *model.Resource) (*model.Resource, error) {
	operation := func() (*model.Resource, error) {
		return c.upsertResource(ctx, data)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "UpsertResource"))
}

func (c *Client) upsertResource(ctx context.Context, data *model.Resource) (*model.Resource, error) {
	if err := c.Master.WithContext(ctx).Exec(insertUpsertResource,
		data.ResourceID, data.ResourceName, data.ProjectID).Error; err != nil {
		return nil, err
	}
	return c.GetResourceByName(ctx, data.ProjectID, data.ResourceName)
}

const selectGetResourceByName = `select * from resource where project_id = ? and resource_name = ?`

func (c *Client) GetResourceByName(ctx context.Context, projectID uint32, resourceName string) (*model.Resource, error) {
	var data model.Resource
	if err := c.Master.WithContext(ctx).Raw(selectGetResourceByName, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResourceTagByKey = `select * from resource_tag where project_id = ? and resource_id = ? and tag = ?`

func (c *Client) GetResourceTagByKey(ctx context.Context, projectID uint32, resourceID uint64, tag string) (*model.ResourceTag, error) {
	var data model.ResourceTag
	if err := c.Master.WithContext(ctx).Raw(selectGetResourceTagByKey, projectID, resourceID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertTagResource = `
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)
`

func (c *Client) TagResource(ctx context.Context, tag *model.ResourceTag) (*model.ResourceTag, error) {
	operation := func() (*model.ResourceTag, error) {
		return c.tagResource(ctx, tag)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "TagResource"))
}

func (c *Client) tagResource(ctx context.Context, tag *model.ResourceTag) (*model.ResourceTag, error) {
	if err := c.Master.WithContext(ctx).Exec(insertTagResource,
		tag.ResourceTagID, tag.ResourceID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return c.GetResourceTagByKey(ctx, tag.ProjectID, tag.ResourceID, tag.Tag)
}

const deleteUntagResource = `delete from resource_tag where project_id = ? and resource_tag_id = ?`

func (c *Client) UntagResource(ctx context.Context, projectID uint32, resourceTagID uint64) error {
	return c.Master.WithContext(ctx).Exec(deleteUntagResource, projectID, resourceTagID).Error
}

func (c *Client) BulkUpsertResource(ctx context.Context, data []*model.Resource) error {
	operation := func() error {
		return c.bulkUpsertResource(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertResource"))
}

func (c *Client) bulkUpsertResource(ctx context.Context, data []*model.Resource) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertResourceSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertResourceSQL(data []*model.Resource) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?),`
		params = append(params, d.ResourceID, d.ResourceName, d.ProjectID)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id),
  updated_at=NOW()`
	return sql, params
}

func (c *Client) BulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error {
	operation := func() error {
		return c.bulkUpsertResourceTag(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertResourceTag"))
}

func (c *Client) bulkUpsertResourceTag(ctx context.Context, data []*model.ResourceTag) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertResourceTagSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
}

func generateBulkUpsertResourceTagSQL(data []*model.ResourceTag) (string, []interface{}) {
	var params []interface{}
	sql := `
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES`
	for _, d := range data {
		sql += `
  (?, ?, ?, ?),`
		params = append(params, d.ResourceTagID, d.ResourceID, d.ProjectID, d.Tag)
	}
	sql = strings.TrimRight(sql, ",")
	sql += `
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)`
	return sql, params
}

const selectGetRecommend = `
select r.* 
from recommend r 
  inner join recommend_finding rf using(recommend_id) 
where rf.project_id=? and rf.finding_id=?
`

func (c *Client) GetRecommend(ctx context.Context, projectID uint32, findingID uint64) (*model.Recommend, error) {
	var data model.Recommend
	if err := c.Slave.WithContext(ctx).Raw(selectGetRecommend, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	operation := func() (*model.Recommend, error) {
		return c.upsertRecommend(ctx, data)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "UpsertRecommend"))
}

func (c *Client) upsertRecommend(ctx context.Context, data *model.Recommend) (*model.Recommend, error) {
	var ret model.Recommend
	if err := c.Master.WithContext(ctx).Where("data_source=? AND type=?", data.DataSource, data.Type).Assign(data).FirstOrCreate(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) UpsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	operation := func() (*model.RecommendFinding, error) {
		return c.upsertRecommendFinding(ctx, data)
	}
	return backoff.RetryNotifyWithData(operation, c.retryer, c.newRetryLogger(ctx, "UpsertRecommendFinding"))
}

func (c *Client) upsertRecommendFinding(ctx context.Context, data *model.RecommendFinding) (*model.RecommendFinding, error) {
	var ret model.RecommendFinding
	if err := c.Master.WithContext(ctx).Where("finding_id=?", data.FindingID).Assign(data).FirstOrCreate(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

const selectGetRecommendByDataSourceType = `
select r.* 
from recommend r 
where r.data_source=? and r.type=?
`

func (c *Client) GetRecommendByDataSourceType(ctx context.Context, dataSource, recommendType string) (*model.Recommend, error) {
	var data model.Recommend
	if err := c.Master.WithContext(ctx).Raw(selectGetRecommendByDataSourceType, dataSource, recommendType).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) BulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	operation := func() error {
		return c.bulkUpsertRecommend(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertRecommend"))
}

func (c *Client) bulkUpsertRecommend(ctx context.Context, data []*model.Recommend) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertRecommendSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
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

func (c *Client) BulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	operation := func() error {
		return c.bulkUpsertRecommendFinding(ctx, data)
	}
	return backoff.RetryNotify(operation, c.retryer, c.newRetryLogger(ctx, "BulkUpsertRecommendFinding"))
}

func (c *Client) bulkUpsertRecommendFinding(ctx context.Context, data []*model.RecommendFinding) error {
	if len(data) == 0 {
		return nil
	}
	sql, params := generateBulkUpsertRecommendFindingSQL(data)
	return c.Master.WithContext(ctx).Exec(sql, params...).Error
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

const selectGetPendFinding = `select * from pend_finding where project_id = ? and finding_id = ?`

func (c *Client) GetPendFinding(ctx context.Context, projectID uint32, findingID uint64) (*model.PendFinding, error) {
	var data model.PendFinding
	if err := c.Master.WithContext(ctx).Raw(selectGetPendFinding, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPendFinding = `
INSERT INTO pend_finding
  (finding_id, project_id, note)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  updated_at = CURRENT_TIMESTAMP()
`

func (c *Client) UpsertPendFinding(ctx context.Context, pend *finding.PendFindingForUpsert) (*model.PendFinding, error) {
	if err := c.Master.WithContext(ctx).Exec(insertPendFinding, pend.FindingId, pend.ProjectId, pend.Note).Error; err != nil {
		return nil, err
	}
	return c.GetPendFinding(ctx, pend.ProjectId, pend.FindingId)
}

const deletePendFinding = `delete from pend_finding where project_id = ? and finding_id = ?`

func (c *Client) DeletePendFinding(ctx context.Context, projectID uint32, findingID uint64) error {
	return c.Master.WithContext(ctx).Exec(deletePendFinding, projectID, findingID).Error
}

type TagName struct {
	Tag string
}
