package model

import "time"

// OrganizationNotification entity model
type OrganizationNotification struct {
	NotificationID uint32 `gorm:"primary_key"`
	Name           string
	OrganizationID uint32
	Type           string
	NotifySetting  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
