package ai

import (
	"context"
	"errors"
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

func (a *AIService) PutRemediationProposal(ctx context.Context, req *aipb.PutRemediationProposalRequest) (*aipb.PutRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := validatePutRemediationProposal(req); err != nil {
		return nil, err
	}

	statusDetail, remediationPlan, generatedAt := remediationProposalOptionalValues(req)

	var result *model.RemediationProposal
	var err error
	if req.RemediationProposalId == 0 {
		result, err = a.repository.CreateRemediationProposal(ctx, &model.RemediationProposal{
			FindingID:       req.FindingId,
			ProjectID:       req.ProjectId,
			Status:          req.Status,
			StatusDetail:    statusDetail,
			RemediationPlan: remediationPlan,
			GeneratedAt:     generatedAt,
		})
	} else {
		result, err = a.repository.UpdateRemediationProposalStatus(ctx, req.ProjectId,
			req.RemediationProposalId, req.Status, statusDetail, remediationPlan, generatedAt)
	}
	if err != nil {
		return nil, err
	}
	return &aipb.PutRemediationProposalResponse{RemediationProposal: convertRemediationProposal(result)}, nil
}

func validatePutRemediationProposal(req *aipb.PutRemediationProposalRequest) error {
	if req.Status != model.RemediationProposalStatusPending && req.RemediationProposalId == 0 {
		return errors.New("remediation_proposal_id is required")
	}
	return nil
}

func remediationProposalOptionalValues(req *aipb.PutRemediationProposalRequest) (*string, *string, *time.Time) {
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
