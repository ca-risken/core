package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
)

const selectGetUserPolicy = `
select
  p.* 
from
  user u
  inner join user_role ur using(user_id)
  inner join role_policy rp using(role_id)
  inner join policy p using(policy_id) 
where
  u.activated = 'true'
  and u.user_id = ?
`

func (i *iamDB) GetUserPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := i.Slave.WithContext(ctx).Raw(selectGetUserPolicy, userID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetTokenPolicy = `
select
  p.* 
from
  access_token at
  inner join access_token_role ur using(access_token_id)
  inner join role_policy rp using(role_id)
  inner join policy p using(policy_id) 
where
  at.expired_at >= NOW()
  and at.access_token_id = ?
`

func (i *iamDB) GetTokenPolicy(ctx context.Context, accessTokenID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := i.Slave.WithContext(ctx).Raw(selectGetTokenPolicy, accessTokenID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAdminPolicy = `
select
  p.* 
from
  user u
	inner join user_role ur using(user_id)
	inner join role r using(role_id)
  inner join role_policy rp using(role_id)
  inner join policy p using(policy_id) 
where
	u.activated = 'true'
	and r.project_id is null
	and p.project_id is null
  and u.user_id = ?
`

func (i *iamDB) GetAdminPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := i.Slave.WithContext(ctx).Raw(selectGetAdminPolicy, userID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (i *iamDB) ListPolicy(ctx context.Context, projectID uint32, name string, roleID uint32) (*[]model.Policy, error) {
	query := `select * from policy p where p.project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and p.name = ?"
		params = append(params, name)
	}
	if !zero.IsZeroVal(roleID) {
		query += " and exists(select * from role_policy rp where rp.policy_id = p.policy_id and rp.role_id = ?)"
		params = append(params, roleID)
	}
	var data []model.Policy
	if err := i.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicy = `select * from policy where project_id = ? and policy_id =?`

func (i *iamDB) GetPolicy(ctx context.Context, projectID, policyID uint32) (*model.Policy, error) {
	var data model.Policy
	if err := i.Master.WithContext(ctx).Raw(selectGetPolicy, projectID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicyByName = `select * from policy where project_id = ? and name =?`

func (i *iamDB) GetPolicyByName(ctx context.Context, projectID uint32, name string) (*model.Policy, error) {
	var data model.Policy
	if err := i.Master.WithContext(ctx).Raw(selectGetPolicyByName, projectID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPutPolicy = `
INSERT INTO policy
  (policy_id, name, project_id, action_ptn, resource_ptn)
VALUES
  (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  project_id=VALUES(project_id),
  action_ptn=VALUES(action_ptn),
  resource_ptn=VALUES(resource_ptn)
`

func (i *iamDB) PutPolicy(ctx context.Context, p *model.Policy) (*model.Policy, error) {
	if err := i.Master.WithContext(ctx).Exec(insertPutPolicy, p.PolicyID, p.Name, p.ProjectID, p.ActionPtn, p.ResourcePtn).Error; err != nil {
		return nil, err
	}
	return i.GetPolicyByName(ctx, p.ProjectID, p.Name)
}

const deleteDeletePolicy = `delete from policy where project_id = ? and policy_id = ?`

func (i *iamDB) DeletePolicy(ctx context.Context, projectID, policyID uint32) error {
	return i.Master.WithContext(ctx).Exec(deleteDeletePolicy, projectID, policyID).Error
}

const selectGetRolePolicy = `select * from role_policy where project_id = ? and role_id = ? and policy_id =?`

func (i *iamDB) GetRolePolicy(ctx context.Context, projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	var data model.RolePolicy
	if err := i.Master.WithContext(ctx).Raw(selectGetRolePolicy, projectID, roleID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertAttachPolicy = `
INSERT INTO role_policy
  (role_id, policy_id, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id)
`

func (i *iamDB) AttachPolicy(ctx context.Context, projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	if !i.roleExists(ctx, projectID, roleID) || !i.policyExists(ctx, projectID, policyID) {
		return nil, fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	if err := i.Master.WithContext(ctx).Exec(insertAttachPolicy, roleID, policyID, projectID).Error; err != nil {
		return nil, err
	}
	return i.GetRolePolicy(ctx, projectID, roleID, policyID)
}

const deleteDetachPolicy = `delete from role_policy where role_id = ? and policy_id = ? and project_id = ?`

func (i *iamDB) DetachPolicy(ctx context.Context, projectID, roleID, policyID uint32) error {
	if !i.roleExists(ctx, projectID, roleID) || !i.policyExists(ctx, projectID, policyID) {
		return fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	return i.Master.WithContext(ctx).Exec(deleteDetachPolicy, roleID, policyID, projectID).Error
}
