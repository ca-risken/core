package model

import "time"

// ReportFinding entity model
type ReportFinding struct {
	ReportFindingID uint32 `gorm:"primary_key"`
	ReportDate      string
	ProjectID       uint32
	ProjectName     string
	DataSource      string
	Score           float32
	Count           uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
