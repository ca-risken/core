package main

import (
	"context"
	"fmt"

	mimosasql "github.com/CyberAgent/mimosa-common/pkg/database/sql"
	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/gorm"
)

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

type projectRepository interface {
	ListProject(ctx context.Context, userID, projectID uint32, name string) (*[]projectWithTag, error)
	CreateProject(ctx context.Context, name string) (*model.Project, error)
	UpdateProject(ctx context.Context, projectID uint32, name string) (*model.Project, error)

	ListProjectTag(ctx context.Context, projectID uint32) (*[]model.ProjectTag, error)
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
