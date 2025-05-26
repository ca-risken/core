package model

import "time"

// OrganizationRole entity model
type OrganizationRole struct {
	RoleID         uint32 `gorm:"primary_key"`
	OrganizationID uint32
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// OrganizationPolicy entity model
type OrganizationPolicy struct {
	PolicyID       uint32 `gorm:"primary_key"`
	OrganizationID uint32
	Name           string
	ActionPtn      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
