package main

import (
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
)

func (f *findingDB) ListFinding(req *finding.ListFindingRequest) (*[]model.Finding, error) {
	query := `
select
  *
from
  finding f
where
  f.project_id = ?
  and f.score between ? and ?
  and f.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, req.FromScore, req.ToScore, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if len(req.DataSource) != 0 {
		query += " and f.data_source in (?)"
		params = append(params, req.DataSource)
	}
	if len(req.ResourceName) != 0 {
		query += " and f.resource_name in (?)"
		params = append(params, req.ResourceName)
	}
	if len(req.Tag) != 0 {
		query += " and exists (select * from finding_tag ft where ft.finding_id=f.finding_id and ft.tag in (?))"
		params = append(params, req.Tag)
	}
	var data []model.Finding
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFinding = `select * from finding where project_id = ? and finding_id = ?`

func (f *findingDB) GetFinding(projectID uint32, findingID uint64) (*model.Finding, error) {
	var data model.Finding
	if err := f.Slave.Raw(selectGetFinding, projectID, findingID).First(&data).Error; err != nil {
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
  data=VALUES(data)
`

func (f *findingDB) UpsertFinding(data *model.Finding) (*model.Finding, error) {
	if err := f.Master.Exec(insertUpsertFinding,
		data.FindingID, data.Description, data.DataSource, data.DataSourceID, data.ResourceName,
		data.ProjectID, data.OriginalScore, data.Score, data.Data).Error; err != nil {
		return nil, err
	}
	return f.GetFindingByDataSource(data.ProjectID, data.DataSource, data.DataSourceID)
}

const selectGetFindingByDataSource = `select * from finding where project_id = ? and data_source = ? and data_source_id = ?`

func (f *findingDB) GetFindingByDataSource(projectID uint32, dataSource, dataSourceID string) (*model.Finding, error) {
	var result model.Finding
	if err := f.Slave.Raw(selectGetFindingByDataSource,
		projectID, dataSource, dataSourceID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

const insertUpsertResource = `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id);
`

func (f *findingDB) UpsertResource(data *model.Resource) (*model.Resource, error) {
	if err := f.Master.Exec(insertUpsertResource,
		data.ResourceID, data.ResourceName, data.ProjectID).Error; err != nil {
		return nil, err
	}
	return f.GetResourceByName(data.ProjectID, data.ResourceName)
}

const selectGetResourceByName = `select * from resource where project_id = ? and resource_name = ?`

func (f *findingDB) GetResourceByName(projectID uint32, resourceName string) (*model.Resource, error) {
	var data model.Resource
	if err := f.Slave.Raw(selectGetResourceByName, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDeleteFinding = `delete from finding where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteFinding(projectID uint32, findingID uint64) error {
	if err := f.Master.Exec(deleteDeleteFinding, projectID, findingID).Error; err != nil {
		return err
	}
	return f.DeleteTagByFindingID(projectID, findingID)
}

const deleteDeleteTagByFindingID = `delete from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteTagByFindingID(projectID uint32, findingID uint64) error {
	return f.Master.Exec(deleteDeleteTagByFindingID, projectID, findingID).Error
}

const selectListFindingTag = `select * from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) ListFindingTag(projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := f.Slave.Raw(selectListFindingTag, projectID, findingID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) ListFindingTagName(req *finding.ListFindingTagNameRequest) (*[]tagName, error) {
	query := `
select
  distinct tag
from
  finding_tag
where
  project_id = ?
  and updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	var data []tagName
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertTagFinding = `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)
`

func (f *findingDB) TagFinding(tag *model.FindingTag) (*model.FindingTag, error) {
	if err := f.Master.Exec(insertTagFinding,
		tag.FindingTagID, tag.FindingID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return f.GetFindingTagByKey(tag.ProjectID, tag.FindingID, tag.Tag)
}

const selectGetFindingTagByKey = `select * from finding_tag where project_id = ? and finding_id = ? and tag = ?`

func (f *findingDB) GetFindingTagByKey(projectID uint32, findingID uint64, tag string) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.Slave.Raw(selectGetFindingTagByKey, projectID, findingID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingTagByID = `select * from finding_tag where project_id = ? and finding_tag_id = ?`

func (f *findingDB) GetFindingTagByID(projectID uint32, findingTagID uint64) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.Slave.Raw(selectGetFindingTagByID, projectID, findingTagID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteUntagFinding = `delete from finding_tag where project_id = ? and finding_tag_id = ?`

func (f *findingDB) UntagFinding(projectID uint32, findingTagID uint64) error {
	return f.Master.Exec(deleteUntagFinding, projectID, findingTagID).Error
}
