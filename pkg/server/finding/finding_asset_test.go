package finding

import (
	"context"
	"testing"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/mock"
)

func TestAssetListService(t *testing.T) {
	cases := []struct {
		name      string
		org, view bool
	}{{"project groups", false, true}, {"org groups", true, true}, {"project expansion", false, false}, {"org expansion", true, false}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := new(mocks.FindingRepository)
			service := &FindingService{repository: repo}
			org := uint32(0)
			if c.org {
				org = 7
			}
			repo.On("ListFindingAssets", mock.Anything, mock.Anything, org).Return([]*finding.FindingAsset{{ProjectId: 1, AssetType: "VM", AssetKey: "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", FindingCount: 2}}, int64(1), uint64(2), nil)
			if c.org {
				if !c.view {
					repo.On("ListFindingForOrg", mock.Anything, mock.Anything).Return(&[]db.FindingWithPend{}, nil)
				}
			}
			if !c.org {
				if !c.view {
					repo.On("ListFinding", mock.Anything, mock.Anything).Return(&[]model.Finding{{FindingID: 9}}, nil)
				}
				res, err := service.ListFinding(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetView: c.view, AssetType: "VM", AssetKey: "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU"})
				if err != nil || res.Total != 1 || (c.view && res.TotalFindingCount != 2) {
					t.Fatalf("res=%v err=%v", res, err)
				}
			}
			if c.org {
				res, err := service.ListFindingForOrg(context.Background(), &finding.ListFindingForOrgRequest{OrganizationId: 7, AssetView: c.view, AssetType: "VM", AssetKey: "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU"})
				if err != nil || res.Total != 1 || (c.view && res.TotalFindingCount != 2) {
					t.Fatalf("res=%v err=%v", res, err)
				}
			}
		})
	}
}
