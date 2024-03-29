package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/DATA-DOG/go-sqlmock"
	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/cenkalti/backoff/v4"
)

type Client struct {
	Master  *gorm.DB
	Slave   *gorm.DB
	logger  logging.Logger
	retryer backoff.BackOff
}

func NewClient(conf *Config, l logging.Logger) (*Client, error) {
	ctx := context.Background()
	m, err := connect(conf, true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database of master: %w", err)
	}
	l.Info(ctx, "Connected to Database of master.")

	s, err := connect(conf, false)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database of slave: %w", err)
	}
	l.Info(ctx, "Connected to Database of slave.")
	return &Client{
		Master:  m,
		Slave:   s,
		logger:  l,
		retryer: backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
	}, nil
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
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open gorm, error: %+w", err)
	}
	return &Client{
		Master:  gormDB,
		Slave:   gormDB,
		logger:  logging.NewLogger(),
		retryer: backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 2),
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

func (c *Client) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		c.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
