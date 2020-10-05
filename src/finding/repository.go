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

type findingRepository interface {
	// Finding
	ListFinding(*finding.ListFindingRequest) (*[]model.Finding, error)
	GetFinding(uint32, uint64) (*model.Finding, error)
	GetFindingByDataSource(uint32, string, string) (*model.Finding, error)
	UpsertFinding(*model.Finding) (*model.Finding, error)
	DeleteFinding(uint32, uint64) error
	ListFindingTag(uint32, uint64) (*[]model.FindingTag, error)
	GetFindingTagByKey(uint32, uint64, string) (*model.FindingTag, error)
	GetFindingTagByID(uint32, uint64) (*model.FindingTag, error)
	TagFinding(*model.FindingTag) (*model.FindingTag, error)
	UntagFinding(uint32, uint64) error

	// Resource
	ListResource(*finding.ListResourceRequest) (*[]model.Resource, error)
	GetResource(uint32, uint64) (*model.Resource, error)
	GetResourceByName(uint32, string) (*model.Resource, error)
	UpsertResource(*model.Resource) (*model.Resource, error)
	DeleteResource(uint32, uint64) error
	ListResourceTag(uint32, uint64) (*[]model.ResourceTag, error)
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

func (f *findingDB) ListFinding(req *finding.ListFindingRequest) (*[]model.Finding, error) {
	query := `
select
  *
from
  finding
where
  project_id = ?
  and score between ? and ?
  and updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, req.FromScore, req.ToScore, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if len(req.DataSource) != 0 {
		query += " and data_source in (?)"
		params = append(params, req.DataSource)
	}
	if len(req.ResourceName) != 0 {
		query += " and resource_name in (?)"
		params = append(params, req.ResourceName)
	}
	var data []model.Finding
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFinding = `select * from finding where project_id = ? and finding_id = ?`

func (f *findingDB) GetFinding(projectID uint32, findingID uint64) (*model.Finding, error) {
	var data model.Finding
	if err := f.Slave.Raw(selectGetFinding, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertFinding = `
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
  data=VALUES(data)
`

func (f *findingDB) UpsertFinding(data *model.Finding) (*model.Finding, error) {
	if err := f.Master.Exec(insertUpsertFinding,
		data.FindingID, data.Description, data.DataSource, data.DataSourceID, data.ResourceName,
		data.ProjectID, data.OriginalScore, data.Score, data.Data).Error; err != nil {
		return nil, err
	}
	return f.GetFindingByDataSource(data.ProjectID, data.DataSource, data.DataSourceID)
}

const selectGetFindingByDataSource = `select * from finding where project_id = ? and data_source = ? and data_source_id = ?`

func (f *findingDB) GetFindingByDataSource(projectID uint32, dataSource, dataSourceID string) (*model.Finding, error) {
	var result model.Finding
	if err := f.Slave.Raw(selectGetFindingByDataSource,
		projectID, dataSource, dataSourceID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

const insertUpsertResource = `
INSERT INTO resource
  (resource_id, resource_name, project_id)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  resource_name=VALUES(resource_name),
  project_id=VALUES(project_id);
`

func (f *findingDB) UpsertResource(data *model.Resource) (*model.Resource, error) {
	if err := f.Master.Exec(insertUpsertResource,
		data.ResourceID, data.ResourceName, data.ProjectID).Error; err != nil {
		return nil, err
	}
	return f.GetResourceByName(data.ProjectID, data.ResourceName)
}

const selectGetResourceByName = `select * from resource where project_id = ? and resource_name = ?`

func (f *findingDB) GetResourceByName(projectID uint32, resourceName string) (*model.Resource, error) {
	var data model.Resource
	if err := f.Slave.Raw(selectGetResourceByName, projectID, resourceName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDeleteFinding = `delete from finding where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteFinding(projectID uint32, findingID uint64) error {
	if err := f.Master.Exec(deleteDeleteFinding, projectID, findingID).Error; err != nil {
		return err
	}
	return f.DeleteTagByFindingID(projectID, findingID)
}

const deleteDeleteTagByFindingID = `delete from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) DeleteTagByFindingID(projectID uint32, findingID uint64) error {
	return f.Master.Exec(deleteDeleteTagByFindingID, projectID, findingID).Error
}

const selectListFindingTag = `select * from finding_tag where project_id = ? and finding_id = ?`

func (f *findingDB) ListFindingTag(projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	var data []model.FindingTag
	if err := f.Slave.Raw(selectListFindingTag, projectID, findingID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertTagFinding = `
INSERT INTO finding_tag
  (finding_tag_id, finding_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)
`

func (f *findingDB) TagFinding(tag *model.FindingTag) (*model.FindingTag, error) {
	if err := f.Master.Exec(insertTagFinding,
		tag.FindingTagID, tag.FindingID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return f.GetFindingTagByKey(tag.ProjectID, tag.FindingID, tag.Tag)
}

const selectGetFindingTagByKey = `select * from finding_tag where project_id = ? and finding_id = ? and tag = ?`

func (f *findingDB) GetFindingTagByKey(projectID uint32, findingID uint64, tag string) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.Slave.Raw(selectGetFindingTagByKey, projectID, findingID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetFindingTagByID = `select * from finding_tag where project_id = ? and finding_tag_id = ?`

func (f *findingDB) GetFindingTagByID(projectID uint32, findingTagID uint64) (*model.FindingTag, error) {
	var data model.FindingTag
	if err := f.Slave.Raw(selectGetFindingTagByID, projectID, findingTagID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteUntagFinding = `delete from finding_tag where project_id = ? and finding_tag_id = ?`

func (f *findingDB) UntagFinding(projectID uint32, findingTagID uint64) error {
	return f.Master.Exec(deleteUntagFinding, projectID, findingTagID).Error
}

func (f *findingDB) ListResource(req *finding.ListResourceRequest) (*[]model.Resource, error) {
	query := `
select 
  r.*
from
  resource r
  left outer join finding f using(resource_name)
where
  r.project_id = ?
  and r.updated_at between ? and ?
`
	var params []interface{}
	params = append(params, req.ProjectId, time.Unix(req.FromAt, 0), time.Unix(req.ToAt, 0))
	if len(req.ResourceName) != 0 {
		query += " and r.resource_name in (?)"
		params = append(params, req.ResourceName)
	}
	query += " group by r.resource_id having sum(COALESCE(f.score, 0)) between ? and ?"
	params = append(params, req.FromSumScore, req.ToSumScore)

	var data []model.Resource
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResource = `select * from resource where project_id = ? and resource_id = ?`

func (f *findingDB) GetResource(projectID uint32, resourceID uint64) (*model.Resource, error) {
	var data model.Resource
	if err := f.Slave.Raw(selectGetResource, projectID, resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteDeleteResource = `delete from resource where project_id = ? and resource_id = ?`

func (f *findingDB) DeleteResource(projectID uint32, resourceID uint64) error {
	if err := f.Master.Exec(deleteDeleteResource, projectID, resourceID).Error; err != nil {
		return err
	}
	return f.DeleteTagByResourceID(projectID, resourceID)
}

const deleteDeleteTagByResourceID = `delete from resource_tag where project_id = ? and resource_id = ?`

func (f *findingDB) DeleteTagByResourceID(projectID uint32, resourceID uint64) error {
	return f.Master.Exec(deleteDeleteTagByResourceID, projectID, resourceID).Error
}

const selectListResourceTag = `select * from resource_tag where project_id = ? and resource_id = ?`

func (f *findingDB) ListResourceTag(projectID uint32, resourceID uint64) (*[]model.ResourceTag, error) {
	var data []model.ResourceTag
	if err := f.Slave.Raw(selectListResourceTag, projectID, resourceID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResourceTagByKey = `select * from resource_tag where project_id = ? and resource_id = ? and tag = ?`

func (f *findingDB) GetResourceTagByKey(projectID uint32, resourceID uint64, tag string) (*model.ResourceTag, error) {
	var data model.ResourceTag
	if err := f.Slave.Raw(selectGetResourceTagByKey, projectID, resourceID, tag).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetResourceTagByID = `select * from resource_tag where project_id = ? and resource_tag_id = ?`

func (f *findingDB) GetResourceTagByID(projectID uint32, resourceID uint64) (*model.ResourceTag, error) {
	var data model.ResourceTag
	if err := f.Slave.Raw(selectGetResourceTagByID, projectID, resourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertTagResource = `
INSERT INTO resource_tag
  (resource_tag_id, resource_id, project_id, tag)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  tag=VALUES(tag)
`

func (f *findingDB) TagResource(tag *model.ResourceTag) (*model.ResourceTag, error) {
	if err := f.Master.Exec(insertTagResource,
		tag.ResourceTagID, tag.ResourceID, tag.ProjectID, tag.Tag).Error; err != nil {
		return nil, err
	}
	return f.GetResourceTagByKey(tag.ProjectID, tag.ResourceID, tag.Tag)
}

const deleteUntagResource = `delete from resource_tag where project_id = ? and resource_tag_id = ?`

func (f *findingDB) UntagResource(projectID uint32, resourceTagID uint64) error {
	return f.Master.Exec(deleteUntagResource, projectID, resourceTagID).Error
}
