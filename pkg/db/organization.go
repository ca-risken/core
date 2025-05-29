package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	// Organization
	ListOrganization(ctx context.Context, organizationID uint32, name string) ([]*model.Organization, error)
	CreateOrganization(ctx context.Context, name, description string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID uint32, name, description string) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID uint32) error

	// OrganizationProject
	CreateOrganizationProject(ctx context.Context, organizationID, projectID uint32) (*model.OrganizationProject, error)
	ListProjectsInOrganization(ctx context.Context, organizationID uint32) ([]*model.Project, error)
	RemoveProjectsInOrganization(ctx context.Context, organizationID, projectID uint32) error
}

var _ OrganizationRepository = (*Client)(nil)

func (c *Client) ListOrganization(ctx context.Context, organizationID uint32, name string) ([]*model.Organization, error) {
	query := `select * from organization o where 1 = 1`
	var params []interface{}
	if organizationID != 0 {
		query += " and o.organization_id = ?"
		params = append(params, organizationID)
	}
	if name != "" {
		query += " and o.name = ?"
		params = append(params, name)
	}
	query += " order by o.organization_id"
	var organizations []*model.Organization
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

const selectGetOrganizationByName = `select * from organization where name = ?`

func (c *Client) GetOrganizationByName(ctx context.Context, name string) (*model.Organization, error) {
	var data model.Organization
	if err := c.Master.WithContext(ctx).Raw(selectGetOrganizationByName, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertCreateOrganization = `insert into organization(name, description) values(?, ?)`

func (c *Client) CreateOrganization(ctx context.Context, name, description string) (*model.Organization, error) {
	// Handling duplicated name error
	if org, err := c.GetOrganizationByName(ctx, name); err == nil {
		return nil, fmt.Errorf("organization name already registered: organization_id=%d, name=%s", org.OrganizationID, org.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("could not get organization data: err=%+v", err)
	}
	if err := c.Master.WithContext(ctx).Exec(insertCreateOrganization, name, description).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationByName(ctx, name)
}

const updateUpdateOrganization = `update organization set name = ?, description = ? where organization_id = ?`

func (c *Client) UpdateOrganization(ctx context.Context, organizationID uint32, name, description string) (*model.Organization, error) {
	if err := c.Master.WithContext(ctx).Exec(updateUpdateOrganization, name, description, organizationID).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationByName(ctx, name)
}

const deleteOrganization = `delete from organization where organization_id = ?`

func (c *Client) DeleteOrganization(ctx context.Context, organizationID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganization, organizationID).Error
}

const createOrganizationProject = `
	INSERT INTO organization_project (
		organization_id,
		project_id
	) VALUES (
		?,
		?
	)
	ON DUPLICATE KEY UPDATE updated_at = NOW()
`

func (c *Client) CreateOrganizationProject(ctx context.Context, organizationID, projectID uint32) (*model.OrganizationProject, error) {
	if err := c.Master.WithContext(ctx).Exec(createOrganizationProject, organizationID, projectID).Error; err != nil {
		return nil, err
	}
	var data model.OrganizationProject
	if err := c.Master.WithContext(ctx).Raw("SELECT * FROM organization_project WHERE organization_id = ? AND project_id = ?", organizationID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const listProjectsInOrganization = `
	SELECT p.*
	FROM project p
	JOIN organization_project op ON p.project_id = op.project_id
	WHERE op.organization_id = ?
`

func (c *Client) ListProjectsInOrganization(ctx context.Context, organizationID uint32) ([]*model.Project, error) {
	var projects []*model.Project
	if err := c.Slave.WithContext(ctx).Raw(listProjectsInOrganization, organizationID).Scan(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

const removeProjectsInOrganization = `
    delete from organization_project 
    where organization_id = ? 
    and project_id = ?
`

func (c *Client) RemoveProjectsInOrganization(ctx context.Context, organizationID, projectID uint32) error {
	return c.Master.WithContext(ctx).Exec(removeProjectsInOrganization, organizationID, projectID).Error
}
