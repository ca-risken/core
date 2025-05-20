package model

import "time"

// Organization entity model
type Organization struct {
	OrganizationID uint32
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
