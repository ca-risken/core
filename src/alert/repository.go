package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/core/src/alert/model"
	"gorm.io/gorm"
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
	db, err := mimosasql.Open(dsn, conf.LogMode, conf.MaxConnection)
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}
