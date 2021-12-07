package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/src/iam/model"
	"github.com/vikyd/zero"
)

func (i *iamDB) ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error) {
	query := `select * from role r where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and r.project_id = ?"
		params = append(params, projectID)
	} else {
		query += " and r.project_id is null"
	}
	if !zero.IsZeroVal(name) {
		query += " and r.name = ?"
		params = append(params, name)
	}
	if !zero.IsZeroVal(userID) {
		query += " and exists (select * from user_role ur where ur.role_id = r.role_id and ur.user_id = ? )"
		params = append(params, userID)
	}
	if !zero.IsZeroVal(accessTokenID) {
		query += " and exists (select * from access_token_role atr inner join access_token at using(access_token_id) where atr.role_id = r.role_id and atr.access_token_id = ?)"
		params = append(params, accessTokenID)
	}
	var data []model.Role
	if err := i.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (i *iamDB) GetRole(ctx context.Context, projectID, roleID uint32) (*model.Role, error) {
	query := `select * from role r where role_id =?`
	var params []interface{}
	params = append(params, roleID)
	if !zero.IsZeroVal(projectID) {
		query += " and r.project_id = ?"
		params = append(params, projectID)
	} else {
		query += " and r.project_id is null"
	}
	var data model.Role
	if err := i.Master.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRoleByName = `select * from role where project_id = ? and name =?`

func (i *iamDB) GetRoleByName(ctx context.Context, projectID uint32, name string) (*model.Role, error) {
	var data model.Role
	if err := i.Master.WithContext(ctx).Raw(selectGetRoleByName, projectID, name).First(&data).Error; err != nil {
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

func (i *iamDB) PutRole(ctx context.Context, r *model.Role) (*model.Role, error) {
	if err := i.Master.WithContext(ctx).Exec(insertPutRole, r.RoleID, r.Name, r.ProjectID).Error; err != nil {
		return nil, err
	}
	return i.GetRoleByName(ctx, r.ProjectID, r.Name)
}

const deleteDeleteRole = `delete from role where project_id = ? and role_id = ?`

func (i *iamDB) DeleteRole(ctx context.Context, projectID, roleID uint32) error {
	return i.Master.WithContext(ctx).Exec(deleteDeleteRole, projectID, roleID).Error
}

const selectGetUserRole = `select * from user_role where project_id = ? and user_id =? and role_id = ?`

func (i *iamDB) GetUserRole(ctx context.Context, projectID, userID, roleID uint32) (*model.UserRole, error) {
	var data model.UserRole
	if err := i.Master.WithContext(ctx).Raw(selectGetUserRole, projectID, userID, roleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const (
	insertAttachAdminRole = `
INSERT INTO user_role
  (user_id, role_id, project_id)
VALUES
  (?, ?, null)
ON DUPLICATE KEY UPDATE
  updated_at=NOW()
`
	insertAttachProjectRole = `
INSERT INTO user_role
  (user_id, role_id, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  updated_at=NOW()
`
)

func (i *iamDB) AttachRole(ctx context.Context, projectID, roleID, userID uint32) (*model.UserRole, error) {
	if !i.userExists(ctx, userID) || !i.roleExists(ctx, projectID, roleID) {
		return nil, fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if zero.IsZeroVal(projectID) {
		if err := i.Master.WithContext(ctx).Exec(insertAttachAdminRole, userID, roleID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := i.Master.WithContext(ctx).Exec(insertAttachProjectRole, userID, roleID, projectID).Error; err != nil {
			return nil, err
		}
	}
	return i.GetUserRole(ctx, projectID, userID, roleID)
}

const insertAttachAllAdminRole = `
INSERT INTO user_role
  (user_id, role_id, project_id)
SELECT
  ?, role_id, null
FROM 
  role
WHERE
  project_id is null
ON DUPLICATE KEY UPDATE
  updated_at=NOW()
`

func (i *iamDB) AttachAllAdminRole(ctx context.Context, userID uint32) error {
	if !i.userExists(ctx, userID) {
		return fmt.Errorf("Not found user: user_id=%d", userID)
	}
	return i.Master.WithContext(ctx).Exec(insertAttachAllAdminRole, userID).Error
}

const (
	deleteDetachAdminRole   = `delete from user_role where user_id = ? and role_id = ? and project_id is null`
	deleteDetachProjectRole = `delete from user_role where user_id = ? and role_id = ? and project_id = ?`
)

func (i *iamDB) DetachRole(ctx context.Context, projectID, roleID, userID uint32) error {
	if !i.userExists(ctx, userID) || !i.roleExists(ctx, projectID, roleID) {
		return fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if zero.IsZeroVal(projectID) {
		return i.Master.WithContext(ctx).Exec(deleteDetachAdminRole, userID, roleID, projectID).Error
	} else {
		return i.Master.WithContext(ctx).Exec(deleteDetachProjectRole, userID, roleID, projectID).Error
	}
}
