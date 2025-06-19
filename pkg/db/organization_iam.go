package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"gorm.io/gorm"
)

type OrganizationIAMRepository interface {
	// OrganizationRole
	ListOrganizationRole(ctx context.Context, organizationID uint32, name string, userID uint32) ([]*model.OrganizationRole, error)
	GetOrganizationRole(ctx context.Context, organizationID, roleID uint32) (*model.OrganizationRole, error)
	GetOrganizationRoleByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationRole, error)
	PutOrganizationRole(ctx context.Context, r *model.OrganizationRole) (*model.OrganizationRole, error)
	DeleteOrganizationRole(ctx context.Context, organizationID, roleID uint32) error
	AttachOrganizationRole(ctx context.Context, organizationID, roleID, userID uint32) (*model.OrganizationRole, error)
	DetachOrganizationRole(ctx context.Context, organizationID, roleID, userID uint32) error

	// OrganizationPolicy
	ListOrganizationPolicy(ctx context.Context, organizationID uint32, name string, roleID uint32) ([]*model.OrganizationPolicy, error)
	GetOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) (*model.OrganizationPolicy, error)
	GetOrganizationPolicyByName(ctx context.Context, organizationID uint32, name string) (*model.OrganizationPolicy, error)
	GetOrganizationPolicyByUserID(ctx context.Context, userID, organizationID uint32) (*[]model.OrganizationPolicy, error)
	PutOrganizationPolicy(ctx context.Context, p *model.OrganizationPolicy) (*model.OrganizationPolicy, error)
	DeleteOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) error
	AttachOrganizationPolicy(ctx context.Context, organizationID, policyID, roleID uint32) (*model.OrganizationPolicy, error)
	DetachOrganizationPolicy(ctx context.Context, organizationID, policyID, roleID uint32) error
}

var _ OrganizationIAMRepository = (*Client)(nil)

const ListOrganizationRole = `select * from organization_role r where 1=1`

