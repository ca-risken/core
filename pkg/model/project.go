package model

import "time"

type Project struct {
	ProjectID uint 32
	Name string
	CreatedAt time.Time
	UpdatedAt time.Time
}
