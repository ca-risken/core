package db

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
)

type AIRepository interface {
	// RemediationProposal
	CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error)
	GetRemediationProposal(ctx context.Context, projectID uint32, requestID string) (*model.RemediationProposal, error)
	ListRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, status []string) ([]*model.RemediationProposal, error)
	GetLatestRemediationProposal(ctx context.Context, projectID uint32, findingID uint64) (*model.RemediationProposal, error)
	UpdateRemediationProposalStatus(ctx context.Context, projectID uint32, requestID, status string, errorMessage, remediationPlan *string, generatedAt *time.Time) (*model.RemediationProposal, error)
	GetActiveRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, createdSince time.Time) (*model.RemediationProposal, error)
}

var _ AIRepository = (*Client)(nil)

const insertCreateRemediationProposal = `
	insert into remediation_proposal (
		request_id,
		finding_id,
		project_id,
		status,
		error_message,
		remediation_plan,
		generated_at
	) values (?, ?, ?, ?, ?, ?, ?)
`

func (c *Client) CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error) {
	if err := c.Master.WithContext(ctx).Exec(insertCreateRemediationProposal,
		data.RequestID, data.FindingID, data.ProjectID, data.Status,
		data.ErrorMessage, data.RemediationPlan, data.GeneratedAt).Error; err != nil {
		return nil, err
	}
	return c.getRemediationProposalMaster(ctx, data.ProjectID, data.RequestID)
}

const selectGetRemediationProposal = `select * from remediation_proposal where project_id = ? and request_id = ?`

func (c *Client) GetRemediationProposal(ctx context.Context, projectID uint32, requestID string) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Slave.WithContext(ctx).Raw(selectGetRemediationProposal, projectID, requestID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) getRemediationProposalMaster(ctx context.Context, projectID uint32, requestID string) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Master.WithContext(ctx).Raw(selectGetRemediationProposal, projectID, requestID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, status []string) ([]*model.RemediationProposal, error) {
	query := `select * from remediation_proposal where project_id = ? and finding_id = ?`
	params := []interface{}{projectID, findingID}
	if len(status) > 0 {
		query += " and status in (?)"
		params = append(params, status)
	}
	query += " order by created_at desc"
	var data []*model.RemediationProposal
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const selectGetLatestRemediationProposal = `
	select *
	from remediation_proposal
	where project_id = ?
	  and finding_id = ?
	order by created_at desc
	limit 1
`

func (c *Client) GetLatestRemediationProposal(ctx context.Context, projectID uint32, findingID uint64) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Slave.WithContext(ctx).Raw(selectGetLatestRemediationProposal, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const updateUpdateRemediationProposalStatus = `
	update remediation_proposal
	set status = ?,
	    error_message = ?,
	    remediation_plan = ?,
	    generated_at = ?
	where project_id = ?
	  and request_id = ?
`

func (c *Client) UpdateRemediationProposalStatus(ctx context.Context, projectID uint32, requestID, status string, errorMessage, remediationPlan *string, generatedAt *time.Time) (*model.RemediationProposal, error) {
	if err := c.Master.WithContext(ctx).Exec(updateUpdateRemediationProposalStatus,
		status, errorMessage, remediationPlan, generatedAt, projectID, requestID).Error; err != nil {
		return nil, err
	}
	return c.getRemediationProposalMaster(ctx, projectID, requestID)
}

const selectGetActiveRemediationProposal = `
	select *
	from remediation_proposal
	where project_id = ?
	  and finding_id = ?
	  and status in (?, ?)
	  and created_at >= ?
	order by created_at desc
	limit 1
`

// GetActiveRemediationProposal returns the latest pending/succeeded proposal created after createdSince.
// Reads from master because the result guards the regeneration cooldown right before a new pending row is created.
func (c *Client) GetActiveRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, createdSince time.Time) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Master.WithContext(ctx).Raw(selectGetActiveRemediationProposal,
		projectID, findingID, model.RemediationProposalStatusPending, model.RemediationProposalStatusSucceeded, createdSince).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
