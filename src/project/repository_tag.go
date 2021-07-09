package main

import (
	"github.com/CyberAgent/mimosa-core/pkg/model"
)

const selectListProjectTag string = `select * from project_tag where project_id=? order by tag`

func (p *projectDB) ListProjectTag(projectID uint32) (*[]model.ProjectTag, error) {
	var data []model.ProjectTag
	if err := p.Slave.Raw(selectListProjectTag, projectID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetProjectTag string = `select * from project_tag where project_id=? and tag=?`

func (p *projectDB) GetProjectTag(projectID uint32, tag string) (*model.ProjectTag, error) {
	var data model.ProjectTag
	if err := p.Master.Raw(selectGetProjectTag, projectID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertTagProject string = `
insert into project_tag
  (project_id, tag, color)
values
  (?, ?, ?)
on duplicate key update
  color=VALUES(color),
  updated_at=NOW()`

func (p *projectDB) TagProject(projectID uint32, tag, color string) (*model.ProjectTag, error) {
	if err := p.Master.Exec(insertTagProject, projectID, tag, color).Error; err != nil {
		return nil, err
	}
	return p.GetProjectTag(projectID, tag)
}

const deleteUntagProject string = `delete from project_tag where project_id=? and tag=?`

func (p *projectDB) UntagProject(projectID uint32, tag string) error {
	if err := p.Master.Exec(deleteUntagProject, projectID, tag).Error; err != nil {
		return err
	}
	return nil
}
