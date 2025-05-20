package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	// Organization
	ListOrganization(ctx context.Context, organizationID uint32, name string) (*[]model.Organization, error)
	CreateOrganization(ctx context.Context, name, description string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, organizationID uint32, name, description string) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, organizationID uint32) error
}

var _ OrganizationRepository = (*Client)(nil)

func (c *Client) ListOrganization(ctx context.Context, organizationID uint32, name string) (*[]model.Organization, error) {
	query := `select * from organization o where 1 = 1`
	var params []interface{}
	if !zero.IsZeroVal(organizationID) {
		query += " and o.organization_id = ?"
		params = append(params, organizationID)
	}
	if !zero.IsZeroVal(name) {
		query += " and o.name = ?"
		params = append(params, name)
	}
	query += " order by o.organization_id"
	var organizations []model.Organization
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&organizations).Error; err != nil {
		return nil, err
	}
	return &organizations, nil
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
