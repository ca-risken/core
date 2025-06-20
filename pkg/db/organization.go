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
	ListOrganization(ctx context.Context, organizationID uint32, name string, userID, projectID uint32) ([]*model.Organization, error)
	CreateOrganization(ctx context.Context, name, description string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID uint32, name, description string) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID uint32) error

	// OrganizationProject
	PutOrganizationProject(ctx context.Context, organizationID, projectID uint32) (*model.OrganizationProject, error)
	ListProjectsInOrganization(ctx context.Context, organizationID uint32) ([]*model.Project, error)
	RemoveProjectsInOrganization(ctx context.Context, organizationID, projectID uint32) error

	// OrganizationInvitation
	ListOrganizationInvitation(ctx context.Context, organizationID, projectID uint32) ([]*model.OrganizationInvitation, error)
	PutOrganizationInvitation(ctx context.Context, organizationID, projectID uint32, status string) (*model.OrganizationInvitation, error)
	DeleteOrganizationInvitation(ctx context.Context, organizationID, projectID uint32) error
}

var _ OrganizationRepository = (*Client)(nil)

func (c *Client) ListOrganization(ctx context.Context, organizationID uint32, name string, userID, projectID uint32) ([]*model.Organization, error) {
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
	if userID != 0 {
		query += " and exists (select 1 from user_organization_role ur inner join organization_role r on ur.role_id = r.role_id where r.organization_id = o.organization_id and ur.user_id = ?)"
		params = append(params, userID)
	}
	if projectID != 0 {
		query += " and exists (select 1 from organization_project op where op.organization_id = o.organization_id and op.project_id = ?)"
		params = append(params, projectID)
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

const (
	putOrganizationProject = `
		insert into organization_project (
			organization_id,
			project_id
		) values (
			?,
			?
		)
		on duplicate key update
			updated_at=NOW()
	`
	selectGetOrganizationProject = `
		select *
		from organization_project
		where organization_id = ? 
		  and project_id = ?
	`
)

func (c *Client) PutOrganizationProject(ctx context.Context, organizationID, projectID uint32) (*model.OrganizationProject, error) {
	if err := c.Master.WithContext(ctx).Exec(putOrganizationProject, organizationID, projectID).Error; err != nil {
		return nil, err
	}
	var data model.OrganizationProject
	if err := c.Master.WithContext(ctx).Raw(selectGetOrganizationProject, organizationID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const listProjectsInOrganization = `
	select p.*
	from project p
	join organization_project op on p.project_id = op.project_id
	where op.organization_id = ?
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

func (c *Client) ListOrganizationInvitation(ctx context.Context, organizationID, projectID uint32) ([]*model.OrganizationInvitation, error) {
	query := `select * from organization_invitation oi where 1=1`
	var params []interface{}
	if organizationID != 0 {
		query += " and oi.organization_id = ?"
		params = append(params, organizationID)
	}
	if projectID != 0 {
		query += " and oi.project_id = ?"
		params = append(params, projectID)
	}
	var invitations []*model.OrganizationInvitation
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

const selectGetOrganizationInvitation = `
	select *
	from organization_invitation
	where organization_id = ?
	and project_id = ?
`

func (c *Client) getOrganizationInvitation(ctx context.Context, organizationID, projectID uint32) (*model.OrganizationInvitation, error) {
	var invitation model.OrganizationInvitation
	if err := c.Master.WithContext(ctx).Raw(selectGetOrganizationInvitation, organizationID, projectID).First(&invitation).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

const putOrganizationInvitation = `
	INSERT INTO organization_invitation (
		organization_id,
		project_id,
		status
	) VALUES (
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		status = VALUES(status)
`

func (c *Client) PutOrganizationInvitation(ctx context.Context, organizationID, projectID uint32, status string) (*model.OrganizationInvitation, error) {
	if err := c.Master.WithContext(ctx).Exec(putOrganizationInvitation, organizationID, projectID, status).Error; err != nil {
		return nil, err
	}
	return c.getOrganizationInvitation(ctx, organizationID, projectID)
}

const deleteOrganizationInvitation = `
    delete from organization_invitation 
    where organization_id = ? 
    and project_id = ?
`

func (c *Client) DeleteOrganizationInvitation(ctx context.Context, organizationID, projectID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganizationInvitation, organizationID, projectID).Error
}
