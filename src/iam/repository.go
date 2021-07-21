package main

import (
	"context"
	"errors"
	"fmt"

	mimosasql "github.com/CyberAgent/mimosa-common/pkg/database/sql"
	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/gorm"
)

type iamRepository interface {
	// User
	ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32) (*[]model.User, error)
	GetUser(context.Context, uint32, string) (*model.User, error)
	GetUserBySub(context.Context, string) (*model.User, error)
	PutUser(context.Context, *model.User) (*model.User, error)

	// Role
	ListRole(context.Context, uint32, string, uint32) (*[]model.Role, error)
	GetRole(context.Context, uint32, uint32) (*model.Role, error)
	GetRoleByName(context.Context, uint32, string) (*model.Role, error)
	PutRole(ctx context.Context, r *model.Role) (*model.Role, error)
	DeleteRole(context.Context, uint32, uint32) error
	AttachRole(context.Context, uint32, uint32, uint32) (*model.UserRole, error)
	DetachRole(context.Context, uint32, uint32, uint32) error

	// Policy
	GetUserPolicy(context.Context, uint32) (*[]model.Policy, error)
	GetAdminPolicy(context.Context, uint32) (*model.Policy, error)
	ListPolicy(context.Context, uint32, string, uint32) (*[]model.Policy, error)
	GetPolicy(context.Context, uint32, uint32) (*model.Policy, error)
	GetPolicyByName(context.Context, uint32, string) (*model.Policy, error)
	PutPolicy(context.Context, *model.Policy) (*model.Policy, error)
	DeletePolicy(context.Context, uint32, uint32) error
	AttachPolicy(context.Context, uint32, uint32, uint32) (*model.RolePolicy, error)
	DetachPolicy(context.Context, uint32, uint32, uint32) error
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

	dsn := fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
		user, pass, host, conf.Port, conf.Schema)
	db, err := mimosasql.Open(dsn, conf.LogMode)
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

func (i *iamDB) userExists(ctx context.Context, userID uint32) bool {
	if _, err := i.GetUser(ctx, userID, ""); errors.Is(err, gorm.ErrRecordNotFound) {
		return false

	} else if err != nil {
		appLogger.Errorf("[userExists]DB error: user_id=%d", userID)
		return false
	}
	return true
}

func (i *iamDB) roleExists(ctx context.Context, projectID, roleID uint32) bool {
	if _, err := i.GetRole(ctx, projectID, roleID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		appLogger.Errorf("[roleExists]DB error: project_id=%d, role_id=%d", projectID, roleID)
		return false
	}
	return true
}

func (i *iamDB) policyExists(ctx context.Context, projectID, policyID uint32) bool {
	if _, err := i.GetPolicy(ctx, projectID, policyID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		appLogger.Errorf("[policyExists]DB error: project_id=%d, policy_id=%d", projectID, policyID)
		return false
	}
	return true
}
