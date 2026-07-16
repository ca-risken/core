package ai

import (
	"context"
	"time"

	"github.com/ca-risken/core/pkg/model"
	aipb "github.com/ca-risken/core/proto/ai"
)

func convertRemediationProposal(r *model.RemediationProposal) *aipb.RemediationProposal {
	data := &aipb.RemediationProposal{
		RemediationProposalId: r.RemediationProposalID,
		FindingId:             r.FindingID,
		ProjectId:             r.ProjectID,
		Status:                r.Status,
		CreatedAt:             r.CreatedAt.Unix(),
		UpdatedAt:             r.UpdatedAt.Unix(),
	}
	if r.StatusDetail != nil {
		data.StatusDetail = *r.StatusDetail
	}
	if r.RemediationPlan != nil {
		data.RemediationPlan = *r.RemediationPlan
	}
	if r.GeneratedAt != nil {
		data.GeneratedAt = r.GeneratedAt.Unix()
	}
	return data
}

func (a *AIService) CreateRemediationProposal(ctx context.Context, req *aipb.CreateRemediationProposalRequest) (*aipb.CreateRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	result, err := a.repository.CreateRemediationProposal(ctx, &model.RemediationProposal{
		FindingID: req.FindingId,
		ProjectID: req.ProjectId,
		Status:    model.RemediationProposalStatusPending,
	})
	if err != nil {
		return nil, err
	}
	return &aipb.CreateRemediationProposalResponse{RemediationProposal: convertRemediationProposal(result)}, nil
}

func (a *AIService) UpdateRemediationProposalStatus(ctx context.Context, req *aipb.UpdateRemediationProposalStatusRequest) (*aipb.UpdateRemediationProposalStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	statusDetail, remediationPlan, generatedAt := remediationProposalOptionalValues(req)

	result, err := a.repository.UpdateRemediationProposalStatus(ctx, req.ProjectId,
		req.RemediationProposalId, req.Status, statusDetail, remediationPlan, generatedAt)
	if err != nil {
		return nil, err
	}
	return &aipb.UpdateRemediationProposalStatusResponse{RemediationProposal: convertRemediationProposal(result)}, nil
}

func remediationProposalOptionalValues(req *aipb.UpdateRemediationProposalStatusRequest) (*string, *string, *time.Time) {
	var statusDetail, remediationPlan *string
	if req.StatusDetail != "" {
		statusDetail = &req.StatusDetail
	}
	if req.RemediationPlan != "" {
		remediationPlan = &req.RemediationPlan
	}
	var generatedAt *time.Time
	if req.GeneratedAt != 0 {
		t := time.Unix(req.GeneratedAt, 0)
		generatedAt = &t
	}
	return statusDetail, remediationPlan, generatedAt
}
