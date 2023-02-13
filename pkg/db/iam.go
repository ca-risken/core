package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type IAMRepository interface {
	// User
	ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32, admin bool) (*[]model.User, error)
	GetUser(ctx context.Context, userID uint32, sub string) (*model.User, error)
	GetUserBySub(ctx context.Context, sub string) (*model.User, error)
	CreateUser(ctx context.Context, u *model.User) (*model.User, error)
	PutUser(ctx context.Context, u *model.User) (*model.User, error)
	GetActiveUserCount(ctx context.Context) (*int, error)
	GetUserByUserIdpKey(ctx context.Context, userIdpKey string) (*model.User, error)

	// Role
	ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error)
	GetRole(ctx context.Context, projectID, roleID uint32) (*model.Role, error)
	GetRoleByName(ctx context.Context, projectID uint32, name string) (*model.Role, error)
	PutRole(ctx context.Context, r *model.Role) (*model.Role, error)
	DeleteRole(ctx context.Context, projectID, roleID uint32) error
	AttachRole(ctx context.Context, projectID, roleID, userID uint32) (*model.UserRole, error)
	AttachAllAdminRole(ctx context.Context, userID uint32) error
	DetachRole(ctx context.Context, projectID, roleID, userID uint32) error

	// Policy
	GetUserPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error)
	GetTokenPolicy(ctx context.Context, accessTokenID uint32) (*[]model.Policy, error)
	GetAdminPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error)
	ListPolicy(ctx context.Context, projectID uint32, name string, roleID uint32) (*[]model.Policy, error)
	GetPolicy(ctx context.Context, projectID, policyID uint32) (*model.Policy, error)
	GetPolicyByName(ctx context.Context, projectID uint32, name string) (*model.Policy, error)
	PutPolicy(ctx context.Context, p *model.Policy) (*model.Policy, error)
	DeletePolicy(ctx context.Context, projectID, policyID uint32) error
	AttachPolicy(ctx context.Context, projectID, roleID, policyID uint32) (*model.RolePolicy, error)
	DetachPolicy(ctx context.Context, projectID, roleID, policyID uint32) error

	// AccessToken
	ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error)
	GetAccessTokenByID(ctx context.Context, projectID, accessTokenID uint32) (*model.AccessToken, error)
	GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error)
	GetActiveAccessTokenHash(ctx context.Context, projectID, accessTokenID uint32, tokenHash string) (*model.AccessToken, error)
	PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error)
	DeleteAccessToken(ctx context.Context, projectID, accessTokenID uint32) error
	AttachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) (*model.AccessTokenRole, error)
	GetAccessTokenRole(ctx context.Context, accessTokenID, roleID uint32) (*model.AccessTokenRole, error)
	DetachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) error
	ExistsAccessTokenMaintainer(ctx context.Context, projectID, accessTokenID uint32) (bool, error)
	ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error)

	// UserReserved
	ListUserReserved(ctx context.Context, projectID uint32, userIdpKey string) (*[]model.UserReserved, error)
	ListUserReservedWithProjectID(ctx context.Context, userIdpKey string) (*[]UserReservedWithProjectID, error)
	PutUserReserved(ctx context.Context, u *model.UserReserved) (*model.UserReserved, error)
	DeleteUserReserved(ctx context.Context, projectID, reservedID uint32) error
}

var _ IAMRepository = (*Client)(nil)

