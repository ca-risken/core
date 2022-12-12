package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]ProjectWithTag, error)
	CreateProject(ctx context.Context, name string) (*model.Project, error)
	UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error)
	DeleteProject(ctx context.Context, projectID uint32) error

	TagProject(ctx context.Context, projectID uint32, tag, color string) (*model.ProjectTag, error)
	UntagProject(ctx context.Context, projectID uint32, tag string) error
}

var _ ProjectRepository = (*Client)(nil)

type ProjectWithTag struct {
	ProjectID uint32
	Name      string
	Tag       *[]model.ProjectTag
	CreatedAt time.Time
	UpdatedAt time.Time
}

type projectTagDenormarize struct {
	ProjectID uint32
	Name      string
	Tag       string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Client) ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]ProjectWithTag, error) {
	query := `
select p.project_id, p.name, pt.tag, pt.color, p.created_at, p.updated_at 
from project p left outer join project_tag pt using(project_id) 
where 1 = 1 `
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
	query += " order by p.project_id, pt.tag"
	denormarize := []projectTagDenormarize{}
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&denormarize).Error; err != nil {
		return nil, err
	}
	normarize := []ProjectWithTag{}
	pjMap := make(map[uint32]int) // key: project_id, value: index number
	for _, pj := range denormarize {
		if idx, ok := pjMap[pj.ProjectID]; ok {
			tags := *normarize[idx].Tag
			tags = append(tags, model.ProjectTag{
				ProjectID: pj.ProjectID,
				Tag:       pj.Tag,
				Color:     pj.Color,
			})
			normarize[idx].Tag = &tags
			continue
		}
		// new project data
		data := ProjectWithTag{
			ProjectID: pj.ProjectID,
			Name:      pj.Name,
			CreatedAt: pj.CreatedAt,
			UpdatedAt: pj.UpdatedAt,
		}
		if pj.Tag != "" {
			data.Tag = &[]model.ProjectTag{
				{ProjectID: pj.ProjectID, Tag: pj.Tag, Color: pj.Color},
			}
		}
		normarize = append(normarize, data)
		pjMap[pj.ProjectID] = len(normarize) - 1
	}
	return &normarize, nil
}

const selectGetProjectByName = `select * from project where name = ?`

func (c *Client) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	var data model.Project
	if err := c.Master.WithContext(ctx).Raw(selectGetProjectByName, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertCreateProject = `insert into project(name) values(?)`

func (c *Client) CreateProject(ctx context.Context, name string) (*model.Project, error) {
	// Handring duplicated name error
	if pr, err := c.GetProjectByName(ctx, name); err == nil {
		return nil, fmt.Errorf("Project name already registerd: project_id=%d, name=%s", pr.ProjectID, pr.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Could not get project data: err=%+v", err)
	}
	if err := c.Master.WithContext(ctx).Exec(insertCreateProject, name).Error; err != nil {
		return nil, err
	}
	return c.GetProjectByName(ctx, name)
}

const updateUpdateProject = `update project set name = ? where project_id = ?`

func (c *Client) UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error) {
	if err := c.Master.WithContext(ctx).Exec(updateUpdateProject, name, projectID).Error; err != nil {
		return nil, err
	}
	return c.GetProjectByName(ctx, name)
}

const selectGetProjectTag string = `select * from project_tag where project_id=? and tag=?`

func (c *Client) GetProjectTag(ctx context.Context, projectID uint32, tag string) (*model.ProjectTag, error) {
	var data model.ProjectTag
	if err := c.Master.WithContext(ctx).Raw(selectGetProjectTag, projectID, tag).First(&data).Error; err != nil {
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

func (c *Client) TagProject(ctx context.Context, projectID uint32, tag, color string) (*model.ProjectTag, error) {
	if err := c.Master.WithContext(ctx).Exec(insertTagProject, projectID, tag, color).Error; err != nil {
		return nil, err
	}
	return c.GetProjectTag(ctx, projectID, tag)
}

const deleteUntagProject string = `delete from project_tag where project_id=? and tag=?`

func (c *Client) UntagProject(ctx context.Context, projectID uint32, tag string) error {
	if err := c.Master.WithContext(ctx).Exec(deleteUntagProject, projectID, tag).Error; err != nil {
		return err
	}
	return nil
}

const deleteProject = `delete from project where project_id=?`

func (c *Client) DeleteProject(ctx context.Context, projectID uint32) error {
	if err := c.Master.WithContext(ctx).Exec(deleteProject, projectID).Error; err != nil {
		return err
	}
	return nil
}
