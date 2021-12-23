package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/core/src/project/model"
	"github.com/gassara-kys/envconfig"
	"gorm.io/gorm"
)

type dbConfig struct {
	MasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	MasterUser     string `split_words:"true" default:"hoge"`
	MasterPassword string `split_words:"true" default:"moge"`
	SlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	SlaveUser      string `split_words:"true" default:"hoge"`
	SlavePassword  string `split_words:"true" default:"moge"`

	Schema        string `required:"true"    default:"mimosa"`
	Port          int    `required:"true"    default:"3306"`
	LogMode       bool   `split_words:"true" default:"false"`
	MaxConnection int    `split_words:"true" default:"10"`
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
	db, err := mimosasql.Open(dsn, conf.LogMode, conf.MaxConnection)
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

type projectRepository interface {
	ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]projectWithTag, error)
	CreateProject(ctx context.Context, name string) (*model.Project, error)
	UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error)

	TagProject(ctx context.Context, projectID uint32, tag, color string) (*model.ProjectTag, error)
	UntagProject(ctx context.Context, projectID uint32, tag string) error
}

type projectDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newProjectRepository() projectRepository {
	return &projectDB{
		Master: initDB(true),
		Slave:  initDB(false),
	}
}
