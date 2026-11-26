package finding

import (
	"context"
	"github.com/ca-risken/core/proto/finding"
)

func (f *FindingService) GetFindingCount(ctx context.Context, req *finding.GetFindingCountRequest) (*finding.GetFindingCountResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	toScore := req.ToScore
	if toScore == 0 {
		toScore = 1.0
	}
	total, err := f.repository.ListFindingCount(ctx, req.ProjectId, 0, req.AlertId,
		req.FromScore, toScore, req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	if err != nil {
		return nil, err
	}
	return &finding.GetFindingCountResponse{Total: convertToUint32(total)}, nil
}
