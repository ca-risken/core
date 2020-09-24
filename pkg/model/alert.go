package model

import "time"

// Alert entity model
type Alert struct {
	AlertID           uint64
	AlertConnditionID uint32
	Description       string
	Severity          string
	ProjectID         uint32
	Activated         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AlertHistory entity model
type AlertHistory struct {
	AlertHistoryID uint64
	HistoryType    string
	AlertID        uint32
	Description    string
	Severity       string
	ProjectID      uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RelAlertFinding entity model
type RelAlertFinding struct {
	AlertID   uint32
	FindingID uint32
	ProjectID uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AlertCondition entity model
type AlertCondition struct {
	AlertConditionID uint32
	Description      string
	Severity         string
	ProjectID        uint32
	AndOr            string
	Enabled          bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// AlertRule entity model
type AlertRule struct {
	AlertRuleID  uint32
	Name         string
	ProjectID    uint32
	Score        float32
	ResourceName string
	Tag          string
	FindingCnt   uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AlertCondRule entity model
type AlertCondRule struct {
	AlertConditionID uint32
	AlertRuleID      uint32
	ProjectID        uint32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Notification entity model
type Notification struct {
	NotificationID uint32
	Name           string
	ProjectID      uint32
	Type           string
	NotifySetting  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// AlertCondNotification entity model
type AlertCondNotification struct {
	AlertConditionID uint32
	NotificationID   uint32
	ProjectID        uint32
	CacheSecond      uint32
	NotifiedAt       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
