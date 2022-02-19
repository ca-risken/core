package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/core/src/project/model"
	"gorm.io/gorm"
)

type DBConfig struct {
	MasterHost     string
	MasterUser     string
	MasterPassword string
	SlaveHost      string
	SlaveUser      string
	SlavePassword  string

	Schema        string
	Port          int
	LogMode       bool
	MaxConnection int
}

func initDB(conf *DBConfig, isMaster bool) *gorm.DB {
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

func newProjectRepository(conf *DBConfig) projectRepository {
	return &projectDB{
		Master: initDB(conf, true),
		Slave:  initDB(conf, false),
	}
}
