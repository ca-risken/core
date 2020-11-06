package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
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

type projectRepository interface {
	ListProject(uint32, uint32, string) (*[]model.Project, error)
	CreateProject(string) (*model.Project, error)
	UpdateProject(uint32, string) (*model.Project, error)
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

func (p *projectDB) ListProject(userID, projectID uint32, name string) (*[]model.Project, error) {
	query := `select p.* from project p where 1 = 1` // プログラム構造をシンプルに保つために必ずtrueとなるwhere条件を入れとく（and条件のみなので一旦これで）
	var params []interface{}
	if !zero.IsZeroVal(userID) {
		query += " and exists (select * from user_role ur where ur.project_id = p.project_id and user_id = ?)"
		params = append(params, userID)
	}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	data := []model.Project{}
	if err := p.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetProjectByName = `select * from project where name = ?`

func (p *projectDB) GetProjectByName(name string) (*model.Project, error) {
	var data model.Project
	if err := p.Master.Raw(selectGetProjectByName, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertCreateProject = `insert into project(name) values(?)`

func (p *projectDB) CreateProject(name string) (*model.Project, error) {
	// Handring duplicated name error
	if pr, err := p.GetProjectByName(name); err == nil {
		return nil, fmt.Errorf("Project name already registerd: project_id=%d, name=%s", pr.ProjectID, pr.Name)
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Could not get project data: err=%+v", err)
	}
	if err := p.Master.Exec(insertCreateProject, name).Error; err != nil {
		return nil, err
	}
	return p.GetProjectByName(name)
}

const updateUpdateProject = `update project set name = ? where project_id = ?`

func (p *projectDB) UpdateProject(projectID uint32, name string) (*model.Project, error) {
	if err := p.Master.Exec(updateUpdateProject, name, projectID).Error; err != nil {
		return nil, err
	}
	return p.GetProjectByName(name)
}
