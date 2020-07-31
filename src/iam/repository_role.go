package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (i *iamDB) ListRole(projectID uint32, name string) (*[]model.Role, error) {
	query := `select * from role where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	var data []model.Role
	if err := i.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRole = `select * from role where project_id = ? and role_id =?`

func (i *iamDB) GetRole(projectID, roleID uint32) (*model.Role, error) {
	var data model.Role
	if err := i.Slave.Raw(selectGetRole, projectID, roleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRoleByName = `select * from role where project_id = ? and name =?`

func (i *iamDB) GetRoleByName(projectID uint32, name string) (*model.Role, error) {
	var data model.Role
	if err := i.Slave.Raw(selectGetRoleByName, projectID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPutRole = `
INSERT INTO role
  (role_id, name, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  project_id=VALUES(project_id)
`

func (i *iamDB) PutRole(r *model.Role) (*model.Role, error) {
	if err := i.Master.Exec(insertPutRole, r.RoleID, r.Name, r.ProjectID).Error; err != nil {
		return nil, err
	}
	return i.GetRoleByName(r.ProjectID, r.Name)
}

const deleteDeleteRole = `delete from role where project_id = ? and role_id = ?`

func (i *iamDB) DeleteRole(projectID, roleID uint32) error {
	return i.Master.Exec(deleteDeleteRole, projectID, roleID).Error
}

const selectGetUserRole = `select * from user_role where project_id = ? and user_id =? and role_id = ?`

func (i *iamDB) GetUserRole(projectID, userID, roleID uint32) (*model.UserRole, error) {
	var data model.UserRole
	if err := i.Slave.Raw(selectGetUserRole, projectID, userID, roleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertAttachRole = `
INSERT INTO user_role
  (user_id, role_id, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id)
`

func (i *iamDB) AttachRole(projectID, roleID, userID uint32) (*model.UserRole, error) {
	if !i.userExists(userID) || !i.roleExists(projectID, roleID) {
		return nil, fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if err := i.Master.Exec(insertAttachRole, userID, roleID, projectID).Error; err != nil {
		return nil, err
	}
	return i.GetUserRole(projectID, userID, roleID)
}

const deleteDetachRole = `delete from user_role where user_id = ? and role_id = ? and project_id = ?`

func (i *iamDB) DetachRole(projectID, roleID, userID uint32) error {
	if !i.userExists(userID) || !i.roleExists(projectID, roleID) {
		return fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	return i.Master.Exec(deleteDetachRole, userID, roleID, projectID).Error
}