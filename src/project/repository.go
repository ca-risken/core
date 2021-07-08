package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
		return nil
	}
	if conf.LogMode {
		db.Logger.LogMode(logger.Info)
	}
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

type projectRepository interface {
	ListProject(userID, projectID uint32, name string) (*[]projectWithTag, error)
	CreateProject(name string) (*model.Project, error)
	UpdateProject(projectID uint32, name string) (*model.Project, error)

	ListProjectTag(projectID uint32) (*[]model.ProjectTag, error)
	TagProject(projectID uint32, tag, color string) (*model.ProjectTag, error)
	UntagProject(projectID uint32, tag string) error
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
