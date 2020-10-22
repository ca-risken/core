package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type alertRepository interface {
	// Alert
	ListAlert(uint32, []string, []string, string, int64, int64) (*[]model.Alert, error)
	GetAlert(uint32, uint32) (*model.Alert, error)
	GetAlertByAlertConditionID(uint32, uint32) (*model.Alert, error)
	UpsertAlert(*model.Alert) (*model.Alert, error)
	DeleteAlert(uint32, uint32) error
	ListAlertHistory(uint32, uint32, []string, []string, int64, int64) (*[]model.AlertHistory, error)
	GetAlertHistory(uint32, uint32) (*model.AlertHistory, error)
	UpsertAlertHistory(*model.AlertHistory) (*model.AlertHistory, error)
	DeleteAlertHistory(uint32, uint32) error
	ListRelAlertFinding(uint32, uint32, uint32, int64, int64) (*[]model.RelAlertFinding, error)
	GetRelAlertFinding(uint32, uint32, uint32) (*model.RelAlertFinding, error)
	UpsertRelAlertFinding(*model.RelAlertFinding) (*model.RelAlertFinding, error)
	DeleteRelAlertFinding(uint32, uint32, uint32) error
	ListAlertCondition(uint32, []string, bool, int64, int64) (*[]model.AlertCondition, error)
	GetAlertCondition(uint32, uint32) (*model.AlertCondition, error)
	UpsertAlertCondition(*model.AlertCondition) (*model.AlertCondition, error)
	DeleteAlertCondition(uint32, uint32) error
	ListAlertRule(uint32, float32, float32, int64, int64) (*[]model.AlertRule, error)
	GetAlertRule(uint32, uint32) (*model.AlertRule, error)
	UpsertAlertRule(*model.AlertRule) (*model.AlertRule, error)
	DeleteAlertRule(uint32, uint32) error
	ListAlertCondRule(uint32, uint32, uint32, int64, int64) (*[]model.AlertCondRule, error)
	GetAlertCondRule(uint32, uint32, uint32) (*model.AlertCondRule, error)
	UpsertAlertCondRule(*model.AlertCondRule) (*model.AlertCondRule, error)
	DeleteAlertCondRule(uint32, uint32, uint32) error
	ListNotification(uint32, string, int64, int64) (*[]model.Notification, error)
	GetNotification(uint32, uint32) (*model.Notification, error)
	UpsertNotification(*model.Notification) (*model.Notification, error)
	DeleteNotification(uint32, uint32) error
	ListAlertCondNotification(uint32, uint32, uint32, int64, int64) (*[]model.AlertCondNotification, error)
	GetAlertCondNotification(uint32, uint32, uint32) (*model.AlertCondNotification, error)
	UpsertAlertCondNotification(*model.AlertCondNotification) (*model.AlertCondNotification, error)
	DeleteAlertCondNotification(uint32, uint32, uint32) error

	// forAnalyze
	ListAlertRuleByAlertConditionID(uint32, uint32) (*[]model.AlertRule, error)
	ListNotificationByAlertConditionID(uint32, uint32) (*[]model.Notification, error)
	DeactivateAlert(*model.Alert) error
	GetAlertByAlertConditionIDStatus(uint32, uint32, []string) (*model.Alert, error)
	ListFinding(uint32) (*[]model.Finding, error)
	ListFindingTag(uint32, uint64) (*[]model.FindingTag, error)
	ListEnabledAlertCondition(uint32, []uint32) (*[]model.AlertCondition, error)
	ListDisabledAlertCondition(uint32, []uint32) (*[]model.AlertCondition, error)
	GetProject(uint32) (*model.Project, error)
}

type alertDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newAlertRepository() alertRepository {
	return &alertDB{
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
