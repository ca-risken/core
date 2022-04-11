package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/src/alert/model"
	mysqldriver "github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type alertRepository interface {
	// Alert
	ListAlert(context.Context, uint32, []string, []string, string, int64, int64) (*[]model.Alert, error)
	GetAlert(context.Context, uint32, uint32) (*model.Alert, error)
	UpsertAlert(context.Context, *model.Alert) (*model.Alert, error)
	DeleteAlert(context.Context, uint32, uint32) error
	ListAlertHistory(context.Context, uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error)
	GetAlertHistory(context.Context, uint32, uint32) (*model.AlertHistory, error)
	UpsertAlertHistory(context.Context, *model.AlertHistory) (*model.AlertHistory, error)
	DeleteAlertHistory(context.Context, uint32, uint32) error
	ListRelAlertFinding(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error)
	GetRelAlertFinding(context.Context, uint32, uint32, uint32) (*model.RelAlertFinding, error)
	UpsertRelAlertFinding(context.Context, *model.RelAlertFinding) (*model.RelAlertFinding, error)
	DeleteRelAlertFinding(context.Context, uint32, uint32, uint32) error
	ListAlertCondition(context.Context, uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error)
	GetAlertCondition(context.Context, uint32, uint32) (*model.AlertCondition, error)
	UpsertAlertCondition(context.Context, *model.AlertCondition) (*model.AlertCondition, error)
	DeleteAlertCondition(context.Context, uint32, uint32) error
	ListAlertRule(context.Context, uint32, float32, float32, int64, int64) (*[]model.AlertRule, error)
	GetAlertRule(context.Context, uint32, uint32) (*model.AlertRule, error)
	UpsertAlertRule(context.Context, *model.AlertRule) (*model.AlertRule, error)
	DeleteAlertRule(context.Context, uint32, uint32) error
	ListAlertCondRule(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error)
	GetAlertCondRule(context.Context, uint32, uint32, uint32) (*model.AlertCondRule, error)
	UpsertAlertCondRule(context.Context, *model.AlertCondRule) (*model.AlertCondRule, error)
	DeleteAlertCondRule(context.Context, uint32, uint32, uint32) error
	ListNotification(context.Context, uint32, string, int64, int64) (*[]model.Notification, error)
	GetNotification(context.Context, uint32, uint32) (*model.Notification, error)
	UpsertNotification(context.Context, *model.Notification) (*model.Notification, error)
	DeleteNotification(context.Context, uint32, uint32) error
	ListAlertCondNotification(context.Context, uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error)
	GetAlertCondNotification(context.Context, uint32, uint32, uint32) (*model.AlertCondNotification, error)
	UpsertAlertCondNotification(context.Context, *model.AlertCondNotification) (*model.AlertCondNotification, error)
	DeleteAlertCondNotification(context.Context, uint32, uint32, uint32) error

	// forAnalyze
	ListAlertRuleByAlertConditionID(context.Context, uint32, uint32) (*[]model.AlertRule, error)
	DeactivateAlert(context.Context, *model.Alert) error
	GetAlertByAlertConditionIDStatus(context.Context, uint32, uint32, []string) (*model.Alert, error)
	ListEnabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)
	ListDisabledAlertCondition(context.Context, uint32, []uint32) (*[]model.AlertCondition, error)
}

type alertDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newAlertRepository(conf *DBConfig) alertRepository {
	return &alertDB{
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
	sqltrace.Register("mysql", &mysqldriver.MySQLDriver{}, sqltrace.WithAnalytics(true))
	dbConn, err := sqltrace.Open("mysql", dsn)
	if err != nil {
		appLogger.Fatalf("Failed to open DB connection. err: %+v", err)
	}
	dbConn.SetMaxOpenConns(conf.MaxConnection)
	dbConn.SetConnMaxLifetime(time.Duration(conf.MaxConnection/2) * time.Second)

	db, err := gormtrace.Open(mysql.New(mysql.Config{
		Conn: dbConn,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}
	if conf.LogMode {
		db.Logger.LogMode(logger.Info)
	}

	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}
