package model

import "time"

// Project entity model
type Project struct {
	ProjectID uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProjectTag entity model
type ProjectTag struct {
	ProjectID uint32
	Tag       string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
