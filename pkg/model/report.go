package model

import "time"

// ReportFinding entity model
type ReportFinding struct {
	ReportFindingID uint32 `gorm:"primary_key"`
	ReportDate      string
	Description     string
	Severity        string
	ProjectID       uint32
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
