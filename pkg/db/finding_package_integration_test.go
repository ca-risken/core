package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/findingremediation"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingPackagesMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	c, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	ctx := context.Background()
	cleanup := func() {
		for _, table := range []string{"finding_tag", "pend_finding", "recommend_finding", "finding"} {
			require.NoError(t, c.Master.Exec("DELETE FROM "+table+" WHERE finding_id BETWEEN 90001 AND 90050").Error)
		}
		require.NoError(t, c.Master.Exec("DELETE FROM recommend WHERE type LIKE 'local-recommend-test-%'").Error)
	}
	cleanup()
	defer cleanup()
	// Two different SCC recommendations, different packages and paths, same method.
	for i := 1; i <= 25; i++ {
		f := &model.Finding{FindingID: uint64(90000 + i), ProjectID: 1, DataSource: "google:scc", DataSourceID: fmt.Sprintf("local-recommend-test-%d", i), ResourceName: "fixture", Score: 0.9, Data: fmt.Sprintf(`{"finding":{"files":[{"path":"/%d/lib/file"}],"vulnerability":{"offending_package":{"package_name":"p%d"}}}}`, i, i), AssetType: "Image", AssetKey: "recommend-test-key", AssetName: "fixture"}
		require.NoError(t, c.BulkUpsertFinding(ctx, []*model.Finding{f}))
		if i <= 23 {
			typ := fmt.Sprintf("local-recommend-test-%d", i%2)
			r, err := c.UpsertRecommend(ctx, &model.Recommend{DataSource: f.DataSource, Type: typ, Risk: "SOFTWARE_VULNERABILITY", Recommendation: "Keep original instructions", RemediationType: findingremediation.Classify(f.DataSource, typ, "SOFTWARE_VULNERABILITY", f.Data)})
			require.NoError(t, err)
			require.NoError(t, c.BulkUpsertRecommendFinding(ctx, []*model.RecommendFinding{{FindingID: f.FindingID, ProjectID: 1, RecommendID: r.RecommendID}}))
		}
	}
	require.NoError(t, c.Master.Exec("INSERT INTO finding_tag(finding_tag_id,finding_id,project_id,tag) VALUES(90001,90001,1,'method-a'),(90002,90001,1,'method-b')").Error)
	require.NoError(t, c.Master.Exec("INSERT INTO pend_finding(finding_id,project_id) VALUES(90002,1)").Error)
	cases := []struct {
		name   string
		tags   []string
		status finding.FindingStatus
		from   float32
		want   uint64
	}{
		{name: "all", want: 25}, {name: "duplicate tags", tags: []string{"method-a", "method-b"}, want: 1},
		{name: "pending", status: finding.FindingStatus_FINDING_PENDING, want: 1}, {name: "active", status: finding.FindingStatus_FINDING_ACTIVE, want: 24}, {name: "score", from: 0.95, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &finding.ListFindingRequest{ProjectId: 1, AssetProjectId: 1, AssetType: "Image", AssetKey: "recommend-test-key", ToScore: 1, FromScore: tc.from, Tag: tc.tags, Status: tc.status, Sort: "finding_id", Direction: "asc", Limit: 20}
			groups, err := c.ListFindingPackages(ctx, req, 0)
			require.NoError(t, err)
			if tc.name == "all" {
				require.Len(t, groups, 2)
				counts := map[string]uint64{}
				for _, g := range groups {
					counts[g.WorkType] += g.FindingCount
				}
				require.Equal(t, uint64(23), counts["library_update"])
				require.Equal(t, uint64(2), counts["review"])
			}
			var total uint64
			for _, g := range groups {
				total += g.FindingCount
				req.PackageKey = g.PackageKey
				var count uint64
				for offset := int32(0); ; offset += 20 {
					req.Offset = offset
					items, err := c.ListFinding(ctx, req)
					require.NoError(t, err)
					require.LessOrEqual(t, len(*items), 20)
					count += uint64(len(*items))
					if len(*items) < 20 {
						break
					}
				}
				require.Equal(t, g.FindingCount, count)
			}
			require.Equal(t, tc.want, total)
		})
	}
	// Bulk relation updates must replace the old recommendation, not just timestamp.
	r, err := c.UpsertRecommend(ctx, &model.Recommend{DataSource: "google:scc", Type: "local-recommend-test-public", RemediationType: "public_access"})
	require.NoError(t, err)
	require.NoError(t, c.BulkUpsertRecommendFinding(ctx, []*model.RecommendFinding{{FindingID: 90001, ProjectID: 1, RecommendID: r.RecommendID}}))
	got, err := c.GetRecommend(ctx, 1, 90001)
	require.NoError(t, err)
	require.Equal(t, r.RecommendID, got.RecommendID)
	// An invalid cross-source relation must not supply another source's method.
	require.NoError(t, c.Master.Exec("UPDATE recommend SET data_source='google:asset' WHERE recommend_id=?", r.RecommendID).Error)
	req := &finding.ListFindingRequest{ProjectId: 1, AssetProjectId: 1, AssetType: "Image", AssetKey: "recommend-test-key", FindingId: 90001, ToScore: 1}
	groups, err := c.ListFindingPackages(ctx, req, 0)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, "review", groups[0].WorkType)
}
