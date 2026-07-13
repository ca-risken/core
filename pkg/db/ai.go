package db

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
)

type AIRepository interface {
	// RemediationProposal
	CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error)
	UpdateRemediationProposalStatus(ctx context.Context, projectID uint32, remediationProposalID uint32, status string, statusDetail, remediationPlan *string, generatedAt *time.Time) (*model.RemediationProposal, error)
}

var _ AIRepository = (*Client)(nil)

func (c *Client) CreateRemediationProposal(ctx context.Context, data *model.RemediationProposal) (*model.RemediationProposal, error) {
	if err := c.Master.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return c.getRemediationProposalMaster(ctx, data.ProjectID, data.RemediationProposalID)
}

const selectRemediationProposalByID = `select * from remediation_proposal where project_id = ? and remediation_proposal_id = ?`

func (c *Client) getRemediationProposalMaster(ctx context.Context, projectID uint32, remediationProposalID uint32) (*model.RemediationProposal, error) {
	var data model.RemediationProposal
	if err := c.Master.WithContext(ctx).Raw(selectRemediationProposalByID, projectID, remediationProposalID).First(&data).Error; err != nil {
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
