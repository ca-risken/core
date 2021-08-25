package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (i *iamDB) ListAccessToken(ctx context.Context, projectID uint32, name string, accessTokenID uint32) (*[]model.AccessToken, error) {
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
	if err := i.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAccessTokenByID = `select * from access_token where project_id=? and access_token_id=?`

func (i *iamDB) GetAccessTokenByID(ctx context.Context, projectID, accessTokenID uint32) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := i.Master.WithContext(ctx).Raw(selectGetAccessTokenByID, projectID, accessTokenID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetActiveAccessTokenHash = `select * from access_token where project_id=? and access_token_id=? and token_hash=? and expired_at >= NOW()`

func (i *iamDB) GetActiveAccessTokenHash(ctx context.Context, projectID, accessTokenID uint32, tokenHash string) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := i.Master.WithContext(ctx).Raw(selectGetActiveAccessTokenHash, projectID, accessTokenID, tokenHash).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAccessTokenByUniqueKey = `select * from access_token where project_id = ? and name =?`

func (i *iamDB) GetAccessTokenByUniqueKey(ctx context.Context, projectID uint32, name string) (*model.AccessToken, error) {
	var data model.AccessToken
	if err := i.Master.WithContext(ctx).Raw(selectGetAccessTokenByUniqueKey, projectID, name).First(&data).Error; err != nil {
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

func (i *iamDB) PutAccessToken(ctx context.Context, r *model.AccessToken) (*model.AccessToken, error) {
	if err := i.Master.WithContext(ctx).Exec(insertPutAccessToken,
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
	return i.GetAccessTokenByUniqueKey(ctx, r.ProjectID, r.Name)
}

const deleteDeleteAccessToken = `delete from access_token where project_id=? and access_token_id=?`

func (i *iamDB) DeleteAccessToken(ctx context.Context, projectID, accessTokenID uint32) error {
	return i.Master.WithContext(ctx).Exec(deleteDeleteAccessToken, projectID, accessTokenID).Error
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

func (i *iamDB) GetAccessTokenRole(ctx context.Context, accessTokenID, roleID uint32) (*model.AccessTokenRole, error) {
	var data model.AccessTokenRole
	if err := i.Master.WithContext(ctx).Raw(selectGetAccessTokenRole, accessTokenID, roleID).First(&data).Error; err != nil {
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

func (i *iamDB) AttachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) (*model.AccessTokenRole, error) {
	if !i.accessTokenExists(ctx, projectID, accessTokenID) || !i.roleExists(ctx, projectID, roleID) {
		return nil, fmt.Errorf(
			"Not found access_token or role: access_token_id=%d, role_id=%d, project_id=%d", accessTokenID, roleID, projectID)
	}
	if err := i.Master.WithContext(ctx).Exec(insertAttachAccessTokenRole, accessTokenID, roleID).Error; err != nil {
		return nil, err
	}
	return i.GetAccessTokenRole(ctx, accessTokenID, roleID)
}

const deleteDetachAccessTokenRole = `delete from access_token_role where access_token_id=? and role_id=?`

func (i *iamDB) DetachAccessTokenRole(ctx context.Context, projectID, roleID, accessTokenID uint32) error {
	if !i.accessTokenExists(ctx, projectID, accessTokenID) || !i.roleExists(ctx, projectID, roleID) {
		return fmt.Errorf(
			"Not found access_token or role: access_token_id=%d, role_id=%d, project_id=%d", accessTokenID, roleID, projectID)
	}
	return i.Master.WithContext(ctx).Exec(deleteDetachAccessTokenRole, accessTokenID, roleID).Error
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

func (i *iamDB) ExistsAccessTokenMaintainer(ctx context.Context, projectID, accessTokenID uint32) (bool, error) {
	var data model.User
	if err := i.Slave.WithContext(ctx).Raw(selectExistsAccessTokenMaintainer, projectID, accessTokenID).First(&data).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

const selectListExpiredAccessToken = `select * from access_token where expired_at < NOW()`

func (i *iamDB) ListExpiredAccessToken(ctx context.Context) (*[]model.AccessToken, error) {
	var data []model.AccessToken
	if err := i.Slave.WithContext(ctx).Raw(selectListExpiredAccessToken).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
