package model

import "time"

// Finding entity model
type Finding struct {
	FindingID     uint64
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
	FindingTagID uint64
	FindingID    uint64
	TagKey       string
	TagValue     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Resource entity model
type Resource struct {
	ResourceID   uint64
	ResourceName string
	ProjectID    uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ResourceTag entity model
type ResourceTag struct {
	ResourceTagID uint64
	ResourceID    uint64
	TagKey        string
	TagValue      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
