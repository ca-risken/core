package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/src/iam/model"
	"github.com/vikyd/zero"
)

func (i *iamDB) ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32, admin bool) (*[]model.User, error) {
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
		query += " and u.name like ?"
		params = append(params, "%"+name+"%")
	}
	if !zero.IsZeroVal(userID) {
		query += " and u.user_id = ?"
		params = append(params, userID)
	}
	if admin {
		query += " and exists (select * from user_role ur where ur.user_id = u.user_id and ur.project_id is null)"
	}
	var data []model.User
	if err := i.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (i *iamDB) GetUser(ctx context.Context, userID uint32, sub string) (*model.User, error) {
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
	if err := i.Master.WithContext(ctx).Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetUserBySub = `select * from user where sub = ?`

func (i *iamDB) GetUserBySub(ctx context.Context, sub string) (*model.User, error) {
	var data model.User
	if err := i.Master.WithContext(ctx).Raw(selectGetUserBySub, sub).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPutUser = `
INSERT INTO user
  (user_id, sub, name, activated)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  activated=VALUES(activated)
`

func (i *iamDB) PutUser(ctx context.Context, u *model.User) (*model.User, error) {
	if err := i.Master.WithContext(ctx).Exec(insertPutUser, u.UserID, u.Sub, u.Name, fmt.Sprintf("%t", u.Activated)).Error; err != nil {
		return nil, err
	}
	return i.GetUserBySub(ctx, u.Sub)
}

const selectGetActiveUserCount = `select count(*) from user where activated = 'true'`

func (i *iamDB) GetActiveUserCount(ctx context.Context) (*int, error) {
	var cnt int
	if err := i.Slave.WithContext(ctx).Raw(selectGetActiveUserCount).Scan(&cnt).Error; err != nil {
		return nil, err
	}
	return &cnt, nil
}
