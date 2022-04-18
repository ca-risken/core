package main

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMockFindingDB() (*findingDB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open mock sql db, error: %+w", err)
	}
	if sqlDB == nil || mock == nil {
		return nil, nil, fmt.Errorf("Failed to create mock db, db: %+v, mock: %+v", sqlDB, mock)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open gorm, error: %+w", err)
	}
	return &findingDB{
		Master: gormDB,
		Slave:  gormDB,
	}, mock, nil
}
