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

// OrganizationProject entity model
type OrganizationProject struct {
	OrganizationID uint32
	ProjectID      uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// OrganizationInvitation entity model
type OrganizationInvitation struct {
	OrganizationID uint32
	ProjectID      uint32
	Status         string // PENDING, ACCEPTED, REJECTED
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
