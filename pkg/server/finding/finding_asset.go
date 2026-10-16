package finding

import (
	"context"

	"github.com/ca-risken/core/proto/finding"
)

func orgAssetRequest(r *finding.ListFindingForOrgRequest) *finding.ListFindingRequest {
	return &finding.ListFindingRequest{DataSource: r.DataSource, ResourceName: r.ResourceName, FromScore: r.FromScore, ToScore: r.ToScore, Tag: r.Tag, Sort: r.Sort, Direction: r.Direction, Offset: r.Offset, Limit: r.Limit, Status: r.Status, FindingId: r.FindingId, AssetView: r.AssetView, AssetType: r.AssetType, AssetKey: r.AssetKey, AssetProjectId: r.AssetProjectId}
}
func (f *FindingService) listFindingAssetView(ctx context.Context, r *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	assets, total, err := f.repository.ListFindingAssets(ctx, r, 0)
	if err != nil {
		return nil, err
	}
	if r.AssetView {
		return &finding.ListFindingResponse{Assets: assets, Count: uint32(len(assets)), Total: convertToUint32(total)}, nil
	}
	list, err := f.repository.ListFinding(ctx, r)
	if err != nil {
		return nil, err
	}
	ids := []uint64{}
	for _, d := range *list {
		ids = append(ids, d.FindingID)
	}
	return &finding.ListFindingResponse{FindingId: ids, Count: uint32(len(ids)), Total: convertToUint32(total)}, nil
}
