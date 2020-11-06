package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/vikyd/zero"
)

func (i *iamDB) ListUser(activated bool, projectID uint32, name string) (*[]model.User, error) {
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
		query += " and exists (select * from user_role ur where ur.user_id = u.user_id and ur.project_id = ?)"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(name) {
		query += " and u.name = ?"
		params = append(params, name)
	}
	var data []model.User
	if err := i.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (i *iamDB) GetUser(userID uint32, sub string) (*model.User, error) {
	query := `select * from	user where activated = 'true'`
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
	if err := i.Slave.Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetUserBySub = `select * from user where sub = ?`

func (i *iamDB) GetUserBySub(sub string) (*model.User, error) {
	var data model.User
	if err := i.Master.Raw(selectGetUserBySub, sub).First(&data).Error; err != nil {
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

func (i *iamDB) PutUser(u *model.User) (*model.User, error) {
	if err := i.Master.Exec(insertPutUser, u.UserID, u.Sub, u.Name, fmt.Sprintf("%t", u.Activated)).Error; err != nil {
		return nil, err
	}
	return i.GetUserBySub(u.Sub)
}
