package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type projectWithTag struct {
	ProjectID uint32
	Name      string
	Tag       *[]model.ProjectTag `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *projectDB) ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]projectWithTag, error) {
	query := `select p.* from project p where 1 = 1` // プログラム構造をシンプルに保つために必ずtrueとなるwhere条件を入れとく（and条件のみなので一旦これで）
	var params []interface{}
	if !zero.IsZeroVal(userID) {
		query += " and exists (select * from user_role ur inner join role r using(project_id, role_id) where ur.project_id = p.project_id and user_id = ?)"
		params = append(params, userID)
	}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	data := []projectWithTag{}
	if err := p.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	for idx, pj := range data {
		tag, err := p.ListProjectTag(ctx, pj.ProjectID)
		if err != nil {
			return nil, err
		}
		data[idx].Tag = tag
	}
	return &data, nil
}

const selectGetProjectByName = `select * from project where name = ?`

func (p *projectDB) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	var data model.Project
	if err := p.Master.WithContext(ctx).Raw(selectGetProjectByName, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertCreateProject = `insert into project(name) values(?)`

func (p *projectDB) CreateProject(ctx context.Context, name string) (*model.Project, error) {
	// Handring duplicated name error
	if pr, err := p.GetProjectByName(ctx, name); err == nil {
		return nil, fmt.Errorf("Project name already registerd: project_id=%d, name=%s", pr.ProjectID, pr.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Could not get project data: err=%+v", err)
	}
	if err := p.Master.WithContext(ctx).Exec(insertCreateProject, name).Error; err != nil {
		return nil, err
	}
	return p.GetProjectByName(ctx, name)
}

const updateUpdateProject = `update project set name = ? where project_id = ?`

func (p *projectDB) UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error) {
	if err := p.Master.WithContext(ctx).Exec(updateUpdateProject, name, projectID).Error; err != nil {
		return nil, err
	}
	return p.GetProjectByName(ctx, name)
}
