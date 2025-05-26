package db

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
)

type OrganizationIAMRepository interface {

	// OrganizationRole
	ListOrganizationRole(ctx context.Context, organizationID uint32, name string, userID uint32) ([]*model.OrganizationRole, error)
	GetOrganizationRole(ctx context.Context, organizationID, roleID uint32) (*model.OrganizationRole, error)
	GetOrganizationRoleByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationRole, error)
	PutOrganizationRole(ctx context.Context, r *model.OrganizationRole) (*model.OrganizationRole, error)
	DeleteOrganizationRole(ctx context.Context, organizationID, roleID uint32) error

	// OrganizationPolicy
	ListOrganizationPolicy(ctx context.Context, organizationID uint32, name string, roleID uint32) ([]*model.OrganizationPolicy, error)
	GetOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) (*model.OrganizationPolicy, error)
	GetOrganizationPolicyByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationPolicy, error)
	PutOrganizationPolicy(ctx context.Context, p *model.OrganizationPolicy) (*model.OrganizationPolicy, error)
	DeleteOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) error
}

var _ OrganizationIAMRepository = (*Client)(nil)

func (c *Client) ListOrganizationRole(ctx context.Context, organizationID uint32, name string, userID uint32) ([]*model.OrganizationRole, error) {
	query := `select * from organization_role or where 1=1`
	var params []interface{}
	if organizationID != 0 {
		query += " and r.organization_id = ?"
		params = append(params, organizationID)
	}
	if name != "" {
		query += " and r.name = ?"
		params = append(params, name)
	}
	if userID != 0 {
		query += " and exists (select * from user_organization_role uor where uor.role_id = r.role_id and uor.user_id = ? )"
		params = append(params, userID)
	}
	var data []*model.OrganizationRole
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetOrganizationRole(ctx context.Context, organizationID, roleID uint32) (*model.OrganizationRole, error) {
	query := `select * from organization_role r where role_id =?`
	var params []interface{}
	params = append(params, roleID)
	if organizationID != 0 {
		query += " and r.organization_id = ?"
		params = append(params, organizationID)
	}
	var data model.OrganizationRole
	if err := c.Slave.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const getOrganizationRoleByName = `
	select * 
	from organization_role 
	where organization_id = ? 
		and name = ?
`

func (c *Client) GetOrganizationRoleByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationRole, error) {
	var data model.OrganizationRole
	if err := c.Master.WithContext(ctx).Raw(getOrganizationRoleByName, organizationID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const putOrganizationRole = `
	INSERT INTO organization_role (
		role_id,
		name,
		organization_id
	) VALUES (
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		organization_id = VALUES(organization_id)
`

func (c *Client) PutOrganizationRole(ctx context.Context, r *model.OrganizationRole) (*model.OrganizationRole, error) {
	if err := c.Master.WithContext(ctx).Exec(putOrganizationRole, r.RoleID, r.Name, r.OrganizationID).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationRoleByName(ctx, r.OrganizationID, r.Name)
}

const deleteOrganizationRole = `
	delete from organization_role 
	where organization_id = ? 
		and role_id = ?
`

func (c *Client) DeleteOrganizationRole(ctx context.Context, organizationID, roleID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganizationRole, organizationID, roleID).Error
}

// OrganizationPolicy
func (c *Client) ListOrganizationPolicy(ctx context.Context, organizationID uint32, name string, roleID uint32) ([]*model.OrganizationPolicy, error) {
	query := `select * from policy p where p.organization_id = ?`
	var params []interface{}
	params = append(params, organizationID)
	if name != "" {
		query += " and p.name = ?"
		params = append(params, name)
	}
	if roleID != 0 {
		query += " and exists(select * from organization_role_policy orp where orp.policy_id = p.policy_id and orp.role_id = ?)"
		params = append(params, roleID)
	}
	var data []*model.OrganizationPolicy
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const getOrganizationPolicy = `
	select * 
	from organization_policy 
	where organization_id = ? 
		and policy_id = ?
`

func (c *Client) GetOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) (*model.OrganizationPolicy, error) {
	var data model.OrganizationPolicy
	if err := c.Slave.WithContext(ctx).Raw(getOrganizationPolicy, organizationID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const getOrganizationPolicyByName = `
	select * 
	from organization_policy 
	where organization_id = ? 
		and name = ?
`

func (c *Client) GetOrganizationPolicyByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationPolicy, error) {
	var data model.OrganizationPolicy
	if err := c.Master.WithContext(ctx).Raw(getOrganizationPolicyByName, organizationID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const putOrganizationPolicy = `
	INSERT INTO organization_policy (
		policy_id,
		name,
		organization_id,
		action_ptn
	) VALUES (
		?,
		?,
		?,
		?
	) ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		organization_id = VALUES(organization_id),
		action_ptn = VALUES(action_ptn)
`

func (c *Client) PutOrganizationPolicy(ctx context.Context, p *model.OrganizationPolicy) (*model.OrganizationPolicy, error) {
	if err := c.Master.WithContext(ctx).Exec(putOrganizationPolicy, p.PolicyID, p.Name, p.OrganizationID, p.ActionPtn).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationPolicyByName(ctx, p.OrganizationID, p.Name)
}

const deleteOrganizationPolicy = `
	delete from organization_policy 
	where organization_id = ? 
		and policy_id = ?
`

func (c *Client) DeleteOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganizationPolicy, organizationID, policyID).Error
}
