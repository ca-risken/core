package model

import "time"

// Finding entity model
type Finding struct {
	FindingID     uint64 `gorm:"primary_key"`
	Description   string
	DataSource    string
	DataSourceID  string
	ResourceName  string
	ProjectID     uint32
	OriginalScore float32
	Score         float32
	Data          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// FindingTag entity model
type FindingTag struct {
	FindingTagID uint64 `gorm:"primary_key"`
	FindingID    uint64
	ProjectID    uint32
	Tag          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Resource entity model
type Resource struct {
	ResourceID   uint64 `gorm:"primary_key"`
	ResourceName string
	ProjectID    uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ResourceTag entity model
type ResourceTag struct {
	ResourceTagID uint64 `gorm:"primary_key"`
	ResourceID    uint64
	ProjectID     uint32
	Tag           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PendFinding entity model
type PendFinding struct {
	FindingID uint64 `gorm:"primary_key"`
	ProjectID uint32
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FindingSetting entity model
type FindingSetting struct {
	FindingSettingID uint32 `gorm:"primary_key"`
	ProjectID        uint32
	ResourceName     string
	Setting          string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Recommend entity model
type Recommend struct {
	RecommendID    uint32 `gorm:"primary_key"`
	DataSource     string
	Type           string
	Risk           string
	Recommendation string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RecommendFinding entity model
type RecommendFinding struct {
	FindingID   uint64 `gorm:"primary_key"`
	RecommendID uint32
	ProjectID   uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
