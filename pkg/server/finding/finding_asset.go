package finding

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/ca-risken/core/proto/finding"
)

func orgAssetRequest(r *finding.ListFindingForOrgRequest) *finding.ListFindingRequest {
	return &finding.ListFindingRequest{DataSource: r.DataSource, ResourceName: r.ResourceName, FromScore: r.FromScore, ToScore: r.ToScore, Tag: r.Tag, Sort: r.Sort, Direction: r.Direction, Offset: r.Offset, Limit: r.Limit, Status: r.Status, FindingId: r.FindingId, AssetView: r.AssetView, AssetType: r.AssetType, AssetKey: r.AssetKey, AssetProjectId: r.AssetProjectId, PackageView: r.PackageView, PackageKey: r.PackageKey}
}
func (f *FindingService) listFindingAssetView(ctx context.Context, r *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	if r.PackageView {
		packages, total, count, err := f.listFindingPackages(ctx, r, 0)
		return &finding.ListFindingResponse{Packages: packages, Total: total, Count: uint32(len(packages)), TotalFindingCount: count}, err
	}
	assets, total, totalFindings, err := f.repository.ListFindingAssets(ctx, r, 0)
	if err != nil {
		return nil, err
	}
	if r.AssetView {
		return &finding.ListFindingResponse{Assets: assets, Count: uint32(len(assets)), Total: convertToUint32(total), TotalFindingCount: totalFindings}, nil
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

func (f *FindingService) listFindingPackages(ctx context.Context, r *finding.ListFindingRequest, orgID uint32) ([]*finding.FindingPackage, uint32, uint64, error) {
	groups, err := f.repository.ListFindingPackages(ctx, r, orgID)
	if err != nil {
		return nil, 0, 0, err
	}
	var count uint64
	for _, g := range groups {
		count += g.FindingCount
	}
	type work struct {
		review   bool
		key      string
		count    uint64
		packages []*finding.FindingPackage
	}
	works := []*work{}
	byWork := map[string]*work{}
	for _, g := range groups {
		location := g.WorkLocation
		if g.WorkType == "location" {
			location = ""
		}
		encoded, _ := json.Marshal([]string{g.WorkType, location})
		key := string(encoded)
		w := byWork[key]
		if w == nil {
			w = &work{key: key, review: g.WorkType == ""}
			byWork[key] = w
			works = append(works, w)
		}
		w.count += g.FindingCount
		w.packages = append(w.packages, g)
	}
	sort.Slice(works, func(i, j int) bool {
		if works[i].review != works[j].review {
			return !works[i].review
		}
		if works[i].count != works[j].count {
			return works[i].count > works[j].count
		}
		return works[i].key < works[j].key
	})
	total := len(works)
	start := min(int(r.Offset), total)
	end := min(start+int(r.Limit), total)
	page := []*finding.FindingPackage{}
	for _, w := range works[start:end] {
		page = append(page, w.packages...)
	}
	return page, uint32(total), count, nil
}
