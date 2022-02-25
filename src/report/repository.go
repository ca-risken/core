package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/core/src/report/model"
	"gorm.io/gorm"
)

type reportRepository interface {
	// Report
	GetReportFinding(context.Context, uint32, []string, string, string, float32) (*[]model.ReportFinding, error)
	GetReportFindingAll(context.Context, []string, string, string, float32) (*[]model.ReportFinding, error)
	CollectReportFinding(ctx context.Context) error
}

type reportDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newReportRepository(conf *DBConfig) reportRepository {
	return &reportDB{
		Master: initDB(conf, true),
		Slave:  initDB(conf, false),
	}
}

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
