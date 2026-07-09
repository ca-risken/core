package db

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
)

type AIRepository interface {
	// RemediationProposal
	CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error)
	GetRemediationProposal(ctx context.Context, projectID uint32, remediationProposalID uint32) (*model.RemediationProposal, error)
	ListRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, status []string) ([]*model.RemediationProposal, error)
	GetLatestRemediationProposal(ctx context.Context, projectID uint32, findingID uint64) (*model.RemediationProposal, error)
	UpdateRemediationProposalStatus(ctx context.Context, projectID uint32, remediationProposalID uint32, status string, statusDetail, remediationPlan *string, generatedAt *time.Time) (*model.RemediationProposal, error)
	GetActiveRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, createdSince time.Time) (*model.RemediationProposal, error)
}

var _ AIRepository = (*Client)(nil)

func (c *Client) CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error) {
	data.RemediationProposalID = 0
	if err := c.Master.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return c.getRemediationProposalMaster(ctx, data.ProjectID, data.RemediationProposalID)
}

const selectGetRemediationProposal = `select * from remediation_proposal where project_id = ? and remediation_proposal_id = ?`

func (c *Client) GetRemediationProposal(ctx context.Context, projectID uint32, remediationProposalID uint32) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Slave.WithContext(ctx).Raw(selectGetRemediationProposal, projectID, remediationProposalID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) getRemediationProposalMaster(ctx context.Context, projectID uint32, remediationProposalID uint32) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Master.WithContext(ctx).Raw(selectGetRemediationProposal, projectID, remediationProposalID).First(&data).Error; err != nil {
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
	    status_detail = ?,
	    remediation_plan = ?,
	    generated_at = ?
	where project_id = ?
	  and remediation_proposal_id = ?
`

func (c *Client) UpdateRemediationProposalStatus(ctx context.Context, projectID uint32, remediationProposalID uint32, status string, statusDetail, remediationPlan *string, generatedAt *time.Time) (*model.RemediationProposal, error) {
	if err := c.Master.WithContext(ctx).Exec(updateUpdateRemediationProposalStatus,
		status, statusDetail, remediationPlan, generatedAt, projectID, remediationProposalID).Error; err != nil {
		return nil, err
	}
	return c.getRemediationProposalMaster(ctx, projectID, remediationProposalID)
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
