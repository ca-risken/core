package main

import (
	"fmt"

	"github.com/cloudflare/cfssl/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/tsenart/nap"
)

type findingRepoInterface interface {
	List() (*[]string, error)
}

type findingRepository struct {
	DB *nap.DB
}

func newFindingRepository() findingRepoInterface {
	return &findingRepository{
		DB: initDB(),
	}
}

type dbConfig struct {
	MasterHost     string `split_words:"true" required:"true"`
	MasterUser     string `split_words:"true" required:"true"`
	MasterPassword string `split_words:"true" required:"true"`
	SlaveHost      string `split_words:"true"`
	SlaveUser      string `split_words:"true"`
	SlavePassword  string `split_words:"true"`

	Schema string `default:"mimosa"`
	Port   int    `default:"3306" required:"true"`
}

func initDB() *nap.DB {
	conf := &dbConfig{}
	if err := envconfig.Process("DB", conf); err != nil {
		appLogger.Fatalf("Failed to load DB config. err: %+v", err)
	}

	dsns := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
		conf.MasterUser, conf.MasterPassword, conf.MasterHost, conf.Port, conf.Schema)

	if conf.SlaveHost != "" {
		dsns += ";"
		dsns += fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
			conf.SlaveUser, conf.SlavePassword, conf.SlaveHost, conf.Port, conf.Schema)
	}
	db, err := nap.Open("mysql", dsns)
	if err != nil {
		appLogger.Fatalf("Failed to open DB. err: %+v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Some physical database is unreachable: %s", err)
	}
	appLogger.Info("Connected to Database.")
	return db
}

func (f *findingRepository) List() (*[]string, error) {
	return &[]string{"0000000001", "0000000002", "0000000003"}, nil
}