func (c *Client) ListOrganizationRole(ctx context.Context, organizationID uint32, name string, userID uint32) ([]*model.OrganizationRole, error) {
	query := ListOrganizationRole
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

const getOrganizationRole = `select * from organization_role r where role_id = ?`

func (c *Client) GetOrganizationRole(ctx context.Context, organizationID, roleID uint32) (*model.OrganizationRole, error) {
	query := getOrganizationRole
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
	insert into organization_role (
		role_id,
		name,
		organization_id
	) values (
		?,
		?,
		?
	) on duplicate key update
		name = values(name),
		organization_id = values(organization_id)
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

const ListOrganizationPolicy = `select * from organization_policy op where op.organization_id = ?`

// OrganizationPolicy
func (c *Client) ListOrganizationPolicy(ctx context.Context, organizationID uint32, name string, roleID uint32) ([]*model.OrganizationPolicy, error) {
	query := ListOrganizationPolicy
	var params []interface{}
	params = append(params, organizationID)
	if name != "" {
		query += " and op.name = ?"
		params = append(params, name)
	}
	if roleID != 0 {
		query += " and exists(select * from organization_role_policy orp where orp.policy_id = op.policy_id and orp.role_id = ?)"
		params = append(params, roleID)
	}
	var data []*model.OrganizationPolicy
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const getOrganizationPolicy = `select * from organization_policy p where policy_id = ?`

func (c *Client) GetOrganizationPolicy(ctx context.Context, organizationID, policyID uint32) (*model.OrganizationPolicy, error) {
	query := getOrganizationPolicy
	var params []interface{}
	params = append(params, policyID)
	if organizationID != 0 {
		query += " and p.organization_id = ?"
		params = append(params, organizationID)
	}
	var data model.OrganizationPolicy
	if err := c.Slave.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
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

const getOrganizationPolicyByUserID = `
select
  op.* 
from
  user u
  inner join user_organization_role uor using(user_id)
  inner join organization_role_policy orp using(role_id)
  inner join organization_policy op using(policy_id) 
where
  u.activated = 'true'
  and u.user_id = ?
  and op.organization_id = ?
`

func (c *Client) GetOrganizationPolicyByUserID(ctx context.Context, userID, organizationID uint32) (*[]model.OrganizationPolicy, error) {
	var data []model.OrganizationPolicy
	if err := c.Slave.WithContext(ctx).Raw(getOrganizationPolicyByUserID, userID, organizationID).Scan(&data).Error; err != nil {
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

const insertAttachOrganizationRole = `
	insert into user_organization_role (
		role_id,
		user_id
	) values (
	 	?, 
		?
	) on duplicate key update
		role_id = values(role_id)`

func (c *Client) AttachOrganizationRole(ctx context.Context, organizationID, roleID, userID uint32) (*model.OrganizationRole, error) {
	// Check if role exists and belongs to the specified organization
	roleExists, err := c.organizationRoleExists(ctx, organizationID, roleID)
	if err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, fmt.Errorf("role not found in organization: organizationID=%d, roleID=%d", organizationID, roleID)
	}

	userExists, err := c.organizationUserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("user not found: userID=%d", userID)
	}

	if err := c.Master.WithContext(ctx).Exec(insertAttachOrganizationRole, roleID, userID).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationRole(ctx, organizationID, roleID)
}

const deleteDetachOrganizationRole = `
	delete from user_organization_role 
	where role_id = ? 
		and user_id = ?
`

func (c *Client) DetachOrganizationRole(ctx context.Context, organizationID, roleID, userID uint32) error {
	roleExists, err := c.organizationRoleExists(ctx, organizationID, roleID)
	if err != nil {
		return err
	}
	if !roleExists {
		return fmt.Errorf("role not found in organization: organizationID=%d, roleID=%d", organizationID, roleID)
	}

	userExists, err := c.organizationUserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !userExists {
		return fmt.Errorf("user not found: userID=%d", userID)
	}

	return c.Master.WithContext(ctx).Exec(deleteDetachOrganizationRole, roleID, userID).Error
}

const insertAttachOrganizationPolicy = `
	insert into organization_role_policy (
		role_id,
		policy_id
	) values (
	 	?, 
		?
	) on duplicate key update
		role_id = values(role_id)`

func (c *Client) AttachOrganizationPolicy(ctx context.Context, organizationID, policyID, roleID uint32) (*model.OrganizationPolicy, error) {
	roleExists, err := c.organizationRoleExists(ctx, organizationID, roleID)
	if err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, fmt.Errorf("role not found in organization: organizationID=%d, roleID=%d", organizationID, roleID)
	}
	policyExists, err := c.organizationPolicyExists(ctx, organizationID, policyID)
	if err != nil {
		return nil, err
	}
	if !policyExists {
		return nil, fmt.Errorf("policy not found in organization: organizationID=%d, policyID=%d", organizationID, policyID)
	}

	if err := c.Master.WithContext(ctx).Exec(insertAttachOrganizationPolicy, roleID, policyID).Error; err != nil {
		return nil, err
	}
	return c.GetOrganizationPolicy(ctx, organizationID, policyID)
}

const deleteDetachOrganizationPolicy = `
	delete from organization_role_policy 
	where role_id = ? 
		and policy_id = ?
`

func (c *Client) DetachOrganizationPolicy(ctx context.Context, organizationID, policyID, roleID uint32) error {
	roleExists, err := c.organizationRoleExists(ctx, organizationID, roleID)
	if err != nil {
		return err
	}
	if !roleExists {
		return fmt.Errorf("role not found in organization: organizationID=%d, roleID=%d", organizationID, roleID)
	}
	policyExists, err := c.organizationPolicyExists(ctx, organizationID, policyID)
	if err != nil {
		return err
	}
	if !policyExists {
		return fmt.Errorf("policy not found in organization: organizationID=%d, policyID=%d", organizationID, policyID)
	}
	return c.Master.WithContext(ctx).Exec(deleteDetachOrganizationPolicy, roleID, policyID).Error
}

func (c *Client) organizationRoleExists(ctx context.Context, organizationID, roleID uint32) (bool, error) {
	if _, err := c.GetOrganizationRole(ctx, organizationID, roleID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get organization role. organization_id=%d, role_id=%d, error: %w", organizationID, roleID, err)
	}
	return true, nil
}

func (c *Client) organizationPolicyExists(ctx context.Context, organizationID, policyID uint32) (bool, error) {
	if _, err := c.GetOrganizationPolicy(ctx, organizationID, policyID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get organization policy. organization_id=%d, policy_id=%d, error: %w", organizationID, policyID, err)
	}
	return true, nil
}

func (c *Client) organizationUserExists(ctx context.Context, userID uint32) (bool, error) {
	if _, err := c.GetUser(ctx, userID, "", ""); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get user. user_id=%d, error: %w", userID, err)
	}
	return true, nil
}
