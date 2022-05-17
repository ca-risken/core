package db

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/common/pkg/logging"
)

type Client struct {
	Master *gorm.DB
	Slave  *gorm.DB
	logger logging.Logger
}

func NewClient(conf *Config, l logging.Logger) *Client {
	ctx := context.Background()
	m, err := connect(conf, true)
	if err != nil {
		l.Fatalf(ctx, "failed to connect database: %w", err)
	}
	l.Infof(ctx, "Connected to Database. isMaster: %t", true)

	s, err := connect(conf, false)
	if err != nil {
		l.Fatalf(ctx, "failed to connect database: %w", err)
	}
	l.Infof(ctx, "Connected to Database. isMaster: %t", false)

	return &Client{
		Master: m,
		Slave:  s,
		logger: l,
	}
}

func newMockClient() (*Client, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open mock sql db, error: %+w", err)
	}
	if sqlDB == nil {
		return nil, nil, fmt.Errorf("failed to create mock db, db: %+v, mock: %+v", sqlDB, mock)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open gorm, error: %+w", err)
	}
	return &Client{
		Master: gormDB,
		Slave:  gormDB,
	}, mock, nil
}

type Config struct {
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

func connect(conf *Config, isMaster bool) (*gorm.DB, error) {
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
		return nil, fmt.Errorf("failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}

	return db, nil
}
