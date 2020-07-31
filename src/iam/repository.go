package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type iamRepository interface {
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

	// Policy
	ListPolicy(uint32, string) (*[]model.Policy, error)
	GetPolicy(uint32, uint32) (*model.Policy, error)
	GetPolicyByName(uint32, string) (*model.Policy, error)
	PutPolicy(*model.Policy) (*model.Policy, error)
	DeletePolicy(uint32, uint32) error
	AttachPolicy(uint32, uint32, uint32) (*model.RolePolicy, error)
	DetachPolicy(uint32, uint32, uint32) error
}

type iamDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newIAMRepository() iamRepository {
	return &iamDB{
		Master: initDB(true),
		Slave:  initDB(false),
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

func (i *iamDB) userExists(userID uint32) bool {
	if _, err := i.GetUser(userID, ""); gorm.IsRecordNotFoundError(err) {
		return false

	} else if err != nil {
		appLogger.Errorf("[userExists]DB error: user_id=%d", userID)
		return false
	}
	return true
}

func (i *iamDB) roleExists(projectID, roleID uint32) bool {
	if _, err := i.GetRole(projectID, roleID); gorm.IsRecordNotFoundError(err) {
		return false
	} else if err != nil {
		appLogger.Errorf("[roleExists]DB error: project_id=%d, role_id=%d", projectID, roleID)
		return false
	}
	return true
}

func (i *iamDB) policyExists(projectID, policyID uint32) bool {
	if _, err := i.GetPolicy(projectID, policyID); gorm.IsRecordNotFoundError(err) {
		return false
	} else if err != nil {
		appLogger.Errorf("[policyExists]DB error: project_id=%d, policy_id=%d", projectID, policyID)
		return false
	}
	return true
}
