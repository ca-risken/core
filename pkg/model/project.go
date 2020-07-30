package model

import "time"

// Project entity model
type Project struct {
	ProjectID uint32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
