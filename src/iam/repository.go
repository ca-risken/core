package main

import (
	"context"
	"errors"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/core/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type iamRepository interface {
	// User
	ListUser(ctx context.Context, activated bool, projectID uint32, name string, userID uint32) (*[]model.User, error)
	GetUser(ctx context.Context, userID uint32, sub string) (*model.User, error)
	GetUserBySub(ctx context.Context, sub string) (*model.User, error)
	PutUser(ctx context.Context, u *model.User) (*model.User, error)
	GetActiveUserCount(ctx context.Context) (*int, error)

	// Role
	ListRole(ctx context.Context, projectID uint32, name string, userID uint32, accessTokenID uint32) (*[]model.Role, error)
	GetRole(ctx context.Context, projectID, roleID uint32) (*model.Role, error)
	GetRoleByName(ctx context.Context, projectID uint32, name string) (*model.Role, error)
	PutRole(ctx context.Context, r *model.Role) (*model.Role, error)
	DeleteRole(ctx context.Context, projectID, roleID uint32) error
	AttachRole(ctx context.Context, projectID, roleID, userID uint32) (*model.UserRole, error)
	AttachAdminRole(ctx context.Context, userID uint32) error
	DetachRole(ctx context.Context, projectID, roleID, userID uint32) error

	// Policy
	GetUserPolicy(ctx context.Context, userID uint32) (*[]model.Policy, error)
	GetTokenPolicy(ctx context.Context, accessTokenID uint32) (*[]model.Policy, error)
	GetAdminPolicy(ctx context.Context, userID uint32) (*model.Policy, error)
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

func (i *iamDB) accessTokenExists(ctx context.Context, projectID, accessTokenID uint32) bool {
	if _, err := i.GetAccessTokenByID(ctx, projectID, accessTokenID); errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		appLogger.Errorf("[accessTokenExists]DB error: project_id=%d, access_token_id=%d", projectID, accessTokenID)
		return false
	}
	return true
}

func convertZeroValueToNull(input interface{}) interface{} {
	if input == nil || zero.IsZeroVal(input) {
		return gorm.Expr("NULL")
	}
	return input
}
