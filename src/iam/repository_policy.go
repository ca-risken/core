package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (i *iamRepository) ListPolicy(projectID uint32, name string) (*[]model.Policy, error) {
	query := `select * from policy where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	var data []model.Policy
	if err := i.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicy = `select * from policy where project_id = ? and policy_id =?`

func (i *iamRepository) GetPolicy(projectID, policyID uint32) (*model.Policy, error) {
	var data model.Policy
	if err := i.SlaveDB.Raw(selectGetPolicy, projectID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicyByName = `select * from policy where project_id = ? and name =?`

func (i *iamRepository) GetPolicyByName(projectID uint32, name string) (*model.Policy, error) {
	var data model.Policy
	if err := i.SlaveDB.Raw(selectGetPolicyByName, projectID, name).First(&data).Error; err != nil {
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

func (i *iamRepository) PutPolicy(p *model.Policy) (*model.Policy, error) {
	if err := i.MasterDB.Exec(insertPutPolicy, p.PolicyID, p.Name, p.ProjectID, p.ActionPtn, p.ResourcePtn).Error; err != nil {
		return nil, err
	}
	return i.GetPolicyByName(p.ProjectID, p.Name)
}

const deleteDeletePolicy = `delete from policy where project_id = ? and policy_id = ?`

func (i *iamRepository) DeletePolicy(projectID, policyID uint32) error {
	return i.MasterDB.Exec(deleteDeletePolicy, projectID, policyID).Error
}

const selectGetRolePolicy = `select * from role_policy where project_id = ? and role_id = ? and policy_id =?`

func (i *iamRepository) GetRolePolicy(projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	var data model.RolePolicy
	if err := i.SlaveDB.Raw(selectGetRolePolicy, projectID, roleID, policyID).First(&data).Error; err != nil {
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

func (i *iamRepository) AttachPolicy(projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	if !i.roleExists(projectID, roleID) || !i.policyExists(projectID, policyID) {
		return nil, fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	if err := i.MasterDB.Exec(insertAttachPolicy, roleID, policyID, projectID).Error; err != nil {
		return nil, err
	}
	return i.GetRolePolicy(projectID, roleID, policyID)
}

const deleteDetachPolicy = `delete from role_policy where role_id = ? and policy_id = ? and project_id = ?`

func (i *iamRepository) DetachPolicy(projectID, roleID, policyID uint32) error {
	if !i.roleExists(projectID, roleID) || !i.policyExists(projectID, policyID) {
		return fmt.Errorf(
			"Not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	return i.MasterDB.Exec(deleteDetachPolicy, roleID, policyID, projectID).Error
}
