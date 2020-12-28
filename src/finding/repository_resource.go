package main

import (
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
)

func (f *findingDB) ListResource(req *finding.ListResourceRequest) (*[]model.Resource, error) {
	query := `
select 
  r.*
from
  resource r
  left outer join finding f using(resource_name)
where
  r.project_id = ?
  and r.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if len(req.ResourceName) != 0 {
		query += " and r.resource_name regexp ?"
		params = append(params, strings.Join(req.ResourceName, "|"))
	}
	if len(req.Tag) != 0 {
		query += " and exists (select * from resource_tag rt where rt.resource_id=r.resource_id and rt.tag in (?) )"
		params = append(params, req.Tag)
	}
	query += " group by r.resource_id having sum(COALESCE(f.score, 0)) between ? and ?"
	params = append(params, req.FromSumScore, req.ToSumScore)

	var data []model.Resource
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResource = `select * from resource where project_id = ? and resource_id = ?`

func (f *findingDB) GetResource(projectID uint32, resourceID uint64) (*model.Resource, error) {
	var data model.Resource
	if err := f.Slave.Raw(selectGetResource, projectID, resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDeleteResource = `delete from resource where project_id = ? and resource_id = ?`

func (f *findingDB) DeleteResource(projectID uint32, resourceID uint64) error {
	if err := f.Master.Exec(deleteDeleteResource, projectID, resourceID).Error; err != nil {
		return err
	}
	return f.DeleteTagByResourceID(projectID, resourceID)
}

const deleteDeleteTagByResourceID = `delete from resource_tag where project_id = ? and resource_id = ?`

func (f *findingDB) DeleteTagByResourceID(projectID uint32, resourceID uint64) error {
	return f.Master.Exec(deleteDeleteTagByResourceID, projectID, resourceID).Error
}

const selectListResourceTag = `select * from resource_tag where project_id = ? and resource_id = ?`

func (f *findingDB) ListResourceTag(projectID uint32, resourceID uint64) (*[]model.ResourceTag, error) {
	var data []model.ResourceTag
	if err := f.Slave.Raw(selectListResourceTag, projectID, resourceID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingDB) ListResourceTagName(req *finding.ListResourceTagNameRequest) (*[]tagName, error) {
	query := `
select
  distinct tag
from
  resource_tag
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
	if err := f.Master.Raw(selectGetResourceByName, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResourceTagByKey = `select * from resource_tag where project_id = ? and resource_id = ? and tag = ?`

func (f *findingDB) GetResourceTagByKey(projectID uint32, resourceID uint64, tag string) (*model.ResourceTag, error) {
	var data model.ResourceTag
	if err := f.Master.Raw(selectGetResourceTagByKey, projectID, resourceID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResourceTagByID = `select * from resource_tag where project_id = ? and resource_tag_id = ?`

func (f *findingDB) GetResourceTagByID(projectID uint32, resourceID uint64) (*model.ResourceTag, error) {
	var data model.ResourceTag
	if err := f.Slave.Raw(selectGetResourceTagByID, projectID, resourceID).First(&data).Error; err != nil {
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

func (f *findingDB) TagResource(tag *model.ResourceTag) (*model.ResourceTag, error) {
	if err := f.Master.Exec(insertTagResource,
		tag.ResourceTagID, tag.ResourceID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return f.GetResourceTagByKey(tag.ProjectID, tag.ResourceID, tag.Tag)
}

const deleteUntagResource = `delete from resource_tag where project_id = ? and resource_tag_id = ?`

func (f *findingDB) UntagResource(projectID uint32, resourceTagID uint64) error {
	return f.Master.Exec(deleteUntagResource, projectID, resourceTagID).Error
}
