package main

import (
	"context"
	"fmt"

	mimosasql "github.com/CyberAgent/mimosa-common/pkg/database/sql"
	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/gorm"
)

type findingRepository interface {
	// Finding
	ListFinding(context.Context, *finding.ListFindingRequest) (*[]model.Finding, error)
	BatchListFinding(context.Context, *finding.BatchListFindingRequest) (*[]model.Finding, error)
	ListFindingCount(
		ctx context.Context,
		projectID uint32,
		fromScore, toScore float32,
		fromAt, toAt int64,
		findingID uint64,
		dataSources, resourceNames, tags []string,
		status finding.FindingStatus) (int64, error)
	GetFinding(context.Context, uint32, uint64) (*model.Finding, error)
	GetFindingByDataSource(context.Context, uint32, string, string) (*model.Finding, error)
	UpsertFinding(context.Context, *model.Finding) (*model.Finding, error)
	DeleteFinding(context.Context, uint32, uint64) error
	ListFindingTag(ctx context.Context, param *finding.ListFindingTagRequest) (*[]model.FindingTag, error)
	ListFindingTagCount(ctx context.Context, param *finding.ListFindingTagRequest) (int64, error)
	ListFindingTagName(ctx context.Context, param *finding.ListFindingTagNameRequest) (*[]tagName, error)
	ListFindingTagNameCount(ctx context.Context, param *finding.ListFindingTagNameRequest) (int64, error)
	GetFindingTagByKey(context.Context, uint32, uint64, string) (*model.FindingTag, error)
	GetFindingTagByID(context.Context, uint32, uint64) (*model.FindingTag, error)
	TagFinding(context.Context, *model.FindingTag) (*model.FindingTag, error)
	UntagFinding(context.Context, uint32, uint64) error
	ClearScoreFinding(ctx context.Context, req *finding.ClearScoreRequest) error

	// Resource
	ListResource(context.Context, *finding.ListResourceRequest) (*[]model.Resource, error)
	ListResourceCount(ctx context.Context, req *finding.ListResourceRequest) (int64, error)
	GetResource(context.Context, uint32, uint64) (*model.Resource, error)
	GetResourceByName(context.Context, uint32, string) (*model.Resource, error)
	UpsertResource(context.Context, *model.Resource) (*model.Resource, error)
	DeleteResource(context.Context, uint32, uint64) error
	ListResourceTag(ctx context.Context, param *finding.ListResourceTagRequest) (*[]model.ResourceTag, error)
	ListResourceTagCount(ctx context.Context, param *finding.ListResourceTagRequest) (int64, error)
	ListResourceTagName(ctx context.Context, param *finding.ListResourceTagNameRequest) (*[]tagName, error)
	ListResourceTagNameCount(ctx context.Context, param *finding.ListResourceTagNameRequest) (int64, error)
	GetResourceTagByKey(context.Context, uint32, uint64, string) (*model.ResourceTag, error)
	GetResourceTagByID(context.Context, uint32, uint64) (*model.ResourceTag, error)
	TagResource(context.Context, *model.ResourceTag) (*model.ResourceTag, error)
	UntagResource(context.Context, uint32, uint64) error

	// PendFinding
	GetPendFinding(ctx context.Context, projectID uint32, findingID uint64) (*model.PendFinding, error)
	UpsertPendFinding(ctx context.Context, pend *finding.PendFindingForUpsert) (*model.PendFinding, error)
	DeletePendFinding(ctx context.Context, projectID uint32, findingID uint64) error

	// FindingSetting
	ListFindingSetting(ctx context.Context, req *finding.ListFindingSettingRequest) (*[]model.FindingSetting, error)
	GetFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) (*model.FindingSetting, error)
	GetFindingSettingByResource(ctx context.Context, projectID uint32, resourceName string) (*model.FindingSetting, error)
	UpsertFindingSetting(ctx context.Context, data *model.FindingSetting) (*model.FindingSetting, error)
	DeleteFindingSetting(ctx context.Context, projectID uint32, findingSettingID uint32) error
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
	db, err := mimosasql.Open(dsn, conf.LogMode)
	if err != nil {
		appLogger.Fatalf("Failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}
	appLogger.Infof("Connected to Database. isMaster: %t", isMaster)
	return db
}

type tagName struct {
	Tag string
}
