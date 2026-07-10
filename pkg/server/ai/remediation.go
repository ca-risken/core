package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/ai"
	"gorm.io/gorm"
)

func convertRemediationProposal(r *model.RemediationProposal) *ai.RemediationProposal {
	data := &ai.RemediationProposal{
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

func (a *AIService) GetRemediationProposal(ctx context.Context, req *ai.GetRemediationProposalRequest) (*ai.GetRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetRemediationProposal(ctx, req.ProjectId, req.RemediationProposalId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ai.GetRemediationProposalResponse{}, nil
		}
		return nil, err
	}
	return &ai.GetRemediationProposalResponse{RemediationProposal: convertRemediationProposal(data)}, nil
}

func (a *AIService) ListRemediationProposal(ctx context.Context, req *ai.ListRemediationProposalRequest) (*ai.ListRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListRemediationProposal(ctx, req.ProjectId, req.FindingId, req.Status)
	if err != nil {
		return nil, err
	}
	data := ai.ListRemediationProposalResponse{}
	for _, r := range list {
		data.RemediationProposal = append(data.RemediationProposal, convertRemediationProposal(r))
	}
	return &data, nil
}

func (a *AIService) PutRemediationProposal(ctx context.Context, req *ai.PutRemediationProposalRequest) (*ai.PutRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := validateRemediationProposalForUpsert(req.ProjectId, req.RemediationProposal); err != nil {
		return nil, err
	}

	exists := req.RemediationProposal.RemediationProposalId != 0
	if exists {
		if _, err := a.repository.GetRemediationProposal(ctx, req.ProjectId, req.RemediationProposal.RemediationProposalId); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			exists = false
		}
	}

	var statusDetail, remediationPlan *string
	if req.RemediationProposal.StatusDetail != "" {
		statusDetail = &req.RemediationProposal.StatusDetail
	}
	if req.RemediationProposal.RemediationPlan != "" {
		remediationPlan = &req.RemediationProposal.RemediationPlan
	}
	var generatedAt *time.Time
	if req.RemediationProposal.GeneratedAt != 0 {
		t := time.Unix(req.RemediationProposal.GeneratedAt, 0)
		generatedAt = &t
	}

	var result *model.RemediationProposal
	var err error
	if exists {
		result, err = a.repository.UpdateRemediationProposalStatus(ctx, req.ProjectId,
			req.RemediationProposal.RemediationProposalId, req.RemediationProposal.Status, statusDetail, remediationPlan, generatedAt)
	} else {
		result, err = a.repository.CreateRemediationProposal(ctx, &model.RemediationProposal{
			RemediationProposalID: req.RemediationProposal.RemediationProposalId,
			FindingID:             req.RemediationProposal.FindingId,
			ProjectID:             req.RemediationProposal.ProjectId,
			Status:                req.RemediationProposal.Status,
			StatusDetail:          statusDetail,
			RemediationPlan:       remediationPlan,
			GeneratedAt:           generatedAt,
		})
	}
	if err != nil {
		return nil, err
	}
	return &ai.PutRemediationProposalResponse{RemediationProposal: convertRemediationProposal(result)}, nil
}

func validateRemediationProposalForUpsert(projectID uint32, data *ai.RemediationProposalForUpsert) error {
	if data.FindingId == 0 {
		return errors.New("finding_id is required")
	}
	if data.ProjectId != projectID {
		return fmt.Errorf("project_id mismatch: request=%d, data=%d", projectID, data.ProjectId)
	}
	switch data.Status {
	case model.RemediationProposalStatusPending, model.RemediationProposalStatusSucceeded, model.RemediationProposalStatusFailed:
	default:
		return fmt.Errorf("invalid status: %s", data.Status)
	}
	if data.Status != model.RemediationProposalStatusPending && data.RemediationProposalId == 0 {
		return errors.New("remediation_proposal_id is required")
	}
	return nil
}
