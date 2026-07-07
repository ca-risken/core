package model

import "time"

// RemediationProposal status
const (
	RemediationProposalStatusPending   = "pending"
	RemediationProposalStatusSucceeded = "succeeded"
	RemediationProposalStatusFailed    = "failed"
)

// RemediationProposal entity model
type RemediationProposal struct {
	RequestID       string `gorm:"primary_key"`
	FindingID       uint64
	ProjectID       uint32
	Status          string
	ErrorMessage    *string
	RemediationPlan *string
	GeneratedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
