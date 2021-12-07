package main

import (
	"context"

	"github.com/ca-risken/core/src/project/model"
)

const selectGetProjectTag string = `select * from project_tag where project_id=? and tag=?`

func (p *projectDB) GetProjectTag(ctx context.Context, projectID uint32, tag string) (*model.ProjectTag, error) {
	var data model.ProjectTag
	if err := p.Master.WithContext(ctx).Raw(selectGetProjectTag, projectID, tag).First(&data).Error; err != nil {
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

func (p *projectDB) TagProject(ctx context.Context, projectID uint32, tag, color string) (*model.ProjectTag, error) {
	if err := p.Master.WithContext(ctx).Exec(insertTagProject, projectID, tag, color).Error; err != nil {
		return nil, err
	}
	return p.GetProjectTag(ctx, projectID, tag)
}

const deleteUntagProject string = `delete from project_tag where project_id=? and tag=?`

func (p *projectDB) UntagProject(ctx context.Context, projectID uint32, tag string) error {
	if err := p.Master.WithContext(ctx).Exec(deleteUntagProject, projectID, tag).Error; err != nil {
		return err
	}
	return nil
}
