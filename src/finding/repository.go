package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type findingRepository interface {
	// Finding
	ListFinding(*finding.ListFindingRequest) (*[]model.Finding, error)
	ListFindingCount(req *finding.ListFindingRequest) (uint32, error)
	GetFinding(uint32, uint64) (*model.Finding, error)
	GetFindingByDataSource(uint32, string, string) (*model.Finding, error)
	UpsertFinding(*model.Finding) (*model.Finding, error)
	DeleteFinding(uint32, uint64) error
	ListFindingTag(param *finding.ListFindingTagRequest) (*[]model.FindingTag, error)
	ListFindingTagCount(param *finding.ListFindingTagRequest) (uint32, error)
	ListFindingTagName(param *finding.ListFindingTagNameRequest) (*[]tagName, error)
	ListFindingTagNameCount(param *finding.ListFindingTagNameRequest) (uint32, error)
	GetFindingTagByKey(uint32, uint64, string) (*model.FindingTag, error)
	GetFindingTagByID(uint32, uint64) (*model.FindingTag, error)
	TagFinding(*model.FindingTag) (*model.FindingTag, error)
	UntagFinding(uint32, uint64) error
	GetPendFinding(projectID uint32, findingID uint64) (*model.PendFinding, error)
	UpsertPendFinding(findingID uint64, projectID uint32) (*model.PendFinding, error)
	DeletePendFinding(projectID uint32, findingID uint64) error

	// Resource
	ListResource(*finding.ListResourceRequest) (*[]model.Resource, error)
	ListResourceCount(req *finding.ListResourceRequest) (uint32, error)
	GetResource(uint32, uint64) (*model.Resource, error)
	GetResourceByName(uint32, string) (*model.Resource, error)
	UpsertResource(*model.Resource) (*model.Resource, error)
	DeleteResource(uint32, uint64) error
	ListResourceTag(param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error)
	ListResourceTagCount(param *finding.ListResourceTagRequest) (uint32, error)
	ListResourceTagName(param *finding.ListResourceTagNameRequest) (*[]tagName, error)
	ListResourceTagNameCount(param *finding.ListResourceTagNameRequest) (uint32, error)
	GetResourceTagByKey(uint32, uint64, string) (*model.ResourceTag, error)
	GetResourceTagByID(uint32, uint64) (*model.ResourceTag, error)
	TagResource(*model.ResourceTag) (*model.ResourceTag, error)
	UntagResource(uint32, uint64) error
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

type tagName struct {
	Tag string
}
