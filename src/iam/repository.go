package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type iamRepoInterface interface {
	// User
	ListUser(activated bool, projectID uint32, name string) (*[]model.User, error)
	GetUser(uint32, string) (*model.User, error)
	GetUserBySub(string) (*model.User, error)
	GetUserPoicy(uint32) (*[]model.Policy, error)
	PutUser(*model.User) (*model.User, error)

	// Role
	ListRole(uint32, string) (*[]model.Role, error)
	GetRole(uint32, uint32) (*model.Role, error)
	GetRoleByName(uint32, string) (*model.Role, error)
	PutRole(r *model.Role) (*model.Role, error)
	DeleteRole(uint32, uint32) error
	AttachRole(uint32, uint32, uint32) (*model.UserRole, error)
	DetachRole(uint32, uint32, uint32) error
}

type iamRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newIAMRepository() iamRepoInterface {
	return &iamRepository{
		MasterDB: initDB(true),
		SlaveDB:  initDB(false),
	}
}

type dbConfig struct {
	MasterHost     string `split_words:"true" required:"true"`
	MasterUser     string `split_words:"true" required:"true"`
	MasterPassword string `split_words:"true" required:"true"`
	SlaveHost      string `split_words:"true"`
	SlaveUser      string `split_words:"true"`
	SlavePassword  string `split_words:"true"`

	Schema  string `required:"true"`
	Port    int    `required:"true"`
	LogMode bool   `split_words:"true" default:"false"`
}

func initDB(isMaster bool) *gorm.DB {
	conf := &dbConfig{}
	if err := envconfig.Process("DB", conf); err != nil {
		appLogger.Fatalf("Failed to load DB config. err: %+v", err)
	}

	var user, pass, host string
	if isMaster {
		user = conf.MasterUser
		pass = conf.MasterPassword
		host = conf.MasterHost
	} else {
		user = conf.SlaveUser
		pass = conf.SlavePassword
		host = conf.SlaveHost
	}

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
			user, pass, host, conf.Port, conf.Schema))
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
		return nil
	}
	db.LogMode(conf.LogMode)
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

func (i *iamRepository) ListUser(activated bool, projectID uint32, name string) (*[]model.User, error) {
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
	if err := i.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (i *iamRepository) GetUser(userID uint32, sub string) (*model.User, error) {
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
	if err := i.SlaveDB.Raw(query, params...).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetUserBySub = `select * from user where sub = ?`

func (i *iamRepository) GetUserBySub(sub string) (*model.User, error) {
	var data model.User
	if err := i.SlaveDB.Raw(selectGetUserBySub, sub).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
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

func (i *iamRepository) GetUserPoicy(userID uint32) (*[]model.Policy, error) {
	var data []model.Policy
	if err := i.SlaveDB.Raw(selectGetUserPolicy, userID).Scan(&data).Error; err != nil {
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

func (i *iamRepository) PutUser(u *model.User) (*model.User, error) {
	if err := i.MasterDB.Exec(insertPutUser, u.UserID, u.Sub, u.Name, fmt.Sprintf("%t", u.Activated)).Error; err != nil {
		return nil, err
	}
	updated, err := i.GetUserBySub(u.Sub)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *iamRepository) ListRole(projectID uint32, name string) (*[]model.Role, error) {
	query := `select * from role where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	var data []model.Role
	if err := i.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRole = `select * from role where project_id = ? and role_id =?`

func (i *iamRepository) GetRole(projectID, roleID uint32) (*model.Role, error) {
	var data model.Role
	if err := i.SlaveDB.Raw(selectGetRole, projectID, roleID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRoleByName = `select * from role where project_id = ? and name =?`

func (i *iamRepository) GetRoleByName(projectID uint32, name string) (*model.Role, error) {
	var data model.Role
	if err := i.SlaveDB.Raw(selectGetRoleByName, projectID, name).First(&data).Error; err != nil {
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

func (i *iamRepository) PutRole(r *model.Role) (*model.Role, error) {
	if err := i.MasterDB.Exec(insertPutRole, r.RoleID, r.Name, r.ProjectID).Error; err != nil {
		return nil, err
	}
	updated, err := i.GetRoleByName(r.ProjectID, r.Name)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

const deleteDeleteRole = `delete from role where project_id = ? and role_id = ?`

func (i *iamRepository) DeleteRole(projectID, roleID uint32) error {
	if err := i.MasterDB.Exec(deleteDeleteRole, projectID, roleID).Error; err != nil {
		return err
	}
	return nil
}

const selectGetUserRole = `select * from user_role where project_id = ? and user_id =? and role_id = ?`

func (i *iamRepository) GetUserRole(projectID, userID, roleID uint32) (*model.UserRole, error) {
	var data model.UserRole
	if err := i.SlaveDB.Raw(selectGetUserRole, projectID, userID, roleID).First(&data).Error; err != nil {
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

func (i *iamRepository) AttachRole(projectID, roleID, userID uint32) (*model.UserRole, error) {
	if !i.userExists(userID) || !i.roleExists(projectID, roleID) {
		return nil, fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if err := i.MasterDB.Exec(insertAttachRole, userID, roleID, projectID).Error; err != nil {
		return nil, err
	}
	return i.GetUserRole(projectID, userID, roleID)
}

const deleteDetachRole = `delete from user_role where user_id = ? and role_id = ? and project_id = ?`

func (i *iamRepository) DetachRole(projectID, roleID, userID uint32) error {
	if !i.userExists(userID) || !i.roleExists(projectID, roleID) {
		return fmt.Errorf(
			"Not found user or role: user_id=%d, role_id=%d, project_id=%d", userID, roleID, projectID)
	}
	if err := i.MasterDB.Exec(deleteDetachRole, userID, roleID, projectID).Error; err != nil {
		return err
	}
	return nil
}

func (i *iamRepository) userExists(userID uint32) bool {
	if _, err := i.GetUser(userID, ""); gorm.IsRecordNotFoundError(err) {
		return false

	} else if err != nil {
		appLogger.Errorf("[userExists]DB error: user_id=%d", userID)
		return false
	}
	return true
}

func (i *iamRepository) roleExists(projectID, roleID uint32) bool {
	if _, err := i.GetRole(projectID, roleID); gorm.IsRecordNotFoundError(err) {
		return false
	} else if err != nil {
		appLogger.Errorf("[roleExists]DB error: project_id=%d, role_id=%d", projectID, roleID)
		return false
	}
	return true
}
