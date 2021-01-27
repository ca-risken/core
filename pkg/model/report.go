package model

import "time"

// ReportFinding entity model
type ReportFinding struct {
	ReportFindingID uint32 `gorm:"primary_key"`
	ReportDate      string
	Description     string
	Score           float32
	Count           uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
