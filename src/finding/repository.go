package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type findingRepository interface {
	// Finding
	ListFinding(*finding.ListFindingRequest) (*[]model.Finding, error)
	BatchListFinding(*finding.BatchListFindingRequest) (*[]model.Finding, error)
	ListFindingCount(
		projectID uint32,
		fromScore, toScore float32,
		fromAt, toAt int64,
		findingID uint64,
		dataSources, resourceNames, tags []string,
		status finding.FindingStatus) (int64, error)
	GetFinding(uint32, uint64) (*model.Finding, error)
	GetFindingByDataSource(uint32, string, string) (*model.Finding, error)
	UpsertFinding(*model.Finding) (*model.Finding, error)
	DeleteFinding(uint32, uint64) error
	ListFindingTag(param *finding.ListFindingTagRequest) (*[]model.FindingTag, error)
	ListFindingTagCount(param *finding.ListFindingTagRequest) (int64, error)
	ListFindingTagName(param *finding.ListFindingTagNameRequest) (*[]tagName, error)
	ListFindingTagNameCount(param *finding.ListFindingTagNameRequest) (int64, error)
	GetFindingTagByKey(uint32, uint64, string) (*model.FindingTag, error)
	GetFindingTagByID(uint32, uint64) (*model.FindingTag, error)
	TagFinding(*model.FindingTag) (*model.FindingTag, error)
	UntagFinding(uint32, uint64) error

	// Resource
	ListResource(*finding.ListResourceRequest) (*[]model.Resource, error)
	ListResourceCount(req *finding.ListResourceRequest) (int64, error)
	GetResource(uint32, uint64) (*model.Resource, error)
	GetResourceByName(uint32, string) (*model.Resource, error)
	UpsertResource(*model.Resource) (*model.Resource, error)
	DeleteResource(uint32, uint64) error
	ListResourceTag(param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error)
	ListResourceTagCount(param *finding.ListResourceTagRequest) (int64, error)
	ListResourceTagName(param *finding.ListResourceTagNameRequest) (*[]tagName, error)
	ListResourceTagNameCount(param *finding.ListResourceTagNameRequest) (int64, error)
	GetResourceTagByKey(uint32, uint64, string) (*model.ResourceTag, error)
	GetResourceTagByID(uint32, uint64) (*model.ResourceTag, error)
	TagResource(*model.ResourceTag) (*model.ResourceTag, error)
	UntagResource(uint32, uint64) error

	// PendFinding
	GetPendFinding(projectID uint32, findingID uint64) (*model.PendFinding, error)
	UpsertPendFinding(pend *finding.PendFindingForUpsert) (*model.PendFinding, error)
	DeletePendFinding(projectID uint32, findingID uint64) error

	// FindingSetting
	ListFindingSetting(req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error)
	GetFindingSetting(projectID uint32, findingSettingID uint32) (*model.FindingSetting, error)
	GetFindingSettingByResource(projectID uint32, resourceName string) (*model.FindingSetting, error)
	UpsertFindingSetting(data *model.FindingSetting) (*model.FindingSetting, error)
	DeleteFindingSetting(projectID uint32, findingSettingID uint32) error
}

type findingDB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func newFindingRepository() findingRepository {
	return &findingDB{
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

type tagName struct {
	Tag string
}
