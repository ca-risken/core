package model

import "time"

// RemediationProposal status
const (
	RemediationProposalStatusPending   = "PENDING"
	RemediationProposalStatusSucceeded = "SUCCEEDED"
	RemediationProposalStatusFailed    = "FAILED"
)

// RemediationProposal entity model
type RemediationProposal struct {
	RemediationProposalID uint32 `gorm:"primary_key"`
	FindingID             uint64
	ProjectID             uint32
	Status                string
	StatusDetail          *string
	RemediationPlan       *string
	GeneratedAt           *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
