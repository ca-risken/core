package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (i *iamDB) ListPolicy(projectID uint32, name string) (*[]model.Policy, error) {
	query := `select * from policy where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	var data []model.Policy
	if err := i.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicy = `select * from policy where project_id = ? and policy_id =?`

func (i *iamDB) GetPolicy(projectID, policyID uint32) (*model.Policy, error) {
	var data model.Policy
	if err := i.Slave.Raw(selectGetPolicy, projectID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicyByName = `select * from policy where project_id = ? and name =?`

func (i *iamDB) GetPolicyByName(projectID uint32, name string) (*model.Policy, error) {
	var data model.Policy
	if err := i.Slave.Raw(selectGetPolicyByName, projectID, name).First(&data).Error; err != nil {
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

func (i *iamDB) PutPolicy(p *model.Policy) (*model.Policy, error) {
	if err := i.Master.Exec(insertPutPolicy, p.PolicyID, p.Name, p.ProjectID, p.ActionPtn, p.ResourcePtn).Error; err != nil {
		return nil, err
	}
	return i.GetPolicyByName(p.ProjectID, p.Name)
}

const deleteDeletePolicy = `delete from policy where project_id = ? and policy_id = ?`

func (i *iamDB) DeletePolicy(projectID, policyID uint32) error {
	return i.Master.Exec(deleteDeletePolicy, projectID, policyID).Error
}

const selectGetRolePolicy = `select * from role_policy where project_id = ? and role_id = ? and policy_id =?`

func (i *iamDB) GetRolePolicy(projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	var data model.RolePolicy
	if err := i.Slave.Raw(selectGetRolePolicy, projectID, roleID, policyID).First(&data).Error; err != nil {
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

func (i *iamDB) AttachPolicy(projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	if !i.roleExists(projectID, roleID) || !i.policyExists(projectID, policyID) {
		return nil, fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	if err := i.Master.Exec(insertAttachPolicy, roleID, policyID, projectID).Error; err != nil {
		return nil, err
	}
	return i.GetRolePolicy(projectID, roleID, policyID)
}

const deleteDetachPolicy = `delete from role_policy where role_id = ? and policy_id = ? and project_id = ?`

func (i *iamDB) DetachPolicy(projectID, roleID, policyID uint32) error {
	if !i.roleExists(projectID, roleID) || !i.policyExists(projectID, policyID) {
		return fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	return i.Master.Exec(deleteDetachPolicy, roleID, policyID, projectID).Error
}