package model

import "time"

// User entity model
type User struct {
	UserID    uint32
	Name      string
	Actevated bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRole entity model
type UserRole struct {
	UserID    uint32
	RoleID    uint32
	ProjectID uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Role entity model
type Role struct {
	RoleID    uint32
	Name      string
	ProjectID uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RolePolicy entity model
type RolePolicy struct {
	RoleID    uint32
	PolicyID  uint32
	ProjectID uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Policy entity model
type Policy struct {
	PolicyID    uint32
	Name        string
	ProjectID   uint32
	ActionPtn   string
	ResourcePtn string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
