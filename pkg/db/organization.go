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
	ListOrganization(ctx context.Context, organizationID uint32, name string) (*[]model.Organization, error)
	GetOrganization(ctx context.Context, organizationID uint32) (*model.Organization, error)
	CreateOrganization(ctx context.Context, name, description string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID uint32, name, description string) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID uint32) error
}

var _ OrganizationRepository = (*Client)(nil)

const selectListOrganization = `
select
  o.*
from
  organization o
where
  (? = 0 or o.organization_id = ?)
  and (? = '' or o.name like ?)
order by
  o.organization_id`

func (c *Client) ListOrganization(ctx context.Context, organizationID uint32, name string) (*[]model.Organization, error) {
	var organizations []model.Organization
	if err := c.Slave.WithContext(ctx).Raw(selectListOrganization,
		organizationID, organizationID,
		name, fmt.Sprintf("%%%s%%", name),
	).Scan(&organizations).Error; err != nil {
		return nil, err
	}
	return &organizations, nil
}

const selectGetOrganization = `select * from organization where organization_id = ?`

func (c *Client) GetOrganization(ctx context.Context, organizationID uint32) (*model.Organization, error) {
	var data model.Organization
	if err := c.Master.WithContext(ctx).Raw(selectGetOrganization, organizationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
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