func (c *Client) ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32, admin bool) (*[]model.User, error) {
	query := `
select
  u.*
from
  user u
where
  activated = ?
`
	var params []interface{}
	params = append(params, fmt.Sprintf("%t", activated))
	if !zero.IsZeroVal(projectID) {
		query += " and exists (select * from user_role ur inner join role r using(role_id, project_id) where ur.user_id = u.user_id and ur.project_id = ?)"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(name) {
		escapedName := escapeLikeParam(name)
		query += fmt.Sprintf(" and (u.name like ? escape '%s' or u.user_idp_key like ? escape '%s' )", escapeString, escapeString)
		params = append(params, "%"+escapedName+"%", "%"+escapedName+"%")
	}
	if !zero.IsZeroVal(userID) {
		query += " and u.user_id = ?"
		params = append(params, userID)
	}
	if admin {
		query += " and exists (select * from user_role ur where ur.user_id = u.user_id and ur.project_id is null)"
	}
	var data []model.User
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetUser(ctx context.Context, userID uint32, sub string) (*model.User, error) {
	query := `select * from user where activated = 'true'`
	var params []interface{}
	if !zero.IsZeroVal(userID) {
		query += " and user_id = ?"
		params = append(params, userID)
	}
	if !zero.IsZeroVal(sub) {
		query += " and sub = ?"
		params = append(params, sub)
	}
	var data model.User
	if err := c.Master.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetUserBySub = `select * from user where sub = ?`

func (c *Client) GetUserBySub(ctx context.Context, sub string) (*model.User, error) {
	var data model.User
	if err := c.Master.WithContext(ctx).Raw(selectGetUserBySub, sub).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetUserByUserIdpKey(ctx context.Context, userIdpKey string) (*model.User, error) {
	var data model.User
	if err := c.Slave.WithContext(ctx).Where("user_idp_key = ?", userIdpKey).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUser = `
INSERT INTO user
  (user_id, sub, name, user_idp_key, activated)
  VALUES (?, ?, ?, ?, ?)
`

func (c *Client) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {

	if err := c.Master.WithContext(ctx).Exec(insertUser, u.UserID, u.Sub, u.Name, convertZeroValueToNull(u.UserIdpKey), fmt.Sprintf("%t", u.Activated)).Error; err != nil {
		return nil, err
	}
	return c.GetUserBySub(ctx, u.Sub)
}

const updateUser = `
UPDATE user
  SET 
    name = ?,
    user_idp_key = ?,
    activated = ?
  WHERE user_id = ?
`

func (c *Client) PutUser(ctx context.Context, u *model.User) (*model.User, error) {
	if err := c.Master.WithContext(ctx).Exec(updateUser, u.Name, convertZeroValueToNull(u.UserIdpKey), fmt.Sprintf("%t", u.Activated), u.UserID).Error; err != nil {
		return nil, err
	}
	return c.GetUserBySub(ctx, u.Sub)
}

const selectGetActiveUserCount = `select count(*) from user where activated = 'true'`

func (c *Client) GetActiveUserCount(ctx context.Context) (*int, error) {
	var cnt int
	if err := c.Slave.WithContext(ctx).Raw(selectGetActiveUserCount).Scan(&cnt).Error; err != nil {
		return nil, err
	}
	return &cnt, nil
}

func (c *Client) ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error) {
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetRole(ctx context.Context, projectID, roleID uint32) (*model.Role, error) {
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
	if err := c.Master.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRoleByName = `select * from role where project_id = ? and name =?`

func (c *Client) GetRoleByName(ctx context.Context, projectID uint32, name string) (*model.Role, error) {
	var data model.Role
	if err := c.Master.WithContext(ctx).Raw(selectGetRoleByName, projectID, name).First(&data).Error; err != nil {
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

func (c *Client) PutRole(ctx context.Context, r *model.Role) (*model.Role, error) {
	if err := c.Master.WithContext(ctx).Exec(insertPutRole, r.RoleID, r.Name, r.ProjectID).Error; err != nil {
		return nil, err
	}
	return c.GetRoleByName(ctx, r.ProjectID, r.Name)
}

const deleteDeleteRole = `delete from role where project_id = ? and role_id = ?`

func (c *Client) DeleteRole(ctx context.Context, projectID, roleID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteDeleteRole, projectID, roleID).Error
}

const (
	selectGetAdminUserRole   = `select * from user_role where project_id is null and user_id = ? and role_id = ?`
	selectGetProjectUserRole = `select * from user_role where project_id = ?     and user_id = ? and role_id = ?`
)

func (c *Client) GetUserRole(ctx context.Context, projectID, userID, roleID uint32) (*model.UserRole, error) {
	var data model.UserRole
	if zero.IsZeroVal(projectID) {
		if err := c.Master.WithContext(ctx).Raw(selectGetAdminUserRole, userID, roleID).First(&data).Error; err != nil {
			return nil, err
		}
	} else {
		if err := c.Master.WithContext(ctx).Raw(selectGetProjectUserRole, projectID, userID, roleID).First(&data).Error; err != nil {
			return nil, err
		}
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

func (c *Client) AttachRole(ctx context.Context, projectID, roleID, userID uint32) (*model.UserRole, error) {
	userExists, err := c.userExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return nil, err
	}
	if !userExists || !roleExists {
		return nil, fmt.Errorf(
			"not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if zero.IsZeroVal(projectID) {
		if err := c.Master.WithContext(ctx).Exec(insertAttachAdminRole, userID, roleID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := c.Master.WithContext(ctx).Exec(insertAttachProjectRole, userID, roleID, projectID).Error; err != nil {
			return nil, err
		}
	}
	return c.GetUserRole(ctx, projectID, userID, roleID)
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

func (c *Client) AttachAllAdminRole(ctx context.Context, userID uint32) error {
	exists, err := c.userExists(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("not found user: user_id=%d", userID)
	}
	return c.Master.WithContext(ctx).Exec(insertAttachAllAdminRole, userID).Error
}

const (
	deleteDetachAdminRole   = `delete from user_role where user_id = ? and role_id = ? and project_id is null`
	deleteDetachProjectRole = `delete from user_role where user_id = ? and role_id = ? and project_id = ?`
)

func (c *Client) DetachRole(ctx context.Context, projectID, roleID, userID uint32) error {
	userExists, err := c.userExists(ctx, userID)
	if err != nil {
		return err
	}
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return err
	}
	if !userExists || !roleExists {
		return fmt.Errorf(
			"not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if zero.IsZeroVal(projectID) {
		return c.Master.WithContext(ctx).Exec(deleteDetachAdminRole, userID, roleID, projectID).Error
	} else {
		return c.Master.WithContext(ctx).Exec(deleteDetachProjectRole, userID, roleID, projectID).Error
	}
}

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

func (c *Client) GetUserPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := c.Slave.WithContext(ctx).Raw(selectGetUserPolicy, userID).Scan(&data).Error; err != nil {
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

func (c *Client) GetTokenPolicy(ctx context.Context, accessTokenID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := c.Slave.WithContext(ctx).Raw(selectGetTokenPolicy, accessTokenID).Scan(&data).Error; err != nil {
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

func (c *Client) GetAdminPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := c.Slave.WithContext(ctx).Raw(selectGetAdminPolicy, userID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListPolicy(ctx context.Context, projectID uint32, name string, roleID uint32) (*[]model.Policy, error) {
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicy = `select * from policy where project_id = ? and policy_id =?`

func (c *Client) GetPolicy(ctx context.Context, projectID, policyID uint32) (*model.Policy, error) {
	var data model.Policy
	if err := c.Master.WithContext(ctx).Raw(selectGetPolicy, projectID, policyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetPolicyByName = `select * from policy where project_id = ? and name =?`

func (c *Client) GetPolicyByName(ctx context.Context, projectID uint32, name string) (*model.Policy, error) {
	var data model.Policy
	if err := c.Master.WithContext(ctx).Raw(selectGetPolicyByName, projectID, name).First(&data).Error; err != nil {
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

func (c *Client) PutPolicy(ctx context.Context, p *model.Policy) (*model.Policy, error) {
	if err := c.Master.WithContext(ctx).Exec(insertPutPolicy, p.PolicyID, p.Name, p.ProjectID, p.ActionPtn, p.ResourcePtn).Error; err != nil {
		return nil, err
	}
	return c.GetPolicyByName(ctx, p.ProjectID, p.Name)
}

const deleteDeletePolicy = `delete from policy where project_id = ? and policy_id = ?`

func (c *Client) DeletePolicy(ctx context.Context, projectID, policyID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteDeletePolicy, projectID, policyID).Error
}

const selectGetRolePolicy = `select * from role_policy where project_id = ? and role_id = ? and policy_id =?`

func (c *Client) GetRolePolicy(ctx context.Context, projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	var data model.RolePolicy
	if err := c.Master.WithContext(ctx).Raw(selectGetRolePolicy, projectID, roleID, policyID).First(&data).Error; err != nil {
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

func (c *Client) AttachPolicy(ctx context.Context, projectID, roleID, policyID uint32) (*model.RolePolicy, error) {
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return nil, err
	}
	policyExists, err := c.policyExists(ctx, projectID, policyID)
	if err != nil {
		return nil, err
	}
	if !roleExists || !policyExists {
		return nil, fmt.Errorf(
			"not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	if err := c.Master.WithContext(ctx).Exec(insertAttachPolicy, roleID, policyID, projectID).Error; err != nil {
		return nil, err
	}
	return c.GetRolePolicy(ctx, projectID, roleID, policyID)
}

const deleteDetachPolicy = `delete from role_policy where role_id = ? and policy_id = ? and project_id = ?`

func (c *Client) DetachPolicy(ctx context.Context, projectID, roleID, policyID uint32) error {
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return err
	}
	policyExists, err := c.policyExists(ctx, projectID, policyID)
	if err != nil {
		return err
	}
	if !roleExists || !policyExists {
		return fmt.Errorf(
			"not found role or policy: role_id=%d, policy_id=%d, project_id=%d", roleID, policyID, projectID)
	}
	return c.Master.WithContext(ctx).Exec(deleteDetachPolicy, roleID, policyID, projectID).Error
}

func (c *Client) ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error) {
	// query := `select * from access_token a where a.project_id=? and a.expired_at >= NOW()`
	query := `select * from access_token a where a.project_id=?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and a.name = ?"
		params = append(params, name)
	}
	if !zero.IsZeroVal(accessTokenID) {
		query += " and a.access_token_id = ?"
		params = append(params, accessTokenID)
	}
	var data []model.AccessToken
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAccessTokenByID = `select * from access_token where project_id=? and access_token_id=?`

func (c *Client) GetAccessTokenByID(ctx context.Context, projectID, accessTokenID uint32) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := c.Master.WithContext(ctx).Raw(selectGetAccessTokenByID, projectID, accessTokenID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetActiveAccessTokenHash = `select * from access_token where project_id=? and access_token_id=? and token_hash=? and expired_at >= NOW()`

func (c *Client) GetActiveAccessTokenHash(ctx context.Context, projectID, accessTokenID uint32, tokenHash string) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := c.Master.WithContext(ctx).Raw(selectGetActiveAccessTokenHash, projectID, accessTokenID, tokenHash).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAccessTokenByUniqueKey = `select * from access_token where project_id = ? and name =?`

func (c *Client) GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := c.Master.WithContext(ctx).Raw(selectGetAccessTokenByUniqueKey, projectID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPutAccessToken = `
INSERT INTO access_token
  (access_token_id, token_hash, name, description, project_id, expired_at, last_updated_user_id)
VALUES
  (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  token_hash=VALUES(token_hash),
  name=VALUES(name),
  description=VALUES(description),
  project_id=VALUES(project_id),
  expired_at=VALUES(expired_at),
  last_updated_user_id=VALUES(last_updated_user_id)
`

func (c *Client) PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error) {
	if err := c.Master.WithContext(ctx).Exec(insertPutAccessToken,
		r.AccessTokenID,
		r.TokenHash,
		r.Name,
		convertZeroValueToNull(r.Description),
		r.ProjectID,
		r.ExpiredAt,
		r.LastUpdatedUserID,
	).Error; err != nil {
		return nil, err
	}
	return c.GetAccessTokenByUniqueKey(ctx, r.ProjectID, r.Name)
}

const deleteDeleteAccessToken = `delete from access_token where project_id=? and access_token_id=?`

func (c *Client) DeleteAccessToken(ctx context.Context, projectID, accessTokenID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteDeleteAccessToken, projectID, accessTokenID).Error
}

const selectGetAccessTokenRole = `
select
  * 
from
  access_token_role atr 
where
  atr.access_token_id=? 
  and atr.role_id=? 
  and exists(
    select * 
    from access_token at 
    where at.access_token_id=atr.access_token_id
  )
`

func (c *Client) GetAccessTokenRole(ctx context.Context, accessTokenID, roleID uint32) (*model.AccessTokenRole, error) {
	var data model.AccessTokenRole
	if err := c.Master.WithContext(ctx).Raw(selectGetAccessTokenRole, accessTokenID, roleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertAttachAccessTokenRole = `
INSERT INTO access_token_role
  (access_token_id, role_id)
VALUES
  (?, ?)
ON DUPLICATE KEY UPDATE
  access_token_id=VALUES(access_token_id),
  role_id=VALUES(role_id)
`

func (c *Client) AttachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) (*model.AccessTokenRole, error) {
	accessTokenExists, err := c.accessTokenExists(ctx, projectID, accessTokenID)
	if err != nil {
		return nil, err
	}
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return nil, err
	}
	if !accessTokenExists || !roleExists {
		return nil, fmt.Errorf(
			"not found access_token or role: access_token_id=%d, role_id=%d, project_id=%d", accessTokenID, roleID, projectID)
	}
	if err := c.Master.WithContext(ctx).Exec(insertAttachAccessTokenRole, accessTokenID, roleID).Error; err != nil {
		return nil, err
	}
	return c.GetAccessTokenRole(ctx, accessTokenID, roleID)
}

const deleteDetachAccessTokenRole = `delete from access_token_role where access_token_id=? and role_id=?`

func (c *Client) DetachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) error {
	accessTokenExists, err := c.accessTokenExists(ctx, projectID, accessTokenID)
	if err != nil {
		return err
	}
	roleExists, err := c.roleExists(ctx, projectID, roleID)
	if err != nil {
		return err
	}
	if !accessTokenExists || !roleExists {
		return fmt.Errorf(
			"not found access_token or role: access_token_id=%d, role_id=%d, project_id=%d", accessTokenID, roleID, projectID)
	}
	return c.Master.WithContext(ctx).Exec(deleteDetachAccessTokenRole, accessTokenID, roleID).Error
}

const selectExistsAccessTokenMaintainer = `
select
  u.user_id 
from
  access_token at
  inner join role r using(project_id)
  inner join user_role ur using(role_id)
  inner join user u using(user_id)
where
  at.project_id=? 
  and at.expired_at >= NOW()
  and at.access_token_id=?
  and u.activated='true'
`

func (c *Client) ExistsAccessTokenMaintainer(ctx context.Context, projectID, accessTokenID uint32) (bool, error) {
	var data model.User
	if err := c.Slave.WithContext(ctx).Raw(selectExistsAccessTokenMaintainer, projectID, accessTokenID).First(&data).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

const selectListExpiredAccessToken = `select * from access_token where expired_at < NOW()`

func (c *Client) ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error) {
	var data []model.AccessToken
	if err := c.Slave.WithContext(ctx).Raw(selectListExpiredAccessToken).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) userExists(ctx context.Context, userID uint32) (bool, error) {
	if _, err := c.GetUser(ctx, userID, ""); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get user. user_id=%d, error: %w", userID, err)
	}
	return true, nil
}

func (c *Client) roleExists(ctx context.Context, projectID, roleID uint32) (bool, error) {
	if _, err := c.GetRole(ctx, projectID, roleID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get role. project_id=%d, role_id=%d, error: %w", projectID, roleID, err)
	}
	return true, nil
}

func (c *Client) policyExists(ctx context.Context, projectID, policyID uint32) (bool, error) {
	if _, err := c.GetPolicy(ctx, projectID, policyID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get policy. project_id=%d, policy_id=%d, error: %w", projectID, policyID, err)
	}
	return true, nil
}

func (c *Client) accessTokenExists(ctx context.Context, projectID, accessTokenID uint32) (bool, error) {
	if _, err := c.GetAccessTokenByID(ctx, projectID, accessTokenID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get access token. project_id=%d, access_token_id=%d, error: %w", projectID, accessTokenID, err)
	}
	return true, nil
}

const listUserReserved = `
select ur.* 
from user_reserved ur inner join role r using(role_id)
where r.project_id = ?
`

func (c *Client) ListUserReserved(ctx context.Context, projectID uint32, userIdpKey string) (*[]model.UserReserved, error) {
	query := listUserReserved
	params := []interface{}{
		projectID,
	}
	if userIdpKey != "" {
		query += " and ur.user_idp_key = ?"
		params = append(params, userIdpKey)
	}
	var data []model.UserReserved
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// For ListUserReservedWithProjectID
type UserReservedWithProjectID struct {
	ProjectID  uint32
	ReservedID uint32
	RoleID     uint32
}

const listUserReservedWithProjectID = `
select ur.reserved_id,ur.role_id,r.project_id 
from user_reserved ur inner join role r using(role_id)
where  ur.user_idp_key = ?
`

func (c *Client) ListUserReservedWithProjectID(ctx context.Context, userIdpKey string) (*[]UserReservedWithProjectID, error) {
	var data []UserReservedWithProjectID
	if err := c.Slave.WithContext(ctx).Raw(listUserReservedWithProjectID, userIdpKey).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) PutUserReserved(ctx context.Context, data *model.UserReserved) (*model.UserReserved, error) {
	var ret *model.UserReserved
	if err := c.Master.WithContext(ctx).Where("reserved_id = ?", data.ReservedID).Assign(data).FirstOrCreate(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

const deleteUserReserved = `
delete ur from user_reserved ur 
where exists 
  (select * from role r where ur.role_id = r.role_id and r.project_id = ?) and ur.reserved_id = ?
`

func (c *Client) DeleteUserReserved(ctx context.Context, projectID, reservedID uint32) error {
	if err := c.Master.WithContext(ctx).Exec(deleteUserReserved, projectID, reservedID).Error; err != nil {
		return err
	}
	return nil
}

func convertZeroValueToNull(input interface{}) interface{} {
	if input == nil || zero.IsZeroVal(input) {
		return gorm.Expr("NULL")
	}
	return input
}
