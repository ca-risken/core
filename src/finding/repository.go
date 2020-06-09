package main

import (
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type findingRepoInterface interface {
	ListFinding(*finding.ListFindingRequest) (*[]findingIds, error)
	GetFinding(uint64) (*model.Finding, error)
	UpsertFinding(*model.Finding) (*model.Finding, error)
	GetFindingByDataSource(uint32, string, string) (*model.Finding, error)
	UpsertResource(*model.Resource) (*model.Resource, error)
	GetResourceByName(uint32, string) (*model.Resource, error)
	DeleteFinding(uint64) error
	ListFindingTag(uint64) (*[]model.FindingTag, error)
	TagFinding(*model.FindingTag) (*model.FindingTag, error)
	GetFindingTagByKey(uint64, string) (*model.FindingTag, error)
}

type findingRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newFindingRepository() findingRepoInterface {
	repo := findingRepository{}
	repo.MasterDB = initDB(true)
	repo.SlaveDB = initDB(false)
	return &repo
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

type findingIds struct {
	FindingID uint64 `gorm:"column:finding_id"`
}

func (f *findingRepository) ListFinding(req *finding.ListFindingRequest) (*[]findingIds, error) {
	query := `
select
	finding_id
from
	finding
where
	score between ? and ?
	and 
	updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.FromScore, req.ToScore, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if len(req.ProjectId) != 0 {
		query += " and project_id in (?)"
		params = append(params, req.ProjectId)
	}
	if len(req.DataSource) != 0 {
		query += " and data_source in (?)"
		params = append(params, req.DataSource)
	}
	if len(req.ResourceName) != 0 {
		query += " and resource_name in (?)"
		params = append(params, req.ResourceName)
	}

	var data []findingIds
	if err := f.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingRepository) GetFinding(findingID uint64) (*model.Finding, error) {
	var data model.Finding
	if err := f.SlaveDB.Raw(`select * from finding where finding_id = ?`, findingID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingRepository) UpsertFinding(data *model.Finding) (*model.Finding, error) {
	err := f.MasterDB.Exec(`
INSERT INTO finding
	(finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	description=VALUES(description),
	resource_name=VALUES(resource_name),
	project_id=VALUES(project_id),
	original_score=VALUES(original_score),
	score=VALUES(score),
	data=VALUES(data);
`,
		data.FindingID, data.Description, data.DataSource, data.DataSourceID, data.ResourceName,
		data.ProjectID, data.OriginalScore, data.Score, data.Data).Error
	if err != nil {
		return nil, err
	}

	updated, err := f.GetFindingByDataSource(data.ProjectID, data.DataSource, data.DataSourceID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (f *findingRepository) GetFindingByDataSource(projectID uint32, dataSource, dataSourceID string) (*model.Finding, error) {
	var result model.Finding
	if err := f.SlaveDB.Raw(`select * from finding where project_id = ? and data_source = ? and data_source_id = ?`,
		projectID, dataSource, dataSourceID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (f *findingRepository) UpsertResource(data *model.Resource) (*model.Resource, error) {
	err := f.MasterDB.Exec(`
INSERT INTO resource
	(resource_id, resource_name, project_id)
VALUES
	(?, ?, ?)
ON DUPLICATE KEY UPDATE
	resource_name=VALUES(resource_name),
	project_id=VALUES(project_id);
`,
		data.ResourceID, data.ResourceName, data.ProjectID).Error
	if err != nil {
		return nil, err
	}

	updated, err := f.GetResourceByName(data.ProjectID, data.ResourceName)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (f *findingRepository) GetResourceByName(projectID uint32, resourceName string) (*model.Resource, error) {
	var data model.Resource
	if err := f.SlaveDB.Raw(`select * from resource where project_id = ? and resource_name = ?`,
		projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingRepository) DeleteFinding(findingID uint64) error {
	if err := f.MasterDB.Exec(`delete from finding where finding_id = ?`, findingID).Error; err != nil {
		return err
	}
	return nil
}

func (f *findingRepository) ListFindingTag(findingID uint64) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := f.SlaveDB.Raw(`select * from finding_tag where finding_id = ?`, findingID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *findingRepository) TagFinding(tag *model.FindingTag) (*model.FindingTag, error) {
	err := f.MasterDB.Exec(`
INSERT INTO finding_tag
	(finding_tag_id, finding_id, tag_key, tag_value)
VALUES
	(?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	tag_value=VALUES(tag_value);
`,
		tag.FindingTagID, tag.FindingID, tag.TagKey, tag.TagValue).Error
	if err != nil {
		return nil, err
	}
	updated, err := f.GetFindingTagByKey(tag.FindingID, tag.TagKey)
	if err != nil {
		return nil, err
	}
	return updated, err
}

func (f *findingRepository) GetFindingTagByKey(findingID uint64, tagKey string) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.SlaveDB.Raw(`select * from finding_tag where finding_id = ? and tag_key = ?`, findingID, tagKey).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
